package entities

import (
	"time"
)

type Aluno struct {
	ID        int64      `json:"id"`
	Nome      string     `json:"nome"`
	Curso     string     `json:"curso"`
	Ano       int        `json:"ano"`
	Unidade   string     `json:"unidade"`
	Rgm       string     `json:"rgm"`
	Senha     string     `json:"senha"`
	Email     *string    `json:"email"`
	Telefone  *string    `json:"telefone"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
