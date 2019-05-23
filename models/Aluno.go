package models

import (
	"database/sql"
	"time"
)

type Aluno struct {
	ID        int64
	Nome      string
	Curso     string
	Ano       int
	Unidade   string
	Rgm       string
	Senha     string
	CreatedAt time.Time
	UpdatedAt *time.Time
}

//GetAll: Busca todos os alunos cadastrados no sistema
func (this Aluno) GetAll(db *sql.DB) ([]*Aluno, error) {
	alunos := make([]*Aluno, 0)
	selDB, err := db.Query("Select * from cadastros.alunos")
	if err != nil {
		return alunos, err
	}
	for selDB.Next() {
		auxAluno := new(Aluno)
		err = selDB.Scan(
			&auxAluno.ID,
			&auxAluno.Nome,
			&auxAluno.Curso,
			&auxAluno.Ano,
			&auxAluno.Unidade,
			&auxAluno.Rgm,
			&auxAluno.Senha,
			&auxAluno.CreatedAt,
			&auxAluno.UpdatedAt,
		)
		if err != nil {
			return alunos, err
		}
		alunos = append(alunos, auxAluno)
	}
	return alunos, err
}
