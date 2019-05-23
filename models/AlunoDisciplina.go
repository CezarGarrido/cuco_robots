package models

import (
	"database/sql"
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

//GetAll: Busca todas as disciplinas cadastradas no sistema
func (this AlunoDisciplina) GetAll(db *sql.DB) ([]*AlunoDisciplina, error) {
	disciplinas := make([]*AlunoDisciplina, 0)
	selDB, err := db.Query("Select * from cadastros.alunos")
	if err != nil {
		return disciplinas, err
	}
	for selDB.Next() {
		auxDisciplina := new(AlunoDisciplina)
		err = selDB.Scan(
			&auxDisciplina.ID,
			&auxDisciplina.IDAluno,
			&auxDisciplina.IDDisciplina,
			&auxDisciplina.IDUEMS,
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
func (this *AlunoDisciplina) Create(db *sql.DB) (int64, error) {
	var id int64
	stmtIns, err := db.Prepare("INSERT INTO cadastros.alunos_disciplinas (id_aluno,id_disciplina,id_uems,created_at,updated_at) VALUES ($1,$2,$3,$4,$5) RETURNING id")
	if err != nil {
		return id, err
	}
	if err := stmtIns.QueryRow(
		this.IDAluno,
		this.IDDisciplina,
		this.IDUEMS,
		this.CreatedAt,
		this.UpdatedAt,
	).Scan(&id); err != nil {
		return id, err
	}
	defer stmtIns.Close()
	return id, nil
}

func (this AlunoDisciplina) IsExist(db *sql.DB) bool {
	var count bool
	_ = db.QueryRow("select exists(select 1 from cadastros.alunos_disciplinas where id_disciplina=$1)", this.IDDisciplina).Scan(&count)
	return count
}
