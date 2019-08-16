CREATE TABLE IF NOT EXISTS cadastros.alunos (
  id bigserial CONSTRAINT pk_id_aluno primary key,
  guid TEXT,
  nome TEXT NOT NULL,
  rgm TEXT NOT NULL,
  senha TEXT NOT NULL,
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

CREATE TABLE  IF NOT EXISTS cadastros.aluno_contatos (
  id bigserial CONSTRAINT pk_id_aluno_contato primary key,
  aluno_id int8 not null,
  Tipo TEXT,
  Valor TEXT,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  FOREIGN KEY (aluno_id) REFERENCES cadastros.alunos(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS cadastros.aluno_enderecos (
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
  FOREIGN KEY (aluno_id) REFERENCES cadastros.alunos(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS cadastros.aluno_disciplinas (
  id bigserial CONSTRAINT pk_id_aluno_disciplina primary key,
  aluno_id int8 not null,
  uems_id int8 not null,
  unidade TEXT,
  curso TEXT,
  disciplina TEXT,
  turma TEXT,
  serie_disciplina TEXT,
  carga_horaria_presencial TEXT,
  maximo_faltas int,
  periodo_letivo TEXT,
  professor TEXT,
  media_avaliacoes numeric(18,2),
  optativa numeric(18,2),
  exame numeric(18,2),
  media_final numeric(18,2),
  faltas int,
  situacao TEXT,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  FOREIGN KEY (aluno_id) REFERENCES cadastros.alunos(id) ON DELETE CASCADE
);

CREATE TABLE  IF NOT EXISTS cadastros.aluno_notas (
  id bigserial CONSTRAINT pk_id_aluno_nota primary key,
  aluno_id int8 not null,
  disciplina_id int8 not null,
  descricao TEXT not null,
  valor numeric(18,2),
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  FOREIGN KEY (aluno_id) REFERENCES cadastros.alunos(id) ON DELETE CASCADE,
  FOREIGN KEY (disciplina_id) REFERENCES cadastros.aluno_disciplinas(id) ON DELETE CASCADE
);