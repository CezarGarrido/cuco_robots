package models

import (
	"database/sql"
	"time"
)

type Nota struct {
	ID           int64
	IDAluno      int64
	IDDisciplina int64
	Documento    string
	CreatedAt    time.Time
	UpdatedAt    *time.Time
}

//Create:Insere o html encontrado
func (this Nota) Create(db *sql.DB) (int64, error) {
	var id int64
	stmt, err := db.Prepare("INSERT INTO cadastros.notas (id_aluno,id_disciplina,documento,created_at) values($1,$2,$3,$4) returning id")
	if err != nil {
		return id, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(this.IDAluno, this.IDDisciplina, this.Documento, this.CreatedAt).Scan(&id)
	if err != nil {
		return id, err
	}
	return id, nil
}

func (this Nota) IsExist(db *sql.DB) bool {
	var count bool
	_ = db.QueryRow("select exists(select 1 from cadastros.notas where id_aluno=$1 and id_disciplina=$2)", this.IDAluno, this.IDDisciplina).Scan(&count)
	return count
}

func (this Nota) Update(db *sql.DB) error {
	stmtIns, err := db.Prepare("UPDATE cadastros.notas set id_aluno=$1, id_disciplina=$2, documento=$3, created_at=$4,updated_at=$5 Where id = $6 ") // ? = placeholder
	if err != nil {
		return err
	}
	defer stmtIns.Close()
	if _, err := stmtIns.Exec(this.IDAluno, this.IDDisciplina, this.Documento, this.CreatedAt, this.UpdatedAt, this.ID); err != nil {
		return err
	}
	return nil
}

func (this *Nota) GetByAluno(id int64, db *sql.DB) error {
	sql := "SELECT id,id_aluno,id_disciplina,documento,created_at,updated_at FROM cadastros.notas where id_aluno=$1 and id_disciplina=$2"
	selDB, err := db.Query(sql, id, this.IDDisciplina)
	if err != nil {
		return err
	}
	for selDB.Next() {
		err = selDB.Scan(
			&this.ID,
			&this.IDAluno,
			&this.IDDisciplina,
			&this.Documento,
			&this.CreatedAt,
			&this.UpdatedAt,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
