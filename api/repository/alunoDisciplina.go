package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	entities "github.com/CezarGarrido/cuco_robots/api/entities"
)

// AlunoDisciplinaRepo explain...
type AlunoDisciplinaRepo interface {
	GetByAlunoID(ctx context.Context, id int64) ([]*entities.AlunoDisciplina, error)
	IsExiste(ctx context.Context, aluno_id, uems_id int64) (bool, error)
	Create(ctx context.Context, alunoDisciplina *entities.AlunoDisciplina) (int64, error)
	Update(ctx context.Context, alunoDisciplina *entities.AlunoDisciplina) (*entities.AlunoDisciplina, error)
	GetByUemsID(ctx context.Context, alunoID, uems_id int64) (*entities.AlunoDisciplina, error)
}

// NewSQLAlunoDisciplinaRepo retunrs implement of AlunoDisciplina repository interface
func NewSQLAlunoDisciplinaRepo(Conn *sql.DB) AlunoDisciplinaRepo {
	return &mysqlAlunoDisciplinaRepo{Conn: Conn}
}

type mysqlAlunoDisciplinaRepo struct {
	Conn *sql.DB
}

func (m *mysqlAlunoDisciplinaRepo) Create(ctx context.Context, alunoDisciplina *entities.AlunoDisciplina) (int64, error) {
	query := `INSERT INTO cadastros.aluno_disciplinas (aluno_id, uems_id, unidade, curso, disciplina, turma, serie_disciplina, carga_horaria_presencial, maximo_faltas, periodo_letivo, professor, media_avaliacoes, optativa, exame, media_final, faltas, situacao, created_at, updated_at) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10,$11, $12, $13, $14, $15, $16, $17, $18, $19) RETURNING id`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return -1, err
	}
	var alunoDisciplinaID int64
	err = stmt.QueryRowContext(ctx,
		alunoDisciplina.AlunoID,
		alunoDisciplina.UemsID,
		alunoDisciplina.Unidade,
		alunoDisciplina.Curso,
		alunoDisciplina.Disciplina,
		alunoDisciplina.Turma,
		alunoDisciplina.SerieDisciplina,
		alunoDisciplina.CargaHorariaPresencial,
		alunoDisciplina.MaximoFaltas,
		alunoDisciplina.PeriodoLetivo,
		alunoDisciplina.Professor,
		alunoDisciplina.MediaAvaliacoes,
		alunoDisciplina.Optativa,
		alunoDisciplina.Exame,
		alunoDisciplina.MediaFinal,
		alunoDisciplina.Faltas,
		alunoDisciplina.Situacao,
		alunoDisciplina.CreatedAt,
		alunoDisciplina.UpdatedAt,
	).Scan(&alunoDisciplinaID)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()
	notaRepo := NewSQLNotaRepo(m.Conn)
	for _, nota := range alunoDisciplina.Notas {
		nota.DisciplinaID = alunoDisciplinaID
		_, err := notaRepo.Create(ctx, nota)
		if err != nil {
			return -1, err
		}
	}
	return alunoDisciplinaID, nil
}

