class Frequencia {
  int id;
  int alunoId;
  int disciplinaId;
  String mes;
  int dia;
  String valor;
  String createdAt;
  String updatedAt;

  Frequencia(
      {this.id,
      this.alunoId,
      this.disciplinaId,
      this.mes,
      this.dia,
      this.valor,
      this.createdAt,
      this.updatedAt});

  Frequencia.fromJson(Map<String, dynamic> json) {
    id = json['id'];
    alunoId = json['aluno_id'];
    disciplinaId = json['disciplina_id'];
    mes = json['mes'];
    dia = json['dia'];
    valor = json['valor'];
    createdAt = json['created_at'];
    updatedAt = json['updated_at'];
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = new Map<String, dynamic>();
    data['id'] = this.id;
    data['aluno_id'] = this.alunoId;
    data['disciplina_id'] = this.disciplinaId;
    data['mes'] = this.mes;
    data['dia'] = this.dia;
    data['valor'] = this.valor;
    data['created_at'] = this.createdAt;
    data['updated_at'] = this.updatedAt;
    return data;
  }
}
