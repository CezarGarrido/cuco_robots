package entities

import (
	"time"
)

type Disciplina struct {
	ID        int64      `json:"id"`
	Descricao string     `json:"descricao"`
	UemsID    int64      `json:"uems_id"`
	Ano       int        `json:"ano"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
