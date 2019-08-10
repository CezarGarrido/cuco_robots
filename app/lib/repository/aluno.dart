import 'dart:async';
import 'package:http/http.dart' show Client;
import 'dart:convert';
import 'package:app/entities/aluno.dart';
import 'package:app/utils/jwt.dart';
import 'package:app/driver/database.dart';

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

  Future<Aluno> findByID(int id) async {
    final response = await client.get("https://api.balta.io/v1/courses/$id");
    if (response.statusCode == 200) {
      var map = json.decode(response.body) as Map<String, dynamic>;
      return new Aluno(
          id: map["id"],
          nome: map["nome"],
          email: map["email"],
          telefone: map["telefone"],
          rgm: map["rgm"],
          senha: map["senha"]);
    } else {
      throw Exception('Deu ruim!');
    }
  }

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
      "http://192.168.0.104:8091/api/v1/login",
      headers: {"content-type": "application/json"},
      body: data,
    );
    var token = json.decode(response.body);
    if (response.statusCode != 200) return null;
    var map = parseJwt(token);
    var auxAluno = map["aluno"] as Map<String, dynamic>;
    print(auxAluno);
    var createdAt = DateTime.parse(auxAluno["created_at"]);

    var aluno = new Aluno(
        id: auxAluno["id"] as int,
        nome: auxAluno["nome"],
        email: auxAluno["email"],
        telefone: auxAluno["telefone"],
        rgm: auxAluno["rgm"],
        senha: auxAluno["senha"],
        curso: auxAluno["curso"],
        ano: auxAluno["ano"] as int,
        unidade: auxAluno["unidade"],
        createdAt: createdAt);
    if (auxAluno["updated_at"] != null) {
      aluno.updatedAt = DateTime.parse(auxAluno["updated_at"]);
    } else {
      aluno.updatedAt = aluno.createdAt;
    }
    return aluno;
  }

  Future<List> getAll() async {
    var conexao = new ConexaoSqlite();
    var db = await conexao.db;
    //var result = await db.query(tableNote, columns: [columnId, columnTitle, columnDescription]);
    var result = await db.rawQuery('SELECT * FROM alunos');
    return result.toList();
  }

  Future<int> createSQL(Aluno aluno) async {
    var db = await conexao.db;
    var table = await db.rawQuery("SELECT MAX(id)+1 as id FROM alunos");
    int id = table.first["id"];
    var raw = await db.rawInsert(
        "INSERT Into alunos (id,nome,email,telefone,rgm,senha,curso,ano,unidade,logado,created_at,updated_at)"
        " VALUES (?,?,?,?,?,?,?,?,?,?,?,?)",
        [
          id,
          aluno.nome,
          aluno.email,
          aluno.telefone,
          aluno.rgm,
          aluno.senha,
          aluno.curso,
          aluno.ano,
          aluno.unidade,
          1,
          aluno.createdAt.toIso8601String(),
          aluno.updatedAt.toIso8601String(),
        ]);
    return raw;
  }

  Future<Aluno> findIsLogado() async {
    var conexao = new ConexaoSqlite();
    var db = await conexao.db;
    //var result = await db.query(tableNote, columns: [columnId, columnTitle, columnDescription]);
    var result = await db.query("alunos", where: "logado = ?", whereArgs: [1]);
    print(result.length);

    if (result.length == 0) {
      return Aluno(logado: 0);
    }

    var auxAluno = result.first;
    var createdAt = DateTime.parse(auxAluno["created_at"]);
    var aluno = new Aluno(
        id: auxAluno["id"] as int,
        nome: auxAluno["nome"],
        email: auxAluno["email"],
        telefone: auxAluno["telefone"],
        rgm: auxAluno["rgm"],
        senha: auxAluno["senha"],
        curso: auxAluno["curso"],
        ano: auxAluno["ano"] as int,
        unidade: auxAluno["unidade"],
        logado: auxAluno["logado"] as int,
        createdAt: createdAt);
    if (auxAluno["updated_at"] != null) {
      aluno.updatedAt = DateTime.parse(auxAluno["updated_at"]);
    } else {
      aluno.updatedAt = aluno.createdAt;
    }
    print(aluno);
    return aluno;
  }
}
