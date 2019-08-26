class Endereco {
  int id;
  int alunoId;
  String logradouro;
  int numero;
  String complemento;
  String bairro;
  String cep;
  String cidade;
  String createdAt;
  Null updatedAt;

  Endereco(
      {this.id,
      this.alunoId,
      this.logradouro,
      this.numero,
      this.complemento,
      this.bairro,
      this.cep,
      this.cidade,
      this.createdAt,
      this.updatedAt});

  Endereco.fromJson(Map<String, dynamic> json) {
    id = json['id'];
    alunoId = json['aluno_id'];
    logradouro = json['logradouro'];
    numero = json['numero'];
    complemento = json['complemento'];
    bairro = json['bairro'];
    cep = json['cep'];
    cidade = json['cidade'];
    createdAt = json['created_at'];
    updatedAt = json['updated_at'];
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = new Map<String, dynamic>();
    data['id'] = this.id;
    data['aluno_id'] = this.alunoId;
    data['logradouro'] = this.logradouro;
    data['numero'] = this.numero;
    data['complemento'] = this.complemento;
    data['bairro'] = this.bairro;
    data['cep'] = this.cep;
    data['cidade'] = this.cidade;
    data['created_at'] = this.createdAt;
    data['updated_at'] = this.updatedAt;
    return data;
  }
}
