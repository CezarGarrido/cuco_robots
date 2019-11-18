import 'dart:async';
import 'package:flutter/material.dart';
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

  AnimationController controller;
  bool _loadingInProgress;
  bool _loadingFailed;
  Periodo selectedPeriodo;
  List<Periodo> periodos = <Periodo>[
    const Periodo(0, 'Selecione a série'),
    const Periodo(1, '1ª Série'),
    const Periodo(2, '2ª Série')
  ];

  Future<Null> _loadListDB() async {
    try {
      List<Disciplina> listDisciplinas =
          await _disciplinaRepository.getDisciplinasDB();
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
//getDisciplinasByPeriodo

  Future<Null> _loadListBySerie(String serie) async {
    try {
      List<Disciplina> listDisciplinas =
          await _disciplinaRepository.getDisciplinasBySerie(serie);

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
    _loadListDB();
    _loadList();
    super.initState();
    controller =
        AnimationController(duration: Duration(milliseconds: 900), vsync: this);
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
                color: Colors.grey[300],
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
                                print(newValue.name);
                                if (newValue.id != 0) {
                                  _loadListBySerie(newValue.name);
                                }
                              },
                            ),
                          )
                        ],
                      )),
                )),
          ],
        ),
        /* Divider(
          height: 0.0,
        ),*/
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
        int notasLength = 0;
        if (data.notas != null) {
          notasLength = data.notas.length;
        }
        Color corMedia = Colors.green;
        final mediaFormated = data.mediaAvaliacoes;
        if (mediaFormated < 6) {
          corMedia = Colors.red; //Color(0xFFF25961);
        }
        String mediaValue = formatValor(mediaFormated);
        return Padding(
          padding: EdgeInsets.all(0.0),
          child: Column(
            children: <Widget>[
              index > 0
                  ? Padding(
                      padding: new EdgeInsets.only(left: 94.0, right: 0.0),
                      child: Divider(),
                    )
                  : Text(''),
              ListTile(
                leading: CircleAvatar(
                  backgroundColor: corMedia,
                  radius: 30.0,
                  child: Text(
                    mediaValue,
                    style: TextStyle(
                      fontWeight: FontWeight.w600,
                      color: Colors.white,
                    ),
                  ),
                ),
                title: Text(
                  data.disciplina,
                  style: new TextStyle(
                    fontWeight: FontWeight.bold,
                  ),
                ),
                subtitle: Text(mediaFormated < 6
                    ? 'Quase lá 👌'
                    : 'Parabéns 👏👏, você passou ✔️'),
                // trailing: Icon(
                //   Icons.keyboard_arrow_right,
                //   color: Colors.black,
                // ),
              ),
              notasLength == 0
                  ? Center(
                      child: Column(
                        children: <Widget>[
                          Padding(
                            padding: EdgeInsets.all(3),
                          ),
                          // FloatingActionButton(
                          //     child: Icon(Icons.cloud_circle),
                          //     onPressed: () {},
                          //     backgroundColor: Colors.blue),
                          // Icon(
                          //   Icons.error,
                          //   color: Colors.blue,
                          // ),
                          // Text('Aguardando lançamentos'),
                          //#6c757d!important
                          /* Text(
                            'Nenhuma nota foi lançada para esta disciplina.',
                            style: TextStyle(
                              //fontWeight: FontWeight.bold,
                              fontSize: 14.0,
                              color: Color(0xFF6C757D),
                            ),
                          ),*/
                        ],
                      ),
                    )
                  : Container(
                      padding: EdgeInsets.only(top: 0.0, left: 85.0),
                      child: ListView.builder(
                        itemCount: notasLength,
                        shrinkWrap: true,
                        physics: ClampingScrollPhysics(),
                        itemBuilder: (BuildContext context, int indexn) {
                          final nota = data.notas[indexn];
                          Color cor = Colors.green; //Color(0xFF31CE36);
                          final intValue = nota.valor;
                          if (intValue > 4 && intValue < 6) {
                            cor = Colors.amber; //LARANJA
                          } else if (intValue < 4) {
                            cor = Colors.red;
                          }
                          //ffad46
                          String anotherValue = formatValor(intValue);
                          DateTime dataAtualizada =
                              DateTime.parse(nota.updatedAt);
                          return ListTile(
                            leading: CircleAvatar(
                              backgroundColor: cor,
                              radius: 20.0,
                              child: Text(
                                anotherValue,
                                style: TextStyle(
                                  fontWeight: FontWeight.bold,
                                  color: Colors.white,
                                ),
                              ),
                            ),
                            contentPadding: EdgeInsets.all(5.0),
                            title: new Row(children: <Widget>[
                              new Expanded(child: new Text(nota.descricao)),
                              new Expanded(child: new Text('')),
                              //new Expanded(child: new Text('')),
                              new Expanded(
                                  child: new Text(
                                new DateFormat.MMMd("pt_BR")
                                    .format(dataAtualizada),
                                style: TextStyle(
                                  fontSize: 12,
                                ),
                              )),
                            ]),
                            /* subtitle: Text(
                              new DateFormat.MMMd("pt_BR")
                                  .format(dataAtualizada),
                              maxLines: 2,
                              overflow: TextOverflow.ellipsis,
                            ),*/
                          );
                        },
                      ),
                    ),
            ],
          ),
        );
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    return new Scaffold(body: _buildBody(context));
  }
}

String formatValor(double valor) {
  if (valor == 10.0) {
    return "10";
  }
  if (valor == 5.0) {
    return "5";
  }
  if (valor == 6.0) {
    return "6";
  }
  if (valor == 7.0) {
    return "7";
  }
  if (valor == 8.0) {
    return "8";
  }
  if (valor == 9.0) {
    return "9";
  }
  if (valor == 4.0) {
    return "4";
  }
  if (valor == 3.0) {
    return "3";
  }
  if (valor == 2.0) {
    return "2";
  }
  if (valor == 1.0) {
    return "7";
  }
  if (valor == 0.0) {
    return "0";
  }
  return "$valor";
}
