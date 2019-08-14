package repository

import (
	"context"
	"database/sql"

	entities "github.com/CezarGarrido/cuco_robots/api/entities"
)

// AlunoDisciplinaRepo explain...
type AlunoDisciplinaRepo interface {
	GetByAlunoID(ctx context.Context, id int64) ([]*entities.AlunoDisciplina, error)
}

// NewSQLAlunoDisciplinaRepo retunrs implement of AlunoDisciplina repository interface
func NewSQLAlunoDisciplinaRepo(Conn *sql.DB) AlunoDisciplinaRepo {
	return &mysqlAlunoDisciplinaRepo{Conn: Conn}
}

type mysqlAlunoDisciplinaRepo struct {
	Conn *sql.DB
}

func (m *mysqlAlunoDisciplinaRepo) GetByAlunoID(ctx context.Context, alunoID int64) ([]*entities.AlunoDisciplina, error) {
	query := "Select id, aluno_id, disciplina_id, uems_id, ano, carga_horaria_presencial, maximo_faltas,periodo_letivo, professor, media_avaliacoes, optativa, exame, faltas, situacao, created_at, updated_at From aluno_disciplinas where aluno_id=$1"
	return m.fetch(ctx, query, alunoID)
}

func (m *mysqlAlunoDisciplinaRepo) fetch(ctx context.Context, query string, args ...interface{}) ([]*entities.AlunoDisciplina, error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	payload := make([]*entities.AlunoDisciplina, 0)
	for rows.Next() {
		AlunoDisciplina := new(entities.AlunoDisciplina)
		err := rows.Scan(
			&AlunoDisciplina.ID,
			&AlunoDisciplina.AlunoID,
			&AlunoDisciplina.DisciplinaID,
			&AlunoDisciplina.UemsID,
			&AlunoDisciplina.Ano,
			&AlunoDisciplina.CargaHorariaPresencial,
			&AlunoDisciplina.MaximoFaltas,
			&AlunoDisciplina.PeriodoLetivo,
			&AlunoDisciplina.Professor,
			&AlunoDisciplina.MediaAvaliacoes,
			&AlunoDisciplina.Optativa,
			&AlunoDisciplina.Exame,
			&AlunoDisciplina.Faltas,
			&AlunoDisciplina.Situacao,
			&AlunoDisciplina.CreatedAt,
			&AlunoDisciplina.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		payload = append(payload, AlunoDisciplina)
	}
	return payload, nil
}
