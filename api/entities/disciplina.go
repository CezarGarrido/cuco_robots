package entities

import (
	"time"
)

type Disciplina struct {
	ID        int64      `json:"id"`
	UemsID    int64      `json:"uems_id"`
	Descricao string     `json:"descricao"`
	Oferta    string     `json:"oferta"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
