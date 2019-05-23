package models

import (
	"database/sql"
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
			&auxDisciplina.ID,
			&auxDisciplina.Descricao,
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
func (this *Disciplina) Create(db *sql.DB) (int64, error) {
	var id int64
	stmtIns, err := db.Prepare("INSERT INTO cadastros.disciplinas (descricao,id_uems,ano,created_at,updated_at) VALUES ($1,$2,$3,$4,$5) RETURNING id")
	if err != nil {
		return id, err
	}
	if err := stmtIns.QueryRow(
		this.Descricao,
		this.IDUEMS,
		this.Ano,
		this.CreatedAt,
		this.UpdatedAt,
	).Scan(&id); err != nil {
		return id, err
	}
	defer stmtIns.Close()
	return id, nil
}

func (this Disciplina) IsExist(db *sql.DB) bool {
	var count bool
	_ = db.QueryRow("select exists(select 1 from cadastros.disciplinas where id_uems=$1 and Descricao=$2)", this.IDUEMS, this.Descricao).Scan(&count)
	return count
}

func (this *Disciplina) GetByIDUEMS(id_uems int64, db *sql.DB) error {
	sql := "SELECT id, descricao, id_uems, ano,created_at,updated_at FROM cadastros.disciplinas where id_uems = $1"
	selDB, err := db.Query(sql, id_uems)
	if err != nil {
		return err
	}
	for selDB.Next() {
		err = selDB.Scan(
			&this.ID,
			&this.Descricao,
			&this.IDUEMS,
			&this.Ano,
			&this.CreatedAt,
			&this.UpdatedAt,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
