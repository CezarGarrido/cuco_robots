package entities

import "time"

type Frequencia struct {
	ID           int64      `json:"id"`
	AlunoID      int64      `json:"aluno_id"`
	DisciplinaID int64      `json:"disciplina_id"`
	Mes          string     `json:"mes"`
	Dia          int        `json:"dia"`
	Valor        string     `json:"valor"`
	CreatedAt    *time.Time `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
}
