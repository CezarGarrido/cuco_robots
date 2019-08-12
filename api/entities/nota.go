package entities

import (
	"time"
)

type Nota struct {
	ID                     int64       `json:"id"`
	IDAluno                int64       `json:"id_aluno"`
	IDDisciplina           int64       `json:"id_disciplina"`
	Unidade                string      `json:"unidade"`
	Curso                  string      `json:"curso"`
	Disciplina             string      `json:"disciplina"`
	Turma                  string      `json:"turma"`
	SerieDisciplina        string      `json:"serieDisciplina"`
	CargaHorariaPresencial string      `json:"cargaHorariaPresencial"`
	MaximoFaltas           string      `json:"maximoFaltas"`
	PeriodoLetivo          string      `json:"periodoLetivo"`
	Professor              string      `json:"professor"`
	MediaAvaliacoes        string      `json:"mediaAvaliacoes"`
	Optativa               string      `json:"optativa"`
	Exame                  string      `json:"exame"`
	MediaFinal             string      `json:"mediaFinal"`
	Faltas                 string      `json:"faltas"`
	Situacao               string      `json:"situacao"`
	Notas                  []ValorNota `json:"notas"`
	CreatedAt              time.Time   `json:"created_at"`
	UpdatedAt              time.Time   `json:"updated_at"`
}
