package repository

import (
	"context"
	"database/sql"

	entities "github.com/CezarGarrido/cuco_robots/api/entities"
)

// PostRepo explain...
type ContatoRepo interface {
	//Fetch(ctx context.Context, num int64) ([]*entities.Contato, error)
	GetByAlunoID(ctx context.Context, aluno_id int64) ([]*entities.Contato, error)
	Create(ctx context.Context, p *entities.Contato) (int64, error)
	//Update(ctx context.Context, p *entities.Contato) (*entities.Contato, error)
	//Delete(ctx context.Context, id int64) (bool, error)
	DeleteAll(ctx context.Context, aluno_id int64) (bool, error)
}

// NewSQLAlunoDisciplinaRepo retunrs implement of AlunoDisciplina repository interface
func NewSQLContatoRepo(Conn *sql.DB) ContatoRepo {
	return &postgresContatoRepo{Conn: Conn}
}

type postgresContatoRepo struct {
	Conn *sql.DB
}
func (m *postgresContatoRepo)GetByAlunoID(ctx context.Context, aluno_id int64) ([]*entities.Contato, error){
	query:="SELECT id, aluno_id, tipo, valor, created_at, updated_at FROM cadastros.aluno_contatos WHERE aluno_id IS $1;"
	return m.fetch(ctx,query, aluno_id)
}
func (m *postgresContatoRepo) Create(ctx context.Context, contato *entities.Contato) (int64, error) {
	query := `INSERT INTO cadastros.aluno_contatos (aluno_id, tipo, valor, created_at, updated_at) 
	VALUES($1, $2, $3, $4, $5) RETURNING id`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return -1, err
	}
	var contatoID int64
	err = stmt.QueryRowContext(ctx,
		contato.AlunoID,
		contato.Tipo,
		contato.Valor,
		contato.CreatedAt,
		contato.UpdatedAt,
	).Scan(&contatoID)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()
	return contatoID, nil
}

func (m *postgresContatoRepo) fetch(ctx context.Context, query string, args ...interface{}) ([]*entities.Contato, error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	payload := make([]*entities.Contato, 0)
	for rows.Next() {
		contato := new(entities.Contato)
		err := rows.Scan(
			&contato.ID,
			&contato.AlunoID,
			&contato.Tipo,
			&contato.Valor,
			&contato.CreatedAt,
			&contato.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		payload = append(payload, contato)
	}
	return payload, nil
}
func (m *postgresContatoRepo) DeleteAll(ctx context.Context, aluno_id int64) (bool, error){
	query := "Delete From cadastros.aluno_contatos Where aluno_id=?"
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