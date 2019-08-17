import 'dart:async';
import 'package:flutter/material.dart';
import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:app/constants.dart';
import 'dart:math' show pi;
import 'package:intl/intl.dart';

class Notas extends StatefulWidget {
  @override
  _NotasState createState() => _NotasState();
}

class _NotasState extends State<Notas> with SingleTickerProviderStateMixin {
  var list = List();
  bool _loadingInProgress;
  Animation<double> _angleAnimation;

  Animation<double> _scaleAnimation;

  AnimationController _controller;

  Future<Null> _loadList() async {
    Map<String, String> headers = {
      "Authorization": "Bearer " +
          "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhbHVubyI6eyJpZCI6NiwiZ3VpZCI6IjYzNDc0NzQxLTU5ZjMtNGNkOC1hZjQxLTViM2NkM2MxNWNiZSIsIm5vbWUiOiJDRVpBUiBHQVJSSURPIEJSSVRFWiIsInJnbSI6IjQwMDg5Iiwic2VuaGEiOiJDMTAyMDMwZyIsImN1cnNvIjoiIiwiZGF0YV9uYXNjaW1lbnRvIjoiMTk5Ny0xMi0yOFQwMDowMDowMFoiLCJzZXhvIjoiTWFzY3VsaW5vIiwibm9tZV9wYWkiOiJWSVRPUiBCUklURVoiLCJub21lX21hZSI6Ik1BUklBTkEgR0FSUklETyIsImVzdGFkb19jaXZpbCI6IlNvbHRlaXJvKGEpIiwibmFjaW9uYWxpZGFkZSI6IkJSQVNJTEVJUk8iLCJuYXR1cmFsaWRhZGUiOiJQQVJBTkhPUy9NUyIsImZlbm90aXBvIjoiIiwiY3BmIjoiMDUwLjQzMy42OTEtNjciLCJyZyI6IjIuMjI1LjIyOCIsInJnX29yZ2FvX2VtaXNzb3IiOiJNRCIsInJnX2VzdGFkb19lbWlzc29yIjoiTVMiLCJyZ19kYXRhX2VtaXNzYW8gIjoiMDAwMS0wMS0wMVQwMDowMDowMFoiLCJjb250YXRvcyI6W3siaWQiOjEsImFsdW5vX2lkIjo2LCJ0aXBvIjoiVGVsZWZvbmUiLCJ2YWxvciI6Iig2NykgOTk2ODItMjQwMiIsImNyZWF0ZWRfYXQiOiIyMDE5LTA4LTE3VDE4OjUzOjI4Ljk4MzcyNloiLCJ1cGRhdGVkX2F0IjpudWxsfSx7ImlkIjoyLCJhbHVub19pZCI6NiwidGlwbyI6IkVtYWlsIiwidmFsb3IiOiJjZXphci5jZ2IxOEBnbWFpbC5jb20iLCJjcmVhdGVkX2F0IjoiMjAxOS0wOC0xN1QxODo1MzoyOC45ODM3MjZaIiwidXBkYXRlZF9hdCI6bnVsbH1dLCJlbmRlcmVjb3MiOlt7ImlkIjoxLCJhbHVub19pZCI6NiwibG9ncmFkb3VybyI6ImNmd2VmIiwibnVtZXJvIjo1MzU0MywiY29tcGxlbWVudG8iOiJzZGZkZnMiLCJiYWlycm8iOiJmc2RmcyIsImNlcCI6IjU0MzUzLTQiLCJjaWRhZGUiOiJBQkFEScOCTklBL0dPIiwiY3JlYXRlZF9hdCI6IjIwMTktMDgtMTdUMTg6NTM6MjguOTgzNzI2WiIsInVwZGF0ZWRfYXQiOm51bGx9LHsiaWQiOjIsImFsdW5vX2lkIjo2LCJsb2dyYWRvdXJvIjoidGVzdGUiLCJudW1lcm8iOjM0MjQsImNvbXBsZW1lbnRvIjoic2ZzZGYiLCJiYWlycm8iOiJmc2Rmc2QiLCJjZXAiOiI3OTgyNC0yMTAiLCJjaWRhZGUiOiJBQkFESUEgREUgR09Jw4FTL0dPIiwiY3JlYXRlZF9hdCI6IjIwMTktMDgtMTdUMTg6NTM6MjguOTgzNzI2WiIsInVwZGF0ZWRfYXQiOm51bGx9LHsiaWQiOjMsImFsdW5vX2lkIjo2LCJsb2dyYWRvdXJvIjoiUlVBIENPTlRJTkVOVEFMIiwibnVtZXJvIjo5ODUsImNvbXBsZW1lbnRvIjoidGVzdGUiLCJiYWlycm8iOiJKQVJESU0gSVRBSVBVIiwiY2VwIjoiNzk4MjQyMTAiLCJjaWRhZGUiOiJET1VSQURPUy9NUyIsImNyZWF0ZWRfYXQiOiIyMDE5LTA4LTE3VDE4OjUzOjI4Ljk4MzcyNloiLCJ1cGRhdGVkX2F0IjpudWxsfV0sImNyZWF0ZWRfYXQiOiIyMDE5LTA4LTE3VDE4OjUzOjI4Ljk4MzcyNloiLCJ1cGRhdGVkX2F0IjpudWxsfX0.sj208Rwdk35lJYsseCdl5anQk4xbRzRSfYjfvJtHTVU",
      'Content-Type': 'application/json; charset=utf-8'
    };
    final response = await http.get(BaseUrl + "/disciplinas", headers: headers);
    if (response.statusCode == 200) {
      await new Future.delayed(const Duration(seconds: 1));
      if (mounted) {
        setState(() {
          list = json.decode(utf8.decode(response.bodyBytes)) as List;
          _dataLoaded();
        });
      }
    } else {
      throw Exception('Failed to load posts');
    }
  }

