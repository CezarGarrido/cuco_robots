package entities

import "time"

type Contato struct {
	ID        int64      `json:"id"`
	AlunoID   int64      `json:"aluno_id"`
	Tipo      string     `json:"tipo"`
	Valor     *string    `json:"valor"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
