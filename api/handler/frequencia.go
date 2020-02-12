package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/CezarGarrido/cuco_robots/api/driver"
	entities "github.com/CezarGarrido/cuco_robots/api/entities"
	repo "github.com/CezarGarrido/cuco_robots/api/repository"
	"github.com/CezarGarrido/cuco_robots/api/utils"
	"github.com/CezarGarrido/cuco_robots/crawler"
	"github.com/gorilla/mux"
)

func NewFrequencia(db *driver.DB) *Frequencia {
	return &Frequencia{
		repo:       repo.NewSQLFrequenciaRepo(db.SQL),
		sessao:     repo.NewSQLSessaoRepo(db.SQL),
		disciplina: repo.NewSQLAlunoDisciplinaRepo(db.SQL),
	}
}

type Frequencia struct {
	repo       repo.FrequenciaRepo
	disciplina repo.AlunoDisciplinaRepo
	sessao     repo.SessaoRepo
}

func (p *Frequencia) Fetch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	select {
	case <-ctx.Done():
		fmt.Fprint(os.Stderr, "request cancelled\n")
	case <-time.After(2 * time.Second):
		params := mux.Vars(r)
		disciplinaID, err := strconv.ParseInt(params["id"], 10, 64)
		if err != nil {
			log.Println(err.Error())
			utils.RespondWithError(w, 500, err.Error())
			return
		}
		creds, err := utils.ValidToken(r)
		if err != nil {
			log.Println(err.Error())
			utils.RespondWithError(w, 500, err.Error())
			return
		}
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
			isValidSession, err := client.ValidSession()
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
			disciplinaAux, err := p.disciplina.GetByID(ctx, creds.Aluno.ID, disciplinaID)
			if err != nil {
				log.Println(err.Error())
				utils.RespondWithError(w, 500, err.Error())
				return
			}
			fmt.Println(disciplinaAux.UemsID)
			faltas, err := client.FindFaltas(strconv.FormatInt(disciplinaAux.UemsID, 10))
			if err != nil {
				log.Println(err.Error())
				utils.RespondWithError(w, 500, err.Error())
				return
			}
			hoje := time.Now()
			for _, falta := range faltas {

				for _, frequencia := range falta.Frequencias {
					freq := &entities.Frequencia{
						AlunoID:      creds.Aluno.ID,
						DisciplinaID: disciplinaID,
						Mes:          falta.Mes,
						Valor:        frequencia.Valor,
						CreatedAt:    &hoje,
					}
					freq.Dia, _ = strconv.Atoi(frequencia.Dia)
					fmt.Println(freq)
					okFreq, _ := p.repo.IsExiste(ctx, creds.Aluno.ID, disciplinaID, falta.Mes, freq.Dia, freq.Valor)
					fmt.Println("Existe->", freq)
					if !okFreq {
						fmt.Println("Não Existe->", freq)
						_, err := p.repo.Create(ctx, freq)
						if err != nil {
							log.Println(err.Error())
							utils.RespondWithError(w, 500, err.Error())
							return
						}
					}

				}
			}

		}
		fmt.Println("aqui")
		faltas, err := p.repo.GetByDisciplinaID(ctx, creds.Aluno.ID, disciplinaID)
		if err != nil {
			log.Println(err.Error())
			utils.RespondWithError(w, 500, err.Error())
			return
		}
		utils.RespondwithJSON(w, 200, faltas)
	}
}
