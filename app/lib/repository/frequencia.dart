import 'dart:async';
import 'package:http/http.dart' show Client;
import 'package:app/entities/frequencia.dart';
import 'package:app/driver/database.dart';
import 'package:app/utils/secure_store.dart';
import 'package:app/constants.dart';
import 'dart:convert';

class FrequenciaRepository {
  Client client = Client();
  ConexaoSqlite conexao = new ConexaoSqlite();

  Future<List<Frequencia>> getFrequenciasApi(int disciplinaid) async {
    try {
      print("# Buscando faltas");
      var key = await getSecureStore("jwt");
      Map<String, String> headers = {
        "Authorization": "Bearer " + key,
        'Content-Type': 'application/json; charset=utf-8'
      };
      final response = await client.get(
          BaseUrl + "/frequencia/disciplina/$disciplinaid",
          headers: headers);
      if (response.statusCode == 200) {
        var list =
            json.decode(utf8.decode(response.bodyBytes)) as List<dynamic>;
        if (list == null || list.length <= 0) {
          return await getFrequenciasDB(disciplinaid);
        }
        for (Map<String, dynamic> item in list) {
          Frequencia frequencia = new Frequencia.fromJson(item);
          save(frequencia);
        }
        return await getFrequenciasDB(disciplinaid);
      } else {
        return await getFrequenciasDB(disciplinaid);
      }
    } catch (_) {
      print('not connected');
      return await getFrequenciasDB(disciplinaid);
    }
  }

  Future<int> save(Frequencia frequencia) async {
         
    var db = await conexao.db;
    var result = await db.rawInsert(
        "INSERT OR REPLACE INTO aluno_frequencia (id, disciplina_id, aluno_id, mes, dia, valor, created_at, updated_at) VALUES(?,?,?,?,?,?,?,?)",
        [
          frequencia.id,
          frequencia.disciplinaId,
          frequencia.alunoId,
          frequencia.mes,
          frequencia.dia,
          frequencia.valor,
          frequencia.createdAt,
          frequencia.updatedAt
        ]);

    //print("# Frequencia criada $result");
    return result;
  }

  Future<List<Frequencia>> getFrequenciasDB(int disciplinaid) async {
    var db = await conexao.db;
    var result = await db.query('aluno_frequencia',
        where: 'disciplina_id = ?', whereArgs: [disciplinaid]);
    var frequencias = new List<Frequencia>();
    for (Map<String, dynamic> item in result) {
      frequencias.add(new Frequencia.fromJson(item));
    }
    return frequencias;
  }
    Future<List<Frequencia>> getFrequenciasByData(int disciplinaid) async {
    var db = await conexao.db;
    var result = await db.query('aluno_frequencia',
        where: 'disciplina_id = ?', whereArgs: [disciplinaid]);
    var frequencias = new List<Frequencia>();
    for (Map<String, dynamic> item in result) {
      frequencias.add(new Frequencia.fromJson(item));
    }
    return frequencias;
  }
}
