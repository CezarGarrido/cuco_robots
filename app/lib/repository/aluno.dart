import 'dart:async';
import 'package:http/http.dart' show Client;
import 'dart:convert';
import 'package:app/entities/aluno.dart';
import 'package:app/utils/jwt.dart';
import 'package:app/driver/database.dart';
import 'package:app/constants.dart';

class AlunoRepository {
  Client client = Client();
  ConexaoSqlite conexao = new ConexaoSqlite();

  Future<List<Aluno>> getAlunos() async {
    final response = await client.get("https://api.balta.io/v1/courses");
    if (response.statusCode == 200) {
      var list = json.decode(response.body) as List<dynamic>;
      var alunos = new List<Aluno>();
      for (dynamic item in list) {
        Aluno aluno = new Aluno(id: item["title"]);
        alunos.add(aluno);
      }
      return alunos;
    } else {
      throw Exception('Deu ruim!');
    }
  }

  // Future<Aluno> findByID(int id) async {
  //   final response = await client.get("https://api.balta.io/v1/courses/$id");
  //   if (response.statusCode == 200) {
  //     var map = json.decode(response.body) as Map<String, dynamic>;
  //     return new Aluno(
  //         id: map["id"],
  //         nome: map["nome"],
  //         email: map["email"],
  //         telefone: map["telefone"],
  //         rgm: map["rgm"],
  //         senha: map["senha"]);
  //   } else {
  //     throw Exception('Deu ruim!');
  //   }
  // }

  Future<bool> createApi(Aluno data) async {
    String jsonData = json.encode(data);
    final response = await client.post(
      "",
      headers: {"content-type": "application/json"},
      body: jsonData,
    );
    if (response.statusCode == 201) {
      return true;
    } else {
      return false;
    }
  }

  Future<bool> update(Aluno data) async {
    String jsonData = json.encode(data);
    final response = await client.put(
      "",
      headers: {"content-type": "application/json"},
      body: jsonData,
    );
    if (response.statusCode == 200) {
      return true;
    } else {
      return false;
    }
  }

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
    var token = json.decode(utf8.decode(response.bodyBytes));
    if (response.statusCode != 200) {
      final String err = token['message'];
      throw ('$err');
    } else {
      var map = parseJwt(token);
      var auxAluno = map["aluno"] as Map<String, dynamic>;
      var aluno = new Aluno.fromJson(auxAluno);
      save(aluno);
      return aluno;
    }
  }

  Future<List> getAll() async {
    var conexao = new ConexaoSqlite();
    var db = await conexao.db;
    //var result = await db.query(tableNote, columns: [columnId, columnTitle, columnDescription]);
    var result = await db.rawQuery('SELECT * FROM alunos');
    return result.toList();
  }

  Future<int> save(Aluno aluno) async {
    var db = await conexao.db;
    print("# Salvando aluno");
    // aluno.id = await db.insert("alunos", aluno.toJson());
    var result = await db.rawInsert(
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

    print("# aluno criado $result");
    return result;
  }

  //Future<Aluno> findIsLogado() async {
  // var conexao = new ConexaoSqlite();
  // var db = await conexao.db;
  // //var result = await db.query(tableNote, columns: [columnId, columnTitle, columnDescription]);
  // var result = await db.query("alunos", where: "logado = ?", whereArgs: [1]);
  // print(result.length);

  // if (result.length == 0) {
  //   return Aluno(logado: 0);
  // }

  // var auxAluno = result.first;
  // var createdAt = DateTime.parse(auxAluno["created_at"]);
  // var aluno = new Aluno(
  //     id: auxAluno["id"] as int,
  //     nome: auxAluno["nome"],
  //     email: auxAluno["email"],
  //     telefone: auxAluno["telefone"],
  //     rgm: auxAluno["rgm"],
  //     senha: auxAluno["senha"],
  //     curso: auxAluno["curso"],
  //     ano: auxAluno["ano"] as int,
  //     unidade: auxAluno["unidade"],
  //     logado: auxAluno["logado"] as int,
  //     createdAt: createdAt);
  // if (auxAluno["updated_at"] != null) {
  //   aluno.updatedAt = DateTime.parse(auxAluno["updated_at"]);
  // } else {
  //   aluno.updatedAt = aluno.createdAt;
  // }
  // print(aluno);
  // return aluno;
  // }
}
