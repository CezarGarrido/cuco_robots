package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/CezarGarrido/cuco_robots/api/driver"
	repo "github.com/CezarGarrido/cuco_robots/api/repository"
	"github.com/CezarGarrido/cuco_robots/api/utils"
	"github.com/CezarGarrido/cuco_robots/crawler"
)

func NewContato(db *driver.DB) *AlunoContato {
	return &AlunoContato{
		repo: repo.NewSQLContatoRepo(db.SQL),
	}
}

type AlunoContato struct {
	repo repo.ContatoRepo
}

func (p *AlunoContato) Fetch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	select {
	case <-ctx.Done():
		fmt.Fprint(os.Stderr, "request cancelled\n")
	case <-time.After(2 * time.Second):
		creds, err := utils.ValidToken(r)
		if err != nil {
			log.Println(err.Error())
			respondWithError(w, 500, err.Error())
			return
		}
		client, err := crawler.NewClient(creds.Aluno.Rgm, creds.Aluno.Senha)
		if err != nil {
			log.Println(err.Error())
			respondWithError(w, 500, err.Error())
			return
		}
		_, _ = client.Logout()
		contatos, err := p.repo.GetByAlunoID(ctx, creds.Aluno.ID)
		if err != nil {
			log.Println(err.Error())
			respondWithError(w, 500, "Erro interno do sistema")
			return
		}
		respondwithJSON(w, 200, contatos)
	}
}
