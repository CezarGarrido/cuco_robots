import 'dart:async';
import 'package:http/http.dart' show Client;

import 'package:app/driver/database.dart';
import 'package:app/entities/endereco.dart';

class EnderecoRepository {
  Client client = Client();
  ConexaoSqlite conexao = new ConexaoSqlite();

  Future<int> save(Endereco endereco) async {
    var db = await conexao.db;
    var result = await db.rawInsert(
        "INSERT OR REPLACE INTO aluno_enderecos (id, aluno_id, logradouro, complemento,bairro,cep,cidade, created_at, updated_at) VALUES(?,?,?,?,?,?,?,?,?)",
        [
          endereco.id,
          endereco.alunoId,
          endereco.logradouro,
          endereco.complemento,
          endereco.bairro,
          endereco.cep,
          endereco.cidade,
          endereco.createdAt,
          endereco.updatedAt,
        ]);

    print("# Endereco criado $result");
    return result;
  }

  Future<List<Endereco>> getContatos(int alunoid) async {
    var db = await conexao.db;
    var result = await db.query('aluno_enderecos',
        where: 'aluno_id = ?', whereArgs: [alunoid]);
    var enderecos = new List<Endereco>();
    for (Map<String, dynamic> item in result) {
      enderecos.add(new Endereco.fromJson(item));
    }
    return enderecos;
  }
}
