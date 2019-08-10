import 'dart:async';
import 'package:path/path.dart';
import 'package:path_provider/path_provider.dart';
import 'package:sqflite/sqflite.dart';
import 'dart:io';

const initScript = [
  '''
create table alunos (
  id integer primary key autoincrement,
  nome text not null,
  email text,
  telefone text,
  rgm text not null,
  senha text not null,
  curso text not null,
  ano int not null,
  unidade text not null,
  logado int not null default 1,
  created_at text not null, 
  updated_at text
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
