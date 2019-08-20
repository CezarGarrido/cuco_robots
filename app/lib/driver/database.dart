import 'dart:async';
import 'package:path/path.dart';
//import 'package:path_provider/path_provider.dart';
import 'package:sqflite/sqflite.dart';
//import 'dart:io';

const initScript = [
  '''
create table if not exists credenciais (
    id integer primary key autoincrement,
    nome text UNIQUE,
    valor text not null
  );
  ''',
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
	created_at timestamp,
	updated_at timestamp
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
	carga_horaria_presencial integer,
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
	CONSTRAINT aluno_disciplinas_aluno_id_fkey 
    FOREIGN KEY (aluno_id) 
    REFERENCES alunos(id)
    ON DELETE CASCADE
  );
''',
  '''
CREATE TABLE if not exists aluno_notas (
	id integer primary key NOT NULL,
	aluno_id integer NOT NULL,
	disciplina_id integer NOT NULL,
	descricao text,
	valor numeric(18,2),
	created_at timestamp,
	updated_at timestamp,
	CONSTRAINT notas_aluno_id_fkey 
    FOREIGN KEY (aluno_id) 
    REFERENCES alunos(id)
    ON DELETE CASCADE,
  CONSTRAINT notas_disciplina_id_fkey 
    FOREIGN KEY (disciplina_id) 
    REFERENCES aluno_disciplinas(id)
    ON DELETE CASCADE
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

  Future<Database> initDb() async {
    String databasesPath = await getDatabasesPath();
    String path = join(databasesPath, 'app_uems.db');
    var db = await openDatabase(
      path,
      version: migrationScripts.length + 1,
      onOpen: (Database db) {
        db.execute("PRAGMA foreign_keys=ON");
      },
      onCreate: _onCreate,
      onUpgrade: _onUpgrade,
    );
    return db;
  }

  void _onCreate(Database db, int newVersion) async {
    print("# Criando tabelas");
    print(initScript);
    initScript.forEach((script) async => await db.execute(script));
  }

  void _onUpgrade(Database db, int oldVersion, int newVersion) async {
    print("# Atualizando tabelas");
    //for (var i = oldVersion - 1; i <= newVersion - 1; i++) {
    migrationScripts.forEach((f) async {
      print(f);
      await db.execute(f);
    });

    //}
    print("# Tabelas");
    _showTables();
  }

  Future close() async {
    var dbClient = await db;
    return dbClient.close();
  }
}

void runMigrate() async {
  print("# Inicializando banco de dados");
  var conexao = new ConexaoSqlite();
  var dbcli = await conexao.db;
  var version = await dbcli.getVersion();
  print("# Versão do bd $version");
  print("# Tabelas");
  _showTables();
}

void _showTables() async {
  var conexao = new ConexaoSqlite();
  var dbcli = await conexao.db;
  var res = await dbcli.rawQuery(
      "SELECT name FROM sqlite_master WHERE type ='table' AND name NOT LIKE 'sqlite_%';");
  print(res);
}
