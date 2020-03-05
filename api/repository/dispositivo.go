package repository

import (
	"context"
	"database/sql"

	entities "github.com/CezarGarrido/cuco_robots/api/entities"
)

type DispositivoRepo interface {
	Find(ctx context.Context, aluno_id int64) (*entities.Dispositivo, bool, error)
	Commit(ctx context.Context, sessao *entities.Dispositivo) error
}

func NewSQLDispositivoRepo(Conn *sql.DB) DispositivoRepo {
	return &postgresDispositivoRepo{Conn: Conn}
}

type postgresDispositivoRepo struct {
	Conn *sql.DB
}

func (m *postgresDispositivoRepo) Commit(ctx context.Context, dispositivo *entities.Dispositivo) error {
	query := `INSERT INTO cadastros.dispositivos (aluno_id, model, platform, uuid, version, manufacter, isvirtual, serial, created_at, updated_at) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx,
		dispositivo.AlunoID,
		dispositivo.Model,
		dispositivo.Platform,
		dispositivo.UUID,
		dispositivo.Version,
		dispositivo.Manufacturer,
		dispositivo.IsVirtual,
		dispositivo.Serial,
		dispositivo.CreatedAt,
		dispositivo.UpdatedAt,
	)
	if err != nil {
		return err
	}
	defer stmt.Close()
	return nil
}
func (m *postgresDispositivoRepo) Find(ctx context.Context, aluno_id int64) (*entities.Dispositivo, bool, error) {
	query := `SELECT id, aluno_id, qtde_login, qtde_req, cookie_name, cookie_value, created_at, updated_at from cadastros.aluno_sessao where aluno_id=$1`
	rows := m.Conn.QueryRowContext(ctx, query, aluno_id)
	defer m.Conn.Close()
	dispositivo := new(entities.Dispositivo)

	err := rows.Scan(
		&dispositivo.ID,
		&dispositivo.AlunoID,
		&dispositivo.Model,
		&dispositivo.Platform,
		&dispositivo.UUID,
		&dispositivo.Version,
		&dispositivo.Manufacturer,
		&dispositivo.IsVirtual,
		&dispositivo.Serial,
		&dispositivo.CreatedAt,
		&dispositivo.UpdatedAt,
	)
	if err != nil {
		return nil, false, err
	}
	return dispositivo, true, nil
}
