import 'package:app/menu_screen.dart';
import 'package:app/zoom_scaffold.dart';
import 'package:flutter/material.dart';
import 'package:app/views/login/login.dart';
import 'package:intl/date_symbol_data_local.dart';
import 'package:app/driver/database.dart' as migrate;
import 'package:app/repository/aluno.dart';
import 'package:app/views/disciplinas/list.dart';
import 'package:app/views/notas/list.dart';

//void main() => runApp(new MyApp());
AlunoRepository _alunoRepository = AlunoRepository();

Future<void> main() async {
  initializeDateFormatting("pt_BR", null).then((onValue) async {
    migrate.runMigrate();
    Widget _defaultHome = new LoginPage();

    bool isLogado = await _alunoRepository.isLoggedIn();
    if (isLogado) {
      _defaultHome = new MyHomePage();
    }

    runApp(new MaterialApp(
      title: 'App',
      debugShowCheckedModeBanner: false,
      home: _defaultHome,
      routes: <String, WidgetBuilder>{
        "/app": (BuildContext context) => new MyHomePage(),
        "/login": (BuildContext context) => new LoginPage(),
      },
    ));
  });
}

class MyHomePage extends StatefulWidget {
  @override
  _MyHomePageState createState() => new _MyHomePageState();
}

class _MyHomePageState extends State<MyHomePage> with TickerProviderStateMixin {
  @override
  Widget build(BuildContext context) {
    return new Scaffold(
      appBar: AppBar(
        title: Text("Uems"),
      ),
      drawer: Drawer(
        child: ListView(
          children: <Widget>[
            UserAccountsDrawerHeader(
              accountName: Text("Cezar Garrido Britez"),
              accountEmail: Text("offline"),
              currentAccountPicture: CircleAvatar(
                backgroundColor:
                    Theme.of(context).platform == TargetPlatform.iOS
                        ? Colors.blue
                        : Colors.white,
                child: Text(
                  "C",
                  style: TextStyle(fontSize: 40.0),
                ),
              ),
            ),
            ListTile(
              leading: Icon(Icons.stars,color: Colors.blue,),
              title: Text("Notas"),
              onTap: () {
                Navigator.of(context).pop();
                Navigator.of(context).push(MaterialPageRoute(
                    builder: (BuildContext context) => Notas()));
              },
            ),
            ListTile(
              leading: Icon(Icons.school,color: Colors.blue,),
              title: Text("Disciplinas"),
              onTap: () {
                Navigator.of(context).pop();
                Navigator.of(context).push(MaterialPageRoute(
                    builder: (BuildContext context) => Notas()));
              },
            ),
            ListTile(
              leading: Icon(Icons.check_circle,color: Colors.blue,),
              title: Text("Frequências"),
              onTap: () {
                Navigator.of(context).pop();
                Navigator.of(context).push(MaterialPageRoute(
                    builder: (BuildContext context) => Notas()));
              },
            ),
            ListTile(
                leading: Icon(Icons.settings,color: Colors.blue,),
                title: Text("Configurações"),
                onTap: () {
                  Navigator.pop(context);
                }),
            Divider(),
            ListTile(
                leading: Icon(Icons.insert_comment),
                title: Text("Críticas ou Sugestões"),
                onTap: () {
                  Navigator.pop(context);
                }),
            ListTile(
                leading: Icon(Icons.info),
                title: Text("Sobre"),
                onTap: () {
                  Navigator.pop(context);
                }),
            ListTile(
              title: Text("Sair"),
              leading: Icon(Icons.power_settings_new),
              onTap: () {
                logout();
              },
            ),
          ],
        ),
      ),
    );
  }

  void logout() async {
    _alunoRepository.delete().then((onValue) {
      print(onValue);
      Navigator.of(context)
          .pushNamedAndRemoveUntil('/login', (Route<dynamic> route) => false);
    });
  }
}
