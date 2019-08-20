import 'package:app/driver/database.dart';

void setSecureStore(String key, String data) async {
  final conexao = new ConexaoSqlite();
  var db = await conexao.db;
  await db.rawInsert(
      "INSERT OR REPLACE INTO credenciais(nome,valor) values(?,?)",
      [key, data]);
}

Future<String> getSecureStore(String key) async {
  final conexao = new ConexaoSqlite();
  var db = await conexao.db;
  var result =
      await db.query("credenciais", where: 'nome = ?', whereArgs: [key]);
  String valor;
  result.forEach((row) => valor = row['valor']);
  return valor;
}

Future<int> removeSecureStore(String key) async {
  final conexao = new ConexaoSqlite();
  final db = await conexao.db;
  return await db.delete("credenciais", where: "nome = ?", whereArgs: [key]);
}

void getCallbackSecureStore(String key, Function callback) async {
  final conexao = new ConexaoSqlite();
  var db = await conexao.db;
  await db
      .rawQuery("select nome,valor from credenciais where nome=$key")
      .then(callback);
}
