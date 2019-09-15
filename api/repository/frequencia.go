package repository

import (
	"context"
	"database/sql"
	"errors"

	entities "github.com/CezarGarrido/cuco_robots/api/entities"
)

// FrequenciaRepo explain...
type FrequenciaRepo interface {
	GetByDisciplinaID(ctx context.Context, alunoID, disciplinaID int64) ([]*entities.Frequencia, error)
	Create(ctx context.Context, frequencia *entities.Frequencia) (int64, error)
	IsExiste(ctx context.Context, aluno_id, aluno_disciplina_id int64, mes string, dia int, valor string) (bool, error)
	GetByDescricao(ctx context.Context, alunoID, disciplinaID int64, descricao string) (*entities.Frequencia, error)
	Update(ctx context.Context, frequencia *entities.Frequencia) (*entities.Frequencia, error)
}

// NewSQLNotaRepo retunrs implement of frequencia repository interface
func NewSQLFrequenciaRepo(Conn *sql.DB) FrequenciaRepo {
	return &mysqlFrequenciaRepo{Conn: Conn}
}

type mysqlFrequenciaRepo struct {
	Conn *sql.DB
}

/*
type Frequencia struct {
	ID           int64
	AlunoID      int64
	DisciplinaID int64
	Mes          string
	Dia          int
	Valor        string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
*/
func (m *mysqlFrequenciaRepo) Create(ctx context.Context, frequencia *entities.Frequencia) (int64, error) {
	query := `INSERT INTO cadastros.aluno_frequencia (aluno_id, disciplina_id, mes, dia,valor, created_at, updated_at) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return -1, err
	}
	var frequenciaID int64
	err = stmt.QueryRowContext(ctx,
		frequencia.AlunoID,
		frequencia.DisciplinaID,
		frequencia.Mes,
		frequencia.Dia,
		frequencia.Valor,
		frequencia.CreatedAt,
		frequencia.UpdatedAt,
	).Scan(&frequenciaID)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()
	return frequenciaID, nil
}

func (m *mysqlFrequenciaRepo) GetByDisciplinaID(ctx context.Context, alunoID, disciplinaID int64) ([]*entities.Frequencia, error) {
	query := "Select id, aluno_id, disciplina_id, mes, dia, valor, created_at, updated_at From cadastros.aluno_frequencia  where aluno_id=$1 and disciplina_id=$2"
	return m.fetch(ctx, query, alunoID, disciplinaID)
}

func (m *mysqlFrequenciaRepo) GetByDescricao(ctx context.Context, alunoID, disciplinaID int64, descricao string) (*entities.Frequencia, error) {
	query := "Select id, aluno_id, disciplina_id, mes, dia, valor, created_at, updated_at From cadastros.aluno_frequencia where aluno_id=$1 and disciplina_id=$2 and descricao=$3"
	rows, err := m.fetch(ctx, query, alunoID, disciplinaID, descricao)
	if err != nil {
		return nil, err
	}
	payload := &entities.Frequencia{}
	if len(rows) > 0 {
		payload = rows[0]
	} else {
		return nil, errors.New("Aluno não encontrado")
	}
	return payload, nil
}

func (m *mysqlFrequenciaRepo) fetch(ctx context.Context, query string, args ...interface{}) ([]*entities.Frequencia, error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	payload := make([]*entities.Frequencia, 0)
	for rows.Next() {
		Frequencia := new(entities.Frequencia)
		err := rows.Scan(
			&Frequencia.ID,
			&Frequencia.AlunoID,
			&Frequencia.DisciplinaID,
			&Frequencia.Mes,
			&Frequencia.Dia,
			&Frequencia.Valor,
			&Frequencia.CreatedAt,
			&Frequencia.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		payload = append(payload, Frequencia)
	}
	return payload, nil
}

func (m *mysqlFrequenciaRepo) IsExiste(ctx context.Context, aluno_id, aluno_disciplina_id int64, mes string, dia int, valor string) (bool, error) {
	var count = 0
	err := m.Conn.QueryRowContext(ctx, "Select count(*) from cadastros.aluno_frequencia where aluno_id=$1 and disciplina_id=$2 and mes=$3 and dia=$4 and valor=$5;", aluno_id, aluno_disciplina_id, mes, dia, valor).Scan(
		&count,
	)
	if err != nil {
		return false, err
	}
	if count < 4 {
		return false, err
	}
	return true, nil
}

func (m *mysqlFrequenciaRepo) Update(ctx context.Context, frequencia *entities.Frequencia) (*entities.Frequencia, error) {
	query := `UPDATE cadastros.aluno_frequencia SET aluno_id=$1, disciplina_id=$2, descricao=$3, valor=$4, updated_at=$5 WHERE id=$6;`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	_, err = stmt.ExecContext(
		ctx,
		frequencia.AlunoID,
		frequencia.DisciplinaID,
		frequencia.Mes,
		frequencia.Dia,
		frequencia.Valor,
		frequencia.UpdatedAt,
		frequencia.ID,
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return frequencia, nil
}
