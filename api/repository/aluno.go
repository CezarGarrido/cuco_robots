package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	entities "github.com/CezarGarrido/cuco_robots/api/entities"
)

// AlunoRepo explain...
type AlunoRepo interface {
	//GetByID(ctx context.Context, id int64) (*entities.Log, error)*/
	IsExiste(ctx context.Context, rgm, senha string) (bool, error)

	Create(ctx context.Context, advogado *entities.Aluno) (int64, error)
	/*Update(ctx context.Context, p *entities.Log) (*entities.Log, error)*/
	Update(ctx context.Context, p *entities.Aluno) (*entities.Aluno, error)
	GetByLogin(ctx context.Context, rgm string) (*entities.Aluno, error)
}

// NewSQLAlunoRepo retunrs implement of Aluno repository interface
func NewSQLAlunoRepo(Conn *sql.DB) AlunoRepo {
	return &mysqlAlunoRepo{Conn: Conn}
}

type mysqlAlunoRepo struct {
	Conn *sql.DB
}

func (m *mysqlAlunoRepo) Create(ctx context.Context, aluno *entities.Aluno) (int64, error) {
	query := "Insert into cadastros.alunos (nome,curso,ano,unidade,rgm,senha,email,created_at,updated_at) values($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id"
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return -1, err
	}
	var alunoID int64
	err = stmt.QueryRowContext(ctx,
		aluno.Nome,
		aluno.Curso,
		aluno.Ano,
		aluno.Unidade,
		aluno.Rgm,
		aluno.Senha,
		aluno.Email,
		aluno.CreatedAt,
		aluno.UpdatedAt,
	).Scan(&alunoID)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()
	return alunoID, nil
}

/*func (m *mysqlAdvogadoRepo) Fetch(ctx context.Context, uf string) ([]*entities.Aluno, error) {
	query := "Select id,publicacao_id,nome, oab, uf, tipo, possui_suplementar, qtde_processos_ms, created_at, updated_at From publicacoes_tjms_advs where uf <> ? and uf <> '' and possui_suplementar is NULL"
	return m.fetch(ctx, query, uf)
}*/

func (m *mysqlAlunoRepo) fetch(ctx context.Context, query string, args ...interface{}) ([]*entities.Aluno, error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	payload := make([]*entities.Aluno, 0)
	for rows.Next() {
		aluno := new(entities.Aluno)
		err := rows.Scan(
			&aluno.ID,
			&aluno.Nome,
			&aluno.Curso,
			&aluno.Ano,
			&aluno.Unidade,
			&aluno.Rgm,
			&aluno.Senha,
			&aluno.Email,
			&aluno.Telefone,
			&aluno.CreatedAt,
			&aluno.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		payload = append(payload, aluno)
	}
	return payload, nil
}

func (m *mysqlAlunoRepo) Update(ctx context.Context, aluno *entities.Aluno) (*entities.Aluno, error) {
	query := "Update alunos set id=$1,nome=$2, curso=$3, ano=$4, unidade=$5, rgm=$6, senha=$7, created_at=$8, updated_at=$9 where id=$1"
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	_, err = stmt.ExecContext(
		ctx,
		&aluno.ID,
		&aluno.Nome,
		&aluno.Curso,
		&aluno.Ano,
		&aluno.Unidade,
		&aluno.Rgm,
		&aluno.Senha,
		&aluno.Email,
		&aluno.Telefone,
		&aluno.CreatedAt,
		&aluno.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return aluno, nil
}

func (m *mysqlAlunoRepo) GetByLogin(ctx context.Context, rgm string) (*entities.Aluno, error) {
	query := "Select id,nome,curso,ano,unidade,rgm,senha,email,telefone,created_at,updated_at FROM cadastros.alunos WHERE rgm=$1"
	rows, err := m.fetch(ctx, query, rgm)
	if err != nil {
		return nil, err
	}
	payload := &entities.Aluno{}
	if len(rows) > 0 {
		payload = rows[0]
	} else {
		return nil, errors.New("Login não encontrado")
	}
	return payload, nil
}

func (m *mysqlAlunoRepo) IsExiste(ctx context.Context, rgm, senha string) (bool, error) {
	var exist = false
	err := m.Conn.QueryRowContext(ctx, "Select exists(Select 1 from cadastros.alunos where rgm=$1 and senha=$2);", rgm, senha).Scan(
		&exist,
	)
	fmt.Println(exist)
	if err != nil {
		return false, err
	}
	return exist, nil
}
