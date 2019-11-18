import 'dart:async';
import 'package:http/http.dart' show Client;
import 'package:app/entities/nota.dart';
import 'package:app/driver/database.dart';

class NotaRepository {
  Client client = Client();
  ConexaoSqlite conexao = new ConexaoSqlite();

  Future<int> save(Nota nota) async {

    var db = await conexao.db;

    var result = await db.rawInsert(
      
        "INSERT OR REPLACE INTO aluno_notas (id, disciplina_id, aluno_id, descricao, valor, created_at, updated_at) VALUES(?,?,?,?,?,?,?)",
        [
          nota.id,
          nota.disciplinaId,
          nota.alunoId,
          nota.descricao,
          nota.valor,
          nota.createdAt,
          nota.updatedAt
        ]);

    print("# Nota criada $result");
    return result;
  }

  Future<List<Nota>> getNotas(int disciplinaid) async {
    var db = await conexao.db;
    var result = await db.query('aluno_notas',
        where: 'disciplina_id = ?', whereArgs: [disciplinaid]);
    var notas = new List<Nota>();
    for (Map<String, dynamic> item in result) {
      notas.add(new Nota.fromJson(item));
    }
    return notas;
  }
}
