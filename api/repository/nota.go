package repository

import (
	"context"
	"database/sql"

	entities "github.com/CezarGarrido/cuco_robots/api/entities"
)

// NotaRepo explain...
type NotaRepo interface {
	GetByDisciplinaID(ctx context.Context, alunoID, disciplinaID int64) ([]*entities.Nota, error)
}

// NewSQLNotaRepo retunrs implement of nota repository interface
func NewSQLNotaRepo(Conn *sql.DB) NotaRepo {
	return &mysqlNotaRepo{Conn: Conn}
}

type mysqlNotaRepo struct {
	Conn *sql.DB
}

func (m *mysqlNotaRepo) GetByDisciplinaID(ctx context.Context, alunoID, disciplinaID int64) ([]*entities.Nota, error) {
	query := "Select id, aluno_id, disciplina_id, descricao, valor, created_at, updated_at From notas where aluno_id=$1 and disciplina_id=$2"
	return m.fetch(ctx, query, alunoID, disciplinaID)
}

func (m *mysqlNotaRepo) fetch(ctx context.Context, query string, args ...interface{}) ([]*entities.Nota, error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	payload := make([]*entities.Nota, 0)
	for rows.Next() {
		Nota := new(entities.Nota)
		err := rows.Scan(
			&Nota.ID,
			&Nota.AlunoID,
			&Nota.DisciplinaID,
			&Nota.Descricao,
			&Nota.Valor,
			&Nota.CreatedAt,
			&Nota.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		payload = append(payload, Nota)
	}
	return payload, nil
}
