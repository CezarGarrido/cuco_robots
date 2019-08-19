import 'package:app/menu_page.dart';
import 'package:flutter/material.dart';
import 'package:app/zoom_scaffold.dart';
import 'package:provider/provider.dart';
import 'package:app/views/login/login.dart';
import 'package:intl/date_symbol_data_local.dart';
import 'package:app/driver/database.dart' as migrate;
import 'package:app/repository/aluno.dart';

//void main() => runApp(new MyApp());
AlunoRepository _alunoRepository = AlunoRepository();

Future<void> main() async {
  Widget _defaultHome = new LoginPage();

  bool isLogado = await _alunoRepository.isLoggedIn();
  if (isLogado) {
    _defaultHome = new MyHomePage();
  }
  migrate.runMigrate();
  initializeDateFormatting("pt_BR", null);

  runApp(new MaterialApp(
    title: 'App',
    debugShowCheckedModeBanner: false,
    home: _defaultHome,
    routes: <String, WidgetBuilder>{
      "/app": (BuildContext context) => new MyHomePage(),
    },
  ));
}

class MyHomePage extends StatefulWidget {
  @override
  _MyHomePageState createState() => new _MyHomePageState();
}

class _MyHomePageState extends State<MyHomePage> with TickerProviderStateMixin {
  MenuController menuController;
  @override
  void initState() {
    super.initState();
    menuController = new MenuController(
      vsync: this,
    )..addListener(() => setState(() {}));
  }

  @override
  void dispose() {
    menuController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return ChangeNotifierProvider(
      builder: (context) => menuController,
      child: ZoomScaffold(
        menuScreen: MenuScreen(),
        contentScreen: Layout(
            contentBuilder: (cc) => Container(
                  color: Colors.grey[200],
                  child: Container(
                    color: Colors.grey[200],
                  ),
                )),
      ),
    );
  }
}