func (m *mysqlAlunoDisciplinaRepo) GetByAlunoID(ctx context.Context, alunoID int64) ([]*entities.AlunoDisciplina, error) {
	query := `SELECT id, aluno_id, uems_id, unidade, curso, disciplina, turma, serie_disciplina, carga_horaria_presencial, maximo_faltas, periodo_letivo, professor, media_avaliacoes, optativa, exame, media_final, faltas, situacao, created_at, updated_at FROM cadastros.aluno_disciplinas WHERE aluno_id=$1 order by id asc;`
	payload, err := m.fetch(ctx, query, alunoID)
	if err != nil {
		return nil, err
	}
	slice_payload := make([]*entities.AlunoDisciplina, 0)
	notaRepo := NewSQLNotaRepo(m.Conn)
	for _, data := range payload {
		notas, err := notaRepo.GetByDisciplinaID(ctx, data.AlunoID, data.ID)
		if err != nil {
			return nil, err
		}
		data.Notas = notas
		slice_payload = append(slice_payload, data)
	}
	return slice_payload, nil
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
			&AlunoDisciplina.UemsID,
			&AlunoDisciplina.Unidade,
			&AlunoDisciplina.Curso,
			&AlunoDisciplina.Disciplina,
			&AlunoDisciplina.Turma,
			&AlunoDisciplina.SerieDisciplina,
			&AlunoDisciplina.CargaHorariaPresencial,
			&AlunoDisciplina.MaximoFaltas,
			&AlunoDisciplina.PeriodoLetivo,
			&AlunoDisciplina.Professor,
			&AlunoDisciplina.MediaAvaliacoes,
			&AlunoDisciplina.Optativa,
			&AlunoDisciplina.Exame,
			&AlunoDisciplina.MediaFinal,
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

func (m *mysqlAlunoDisciplinaRepo) IsExiste(ctx context.Context, aluno_id, uems_id int64) (bool, error) {
	var exist = false
	err := m.Conn.QueryRowContext(ctx, "Select exists(Select 1 from cadastros.aluno_disciplinas where aluno_id=$1 and uems_id=$2);", aluno_id, uems_id).Scan(
		&exist,
	)
	if err != nil {
		return false, err
	}
	return exist, nil
}

func (m *mysqlAlunoDisciplinaRepo) Update(ctx context.Context, alunoDisciplina *entities.AlunoDisciplina) (*entities.AlunoDisciplina, error) {
	query := `UPDATE cadastros.aluno_disciplinas SET aluno_id=$1, uems_id=$2, unidade=$3, curso=$4, disciplina=$5, turma=$6, serie_disciplina=$7, carga_horaria_presencial=$8, maximo_faltas=$9, periodo_letivo=$10, professor=$11, media_avaliacoes=$12, optativa=$13, exame=$14, media_final=$15, faltas=$16, situacao=$17, updated_at=$18 WHERE id=$19;
	`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	_, err = stmt.ExecContext(
		ctx,
		alunoDisciplina.AlunoID,
		alunoDisciplina.UemsID,
		alunoDisciplina.Unidade,
		alunoDisciplina.Curso,
		alunoDisciplina.Disciplina,
		alunoDisciplina.Turma,
		alunoDisciplina.SerieDisciplina,
		alunoDisciplina.CargaHorariaPresencial,
		alunoDisciplina.MaximoFaltas,
		alunoDisciplina.PeriodoLetivo,
		alunoDisciplina.Professor,
		alunoDisciplina.MediaAvaliacoes,
		alunoDisciplina.Optativa,
		alunoDisciplina.Exame,
		alunoDisciplina.MediaFinal,
		alunoDisciplina.Faltas,
		alunoDisciplina.Situacao,
		alunoDisciplina.UpdatedAt,
		alunoDisciplina.ID,
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	notaRepo := NewSQLNotaRepo(m.Conn)
	for _, nota := range alunoDisciplina.Notas {
		nota.DisciplinaID = alunoDisciplina.ID
		exist, err := notaRepo.IsExiste(ctx, alunoDisciplina.AlunoID, alunoDisciplina.ID, nota.Descricao)
		if !exist {
			_, err = notaRepo.Create(ctx, nota)
			if err != nil {
				return nil, err
			}
		} else {
			notaAnterior, err := notaRepo.GetByDescricao(ctx, alunoDisciplina.AlunoID, alunoDisciplina.ID, nota.Descricao)
			if err != nil {
				return nil, err
			}

			if nota.Valor != nil {
				notaAnterior.Valor = nota.Valor
			}
			h := time.Now()
			notaAnterior.UpdatedAt = &h
			_, err = notaRepo.Update(ctx, notaAnterior)
			if err != nil {
				return nil, err
			}
		}
	}
	return alunoDisciplina, nil
}

func (m *mysqlAlunoDisciplinaRepo) GetByUemsID(ctx context.Context, alunoID, uems_id int64) (*entities.AlunoDisciplina, error) {
	query := `SELECT id, aluno_id, uems_id, unidade, curso, disciplina, turma, serie_disciplina, carga_horaria_presencial, maximo_faltas, periodo_letivo, professor, media_avaliacoes, optativa, exame, media_final, faltas, situacao, created_at, updated_at FROM cadastros.aluno_disciplinas WHERE aluno_id=$1 and uems_id=$2;
	`
	rows, err := m.fetch(ctx, query, alunoID, uems_id)
	if err != nil {
		return nil, err
	}
	payload := &entities.AlunoDisciplina{}
	if len(rows) > 0 {
		payload = rows[0]
	} else {
		return nil, errors.New("Disciplina não encontrada.")
	}
	return payload, nil
}
