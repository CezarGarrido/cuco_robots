package entities

import (
	"time"
)

type Disciplina struct {
	ID        int64
	Descricao string
	IDUEMS    int64
	Ano       int
	CreatedAt time.Time
	UpdatedAt *time.Time
}
