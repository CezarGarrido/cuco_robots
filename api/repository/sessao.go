package repository

import (
	"context"
	"database/sql"
	"fmt"

	entities "github.com/CezarGarrido/cuco_robots/api/entities"
)

// NotaRepo explain...
type SessaoRepo interface {
	Find(ctx context.Context, aluno_id int64) (*entities.Sessao, bool, error)
	Commit(ctx context.Context, sessao *entities.Sessao) error
}

// NewSQLNotaRepo retunrs implement of nota repository interface
func NewSQLSessaoRepo(Conn *sql.DB) SessaoRepo {
	return &postgressqlRepo{Conn: Conn}
}

type postgressqlRepo struct {
	Conn *sql.DB
}

func (m *postgressqlRepo) Commit(ctx context.Context, sessao *entities.Sessao) error {
	query := `INSERT INTO cadastros.aluno_sessao (aluno_id, qtde_login, qtde_req, cookie_name, cookie_value, created_at, updated_at) VALUES($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (aluno_id) DO UPDATE SET qtde_login = EXCLUDED.qtde_login, qtde_req = EXCLUDED.qtde_req,cookie_name = EXCLUDED.cookie_name, cookie_value = EXCLUDED.cookie_value, updated_at = EXCLUDED.updated_at `
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx,
		sessao.AlunoID,
		sessao.QtdeLogin,
		sessao.QtdeRequest,
		sessao.CookieName,
		sessao.CookieValue,
		sessao.CreatedAt,
		sessao.UpdatedAt,
	)
	if err != nil {
		return err
	}
	defer stmt.Close()
	return nil
}
func (m *postgressqlRepo) Find(ctx context.Context, aluno_id int64) (*entities.Sessao, bool, error) {
	query := `SELECT id, aluno_id, qtde_login, qtde_req, cookie_name, cookie_value, created_at, updated_at from cadastros.aluno_sessao where aluno_id=$1`
	rows := m.Conn.QueryRowContext(ctx,
		query, aluno_id)

	sessao := new(entities.Sessao)

	err := rows.Scan(
		&sessao.ID,
		&sessao.AlunoID,
		&sessao.QtdeLogin,
		&sessao.QtdeRequest,
		&sessao.CookieName,
		&sessao.CookieValue,
		&sessao.CreatedAt,
		&sessao.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, false, err
	} else if err != nil {
		return nil, false, err
	}

	//defer rows.Close()
	fmt.Println(sessao)
	return sessao, true, nil
}
