import 'package:app/menu_screen.dart';
import 'package:app/zoom_scaffold.dart';
import 'package:flutter/material.dart';
import 'package:app/views/login/login.dart';
import 'package:intl/date_symbol_data_local.dart';
import 'package:app/driver/database.dart' as migrate;
import 'package:app/repository/aluno.dart';
import 'package:app/entities/aluno.dart';
import 'package:app/views/disciplinas/list.dart';
import 'package:app/views/notas/list.dart';
import 'package:app/views/aluno/perfil.dart';
import 'package:app/views/about/about.dart';
import 'package:app/repository/disciplina.dart';
import 'package:app/entities/disciplina.dart';

//void main() => runApp(new MyApp());
AlunoRepository _alunoRepository = AlunoRepository();
DisciplinaRepository _disciplinaRepository = DisciplinaRepository();

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
        "/perfil":(BuildContext context)=> new ContactsDemo(),
      },
    ));
  });
}

class MyHomePage extends StatefulWidget {
  @override
  _MyHomePageState createState() => new _MyHomePageState();
}

class _MyHomePageState extends State<MyHomePage> with TickerProviderStateMixin {
  Aluno aluno = new Aluno();
  var nomeAluno = "";

  Future<Null> loadAluno() async {
    final alunoDB = await _alunoRepository.getAluno();
    await new Future.delayed(const Duration(seconds: 1));
    if (mounted) {
      setState(() {
        aluno = alunoDB;
        nomeAluno = aluno.nome;
      });
    }
  }
  Future<Null> sincronize() async {
     await _disciplinaRepository.getDisciplinas();
  }

  @override
  void initState() {
    loadAluno();
    sincronize();
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return new Scaffold(
      appBar: AppBar(
        title: Text("Inicio"),
      ),
      drawer: Drawer(
        child: ListView(
          children: <Widget>[
            UserAccountsDrawerHeader(
              accountName: Text(nomeAluno),
              accountEmail: Text("offline"),
              currentAccountPicture: CircleAvatar(
                backgroundColor:
                    Theme.of(context).platform == TargetPlatform.iOS
                        ? Colors.blue
                        : Colors.white,
                child:  Center( // Replace with a Row for horizontal icon + text
                  child: Icon(Icons.photo_camera, color: Colors.grey),
                ),
              ),
            ),
            ListTile(
              leading: Icon(
                Icons.star,
              ),
              title: Text("Notas"),
              onTap: () {
                Navigator.of(context).pop();
                Navigator.of(context).push(MaterialPageRoute(
                    builder: (BuildContext context) => Notas()));
              },
            ),
            ListTile(
              leading: Icon(
                Icons.school,
              ),
              title: Text("Disciplinas"),
              onTap: () {
                Navigator.of(context).pop();
                Navigator.of(context).push(MaterialPageRoute(
                    builder: (BuildContext context) => Notas()));
              },
            ),
            ListTile(
              leading: Icon(
                Icons.check_circle,
              ),
              title: Text("Frequências"),
              onTap: () {
                Navigator.of(context).pop();
                Navigator.of(context).push(MaterialPageRoute(
                    builder: (BuildContext context) => Notas()));
              },
            ),
            ListTile(
                leading: Icon(
                  Icons.settings,
                ),
                title: Text("Configurações"),
                onTap: () {
                  Navigator.pop(context);
                                Navigator.of(context).push(MaterialPageRoute(
                    builder: (BuildContext context) => ContactsDemo()));
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
                  showGalleryAboutDialog(context);
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
