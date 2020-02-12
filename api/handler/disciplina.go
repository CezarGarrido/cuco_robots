package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/CezarGarrido/cuco_robots/api/driver"
	entities "github.com/CezarGarrido/cuco_robots/api/entities"
	repo "github.com/CezarGarrido/cuco_robots/api/repository"
	"github.com/CezarGarrido/cuco_robots/api/utils"
	"github.com/CezarGarrido/cuco_robots/crawler"
)

func NewAlunoDisciplina(db *driver.DB) *AlunoDisciplina {
	return &AlunoDisciplina{
		repo:   repo.NewSQLAlunoDisciplinaRepo(db.SQL),
		sessao: repo.NewSQLSessaoRepo(db.SQL),
		nota: repo.NewSQLNotaRepo(db.SQL),
	}
}

type AlunoDisciplina struct {
	repo   repo.AlunoDisciplinaRepo
	sessao repo.SessaoRepo
	nota  repo.NotaRepo
}

func (p *AlunoDisciplina) Fetch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	select {
	case <-ctx.Done():
		fmt.Fprint(os.Stderr, "request cancelled\n")
	case <-time.After(2 * time.Second):
		creds, err := utils.ValidToken(r)
		if err != nil {
			log.Println(err.Error())
			utils.RespondWithError(w, 500, err.Error())
			return
		}
		var t interface{}
		sessao, ok, _ := p.sessao.Find(ctx, creds.Aluno.ID)
		if ok {
			var cookies []*http.Cookie
			cookie := &http.Cookie{
				Name:   sessao.CookieName,
				Value:  sessao.CookieValue,
				Path:   "/",
				Domain: "sistemas.uems.br",
			}
			fmt.Println(cookie.Name, cookie.Value)
			cookies = append(cookies, cookie)
			client, err := crawler.NewSetCookieClient(cookies)
			if err != nil {
				log.Println(err.Error())
				utils.RespondWithError(w, 500, err.Error())
				return
			}
			isValidSession, err :=  client.ValidSession()
			if err != nil {
				log.Println(err.Error())
				utils.RespondWithError(w, 500, "Não foi possivel estabelecer uma conexão")
				return
			}
			if !isValidSession {
				client, err = crawler.NewClientCtx(ctx, creds.Rgm, creds.Senha)
				if err != nil {
					log.Println(err.Error())
					utils.RespondWithError(w, 500, err.Error())
					return
				}
				cookie := client.GetCookies()[0]
				fmt.Println(cookie.Value, cookie.Name)
				hoje := time.Now()
				newSessao := &entities.Sessao{
					QtdeLogin:   1,
					QtdeRequest: 1,
					CookieName:  cookie.Name,
					CookieValue: cookie.Value,
					CreatedAt:   hoje,
					UpdatedAt:   &hoje,
				}
				newSessao.AlunoID = creds.ID

				err = p.sessao.Commit(ctx, newSessao)
				if err != nil {
					log.Println(err.Error())
					utils.RespondWithError(w, 500, "Erro interno do sistema")
					return
				}
			}
			aux_disciplinas, err := client.FindDisciplinas()
			if err != nil {
				log.Println(err.Error())
				utils.RespondWithError(w, 500, "Sistema indisponivel")
				return
			}
			fmt.Println(aux_disciplinas)
			for _, disciplina := range aux_disciplinas {
				isExists, err := p.repo.IsExiste(ctx, creds.Aluno.ID, disciplina.UemsID)
				if err != nil {
					log.Println(err.Error())
					utils.RespondWithError(w, 500, "Erro interno do sistema")
					return
				}
				detalhe, err := client.FindNotasByDisciplina(strconv.FormatInt(disciplina.UemsID, 10))
				if err != nil {
					log.Println(err.Error())
					utils.RespondWithError(w, 500, "Erro interno do sistema")
					return
				}

				hoje := time.Now()
				newAlunoDisciplina := &entities.AlunoDisciplina{
					AlunoID:         creds.Aluno.ID,
					UemsID:          disciplina.UemsID,
					Unidade:         &detalhe.Unidade,
					Curso:           &detalhe.Curso,
					Disciplina:      &disciplina.Descricao,
					Turma:           &detalhe.Turma,
					SerieDisciplina: &detalhe.SerieDisciplina,
					PeriodoLetivo:   &detalhe.PeriodoLetivo,
					Professor:       &detalhe.Professor,
					Situacao:        &detalhe.Situacao,
					CreatedAt:       &hoje,
				}
				notas := make([]*entities.Nota, 0)
				for _, nota := range detalhe.Notas {
					newNota := &entities.Nota{
						AlunoID:   creds.Aluno.ID,
						Descricao: nota.Descricao,
						CreatedAt: &hoje,
					}
					valorNormalized := strings.Replace(nota.Valor, ",", ".", -1)
					Valor, _ := strconv.ParseFloat(valorNormalized, 64)
					newNota.Valor = &Valor
					notas = append(notas, newNota)
				}
				newAlunoDisciplina.Notas = notas
				CargaHorariaPresencial, _ := strconv.Atoi(detalhe.CargaHorariaPresencial)
				MaximoFaltas, _ := strconv.Atoi(detalhe.MaximoFaltas)
				Faltas, _ := strconv.Atoi(detalhe.Faltas)
				MediaAvaliacoesNormalized := strings.Replace(detalhe.MediaAvaliacoes, ",", ".", -1)
				MediaFinalNormalized := strings.Replace(detalhe.MediaFinal, ",", ".", -1)
				OptativaNormalized := strings.Replace(detalhe.Optativa, ",", ".", -1)
				ExameNormalized := strings.Replace(detalhe.Exame, ",", ".", -1)
				MediaAvaliacoes, _ := strconv.ParseFloat(MediaAvaliacoesNormalized, 64)
				MediaFinal, _ := strconv.ParseFloat(MediaFinalNormalized, 64)
				Optativa, _ := strconv.ParseFloat(OptativaNormalized, 64)
				Exame, _ := strconv.ParseFloat(ExameNormalized, 64)
				newAlunoDisciplina.CargaHorariaPresencial = &CargaHorariaPresencial
				newAlunoDisciplina.MaximoFaltas = &MaximoFaltas
				newAlunoDisciplina.MediaAvaliacoes = &MediaAvaliacoes
				newAlunoDisciplina.MediaFinal = &MediaFinal
				newAlunoDisciplina.Optativa = &Optativa
				newAlunoDisciplina.Exame = &Exame
				newAlunoDisciplina.Faltas = &Faltas
				if !isExists {
					_, err = p.repo.Create(ctx, newAlunoDisciplina)
					if err != nil {
						log.Println(err.Error())
						utils.RespondWithError(w, 500, "Erro interno do sistema")
						return
					}
				} else {
					disciplinaAnterior, err := p.repo.GetByUemsID(ctx, creds.Aluno.ID, disciplina.UemsID)
					if err != nil {
						log.Println(err.Error())
						utils.RespondWithError(w, 500, "Erro interno do sistema")
						return
					}
					id_aux := disciplinaAnterior.ID
					disciplinaAnterior = newAlunoDisciplina
					disciplinaAnterior.ID = id_aux
					disciplinaAnterior.UpdatedAt = &hoje
					_, err = p.repo.Update(ctx, disciplinaAnterior)
					if err != nil {
						log.Println(err.Error())
						utils.RespondWithError(w, 500, "Erro interno do sistema")
						return
					}
				}

			}
			disciplinas, err := p.repo.GetByAlunoID(ctx, creds.Aluno.ID)
			if err != nil {
				log.Println(err.Error())
				utils.RespondWithError(w, 500, "Erro interno do sistema")
				return
			}
			utils.RespondwithJSON(w, 200, disciplinas)
			return
		} else {
			client, err := crawler.NewClientCtx(ctx, creds.Aluno.Rgm, creds.Aluno.Senha)
			if err != nil {
				log.Println(err.Error())
				utils.RespondWithError(w, 500, err.Error())
				return
			}
			aux_disciplinas, err := client.FindDisciplinas()
			if err != nil {
				log.Println(err.Error())
				utils.RespondWithError(w, 500, "Sistema indisponivel")
				return
			}

			for _, disciplina := range aux_disciplinas {
				isExists, err := p.repo.IsExiste(ctx, creds.Aluno.ID, disciplina.UemsID)
				if err != nil {
					log.Println(err.Error())
					utils.RespondWithError(w, 500, "Erro interno do sistema")
					return
				}
				detalhe, err := client.FindNotasByDisciplina(strconv.FormatInt(disciplina.UemsID, 10))
				if err != nil {
					log.Println(err.Error())
					utils.RespondWithError(w, 500, "Erro interno do sistema")
					return
				}

				hoje := time.Now()
				newAlunoDisciplina := &entities.AlunoDisciplina{
					AlunoID:         creds.Aluno.ID,
					UemsID:          disciplina.UemsID,
					Unidade:         &detalhe.Unidade,
					Curso:           &detalhe.Curso,
					Disciplina:      &disciplina.Descricao,
					Turma:           &detalhe.Turma,
					SerieDisciplina: &detalhe.SerieDisciplina,
					PeriodoLetivo:   &detalhe.PeriodoLetivo,
					Professor:       &detalhe.Professor,
					Situacao:        &detalhe.Situacao,
					CreatedAt:       &hoje,
				}
				notas := make([]*entities.Nota, 0)
				for _, nota := range detalhe.Notas {
					newNota := &entities.Nota{
						AlunoID:   creds.Aluno.ID,
						Descricao: nota.Descricao,
						CreatedAt: &hoje,
					}
					valorNormalized := strings.Replace(nota.Valor, ",", ".", -1)
					Valor, _ := strconv.ParseFloat(valorNormalized, 64)
					newNota.Valor = &Valor
					notas = append(notas, newNota)
				}
				newAlunoDisciplina.Notas = notas
				CargaHorariaPresencial, _ := strconv.Atoi(detalhe.CargaHorariaPresencial)
				MaximoFaltas, _ := strconv.Atoi(detalhe.MaximoFaltas)
				Faltas, _ := strconv.Atoi(detalhe.Faltas)

				MediaAvaliacoesNormalized := strings.Replace(detalhe.MediaAvaliacoes, ",", ".", -1)
				MediaFinalNormalized := strings.Replace(detalhe.MediaFinal, ",", ".", -1)
				OptativaNormalized := strings.Replace(detalhe.Optativa, ",", ".", -1)
				ExameNormalized := strings.Replace(detalhe.Exame, ",", ".", -1)

				MediaAvaliacoes, _ := strconv.ParseFloat(MediaAvaliacoesNormalized, 64)
				MediaFinal, _ := strconv.ParseFloat(MediaFinalNormalized, 64)
				Optativa, _ := strconv.ParseFloat(OptativaNormalized, 64)
				Exame, _ := strconv.ParseFloat(ExameNormalized, 64)
				newAlunoDisciplina.CargaHorariaPresencial = &CargaHorariaPresencial
				newAlunoDisciplina.MaximoFaltas = &MaximoFaltas
				newAlunoDisciplina.MediaAvaliacoes = &MediaAvaliacoes
				newAlunoDisciplina.MediaFinal = &MediaFinal
				newAlunoDisciplina.Optativa = &Optativa
				newAlunoDisciplina.Exame = &Exame
				newAlunoDisciplina.Faltas = &Faltas
				if !isExists {
					_, err = p.repo.Create(ctx, newAlunoDisciplina)
					if err != nil {
						log.Println(err.Error())
						utils.RespondWithError(w, 500, "Erro interno do sistema")
						return
					}
				} else {
					disciplinaAnterior, err := p.repo.GetByUemsID(ctx, creds.Aluno.ID, disciplina.UemsID)
					if err != nil {
						log.Println(err.Error())
						utils.RespondWithError(w, 500, "Erro interno do sistema")
						return
					}
					id_aux := disciplinaAnterior.ID
					disciplinaAnterior = newAlunoDisciplina
					disciplinaAnterior.ID = id_aux
					disciplinaAnterior.UpdatedAt = &hoje
					_, err = p.repo.Update(ctx, disciplinaAnterior)
					if err != nil {
						log.Println(err.Error())
						utils.RespondWithError(w, 500, "Erro interno do sistema")
						return
					}
				}
			}
			fmt.Println("disciplinas")
			disciplinas, err := p.repo.GetByAlunoID(ctx, creds.Aluno.ID)
			if err != nil {
				log.Println(err.Error())
				utils.RespondWithError(w, 500, "Erro interno do sistema")
				return
			}
			fmt.Println(disciplinas)
			t = disciplinas
		}
		utils.RespondwithJSON(w, 200, t)
	}
}
