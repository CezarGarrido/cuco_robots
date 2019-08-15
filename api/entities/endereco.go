package entities

import "time"

type Endereco struct {
	ID          int64      `json:"id"`
	AlunoID     int64      `json:"aluno_id"`
	Logradouro  *string    `json:"logradouro"`
	Numero      *int       `json:"numero"`
	Complemento *string    `json:"complemento"`
	Bairro      *string    `json:"bairro"`
	CEP         *string    `json:"cep"`
	Cidade      *string    `json:"cidade"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}