  Widget get _loadingView {
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
    _loadingInProgress = true;
    _controller = new AnimationController(
        duration: const Duration(milliseconds: 2000), vsync: this);
    _angleAnimation = new Tween(begin: 0.0, end: 360.0).animate(_controller)
      ..addListener(() {
        setState(() {
          // the state that has changed here is the animation object’s value
        });
      });
    _scaleAnimation = new Tween(begin: 1.0, end: 6.0).animate(_controller)
      ..addListener(() {
        setState(() {
          // the state that has changed here is the animation object’s value
        });
      });

    _angleAnimation.addStatusListener((status) {
      if (status == AnimationStatus.completed) {
        if (_loadingInProgress) {
          _controller.reverse();
        }
      } else if (status == AnimationStatus.dismissed) {
        if (_loadingInProgress) {
          _controller.forward();
        }
      }
    });

    _controller.forward();
    _loadList();
    super.initState();
  }

  Widget _buildBody(BuildContext context) {
    if (_loadingInProgress) {
      return new Center(
        child: _buildAnimation(),
      );
    } else {
      return new RefreshIndicator(
        child: _buildListview(context),
        onRefresh: _loadList,
      );
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
                  ? Center(child: Text('Aguardando atualização...'))
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

  Widget _buildAnimation() {
    double circleWidth = 10.0 * _scaleAnimation.value;
    Widget circles = new Container(
      width: circleWidth * 2.0,
      height: circleWidth * 2.0,
      child: new Column(
        children: <Widget>[
          new Row(
            children: <Widget>[
              _buildCircle(circleWidth, Colors.blue),
              _buildCircle(circleWidth, Colors.red),
            ],
          ),
          new Row(
            children: <Widget>[
              _buildCircle(circleWidth, Colors.yellow),
              _buildCircle(circleWidth, Colors.green),
            ],
          ),
        ],
      ),
    );

    double angleInDegrees = _angleAnimation.value;
    return new Transform.rotate(
      angle: angleInDegrees / 360 * 2 * pi,
      child: new Container(
        child: circles,
      ),
    );
  }

  Widget _buildCircle(double circleWidth, Color color) {
    return new Container(
      width: circleWidth,
      height: circleWidth,
      decoration: new BoxDecoration(
        color: color,
        shape: BoxShape.circle,
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return _buildBody(context);
  }
}
