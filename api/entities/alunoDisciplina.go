package entities

import (
	"time"
)

type AlunoDisciplina struct {
	ID                     int64      `json:"id"`
	AlunoID                int64      `json:"aluno_id"`
	DisciplinaID           int64      `json:"disciplina_id"`
	UemsID                 int64      `json:"uems_id"`
	Ano                    int64      `json:"ano"`
	CargaHorariaPresencial int        `json:"cargaHorariaPresencial"`
	MaximoFaltas           int        `json:"maximoFaltas"`
	PeriodoLetivo          string     `json:"periodoLetivo"`
	Professor              string     `json:"professor"`
	MediaAvaliacoes        float64    `json:"mediaAvaliacoes"`
	Optativa               float64    `json:"optativa"`
	Exame                  float64    `json:"exame"`
	Faltas                 int        `json:"faltas"`
	Situacao               string     `json:"situacao"`
	Notas                  []Nota     `json:"notas"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              *time.Time `json:"updated_at"`
}

/*
type Nota struct {
	ID                     int64       `json:"id"`
	AlunoID                int64       `json:"id_aluno"`
	DisciplinaID           int64       `json:"id_disciplina"`
	//Unidade                string      `json:"unidade"`
	//Curso                  string      `json:"curso"`
	//Disciplina             string      `json:"disciplina"`
	//Turma                  string      `json:"turma"`
	//SerieDisciplina        string      `json:"serieDisciplina"`
	//CargaHorariaPresencial string      `json:"cargaHorariaPresencial"`
	//MaximoFaltas           string      `json:"maximoFaltas"`
    //	PeriodoLetivo          string      `json:"periodoLetivo"`
	//Professor              string      `json:"professor"` //descrição do professor
	//MediaAvaliacoes        string      `json:"mediaAvaliacoes"`
	//Optativa               string      `json:"optativa"`
	//Exame                  string      `json:"exame"`
	//MediaFinal             string      `json:"mediaFinal"`
	//Faltas                 string      `json:"faltas"`
	//Situacao               string      `json:"situacao"`
//	Notas                  []ValorNota `json:"notas"`
	CreatedAt              time.Time   `json:"created_at"`
	UpdatedAt              time.Time   `json:"updated_at"`
}


*/
