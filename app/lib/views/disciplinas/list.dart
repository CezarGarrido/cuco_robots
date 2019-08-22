import 'dart:async';
import 'package:flutter/material.dart';
import 'package:app/repository/disciplina.dart';
import 'package:app/entities/disciplina.dart';
import 'package:app/zoom_scaffold.dart';

final Screen disciplinasScreen = new Screen(
    title: 'Disciplinas',
    contentBuilder: (BuildContext context) {
      return new DisciplinasPage();
    });

class DisciplinasPage extends StatefulWidget {
  @override
  _DiscplinasState createState() => _DiscplinasState();
}

class Periodo {
  const Periodo(this.id, this.name);
  final String name;
  final int id;
}

DisciplinaRepository _disciplinaRepository = DisciplinaRepository();

class _DiscplinasState extends State<DisciplinasPage>
    with SingleTickerProviderStateMixin {
  List<Disciplina> listDisc = List();

  bool _loadingInProgress;
  bool _loadingFailed;
  Periodo selectedPeriodo;
  List<Periodo> periodos = <Periodo>[
    const Periodo(1, '1º Periodo'),
    const Periodo(2, '2º Periodo')
  ];

  Future<Null> _loadList() async {
    try {
      List<Disciplina> listDisciplinas =
          await _disciplinaRepository.getDisciplinas();
      await new Future.delayed(const Duration(seconds: 1));
      if (mounted) {
        setState(() {
          listDisc = listDisciplinas;
          _loadingFailed = false;
          _dataLoaded();
        });
      }
    } on TimeoutException catch (_) {
      setState(() {
        _loadingFailed = true;
        _dataLoaded();
      });
    }
  }

  Widget _loadingView() {
    return new Center(
      child: new CircularProgressIndicator(),
    );
  }

  void _dataLoaded() {
    setState(() {
      _loadingInProgress = false;
    });
  }

  @override
  void initState() {
    selectedPeriodo = periodos[0];
    _loadingInProgress = true;
    _loadingFailed = false;
    _loadList();
    super.initState();
  }

  Widget _myListNoData(BuildContext context) {
    return ListView(
      children: <Widget>[
        Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: <Widget>[
              Padding(
                padding: EdgeInsets.all(5),
              ),
              Icon(
                Icons.error,
              ),
              Text('Sistema indiponível, tente novamente mais tarde'),
            ],
          ),
        )
      ],
    );
  }

  Widget _buildBody(BuildContext context) {
    if (_loadingFailed) {
      return Center(
          child: RefreshIndicator(
        child: _myListNoData(context),
        onRefresh: _loadList,
      ));
    }
    if (_loadingInProgress && !_loadingFailed) {
      return new Center(
        child: _loadingView(),
      );
    } else {
      return Container(
          child: Column(children: <Widget>[
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: <Widget>[
            Container(
                padding: EdgeInsets.only(left: 10.0),
                child: InkWell(
                  onTap: () => {},
                  child: Container(
                      margin: EdgeInsets.only(top: 5.0, bottom: 5.0),
                      child: Row(
                        children: <Widget>[
                          new Container(
                            padding: new EdgeInsets.all(5.0),
                          ),
                          new Icon(Icons.event),
                          new Container(
                            padding: new EdgeInsets.all(5.0),
                          ),
                          new DropdownButtonHideUnderline(
                            child: new DropdownButton<Periodo>(
                              value: selectedPeriodo,
                              items: periodos.map((Periodo periodo) {
                                return new DropdownMenuItem<Periodo>(
                                  value: periodo,
                                  child: new Text(
                                    periodo.name,
                                    style: new TextStyle(color: Colors.black),
                                  ),
                                );
                              }).toList(),
                              onChanged: (Periodo newValue) {
                                setState(() {
                                  selectedPeriodo = newValue;
                                });
                              },
                            ),
                          )
                        ],
                      )),
                )),
          ],
        ),
        Divider(
          height: 0.0,
        ),
        Expanded(
            child: RefreshIndicator(
          child: _buildListview(context),
          onRefresh: _loadList,
        )),
      ]));
    }
  }

  Widget _buildListview(BuildContext context) {
    return ListView.builder(
      itemCount: listDisc.length,
      itemBuilder: (BuildContext context, int index) {
        final data = listDisc[index];
  return Card(
          elevation: 4,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.all(
              Radius.circular(10),
            ),
          ),
          child: ListTile(
            leading: CircleAvatar(
               backgroundImage: AssetImage(
               'https://github.com/JideGuru/FlutterCryptoUI/blob/master/assets/cm1.jpeg',
               ),
              radius: 25,
            ),
            title: Text(data.disciplina),
            subtitle: Text('teste'),
            trailing: Text('102112',
              style: TextStyle(
                color:
                    Colors.green,
                fontWeight: FontWeight.bold,
              ),
            ),
          ),
        );
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    return _buildBody(context);
  }
}
