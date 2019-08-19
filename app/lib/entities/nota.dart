class Nota {
  int id;
  int disciplinaId;
  int alunoId;
  String descricao;
  double valor;
  String createdAt;
  String updatedAt;

  Nota(
      {this.id,
      this.disciplinaId,
      this.alunoId,
      this.descricao,
      this.valor,
      this.createdAt,
      this.updatedAt});

  Nota.fromJson(Map<String, dynamic> json) {
    id = json['id'];
    disciplinaId = json['disciplina_id'];
    alunoId = json['aluno_id'];
    descricao = json['descricao'];
    valor = json['valor'];
    createdAt = json['created_at'];
    updatedAt = json['updated_at'];
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = new Map<String, dynamic>();
    data['id'] = this.id;
    data['disciplina_id'] = this.disciplinaId;
    data['aluno_id'] = this.alunoId;
    data['descricao'] = this.descricao;
    data['valor'] = this.valor;
    data['created_at'] = this.createdAt;
    data['updated_at'] = this.updatedAt;
    return data;
  }
}