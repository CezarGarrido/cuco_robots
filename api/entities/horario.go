package entities

import "time"

type Horario struct {
	ID            int64      `json:"id"`
	Curso         string     `json:"curso"`
	AnoLetivo     string     `json:"ano_letivo"`
	Serie         string     `json:"serie"`
	Periodo       string     `json:"periodo"`
	Horario       string     `json:"horario"`
	ProfessorNome string     `json:"professor_nome"`
	Disciplina    string     `json:"disciplina"`
	DiaSemana     string     `json:"dia_semana"`
	CreatedAt     *time.Time `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
}