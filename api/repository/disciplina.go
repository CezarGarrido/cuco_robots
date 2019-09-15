package repository

import (
	"context"
	"database/sql"
	"errors"

	entities "github.com/CezarGarrido/cuco_robots/api/entities"
)

// PostRepo explain...
type DisciplinaRepo interface {
	//Fetch(ctx context.Context, num int64) ([]*entities.Contato, error)
	GetByUemsID(ctx context.Context, aluno_id int64) (*entities.Disciplina, error)
	GetByID(ctx context.Context, aluno_id, id int64) (*entities.Disciplina, error)
	Create(ctx context.Context, p *entities.Disciplina) (int64, error)
	//Update(ctx context.Context, p *entities.Contato) (*entities.Contato, error)
	//Delete(ctx context.Context, id int64) (bool, error)
	IsExiste(ctx context.Context, uems_id int64) (bool, error)
	DeleteAll(ctx context.Context, aluno_id int64) (bool, error)
}

// NewSQLAlunoDisciplinaRepo retunrs implement of AlunoDisciplina repository interface
func NewSQLDisciplinaRepo(Conn *sql.DB) DisciplinaRepo {
	return &postgresDisciplinaRepo{Conn: Conn}
}

type postgresDisciplinaRepo struct {
	Conn *sql.DB
}

func (m *postgresDisciplinaRepo) GetByID(ctx context.Context, aluno_id, id int64) (*entities.Disciplina, error) {
	query := "SELECT id, uems_id, descricao, oferta, created_at, updated_at FROM cadastros.disciplinas WHERE aluno_id=$1 and id=$2"
	rows, err := m.fetch(ctx, query, aluno_id, id)
	if err != nil {
		return nil, err
	}
	payload := &entities.Disciplina{}
	if len(rows) > 0 {
		payload = rows[0]
	} else {
		return nil, errors.New("Disciplina não encontada.")
	}
	return payload, nil
}
func (m *postgresDisciplinaRepo) GetByUemsID(ctx context.Context, uems_id int64) (*entities.Disciplina, error) {
	query := "SELECT id, uems_id, descricao, oferta, created_at, updated_at FROM cadastros.disciplinas WHERE id IS $1"
	rows, err := m.fetch(ctx, query, uems_id)
	if err != nil {
		return nil, err
	}
	payload := &entities.Disciplina{}
	if len(rows) > 0 {
		payload = rows[0]
	} else {
		return nil, errors.New("Disciplina não encontada.")
	}
	return payload, nil
}

func (m *postgresDisciplinaRepo) Create(ctx context.Context, disciplina *entities.Disciplina) (int64, error) {
	query := `INSERT INTO cadastros.disciplinas (uems_id, descricao, oferta, created_at, updated_at) VALUES($1, $2, $3, $4, $5) RETURNING id`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return -1, err
	}
	var enderecoID int64
	err = stmt.QueryRowContext(ctx,
		disciplina.UemsID,
		disciplina.Descricao,
		disciplina.Oferta,
		disciplina.CreatedAt,
		disciplina.UpdatedAt,
	).Scan(&enderecoID)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()
	return enderecoID, nil
}

func (m *postgresDisciplinaRepo) fetch(ctx context.Context, query string, args ...interface{}) ([]*entities.Disciplina, error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	payload := make([]*entities.Disciplina, 0)
	for rows.Next() {
		disciplina := new(entities.Disciplina)
		err := rows.Scan(
			&disciplina.ID,
			&disciplina.UemsID,
			&disciplina.Descricao,
			&disciplina.Oferta,
			&disciplina.CreatedAt,
			&disciplina.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		payload = append(payload, disciplina)
	}
	return payload, nil
}
func (m *postgresDisciplinaRepo) DeleteAll(ctx context.Context, aluno_id int64) (bool, error) {
	query := "Delete From cadastros.disciplinas Where aluno_id=?"
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return false, err
	}
	_, err = stmt.ExecContext(ctx, aluno_id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (m *postgresDisciplinaRepo) IsExiste(ctx context.Context, uems_id int64) (bool, error) {
	var exist = false
	err := m.Conn.QueryRowContext(ctx, "Select exists(Select 1 from cadastros.disciplinas where uems_id=$1);", uems_id).Scan(
		&exist,
	)
	if err != nil {
		return false, err
	}
	return exist, nil
}
