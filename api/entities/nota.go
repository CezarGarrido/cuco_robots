package entities

import "time"

type Nota struct {
	ID           int64      `json:"id"`
	DisciplinaID int64      `json:"disciplina_id"`
	AlunoID      int64      `json:"aluno_id"`
	Descricao    string     `json:"descricao"`
	Valor        *float64   `json:"valor"`
	CreatedAt    *time.Time `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
}
