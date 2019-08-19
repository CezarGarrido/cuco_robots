import 'dart:async';
import 'package:http/http.dart' show Client;
import 'dart:convert';
import 'package:app/entities/disciplina.dart';
import 'package:app/utils/secure_store.dart';
import 'package:app/driver/database.dart';
import 'package:app/constants.dart';

class DisciplinaRepository {
  Client client = Client();

  ConexaoSqlite conexao = new ConexaoSqlite();
  Future<List<Disciplina>> getDisciplinas() async {
    var key = await getSecureStore("jwt");
    Map<String, String> headers = {
      "Authorization": "Bearer " + key,
      'Content-Type': 'application/json; charset=utf-8'
    };
    final response =
        await client.get(BaseUrl + "/disciplinas", headers: headers);

        
    if (response.statusCode == 200) {
      var list = json.decode(utf8.decode(response.bodyBytes)) as List<dynamic>;
      if (list.length <= 0) {
        return await _getDisciplinas();
      }
      var disciplinas = new List<Disciplina>();
      for (Map<String, dynamic> item in list) {
        Disciplina disciplina = new Disciplina.fromJson(item);
        _save(disciplina);
        disciplinas.add(disciplina);
      }
      return disciplinas;
    } else {
      throw ('Errrou!');
    }
  }

  Future<int> _save(Disciplina disciplina) async {
    var db = await conexao.db;
    var result = await db.rawInsert(
        "INSERT OR REPLACE INTO aluno_disciplinas (id, aluno_id, uems_id, unidade, curso, disciplina, turma, serie_disciplina, carga_horaria_presencial,maximo_faltas, periodo_letivo, professor, media_avaliacoes, optativa, exame, media_final, faltas, situacao,created_at, updated_at) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
        [
          disciplina.id,
          disciplina.alunoId,
          disciplina.uemsId,
          disciplina.unidade,
          disciplina.curso,
          disciplina.disciplina,
          disciplina.turma,
          disciplina.serieDisciplina,
          disciplina.cargaHorariaPresencial,
          disciplina.maximoFaltas,
          disciplina.periodoLetivo,
          disciplina.professor,
          disciplina.mediaAvaliacoes,
          disciplina.optativa,
          disciplina.exame,
          disciplina.mediaFinal,
          disciplina.faltas,
          disciplina.situacao,
          disciplina.createdAt,
          disciplina.updatedAt
        ]);

    print("# Disciplina criada $result");
    return result;
  }

  Future<List<Disciplina>> _getDisciplinas() async {
    var db = await conexao.db;
    var result = await db.rawQuery('SELECT * FROM aluno_disciplinas');
    var disciplinas = new List<Disciplina>();
    for (Map<String, dynamic> item in result) {
      disciplinas.add(new Disciplina.fromJson(item));
    }
    return disciplinas;
  }
}
