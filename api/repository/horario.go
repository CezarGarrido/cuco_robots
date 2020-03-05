package repository

import (
	"context"
	"database/sql"

	entities "github.com/CezarGarrido/cuco_robots/api/entities"
)

type HorarioRepo interface {
	GetByCurso(ctx context.Context, descricao string) ([]*entities.Horario, error)
}

func NewSQLHorarioRepo(Conn *sql.DB) HorarioRepo {
	return &postgresHorarioRepo{Conn: Conn}
}

type postgresHorarioRepo struct {
	Conn *sql.DB
}

func (m *postgresHorarioRepo) GetByCurso(ctx context.Context, descricao string) ([]*entities.Horario, error) {
	query := "SELECT id, curso, ano_letivo, serie, periodo, horario, professor_nome, disciplina, dia_semana,created_at, updated_at FROM cadastros.horarios WHERE curso=$1;"
	return m.fetch(ctx, query, descricao)
}
func (m *postgresHorarioRepo) fetch(ctx context.Context, query string, args ...interface{}) ([]*entities.Horario, error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	payload := make([]*entities.Horario, 0)
	for rows.Next() {
		horario := new(entities.Horario)
		err := rows.Scan(
			&horario.ID,
			&horario.Curso,
			&horario.AnoLetivo,
			&horario.Serie,
			&horario.Periodo,
			&horario.Horario,
			&horario.ProfessorNome,
			&horario.Disciplina,
			&horario.DiaSemana,
			&horario.CreatedAt,
			&horario.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		payload = append(payload, horario)
	}
	return payload, nil
}
