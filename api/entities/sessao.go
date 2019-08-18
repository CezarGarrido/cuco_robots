package entities

import "time"

type Sessao struct {
	ID          int64      `json:"id"`
	AlunoID     int64      `json:"aluno_id"`
	QtdeLogin   int64      `json:"qtde_login"`
	QtdeRequest int64      `json:"qtde_request"`
	CookieName  string     `json:"cookie_name"`
	CookieValue string     `json:"cookie_value"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}
