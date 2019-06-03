package models

import (
	"database/sql"
	"time"
)

type AlunoNota struct {
	ID                     int64       `json:"id"`
	IDAluno                int64       `json:"id_aluno"`
	IDDisciplina           int64       `json:"id_disciplina"`
	Unidade                string      `json:"unidade"`
	Curso                  string      `json:"curso"`
	Disciplina             string      `json:"disciplina"`
	Turma                  string      `json:"turma"`
	SerieDisciplina        string      `json:"serieDisciplina"`
	CargaHorariaPresencial string      `json:"cargaHorariaPresencial"`
	MaximoFaltas           string      `json:"maximoFaltas"`
	PeriodoLetivo          string      `json:"periodoLetivo"`
	Professor              string      `json:"professor"`
	MediaAvaliacoes        string      `json:"mediaAvaliacoes"`
	Optativa               string      `json:"optativa"`
	Exame                  string      `json:"exame"`
	MediaFinal             string      `json:"mediaFinal"`
	Faltas                 string      `json:"faltas"`
	Situacao               string      `json:"situacao"`
	Notas                  []ValorNota `json:"notas"`
	CreatedAt              time.Time   `json:"created_at"`
	UpdatedAt              time.Time   `json:"updated_at"`
}

func (this *AlunoNota) Create(db *sql.DB) (int64, error) {
	var id int64
	stmtIns, err := db.Prepare("INSERT INTO cadastros.alunos_notas (id_aluno,id_disciplina,unidade,curso,disciplina,turma," +
		"seriedisciplina,cargahorariapresencial,maximofaltas,periodoletivo,professor,mediaavaliacoes,optativa,exame,mediafinal," +
		"faltas,situacao,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19) RETURNING id")

	if err != nil {
		return id, err
	}
	if err := stmtIns.QueryRow(
		this.IDAluno,
		this.IDDisciplina,
		this.Unidade,
		this.Curso,
		this.Disciplina,
		this.Turma,
		this.SerieDisciplina,
		this.CargaHorariaPresencial,
		this.MaximoFaltas,
		this.PeriodoLetivo,
		this.Professor,
		this.MediaAvaliacoes,
		this.Optativa,
		this.Exame,
		this.MediaFinal,
		this.Faltas,
		this.Situacao,
		this.CreatedAt,
		this.UpdatedAt,
	).Scan(&id); err != nil {
		return id, err
	}
	defer stmtIns.Close()
	return id, nil
}
func (this AlunoNota) IsExist(db *sql.DB) bool {
	var count bool
	_ = db.QueryRow("select exists(select 1 from cadastros.alunos_notas where id_disciplina=$1 and id_aluno=$2)", this.IDDisciplina, this.IDAluno).Scan(&count)
	return count
}
func (this AlunoNota) Update(db *sql.DB) error {
	stmtIns, err := db.Prepare("UPDATE cadastros.alunos_notas set id_aluno=$1, id_disciplina=$2," +
		"unidade=$3, curso=$4, disciplina=$5, turma=$6, seriedisciplina=$7, cargahorariapresencial=$8," +
		"maximofaltas=$9, periodoletivo=$10, professor=$11, mediaavaliacoes=$12, optativa=$13, exame=$14, mediafinal=$15," +
		"faltas=$16, situacao=$17, created_at=$18, updated_at=$19 Where id=$20 ") // ? = placeholder
	if err != nil {
		return err
	}
	defer stmtIns.Close()
	if _, err := stmtIns.Exec(
		this.IDAluno,
		this.IDDisciplina,
		this.Unidade,
		this.Curso,
		this.Disciplina,
		this.Turma,
		this.SerieDisciplina,
		this.CargaHorariaPresencial,
		this.MaximoFaltas,
		this.PeriodoLetivo,
		this.Professor,
		this.MediaAvaliacoes,
		this.Optativa,
		this.Exame,
		this.MediaFinal,
		this.Faltas,
		this.Situacao,
		this.CreatedAt,
		this.UpdatedAt,
		this.ID,
	); err != nil {
		return err
	}
	return nil
}

func (this *AlunoNota) GetByDisciplina(id int64, db *sql.DB) error {
	sql := "SELECT id, id_aluno, id_disciplina, " +
		"unidade, curso, disciplina, turma, seriedisciplina, cargahorariapresencial, " +
		"maximofaltas, periodoletivo, professor, mediaavaliacoes, optativa, exame, mediafinal, " +
		"faltas, situacao, created_at, updated_at FROM cadastros.alunos_notas Where id_disciplina=$1 and id_aluno=$2"
	selDB, err := db.Query(sql, id, this.IDAluno)
	if err != nil {
		return err
	}
	for selDB.Next() {
		err = selDB.Scan(
			&this.ID,
			&this.IDAluno,
			&this.IDDisciplina,
			&this.Unidade,
			&this.Curso,
			&this.Disciplina,
			&this.Turma,
			&this.SerieDisciplina,
			&this.CargaHorariaPresencial,
			&this.MaximoFaltas,
			&this.PeriodoLetivo,
			&this.Professor,
			&this.MediaAvaliacoes,
			&this.Optativa,
			&this.Exame,
			&this.MediaFinal,
			&this.Faltas,
			&this.Situacao,
			&this.CreatedAt,
			&this.UpdatedAt,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

type ValorNota struct {
	ID        int64     `json:"id"`
	IDNota    int64     `json:"id_nota"`
	Descricao string    `json:"descricao"`
	Valor     string    `json:"valor"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (this *ValorNota) Create(db *sql.DB) (int64, error) {
	var id int64
	stmtIns, err := db.Prepare("INSERT INTO cadastros.valores_notas (id_nota,descricao,valor,created_at,updated_at) VALUES ($1,$2,$3,$4,$5) RETURNING id")

	if err != nil {
		return id, err
	}
	if err := stmtIns.QueryRow(
		this.IDNota,
		this.Descricao,
		this.Valor,
		this.CreatedAt,
		this.UpdatedAt,
	).Scan(&id); err != nil {
		return id, err
	}
	defer stmtIns.Close()
	return id, nil
}

func (this ValorNota) IsExist(db *sql.DB) bool {
	var count bool
	_ = db.QueryRow("select exists(select 1 from cadastros.valores_notas where id_nota=$1 and descricao=$2)", this.IDNota, this.Descricao).Scan(&count)
	return count
}

func (this ValorNota) Update(db *sql.DB) error {
	stmtIns, err := db.Prepare("UPDATE cadastros.valores_notas set id_nota=$1,descricao=$2,valor=$3,created_at=$4,updated_at=$5 Where id=$6 ") // ? = placeholder
	if err != nil {
		return err
	}
	defer stmtIns.Close()
	if _, err := stmtIns.Exec(
		this.IDNota,
		this.Descricao,
		this.Valor,
		this.CreatedAt,
		this.UpdatedAt,
		this.ID,
	); err != nil {
		return err
	}
	return nil
}

func (this *ValorNota) GetByDescricao(descricao string, db *sql.DB) error {
	sql := "SELECT id, id_nota, descricao, valor, created_at, updated_at FROM cadastros.valores_notas Where id_nota=$1 and descricao=$2"
	selDB, err := db.Query(sql, this.IDNota, descricao)
	if err != nil {
		return err
	}
	for selDB.Next() {
		err = selDB.Scan(
			&this.ID,
			&this.IDNota,
			&this.Descricao,
			&this.Valor,
			&this.CreatedAt,
			&this.UpdatedAt,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
