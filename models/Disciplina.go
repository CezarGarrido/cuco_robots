package models

import (
	"database/sql"
	"time"
)

type Disciplina struct {
	Id        int64
	Descricao string
	IdUems    int64
	Ano       int
	CreatedAt time.Time
	UpdatedAt *time.Time
}
//GetAll: Busca todas as disciplinas cadastradas no sistema
func (this Disciplina) GetAll(db *sql.DB) ([]*Disciplina, error) {
	disciplinas := make([]*Disciplina, 0)
	selDB, err := db.Query("Select * from cadastros.alunos")
	if err != nil {
		return disciplinas, err
	}
	for selDB.Next() {
		auxDisciplina := new(Disciplina)
		err = selDB.Scan(
			&auxDisciplina.Id,
			&auxDisciplina.Descricao,
			&auxDisciplina.IdUems,
			&auxDisciplina.Ano,
			&auxDisciplina.CreatedAt,
			&auxDisciplina.UpdatedAt,
		)
		if err != nil {
			return disciplinas, err
		}
		disciplinas = append(disciplinas, auxDisciplina)
	}
	return disciplinas, err
}
