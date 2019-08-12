package entities

import "time"

type ValorNota struct {
	ID        int64     `json:"id"`
	IDNota    int64     `json:"id_nota"`
	Descricao string    `json:"descricao"`
	Valor     string    `json:"valor"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
