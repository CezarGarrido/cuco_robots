import 'dart:async';
import 'package:http/http.dart' show Client;

import 'package:app/driver/database.dart';
import 'package:app/entities/contato.dart';

class ContatoRepository {
  Client client = Client();
  ConexaoSqlite conexao = new ConexaoSqlite();

  Future<int> save(Contato contact) async {
    var db = await conexao.db;
    var result = await db.rawInsert(
        "INSERT OR REPLACE INTO aluno_contatos (id, aluno_id, tipo, valor, created_at, updated_at) VALUES(?,?,?,?,?,?)",
        [
          contact.id,
          contact.alunoId,
          contact.tipo,
          contact.valor,
          contact.createdAt,
          contact.updatedAt
        ]);

    print("# Contato criado $result");
    return result;
  }

  Future<List<Contato>> getContatos(int alunoid) async {
    var db = await conexao.db;
    var result = await db.query('aluno_contatos',
        where: 'aluno_id = ?', whereArgs: [alunoid]);
    var contatos = new List<Contato>();
    for (Map<String, dynamic> item in result) {
      contatos.add(new Contato.fromJson(item));
    }
    return contatos;
  }
}
            