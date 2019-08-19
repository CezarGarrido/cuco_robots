import 'dart:async';
import 'package:flutter/material.dart';
import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:app/constants.dart';
import 'package:intl/intl.dart';
import 'package:app/repository/disciplina.dart';
import 'package:app/entities/disciplina.dart';

class Notas extends StatefulWidget {
  @override
  _NotasState createState() => _NotasState();
}

class Periodo {
  const Periodo(this.id, this.name);
  final String name;
  final int id;
}

DisciplinaRepository _disciplinaRepository = DisciplinaRepository();

class _NotasState extends State<Notas> with SingleTickerProviderStateMixin {
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
        int notasLength = data.notas.length;
        return Padding(
          padding: EdgeInsets.all(16.0),
          child: Column(
            children: <Widget>[
              ListTile(
                title: Text(
                  data.disciplina,
                ),
                trailing: Icon(Icons.keyboard_arrow_right),
              ),
              notasLength == 0
                  ? Center(
                      child: Column(
                        children: <Widget>[
                          Padding(
                            padding: EdgeInsets.all(5),
                          ),
                          Icon(
                            Icons.error,
                          ),
                          Text('Não há notas'),
                        ],
                      ),
                    )
                  : Container(
                      padding: EdgeInsets.only(top: 2.0),
                      child: ListView.builder(
                        itemCount: notasLength,
                        shrinkWrap: true,
                        physics: ClampingScrollPhysics(),
                        itemBuilder: (BuildContext context, int indexn) {
                          final nota = data.notas[indexn];
                          Color cor = Colors.green[300];
                          final intValue = nota.valor;
                          if (intValue < 6) {
                            cor = Colors.redAccent;
                          }
                          String anotherValue = '$intValue';
                          DateTime dataAtualizada =
                              DateTime.parse(nota.updatedAt);
                          return ListTile(
                            leading: CircleAvatar(
                              backgroundColor: cor,
                              child: Text(
                                anotherValue,
                                style: TextStyle(
                                  color: Colors.white,
                                ),
                              ),
                            ),
                            contentPadding: EdgeInsets.all(10.0),
                            title: new Row(children: <Widget>[
                              new Expanded(child: new Text(nota.descricao)),
                              new Expanded(
                                  child: new Text(
                                new DateFormat.MMMd("pt_BR")
                                    .format(dataAtualizada),
                                style: TextStyle(
                                  fontSize: 12,
                                ),
                              )),
                            ]),
                            subtitle: Text(
                              anotherValue,
                              maxLines: 2,
                              overflow: TextOverflow.ellipsis,
                            ),
                          );
                        },
                      ),
                    )
            ],
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
