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
  migrate.runMigrate();
  initializeDateFormatting("pt_BR", null);

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
}

class MyHomePage extends StatefulWidget {
  @override
  _MyHomePageState createState() => new _MyHomePageState();
}

class _MyHomePageState extends State<MyHomePage> with TickerProviderStateMixin {
  final menu = new Menu(
    items: [
      new MenuItem(
        id: 'notas',
        title: 'Notas',
      ),
      new MenuItem(
        id: 'disciplinas',
        title: 'Disciplinas',
      ),
      new MenuItem(
        id: 'logout',
        title: 'Sair',
      ),
    ],
  );

  var selectedMenuItemId = 'notas';
  var activeScreen = notasScreen;

  @override
  Widget build(BuildContext context) {
    return new ZoomScaffold(
      menuScreen: new MenuScreen(
        menu: menu,
        selectedItemId: selectedMenuItemId,
        onMenuItemSelected: (String itemId) {
          selectedMenuItemId = itemId;
          if (itemId == 'disciplinas') {
            setState(() => activeScreen = disciplinasScreen);
          } else if (itemId == 'notas') {
            setState(() => activeScreen = notasScreen);
          } else if (itemId == 'logout') {
            logout();
          } else {
            setState(() => activeScreen = notasScreen);
          }
        },
      ),
      contentScreen: activeScreen,
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
