package repository

import (
	"context"
	"database/sql"

	entities "github.com/CezarGarrido/cuco_robots/api/entities"
)

// PostRepo explain...
type EnderecoRepo interface {
	//Fetch(ctx context.Context, num int64) ([]*entities.Contato, error)
	GetByAlunoID(ctx context.Context, aluno_id int64) ([]*entities.Endereco, error)
	Create(ctx context.Context, p *entities.Endereco) (int64, error)
	//Update(ctx context.Context, p *entities.Contato) (*entities.Contato, error)
	//Delete(ctx context.Context, id int64) (bool, error)
	DeleteAll(ctx context.Context, aluno_id int64) (bool, error)
}

// NewSQLAlunoDisciplinaRepo retunrs implement of AlunoDisciplina repository interface
func NewSQLEnderecoRepo(Conn *sql.DB) EnderecoRepo {
	return &postgresEnderecoRepo{Conn: Conn}
}

type postgresEnderecoRepo struct {
	Conn *sql.DB
}

func (m *postgresEnderecoRepo) GetByAlunoID(ctx context.Context, aluno_id int64) ([]*entities.Endereco, error) {
	query := `SELECT id, aluno_id, logradouro, numero, complemento, bairro, cep, cidade, created_at, updated_at FROM cadastros.aluno_enderecos WHERE aluno_id IS $1;
	`
	return m.fetch(ctx, query, aluno_id)
}
func (m *postgresEnderecoRepo) Create(ctx context.Context, endereco *entities.Endereco) (int64, error) {
	query := `INSERT INTO cadastros.aluno_enderecos (aluno_id, logradouro, numero, complemento, bairro, cep, cidade, created_at, updated_at) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return -1, err
	}
	var enderecoID int64
	err = stmt.QueryRowContext(ctx,
		endereco.AlunoID,
		endereco.Logradouro,
		endereco.Numero,
		endereco.Complemento,
		endereco.Bairro,
		endereco.CEP,
		endereco.Cidade,
		endereco.CreatedAt,
		endereco.UpdatedAt,
	).Scan(&enderecoID)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()
	return enderecoID, nil
}

func (m *postgresEnderecoRepo) fetch(ctx context.Context, query string, args ...interface{}) ([]*entities.Endereco, error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	payload := make([]*entities.Endereco, 0)
	for rows.Next() {
		endereco := new(entities.Endereco)
		err := rows.Scan(
			&endereco.ID,
			&endereco.AlunoID,
			&endereco.Logradouro,
			&endereco.Numero,
			&endereco.Complemento,
			&endereco.Bairro,
			&endereco.CEP,
			&endereco.Cidade,
			&endereco.CreatedAt,
			&endereco.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		payload = append(payload, endereco)
	}
	return payload, nil
}
func (m *postgresEnderecoRepo) DeleteAll(ctx context.Context, aluno_id int64) (bool, error) {
	query := "Delete From cadastros.aluno_enderecos Where aluno_id=$1"
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
