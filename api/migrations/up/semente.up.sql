CREATE TABLE cadastros.alunos_teste (
  id bigserial CONSTRAINT pk_id_aluno primary key,
  guid TEXT,
  nome TEXT NOT NULL,
  data_nascimento TIMESTAMP,
  sexo TEXT,
  nome_pai TEXT,
  nome_mae TEXT,
  estado_civil TEXT,
  nacionalidade TEXT,
  naturalidade TEXT,
  fenotipo TEXT,
  cpf TEXT,
  rg TEXT,
  rg_orgao_emissor TEXT,
  rg_estado_emissor TEXT,
  rg_data_emissao TIMESTAMP,
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);

CREATE TABLE cadastros.aluno_contatos (
  id bigserial CONSTRAINT pk_id_aluno_contato primary key,
  aluno_id int8 not null,
  Tipo TEXT,
  Valor TEXT,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  FOREIGN KEY (aluno_id) REFERENCES cadastros.alunos_teste(id) ON DELETE CASCADE
);

CREATE TABLE cadastros.aluno_enderecos (
  id bigserial CONSTRAINT pk_id_aluno_endereco primary key,
  aluno_id int8 not null,
  Logradouro TEXT,
  Numero INT,
  Complemento TEXT,
  Bairro TEXT,
  CEP TEXT,
  Cidade TEXT,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  FOREIGN KEY (aluno_id) REFERENCES cadastros.alunos_teste(id) ON DELETE CASCADE
);