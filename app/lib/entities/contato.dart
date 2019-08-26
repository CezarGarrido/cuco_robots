class Contato {
  int id;
  int alunoId;
  String tipo;
  String valor;
  String createdAt;
  Null updatedAt;

  Contato(
      {this.id,
      this.alunoId,
      this.tipo,
      this.valor,
      this.createdAt,
      this.updatedAt});

  Contato.fromJson(Map<String, dynamic> json) {
    id = json['id'];
    alunoId = json['aluno_id'];
    tipo = json['tipo'];
    valor = json['valor'];
    createdAt = json['created_at'];
    updatedAt = json['updated_at'];
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = new Map<String, dynamic>();
    data['id'] = this.id;
    data['aluno_id'] = this.alunoId;
    data['tipo'] = this.tipo;
    data['valor'] = this.valor;
    data['created_at'] = this.createdAt;
    data['updated_at'] = this.updatedAt;
    return data;
  }
}