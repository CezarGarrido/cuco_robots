package entities

import (
	"time"
)

type AlunoDisciplina struct {
	ID                     int64      `json:"id"`
	AlunoID                int64      `json:"aluno_id"`
	UemsID                 int64      `json:"uems_id"`
	Ano                    *int64     `json:"ano"`
	Unidade                *string    `json:"unidade"`
	Curso                  *string    `json:"curso"`
	Disciplina             *string    `json:"disciplina"`
	Turma                  *string    `json:"turma"`
	SerieDisciplina        *string    `json:"serie_disciplina"`
	CargaHorariaPresencial *int       `json:"carga_horaria_presencial"`
	MaximoFaltas           *int       `json:"maximo_faltas"`
	PeriodoLetivo          *string    `json:"periodo_letivo"`
	Professor              *string    `json:"professor"`
	MediaAvaliacoes        *float64   `json:"media_avaliacoes"`
	Optativa               *float64   `json:"optativa"`
	Exame                  *float64   `json:"exame"`
	MediaFinal             *float64   `json:"media_final"`
	Faltas                 *int       `json:"faltas"`
	Situacao               *string    `json:"situacao"`
	Notas                  []*Nota    `json:"notas"`
	CreatedAt              *time.Time `json:"created_at"`
	UpdatedAt              *time.Time `json:"updated_at"`
}
