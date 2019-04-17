package models
import (
	"database/sql"
	"time"
)
type Nota struct {
	Id           int64
	IdAluno      int64
	IdDisciplina int64
	Documento    string
	CreatedAt    time.Time
	UpdatedAt    *time.Time
}
//Create:Insere o html encontrado
func (this Nota) Create(db *sql.DB) error {
	stmt, err := db.Prepare("INSERT INTO cadastros.notas (id_aluno,id_disciplina,documento,created_at) values($1,$2,$3,$4)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(this.IdAluno, this.IdDisciplina, this.Documento, time.Now())
	if err != nil {
		return err
	}
	return nil
}
