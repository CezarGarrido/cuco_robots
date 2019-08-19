import 'dart:async';
import 'package:http/http.dart' show Client;
import 'dart:convert';
import 'package:app/entities/aluno.dart';
import 'package:app/utils/jwt.dart';
import 'package:app/utils/secure_store.dart';
import 'package:app/driver/database.dart';
import 'package:app/constants.dart';

class AlunoRepository {
  Client client = Client();
  ConexaoSqlite conexao = new ConexaoSqlite();

  Future<Aluno> login(String rgm, String senha) async {
    var data = json.encode({
      "rgm": rgm,
      "senha": senha,
    });
    final response = await client.post(
      BaseUrl + "/login",
      headers: {"content-type": "application/json"},
      body: data,
    );
    var res = json.decode(utf8.decode(response.bodyBytes));
    if (response.statusCode != 200) {
      final String err = res['error_message'];
      throw ('$err');
    } else {
      var token = parseJwt(res);
      var auxAluno = token["aluno"] as Map<String, dynamic>;
      var aluno = new Aluno.fromJson(auxAluno);
      save(aluno);
      setSecureStore("jwt", res);
      return aluno;
    }
  }

  Future<Aluno> getAluno(int id) async {
    var db = await conexao.db;
    var result = await db.query("alunos", where: "id = ?", whereArgs: [id]);
    return result.isNotEmpty ? Aluno.fromJson(result.first) : Null;
  }

  Future<int> save(Aluno aluno) async {
    var db = await conexao.db;
    return await db.rawInsert(
        "INSERT OR REPLACE INTO alunos (id, guid, nome, rgm, senha, curso, data_nascimento, sexo, nome_pai, nome_mae, estado_civil, nacionalidade,naturalidade, fenotipo, cpf, rg, rg_orgao_emissor, rg_estado_emissor, rg_data_emissao , created_at, updated_at) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
        [
          aluno.id,
          aluno.guid,
          aluno.nome,
          aluno.rgm,
          aluno.senha,
          aluno.curso,
          aluno.dataNascimento,
          aluno.sexo,
          aluno.nomePai,
          aluno.nomeMae,
          aluno.estadoCivil,
          aluno.nacionalidade,
          aluno.naturalidade,
          aluno.fenotipo,
          aluno.cpf,
          aluno.rg,
          aluno.rgOrgaoEmissor,
          aluno.rgEstadoEmissor,
          aluno.rgDataEmissao,
          aluno.createdAt,
          aluno.updatedAt
        ]);
  }

  Future<bool> isLoggedIn() async {
    var db = await conexao.db;
    var res = await db.query("alunos");
    return res.length > 0 ? true : false;
  }
}
