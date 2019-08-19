import 'dart:async';
import 'package:path/path.dart';
import 'package:path_provider/path_provider.dart';
import 'package:sqflite/sqflite.dart';
import 'dart:io';

const initScript = [
  '''
create table if not exists alunos (
  id integer primary key autoincrement,
	guid text NULL,
	nome text NOT NULL,
	rgm text unique,
	senha text NOT NULL,
  curso text,
	data_nascimento text,
	sexo text,
	nome_pai text,
	nome_mae text,
	estado_civil text,
	nacionalidade text,
	naturalidade text,
	fenotipo text,
	cpf text,
	rg text,
	rg_orgao_emissor text,
	rg_estado_emissor text ,
	rg_data_emissao text,
	created_at text,
	updated_at text
  );
 ''',
  '''
 CREATE TABLE if not exists aluno_disciplinas (
	id integer primary key NOT NULL,
	aluno_id integer NOT NULL,
	uems_id integer NOT NULL,
	unidade text,
	curso text,
	disciplina text,
	turma text,
	serie_disciplina text,
	carga_horaria_presencial text,
	maximo_faltas integer,
	periodo_letivo text,
	professor text,
	media_avaliacoes numeric(18,2),
	optativa numeric(18,2),
	exame numeric(18,2),
	media_final numeric(18,2),
	faltas integer,
	situacao text,
	created_at timestamp,
	updated_at timestamp,
	CONSTRAINT aluno_disciplinas_aluno_id_fkey FOREIGN KEY (aluno_id) REFERENCES alunos(id)
  );
 '''
];
const migrationScripts = [];

class ConexaoSqlite {
  static final ConexaoSqlite _instance = new ConexaoSqlite.internal();
  factory ConexaoSqlite() => _instance;
  static Database _db;
  ConexaoSqlite.internal();

  Future<Database> get db async {
    if (_db != null) {
      return _db;
    }
    _db = await initDb();
    return _db;
  }

  initDb() async {
    String databasesPath = await getDatabasesPath();
    String path = join(databasesPath, 'hello_uems.db');
    var db = await openDatabase(
      path,
      version: migrationScripts.length + 1,
      onCreate: _onCreate,
      onUpgrade: _onUpgrade,
    );
    return db;
  }

  void _onCreate(Database db, int newVersion) async {
    print(initScript);
    initScript.forEach((script) async => await db.execute(script));
  }

  void _onUpgrade(Database db, int oldVersion, int newVersion) async {
    for (var i = oldVersion - 1; i <= newVersion - 1; i++) {
      await db.execute(migrationScripts[i]);
    }
  }

  Future close() async {
    var dbClient = await db;
    return dbClient.close();
  }
}

void runMigrate() async {
  var conexao = new ConexaoSqlite();
  var dbcli = await conexao.db;
  var version = await dbcli.getVersion();

  print("# DB version");
  print(version);
  print("# Executando migrations");
  initScript.forEach((script) async => await dbcli.execute(script));
  showTables();
}

void showTables() async {
  print("# Show tables");
  var conexao = new ConexaoSqlite();
  var dbcli = await conexao.db;
   var res = await dbcli.rawQuery(
      "SELECT name FROM sqlite_master WHERE type ='table' AND name NOT LIKE 'sqlite_%';");
  print(res);

  /* */
}
