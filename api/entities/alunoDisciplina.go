package entities

import (
	"time"
)

type AlunoDisciplina struct {
	ID           int64
	IDAluno      int64
	IDDisciplina int64
	IDUEMS       int64
	Ano          int64
	CreatedAt    time.Time
	UpdatedAt    *time.Time
}
