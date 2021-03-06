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
  periodo_letivo TEXT,
  descricao TEXT not null,
  valor numeric(18,2),
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  FOREIGN KEY (aluno_id) REFERENCES cadastros.alunos(id) ON DELETE CASCADE,
  FOREIGN KEY (disciplina_id) REFERENCES cadastros.aluno_disciplinas(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS cadastros.aluno_frequencia (
  id bigserial CONSTRAINT pk_id_aluno_frequencia primary key,
  aluno_id int8 not null,
  disciplina_id int8 not null,
  mes TEXT not null,
  dia int not null,
  valor TEXT not null, 
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  FOREIGN KEY (aluno_id) REFERENCES cadastros.alunos(id) ON DELETE CASCADE,
  FOREIGN KEY (disciplina_id) REFERENCES cadastros.aluno_disciplinas(id) ON DELETE CASCADE
)

CREATE TABLE IF NOT EXISTS cadastros.dispositivos (
  id bigserial CONSTRAINT pk_id_dispositivo primary key,
  aluno_id int8 not null,
  model TEXT not null,
  platform TEXT not null,
  uuid TEXT not null,
  `version` TEXT not null,
  manufacturer TEXT not null,
  is_virtual TEXT not null,
  `serial` TEXT not null, 
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  FOREIGN KEY (aluno_id) REFERENCES cadastros.alunos(id) ON DELETE CASCADE
)

CREATE TABLE IF NOT EXISTS cadastros.horarios (
  id bigserial CONSTRAINT pk_id_horario primary key,
  curso varchar(100) NOT NULL,
  ano_letivo varchar(50) NOT NULL,
  serie varchar(50) NOT NULL,
  periodo varchar(50) NOT NULL,
  horario varchar(50) NOT NULL,
  professor_nome varchar(50) NOT NULL,
  disciplina varchar(50) NOT NULL,
  dia_semana varchar(50) NOT NULL,
  created_at TIMESTAMP,
  updated_at TIMESTAMP
)