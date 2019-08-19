import 'nota.dart';

class Disciplina {
  int id;
  int alunoId;
  int uemsId;
  Null ano;
  String unidade;
  String curso;
  String disciplina;
  String turma;
  String serieDisciplina;
  int cargaHorariaPresencial;
  int maximoFaltas;
  String periodoLetivo;
  String professor;
  double mediaAvaliacoes;
  double optativa;
  double exame;
  double mediaFinal;
  int faltas;
  String situacao;
  List<Nota> notas;
  String createdAt;
  String updatedAt;

  Disciplina(
      {this.id,
      this.alunoId,
      this.uemsId,
      this.ano,
      this.unidade,
      this.curso,
      this.disciplina,
      this.turma,
      this.serieDisciplina,
      this.cargaHorariaPresencial,
      this.maximoFaltas,
      this.periodoLetivo,
      this.professor,
      this.mediaAvaliacoes,
      this.optativa,
      this.exame,
      this.mediaFinal,
      this.faltas,
      this.situacao,
      this.createdAt,
      this.updatedAt});

  Disciplina.fromJson(Map<String, dynamic> json) {
    id = json['id'];
    alunoId = json['aluno_id'];
    uemsId = json['uems_id'];
    ano = json['ano'];
    unidade = json['unidade'];
    curso = json['curso'];
    disciplina = json['disciplina'];
    turma = json['turma'];
    serieDisciplina = json['serie_disciplina'];
    cargaHorariaPresencial = json['carga_horaria_presencial'];
    maximoFaltas = json['maximo_faltas'];
    periodoLetivo = json['periodo_letivo'];
    professor = json['professor'];
    mediaAvaliacoes = json['media_avaliacoes'].toDouble();
    optativa = json['optativa'].toDouble();
    exame = json['exame'].toDouble();
    mediaFinal = json['media_final'].toDouble();
    faltas = json['faltas'];
    situacao = json['situacao'];
    if (json['notas'] != null) {
      notas = new List<Nota>();
      json['notas'].forEach((v) {
        notas.add(new Nota.fromJson(v));
      });
    }
    createdAt = json['created_at'];
    updatedAt = json['updated_at'];
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = new Map<String, dynamic>();
    data['id'] = this.id;
    data['aluno_id'] = this.alunoId;
    data['uems_id'] = this.uemsId;
    data['ano'] = this.ano;
    data['unidade'] = this.unidade;
    data['curso'] = this.curso;
    data['disciplina'] = this.disciplina;
    data['turma'] = this.turma;
    data['serie_disciplina'] = this.serieDisciplina;
    data['carga_horaria_presencial'] = this.cargaHorariaPresencial;
    data['maximo_faltas'] = this.maximoFaltas;
    data['periodo_letivo'] = this.periodoLetivo;
    data['professor'] = this.professor;
    data['media_avaliacoes'] = this.mediaAvaliacoes;
    data['optativa'] = this.optativa;
    data['exame'] = this.exame;
    data['media_final'] = this.mediaFinal;
    data['faltas'] = this.faltas;
    data['situacao'] = this.situacao;
    if (this.notas != null) {
      data['notas'] = this.notas.map((v) => v.toJson()).toList();
    }
    data['created_at'] = this.createdAt;
    data['updated_at'] = this.updatedAt;
    return data;
  }
}
