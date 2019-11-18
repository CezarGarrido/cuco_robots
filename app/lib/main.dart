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
import 'package:app/views/faltas/faltas.dart';
import 'package:app/repository/disciplina.dart';
import 'package:app/entities/disciplina.dart';
import 'package:charts_flutter/flutter.dart' as charts;
import 'dart:math';

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
        "/perfil": (BuildContext context) => new ContactsDemo(),
      },
    ));
  });
}

class MyHomePage extends StatefulWidget {
  @override
  _MyHomePageState createState() => new _MyHomePageState();
}

class _MyHomePageState extends State<MyHomePage> with TickerProviderStateMixin {
  final GlobalKey<ScaffoldState> _scaffoldKey = new GlobalKey<ScaffoldState>();

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
      key: _scaffoldKey,
      appBar: AppBar(
          elevation: 0.5,
          backgroundColor: Colors.white, //Color(0xFF1572E8),
          title: Text(
            "Notas",
            style: TextStyle(
              fontWeight: FontWeight.bold,
              //fontSize: 14.0,
              color: Colors.black,
            ),
          ),
          automaticallyImplyLeading: true,
          //`true` if you want Flutter to automatically add Back Button when needed,
          //or `false` if you want to force your own back button every where

          leading: IconButton(
            color: Colors.black,
            icon: Icon(Icons.menu),
            onPressed: () => _scaffoldKey.currentState.openDrawer(),
          )),
      drawer: Drawer(
        child: ListView(
          children: <Widget>[
            UserAccountsDrawerHeader(
              decoration: BoxDecoration(
                color: Colors.white,
              ),
              accountName: Text(
                nomeAluno,
                style: new TextStyle(color: Colors.black),
              ),
              accountEmail: Text(
                "",
                style: new TextStyle(color: Colors.black),
              ),
              currentAccountPicture: CircleAvatar(
                backgroundColor:
                    Theme.of(context).platform == TargetPlatform.iOS
                        ? Color(0xFF1572E8)
                        : Colors.grey,
                child: Center(
                  // Replace with a Row for horizontal icon + text
                  child: Icon(Icons.photo_camera, color: Colors.black),
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
                    builder: (BuildContext context) => MyHomePage()));
              },
            ),
            /*ListTile(
              leading: Icon(
                Icons.school,
              ),
              title: Text("Disciplinas"),
              onTap: () {
                Navigator.of(context).pop();
                Navigator.of(context).push(MaterialPageRoute(
                    builder: (BuildContext context) => Disciplinas()));
              },
            ),*/
            /*ListTile(
              leading: Icon(
                Icons.check_circle,
              ),
              title: Text("Faltas"),
              onTap: () {
                Navigator.of(context).pop();
                Navigator.of(context).push(MaterialPageRoute(
                    builder: (BuildContext context) => FaltasView()));
              },
            ),*/
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
      body: Notas(),
      /* body: new Padding(
          padding: const EdgeInsets.all(8.0),
          child: new Column(children: <Widget>[
            new Card(
                child: Column(
              mainAxisSize: MainAxisSize.min,
              children: <Widget>[
                new ListTile(
                  // leading: Icon(Icons.album),
                  title: Text('Suas notas'),
                  subtitle: Text('Em breve...'),
                ),
                new SizedBox(
                    height: 250.0,
                    child: Padding(
                      padding: const EdgeInsets.all(8.0),
                      child: new SimpleBarChart.withRandomData(),
                    )),
              ],
            ))
          ])),*/
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

class SimpleBarChart extends StatelessWidget {
  final List<charts.Series> seriesList;
  final bool animate;

  SimpleBarChart(this.seriesList, {this.animate});

  /// Creates a [BarChart] with sample data and no transition.
  factory SimpleBarChart.withSampleData() {
    return new SimpleBarChart(
      _createSampleData(),
      // Disable animations for image tests.
      animate: false,
    );
  }

  // EXCLUDE_FROM_GALLERY_DOCS_START
  // This section is excluded from being copied to the gallery.
  // It is used for creating random series data to demonstrate animation in
  // the example app only.
  factory SimpleBarChart.withRandomData() {
    return new SimpleBarChart(_createRandomData());
  }

  /// Create random data.
  static List<charts.Series<OrdinalSales, String>> _createRandomData() {
    final random = new Random();

    final data = [
      new OrdinalSales('Bio', random.nextInt(100)),
      new OrdinalSales('Hist', random.nextInt(100)),
      new OrdinalSales('Quim', random.nextInt(100)),
      new OrdinalSales('Port', random.nextInt(100)),
      new OrdinalSales('Mat', random.nextInt(100)),
      new OrdinalSales('Geo', random.nextInt(100)),
    ];

    return [
      new charts.Series<OrdinalSales, String>(
          id: 'Notas',
          colorFn: (_, __) => charts.MaterialPalette.blue.shadeDefault,
          domainFn: (OrdinalSales sales, _) => sales.year,
          measureFn: (OrdinalSales sales, _) => sales.sales,
          data: data,
          labelAccessorFn: (OrdinalSales sales, _) =>
              '${sales.sales.toString()}%')
    ];
  }
  // EXCLUDE_FROM_GALLERY_DOCS_END

  @override
  Widget build(BuildContext context) {
    return new charts.BarChart(
      seriesList,
      animate: animate,
      barRendererDecorator: new charts.BarLabelDecorator<String>(),
      domainAxis: new charts.OrdinalAxisSpec(),
    );
  }

  /// Create one series with sample hard coded data.
  static List<charts.Series<OrdinalSales, String>> _createSampleData() {
    final data = [
      new OrdinalSales('2014', 5),
      new OrdinalSales('2015', 25),
      new OrdinalSales('2016', 100),
      new OrdinalSales('2017', 75),
    ];

    return [
      new charts.Series<OrdinalSales, String>(
          id: 'Notas',
          colorFn: (_, __) => charts.MaterialPalette.blue.shadeDefault,
          domainFn: (OrdinalSales sales, _) => sales.year,
          measureFn: (OrdinalSales sales, _) => sales.sales,
          data: data,
          labelAccessorFn: (OrdinalSales sales, _) =>
              '${sales.sales.toString()}%')
    ];
  }
}

/// Sample ordinal data type.
class OrdinalSales {
  final String year;
  final int sales;

  OrdinalSales(this.year, this.sales);
}
