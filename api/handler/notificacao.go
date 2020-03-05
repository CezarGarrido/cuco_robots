package handler

import (
	"context"
	"fmt"
	"github.com/CezarGarrido/cuco_robots/api/driver"
	entities "github.com/CezarGarrido/cuco_robots/api/entities"
	repo "github.com/CezarGarrido/cuco_robots/api/repository"
	"github.com/CezarGarrido/cuco_robots/api/utils"
	"github.com/CezarGarrido/cuco_robots/crawler"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func NewNotificacao(db *driver.DB) *Notificacao {
	return &Notificacao{
		sessao:     repo.NewSQLSessaoRepo(db.SQL),
		repo:       repo.NewSQLNotificacaoRepo(db.SQL),
		aluno:      repo.NewSQLAlunoRepo(db.SQL),
		disciplina: repo.NewSQLAlunoDisciplinaRepo(db.SQL),
		nota:       repo.NewSQLNotaRepo(db.SQL),
	}
}

type Notificacao struct {
	repo       repo.NotificacaoRepo
	sessao     repo.SessaoRepo
	aluno      repo.AlunoRepo
	disciplina repo.AlunoDisciplinaRepo
	nota       repo.NotaRepo
}

func (this *Notificacao) StartVascullhador(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	select {
	case <-ctx.Done():
		fmt.Fprint(os.Stderr, "request cancelled\n")
	case <-time.After(3 * time.Second):
		go this.startCrawler()
		utils.RespondwithJSON(w, 200, "OK")
	}
}

func (this *Notificacao) startCrawler() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("recovered error %v", r)
		}
	}()

	ctx := context.Background()
	log.Println("# iniciando vasculhador")
	alunos, err := this.aluno.GetAll(ctx)
	if err != nil {
		log.Println(err.Error())
		return
	}
	for _, aluno := range alunos {
		log.Println("# Buscando notas do aluno ->", aluno.ID, aluno.Nome)
		sessao, ok, _ := this.sessao.Find(ctx, aluno.ID)
		if ok {
			var cookies []*http.Cookie
			cookie := &http.Cookie{
				Name:   sessao.CookieName,
				Value:  sessao.CookieValue,
				Path:   "/",
				Domain: "sistemas.uems.br",
			}
			cookies = append(cookies, cookie)
			client, err := crawler.NewSetCookieClient(cookies)
			if err != nil {
				log.Println(err.Error())
				return
			}
			isValidSession, err := client.ValidSession()
			if err != nil {
				log.Println(err.Error())
				return
			}
			if !isValidSession {
				client, err = crawler.NewClientCtx(ctx, aluno.Rgm, aluno.Senha)
				if err != nil {
					log.Println(err.Error())
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
				newSessao.AlunoID = aluno.ID

				err = this.sessao.Commit(ctx, newSessao)
				if err != nil {
					log.Println(err.Error())
					return
				}
			}
			aux_disciplinas, err := client.FindDisciplinas()
			if err != nil {
				log.Println(err.Error())
				return
			}
			for _, disciplina := range aux_disciplinas {
				log.Println("# Buscando notas da disciplina ->", disciplina.UemsID, disciplina.Descricao)
				isExists, err := this.disciplina.IsExiste(ctx, aluno.ID, disciplina.UemsID)
				if err != nil {
					log.Println(err.Error())
					return
				}
				detalhe, err := client.FindNotasByDisciplina(strconv.FormatInt(disciplina.UemsID, 10))
				if err != nil {
					log.Println(err.Error())
					return
				}
				fmt.Println("# Detalhes ->", detalhe)
				hoje := time.Now()
				newAlunoDisciplina := &entities.AlunoDisciplina{
					AlunoID:         aluno.ID,
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
						AlunoID:   aluno.ID,
						Descricao: nota.Descricao,
						CreatedAt: &hoje,
					}
					if nota.Valor != "" {
						valorNormalized := strings.Replace(nota.Valor, ",", ".", -1)
						Valor, _ := strconv.ParseFloat(valorNormalized, 64)
						newNota.Valor = &Valor
					}
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
					_, err = this.disciplina.Create(ctx, newAlunoDisciplina)
					if err != nil {
						log.Println(err.Error())
						return
					}
				} else {
					disciplinaAnterior, err := this.disciplina.GetByUemsID(ctx, aluno.ID, disciplina.UemsID)
					if err != nil {
						log.Println(err.Error())
						return
					}
					id_aux := disciplinaAnterior.ID
					disciplinaAnterior = newAlunoDisciplina
					disciplinaAnterior.ID = id_aux
					disciplinaAnterior.UpdatedAt = &hoje
					_, err = this.disciplina.Update(ctx, disciplinaAnterior)
					if err != nil {
						log.Println(err.Error())
						return
					}

					for _, nota := range newAlunoDisciplina.Notas {
						nota.DisciplinaID = disciplinaAnterior.ID
						exist, err := this.nota.IsExiste(ctx, disciplinaAnterior.AlunoID, disciplinaAnterior.ID, nota.Descricao)
						if !exist {
							_, err = this.nota.Create(ctx, nota)
							if err != nil {
								log.Println(err.Error())
								return
							}
						} else {
							notaAnterior, err := this.nota.GetByDescricao(ctx, disciplinaAnterior.AlunoID, disciplinaAnterior.ID, nota.Descricao)
							if err != nil {
								log.Println(err.Error())
								return
							}
							if nota.Valor != nil {
								if notaAnterior.Valor != nota.Valor {
									log.Println("# Nova nota lançada:", nota.Descricao, nota.Valor)
									notaAnterior.Valor = nota.Valor
									this.enviaNotificacao(*disciplinaAnterior.Disciplina+" - Nova nota lançada.", "Valor: "+fmt.Sprintf("%.2f", *nota.Valor))
								}
								notaAnterior.Valor = nota.Valor
							}
							h := time.Now()
							notaAnterior.UpdatedAt = &h
							_, err = this.nota.Update(ctx, notaAnterior)
							if err != nil {
								log.Println(err.Error())
								return
							}
						}
					}
				}
			}
		}
	}
	log.Println("# fim da execução")
}

func (this *Notificacao) enviaNotificacao(titulo, body string) {
	log.Println("# Enviando notificação ->", titulo, body)
}
