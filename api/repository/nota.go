package repository

import (
	"context"
	"database/sql"
	"errors"

	entities "github.com/CezarGarrido/cuco_robots/api/entities"
)

// NotaRepo explain...
type NotaRepo interface {
	GetByDisciplinaID(ctx context.Context, alunoID, disciplinaID int64) ([]*entities.Nota, error)
	Create(ctx context.Context, nota *entities.Nota) (int64, error)
	IsExiste(ctx context.Context, aluno_id, aluno_disciplina_id int64, descricao string) (bool, error)
	GetByDescricao(ctx context.Context, alunoID, disciplinaID int64, descricao string) (*entities.Nota, error)
	Update(ctx context.Context, nota *entities.Nota) (*entities.Nota, error)
}

// NewSQLNotaRepo retunrs implement of nota repository interface
func NewSQLNotaRepo(Conn *sql.DB) NotaRepo {
	return &mysqlNotaRepo{Conn: Conn}
}

type mysqlNotaRepo struct {
	Conn *sql.DB
}

func (m *mysqlNotaRepo) Create(ctx context.Context, nota *entities.Nota) (int64, error) {
	query := `INSERT INTO cadastros.aluno_notas (aluno_id, disciplina_id, descricao, valor, created_at, updated_at) VALUES($1, $2, $3, $4, $5, $6) RETURNING id`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return -1, err
	}
	var contatoID int64
	err = stmt.QueryRowContext(ctx,
		nota.AlunoID,
		nota.DisciplinaID,
		nota.Descricao,
		nota.Valor,
		nota.CreatedAt,
		nota.UpdatedAt,
	).Scan(&contatoID)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()
	return contatoID, nil
}

func (m *mysqlNotaRepo) GetByDisciplinaID(ctx context.Context, alunoID, disciplinaID int64) ([]*entities.Nota, error) {
	query := "Select id, aluno_id, disciplina_id, descricao, valor, created_at, updated_at From cadastros.aluno_notas  where aluno_id=$1 and disciplina_id=$2"
	return m.fetch(ctx, query, alunoID, disciplinaID)
}

func (m *mysqlNotaRepo) GetByDescricao(ctx context.Context, alunoID, disciplinaID int64, descricao string) (*entities.Nota, error) {
	query := "Select id, aluno_id, disciplina_id, descricao, valor, created_at, updated_at From cadastros.aluno_notas where aluno_id=$1 and disciplina_id=$2 and descricao=$3"
	rows, err := m.fetch(ctx, query, alunoID, disciplinaID, descricao)
	if err != nil {
		return nil, err
	}
	payload := &entities.Nota{}
	if len(rows) > 0 {
		payload = rows[0]
	} else {
		return nil, errors.New("Aluno não encontrado")
	}
	return payload, nil
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

func (m *mysqlNotaRepo) IsExiste(ctx context.Context, aluno_id, aluno_disciplina_id int64, descricao string) (bool, error) {
	var exist = false
	err := m.Conn.QueryRowContext(ctx, "Select exists(Select 1 from cadastros.aluno_notas where aluno_id=$1 and disciplina_id=$2 and descricao=$3);", aluno_id, aluno_disciplina_id, descricao).Scan(
		&exist,
	)
	if err != nil {
		return false, err
	}
	return exist, nil
}

func (m *mysqlNotaRepo) Update(ctx context.Context, nota *entities.Nota) (*entities.Nota, error) {
	query := `UPDATE cadastros.aluno_notas SET aluno_id=$1, disciplina_id=$2, descricao=$3, valor=$4, updated_at=$5 WHERE id=$6;`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	_, err = stmt.ExecContext(
		ctx,
		nota.AlunoID,
		nota.DisciplinaID,
		nota.Descricao,
		nota.Valor,
		nota.UpdatedAt,
		nota.ID,
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return nota, nil
}
