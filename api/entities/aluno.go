package entities

import (
	"time"
)

type Aluno struct {
	ID              int64       `json:"id"`
	Guid            string      `json:"guid"`
	Nome            string      `json:"nome"`
	Rgm             string      `json:"rgm"`
	Senha           string      `json:"senha"`
	Curso           string      `json:"curso"`
	DataNascimento  *time.Time  `json:"data_nascimento"`
	Sexo            *string     `json:"sexo"`
	NomePai         *string     `json:"nome_pai"`
	NomeMae         *string     `json:"nome_mae"`
	EstadoCivil     *string     `json:"estado_civil"`
	Nacionalidade   *string     `json:"nacionalidade"`
	Naturalidade    *string     `json:"naturalidade"`
	Fenotipo        *string     `json:"fenotipo"`
	CPF             *string     `json:"cpf"`
	RG              *string     `json:"rg"`
	RGOrgaoEmissor  *string     `json:"rg_orgao_emissor"`
	RGEstadoEmissor *string     `json:"rg_estado_emissor"`
	RGDataEmissao   *time.Time  `json:"rg_data_emissao "`
	Contatos        []*Contato  `json:"contatos"`
	Enderecos       []*Endereco `json:"enderecos"`
	CreatedAt       *time.Time  `json:"created_at"`
	UpdatedAt       *time.Time  `json:"updated_at"`
}
