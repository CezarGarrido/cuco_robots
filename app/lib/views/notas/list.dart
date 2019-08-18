import 'dart:async';
import 'package:flutter/material.dart';
import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:app/constants.dart';
import 'package:intl/intl.dart';

class Notas extends StatefulWidget {
  @override
  _NotasState createState() => _NotasState();
}

class Periodo {
  const Periodo(this.id, this.name);
  final String name;
  final int id;
}

class _NotasState extends State<Notas> with SingleTickerProviderStateMixin {
  var list = List();
  bool _loadingInProgress;
  bool _loadingFailed;
  Periodo selectedPeriodo;
  List<Periodo> periodos = <Periodo>[
    const Periodo(1, '1º Periodo'),
    const Periodo(2, '2º Periodo')
  ];

  Future<Null> _loadList() async {
    try {
      Map<String, String> headers = {
        "Authorization": "Bearer " +
            "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhbHVubyI6eyJpZCI6NiwiZ3VpZCI6IjYzNDc0NzQxLTU5ZjMtNGNkOC1hZjQxLTViM2NkM2MxNWNiZSIsIm5vbWUiOiJDRVpBUiBHQVJSSURPIEJSSVRFWiIsInJnbSI6IjQwMDg5Iiwic2VuaGEiOiJDMTAyMDMwZyIsImN1cnNvIjoiIiwiZGF0YV9uYXNjaW1lbnRvIjoiMTk5Ny0xMi0yOFQwMDowMDowMFoiLCJzZXhvIjoiTWFzY3VsaW5vIiwibm9tZV9wYWkiOiJWSVRPUiBCUklURVoiLCJub21lX21hZSI6Ik1BUklBTkEgR0FSUklETyIsImVzdGFkb19jaXZpbCI6IlNvbHRlaXJvKGEpIiwibmFjaW9uYWxpZGFkZSI6IkJSQVNJTEVJUk8iLCJuYXR1cmFsaWRhZGUiOiJQQVJBTkhPUy9NUyIsImZlbm90aXBvIjoiIiwiY3BmIjoiMDUwLjQzMy42OTEtNjciLCJyZyI6IjIuMjI1LjIyOCIsInJnX29yZ2FvX2VtaXNzb3IiOiJNRCIsInJnX2VzdGFkb19lbWlzc29yIjoiTVMiLCJyZ19kYXRhX2VtaXNzYW8gIjoiMDAwMS0wMS0wMVQwMDowMDowMFoiLCJjb250YXRvcyI6W3siaWQiOjEsImFsdW5vX2lkIjo2LCJ0aXBvIjoiVGVsZWZvbmUiLCJ2YWxvciI6Iig2NykgOTk2ODItMjQwMiIsImNyZWF0ZWRfYXQiOiIyMDE5LTA4LTE3VDE4OjUzOjI4Ljk4MzcyNloiLCJ1cGRhdGVkX2F0IjpudWxsfSx7ImlkIjoyLCJhbHVub19pZCI6NiwidGlwbyI6IkVtYWlsIiwidmFsb3IiOiJjZXphci5jZ2IxOEBnbWFpbC5jb20iLCJjcmVhdGVkX2F0IjoiMjAxOS0wOC0xN1QxODo1MzoyOC45ODM3MjZaIiwidXBkYXRlZF9hdCI6bnVsbH1dLCJlbmRlcmVjb3MiOlt7ImlkIjoxLCJhbHVub19pZCI6NiwibG9ncmFkb3VybyI6ImNmd2VmIiwibnVtZXJvIjo1MzU0MywiY29tcGxlbWVudG8iOiJzZGZkZnMiLCJiYWlycm8iOiJmc2RmcyIsImNlcCI6IjU0MzUzLTQiLCJjaWRhZGUiOiJBQkFEScOCTklBL0dPIiwiY3JlYXRlZF9hdCI6IjIwMTktMDgtMTdUMTg6NTM6MjguOTgzNzI2WiIsInVwZGF0ZWRfYXQiOm51bGx9LHsiaWQiOjIsImFsdW5vX2lkIjo2LCJsb2dyYWRvdXJvIjoidGVzdGUiLCJudW1lcm8iOjM0MjQsImNvbXBsZW1lbnRvIjoic2ZzZGYiLCJiYWlycm8iOiJmc2Rmc2QiLCJjZXAiOiI3OTgyNC0yMTAiLCJjaWRhZGUiOiJBQkFESUEgREUgR09Jw4FTL0dPIiwiY3JlYXRlZF9hdCI6IjIwMTktMDgtMTdUMTg6NTM6MjguOTgzNzI2WiIsInVwZGF0ZWRfYXQiOm51bGx9LHsiaWQiOjMsImFsdW5vX2lkIjo2LCJsb2dyYWRvdXJvIjoiUlVBIENPTlRJTkVOVEFMIiwibnVtZXJvIjo5ODUsImNvbXBsZW1lbnRvIjoidGVzdGUiLCJiYWlycm8iOiJKQVJESU0gSVRBSVBVIiwiY2VwIjoiNzk4MjQyMTAiLCJjaWRhZGUiOiJET1VSQURPUy9NUyIsImNyZWF0ZWRfYXQiOiIyMDE5LTA4LTE3VDE4OjUzOjI4Ljk4MzcyNloiLCJ1cGRhdGVkX2F0IjpudWxsfV0sImNyZWF0ZWRfYXQiOiIyMDE5LTA4LTE3VDE4OjUzOjI4Ljk4MzcyNloiLCJ1cGRhdGVkX2F0IjpudWxsfX0.sj208Rwdk35lJYsseCdl5anQk4xbRzRSfYjfvJtHTVU",
        'Content-Type': 'application/json; charset=utf-8'
      };
      final response =
          await http.get(BaseUrl + "/disciplinas", headers: headers);
      if (response.statusCode == 200) {
        await new Future.delayed(const Duration(seconds: 10));
        if (mounted) {
          setState(() {
            list = json.decode(utf8.decode(response.bodyBytes)) as List;
            _dataLoaded();
          });
        }
      } else {
        setState(() {
          _loadingFailed = true;
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
            children: <Widget>[
              Padding(
                padding: EdgeInsets.all(5),
              ),
              Icon(
                Icons.error,
              ),
              Text('Sistema indiponível, tente novamente'),
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
      itemCount: list.length,
      itemBuilder: (BuildContext context, int index) {
        final data = list[index];
        final notas = data['notas'];
        int notasLength = notas.length;
        return Padding(
          padding: EdgeInsets.all(16.0),
          child: Column(
            children: <Widget>[
              ListTile(
                title: Text(
                  data['disciplina'],
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
                          final nota = notas[indexn];
                          Color cor = Colors.green[300];
                          final intValue = nota['valor'];
                          if (intValue < 6) {
                            cor = Colors.redAccent;
                          }
                          String anotherValue = '$intValue';
                          DateTime dataAtualizada =
                              DateTime.parse(nota['updated_at']);
                          //  final f = new DateFormat('yyyy-MM-dd hh:mm');
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
                              new Expanded(child: new Text(nota['descricao'])),
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
