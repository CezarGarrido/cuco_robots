class Aluno {
  int id;
  String guid;
  String nome;
  String rgm;
  String senha;
  String curso;
  String dataNascimento;
  String sexo;
  String nomePai;
  String nomeMae;
  String estadoCivil;
  String nacionalidade;
  String naturalidade;
  String fenotipo;
  String cpf;
  String rg;
  String rgOrgaoEmissor;
  String rgEstadoEmissor;
  String rgDataEmissao;
  Null contatos;
  Null enderecos;
  String createdAt;
  String updatedAt;

  Aluno(
      {this.id,
      this.guid,
      this.nome,
      this.rgm,
      this.senha,
      this.curso,
      this.dataNascimento,
      this.sexo,
      this.nomePai,
      this.nomeMae,
      this.estadoCivil,
      this.nacionalidade,
      this.naturalidade,
      this.fenotipo,
      this.cpf,
      this.rg,
      this.rgOrgaoEmissor,
      this.rgEstadoEmissor,
      this.rgDataEmissao,
      this.contatos,
      this.enderecos,
      this.createdAt,
      this.updatedAt});

  Aluno.fromJson(Map<String, dynamic> json) {
    id = json['id'];
    guid = json['guid'];
    nome = json['nome'];
    rgm = json['rgm'];
    senha = json['senha'];
    curso = json['curso'];
    dataNascimento = json['data_nascimento'];
    sexo = json['sexo'];
    nomePai = json['nome_pai'];
    nomeMae = json['nome_mae'];
    estadoCivil = json['estado_civil'];
    nacionalidade = json['nacionalidade'];
    naturalidade = json['naturalidade'];
    fenotipo = json['fenotipo'];
    cpf = json['cpf'];
    rg = json['rg'];
    rgOrgaoEmissor = json['rg_orgao_emissor'];
    rgEstadoEmissor = json['rg_estado_emissor'];
    rgDataEmissao = json['rg_data_emissao '];
    createdAt = json['created_at'];
    updatedAt = json['updated_at'];
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = new Map<String, dynamic>();
    data['id'] = this.id;
    data['guid'] = this.guid;
    data['nome'] = this.nome;
    data['rgm'] = this.rgm;
    data['senha'] = this.senha;
    data['curso'] = this.curso;
    data['data_nascimento'] = this.dataNascimento;
    data['sexo'] = this.sexo;
    data['nome_pai'] = this.nomePai;
    data['nome_mae'] = this.nomeMae;
    data['estado_civil'] = this.estadoCivil;
    data['nacionalidade'] = this.nacionalidade;
    data['naturalidade'] = this.naturalidade;
    data['fenotipo'] = this.fenotipo;
    data['cpf'] = this.cpf;
    data['rg'] = this.rg;
    data['rg_orgao_emissor'] = this.rgOrgaoEmissor;
    data['rg_estado_emissor'] = this.rgEstadoEmissor;
    data['rg_data_emissao '] = this.rgDataEmissao;
    data['created_at'] = this.createdAt;
    data['updated_at'] = this.updatedAt;
    return data;
  }
}