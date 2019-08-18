package repository

import (
	"context"
	"database/sql"
	"errors"

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
	query := `INSERT INTO cadastros.alunos (guid, nome, rgm, senha, data_nascimento, sexo, nome_pai, nome_mae, estado_civil, nacionalidade,naturalidade, fenotipo, cpf, rg, rg_orgao_emissor, rg_estado_emissor, rg_data_emissao, created_at, updated_at) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19) RETURNING id`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return -1, err
	}
	var alunoID int64
	err = stmt.QueryRowContext(ctx,
		aluno.Guid,
		aluno.Nome,
		aluno.Rgm,
		aluno.Senha,
		aluno.DataNascimento,
		aluno.Sexo,
		aluno.NomePai,
		aluno.NomeMae,
		aluno.EstadoCivil,
		aluno.Nacionalidade,
		aluno.Naturalidade,
		aluno.Fenotipo,
		aluno.CPF,
		aluno.RG,
		aluno.RGOrgaoEmissor,
		aluno.RGEstadoEmissor,
		aluno.RGDataEmissao,
		aluno.CreatedAt,
		aluno.UpdatedAt,
	).Scan(&alunoID)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()
	contatoRepo := NewSQLContatoRepo(m.Conn)
	for _, contato := range aluno.Contatos {
		contato.AlunoID = alunoID
		_, err = contatoRepo.Create(ctx, contato)
		if err != nil {
			return -1, err
		}
	}
	enderecoRepo := NewSQLEnderecoRepo(m.Conn)
	for _, endereco := range aluno.Enderecos {
		endereco.AlunoID = alunoID
		_, err = enderecoRepo.Create(ctx, endereco)
		if err != nil {
			return -1, err
		}
	}
	sessaoRepo := NewSQLSessaoRepo(m.Conn)
	aluno.Sessao.AlunoID = alunoID
	err = sessaoRepo.Commit(ctx, aluno.Sessao)
	if err != nil {
		return -1, err
	}
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
			&aluno.Guid,
			&aluno.Nome,
			&aluno.Rgm,
			&aluno.Senha,
			&aluno.DataNascimento,
			&aluno.Sexo,
			&aluno.NomePai,
			&aluno.NomeMae,
			&aluno.EstadoCivil,
			&aluno.Nacionalidade,
			&aluno.Naturalidade,
			&aluno.Fenotipo,
			&aluno.CPF,
			&aluno.RG,
			&aluno.RGOrgaoEmissor,
			&aluno.RGEstadoEmissor,
			&aluno.RGDataEmissao,
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
	query := `UPDATE cadastros.alunos SET guid=$1, nome=$2, rgm=$3, senha=$4, data_nascimento=$5, sexo=$6, nome_pai=$7, nome_mae=$7, estado_civil=$8, nacionalidade=$9, naturalidade=$10, fenotipo=$11, cpf=$12, rg=$13, rg_orgao_emissor=$14, rg_estado_emissor=$15, rg_data_emissao=$16,updated_at=$17 WHERE id IS $18;`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	_, err = stmt.ExecContext(
		ctx,
		aluno.Guid,
		aluno.Nome,
		aluno.Rgm,
		aluno.Senha,
		aluno.Curso,
		aluno.DataNascimento,
		aluno.Sexo,
		aluno.NomePai,
		aluno.NomeMae,
		aluno.EstadoCivil,
		aluno.Nacionalidade,
		aluno.Naturalidade,
		aluno.Fenotipo,
		aluno.CPF,
		aluno.RG,
		aluno.RGOrgaoEmissor,
		aluno.RGEstadoEmissor,
		aluno.RGDataEmissao,
		aluno.UpdatedAt,
		aluno.ID,
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	contatoRepo := NewSQLContatoRepo(m.Conn)
	for _, contato := range aluno.Contatos {
		contato.AlunoID = aluno.ID
		deletedOk, err := contatoRepo.DeleteAll(ctx, aluno.ID)
		if err != nil {
			return nil, err
		}
		if deletedOk {
			_, err = contatoRepo.Create(ctx, contato)
			if err != nil {
				return nil, err
			}
		}
	}
	enderecoRepo := NewSQLEnderecoRepo(m.Conn)
	for _, endereco := range aluno.Enderecos {
		endereco.AlunoID = aluno.ID
		deletedOk, err := enderecoRepo.DeleteAll(ctx, aluno.ID)
		if err != nil {
			return nil, err
		}
		if deletedOk {
			_, err = enderecoRepo.Create(ctx, endereco)
			if err != nil {
				return nil, err
			}
		}
	}
	return aluno, nil
}

func (m *mysqlAlunoRepo) GetByLogin(ctx context.Context, rgm string) (*entities.Aluno, error) {

	query := "SELECT id, guid, nome, rgm, senha, data_nascimento, sexo, nome_pai, nome_mae, estado_civil, nacionalidade, naturalidade, fenotipo, cpf, rg, rg_orgao_emissor, rg_estado_emissor, rg_data_emissao, created_at, updated_at FROM cadastros.alunos WHERE rgm=$1"

	rows, err := m.fetch(ctx, query, rgm)
	if err != nil {
		return nil, err
	}
	payload := &entities.Aluno{}
	if len(rows) > 0 {
		payload = rows[0]
	} else {
		return nil, errors.New("Aluno não encontrado")
	}

	enderecoRepo := NewSQLEnderecoRepo(m.Conn)
	payload.Enderecos, err = enderecoRepo.GetByAlunoID(ctx, payload.ID)
	if err != nil {
		return nil, err
	}
	contatoRepo := NewSQLContatoRepo(m.Conn)
	payload.Contatos, err = contatoRepo.GetByAlunoID(ctx, payload.ID)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func (m *mysqlAlunoRepo) IsExiste(ctx context.Context, rgm, senha string) (bool, error) {
	var exist = false
	err := m.Conn.QueryRowContext(ctx, "Select exists(Select 1 from cadastros.alunos where rgm=$1 and senha=$2);", rgm, senha).Scan(
		&exist,
	)
	if err != nil {
		return false, err
	}
	return exist, nil
}
