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
	"github.com/gorilla/mux"
)

func NewHorario(db *driver.DB) *Horario {
	return &Horario{
		repo: repo.NewSQLHorarioRepo(db.SQL),
	}
}

type Horario struct {
	repo repo.HorarioRepo
}

func (p *Horario) GetByCurso(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	select {
	case <-ctx.Done():
		fmt.Fprint(os.Stderr, "request cancelled\n")
	case <-time.After(2 * time.Second):
		horarios, err := p.repo.GetByCurso(ctx, mux.Vars(r)["curso"])
		if err != nil {
			log.Println(err.Error())
			utils.RespondWithError(w, 500, "Erro interno do sistema")
			return
		}
		utils.RespondwithJSON(w, 200, horarios)
	}
}
