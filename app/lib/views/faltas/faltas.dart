import 'dart:async';
import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import 'package:app/repository/frequencia.dart';
import 'package:app/entities/frequencia.dart';
import 'package:app/utils/collection.dart';

class FaltasView extends StatefulWidget {
  @override
  _FaltasState createState() => _FaltasState();
}

class Periodo {
  const Periodo(this.id, this.name);
  final String name;
  final int id;
}

FrequenciaRepository _frequenciaRepository = FrequenciaRepository();

class _FaltasState extends State<FaltasView>
    with SingleTickerProviderStateMixin {
  List<Frequencia> listDisc = List();
  var listaPorMes = List();

  AnimationController controller;
  bool _loadingInProgress;
  bool _loadingFailed;
  Periodo selectedPeriodo;

  List<Periodo> periodos = <Periodo>[
    const Periodo(1, '1º Periodo'),
    const Periodo(2, '2º Periodo')
  ];

  Future<Null> _loadListDB() async {
    try {
      List<Frequencia> listDisciplinas =
          await _frequenciaRepository.getFrequenciasDB(62);
      await new Future.delayed(const Duration(seconds: 1));
      if (mounted) {
        setState(() {
          listDisc = listDisciplinas;
          listaPorMes = [
            {
              "Fevereiro": listDisc
                  .where((o) => o.mes == "Fevereiro")
                  .toList()
                  .map((f) => {"dia": f.dia, "valor": f.valor})
                  .toList()
            },
          ];

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
      List<Frequencia> listDisciplinas =
          await _frequenciaRepository.getFrequenciasApi(62);
      await new Future.delayed(const Duration(seconds: 1));
      if (mounted) {
        setState(() {
          listDisc = listDisciplinas;

          listaPorMes = [
            {
              "Fevereiro": listDisc
                  .where((o) => o.mes == "Fevereiro")
                  .toList()
                  .map((f) => {"dia": f.dia, "valor": f.valor})
                  .toList()
            },
          ];

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
      itemCount: listaPorMes.length,
      itemBuilder: (BuildContext context, int index) {

        teste(listDisc);

        final mes = listaPorMes[index]["Fevereiro"];
        int faltasLength = 0;
        if (mes != null) {
          faltasLength = mes.length;
        }

        return Padding(
          padding: EdgeInsets.all(10.0),
          child: Column(
            children: <Widget>[
              index > 0 ? Divider() : Text(''),
              ListTile(
                title: Text(
                  "Fevereiro",
                  style: new TextStyle(
                      //  fontWeight: FontWeight.bold,
                      ),
                ),
                subtitle: Text(""),
                // trailing: Icon(
                //   Icons.keyboard_arrow_right,
                //   color: Colors.black,
                // ),
              ),
              Container(
                padding: EdgeInsets.only(top: 2.0, left: 40.0),
                child: ListView.builder(
                  itemCount: faltasLength,
                  shrinkWrap: true,
                  physics: ClampingScrollPhysics(),
                  itemBuilder: (BuildContext context, int indexn) {
                    final data = mes[indexn];
                    //print("# data");
                    var dia = data["dia"] as int;
                    var valor = data["valor"] as String;

                    //int dia = data.dia;
                    return ListTile(
                      leading: CircleAvatar(
                        backgroundColor: Colors.red,
                        radius: 20.0,
                        child: Text(
                          "$valor",
                          style: TextStyle(
                            fontWeight: FontWeight.bold,
                            color: Colors.white,
                          ),
                        ),
                      ),
                      contentPadding: EdgeInsets.all(10.0),
                      title: new Row(children: <Widget>[
                        new Expanded(child: new Text("$valor")),
                        // new Expanded(child: new Text('')),
                        // new Expanded(child: new Text('')),
                        /*new Expanded(
                                  child: new Text(
                                new DateFormat.MMMd("pt_BR")
                                    .format(dataAtualizada),
                                style: TextStyle(
                                  fontSize: 12,
                                ),
                              )),*/
                      ]),
                      subtitle: Text(
                        '$dia de ' + 'Fevereiro',
                        maxLines: 2,
                        overflow: TextOverflow.ellipsis,
                      ),
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
    return new Scaffold(
        appBar: AppBar(
          backgroundColor: Color(0xFF1572E8),
          title: Text("Faltas"),
        ),
        body: _buildBody(context));
  }
}

var meses = [
  "Janeiro",
  "Fevereiro",
  "Março",
  "Abril",
  "Maio",
  "Junho",
  "Julho",
  "Agosto",
  "Setembro",
  "Outubro",
  "Novembro",
  "Dezembro"
];

class Filter {
  String mes;
  List faltas;
}

class Filtros {
  List<Filter> frequencias;
}

void teste(List<Frequencia> array) {
  List lista = new List();

   print("array");
  print(array);
  
  for (var item in meses) {
    List listaPorMes = [
      {
        item: array
            .where((o) => o.mes == item)
            .toList()
            .map((f) => {"dia": f.dia, "valor": f.valor})
            .toList()
      },
    ];
    if(listaPorMes!=null){
        lista.add(listaPorMes);
    }
    
  }
  print(lista);
}
