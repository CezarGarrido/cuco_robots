// Copyright 2015 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:app/repository/aluno.dart';
import 'package:app/repository/contato.dart';
import 'package:app/entities/aluno.dart';
import 'package:app/entities/contato.dart';
import 'package:app/repository/endereco.dart';
import 'package:app/entities/endereco.dart';
import 'package:intl/intl.dart';

class _ContactCategory extends StatelessWidget {
  const _ContactCategory({Key key, this.icon, this.children}) : super(key: key);

  final IconData icon;
  final List<Widget> children;

  @override
  Widget build(BuildContext context) {
    final ThemeData themeData = Theme.of(context);
    return Container(
      padding: const EdgeInsets.symmetric(vertical: 16.0),
      decoration: BoxDecoration(
          border: Border(bottom: BorderSide(color: themeData.dividerColor))),
      child: DefaultTextStyle(
        style: Theme.of(context).textTheme.subhead,
        child: SafeArea(
          top: false,
          bottom: false,
          child: Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: <Widget>[
              Container(
                padding: const EdgeInsets.symmetric(vertical: 24.0),
                width: 72.0,
                child: Icon(icon, color: themeData.primaryColor),
              ),
              Expanded(child: Column(children: children)),
            ],
          ),
        ),
      ),
    );
  }
}

class _ContactItem extends StatelessWidget {
  const _ContactItem(
      {Key key, this.icon, this.lines, this.tooltip, this.onPressed})
      : assert(lines.length > 1),
        super(key: key);

  final IconData icon;
  final List<String> lines;
  final String tooltip;
  final VoidCallback onPressed;

  @override
  Widget build(BuildContext context) {
    final ThemeData themeData = Theme.of(context);
    final List<Widget> columnChildren = lines
        .sublist(0, lines.length - 1)
        .map<Widget>((String line) => Text(line))
        .toList();
    columnChildren.add(Text(lines.last, style: themeData.textTheme.caption));

    final List<Widget> rowChildren = <Widget>[
      Expanded(
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: columnChildren,
        ),
      ),
    ];
    if (icon != null) {
      rowChildren.add(SizedBox(
        width: 72.0,
        child: IconButton(
          icon: Icon(icon),
          color: themeData.primaryColor,
          onPressed: onPressed,
        ),
      ));
    }
    return MergeSemantics(
      child: Padding(
        padding: const EdgeInsets.symmetric(vertical: 16.0),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: rowChildren,
        ),
      ),
    );
  }
}

class ContactsDemo extends StatefulWidget {
  static const String routeName = '/contacts';

  @override
  ContactsDemoState createState() => ContactsDemoState();
}

enum AppBarBehavior { normal, pinned, floating, snapping }

class ContactsDemoState extends State<ContactsDemo> {
  static final GlobalKey<ScaffoldState> _scaffoldKey =
      GlobalKey<ScaffoldState>();
  final double _appBarHeight = 256.0;

  AppBarBehavior _appBarBehavior = AppBarBehavior.pinned;
  AlunoRepository _alunoRepository = AlunoRepository();

  ContatoRepository _contatoRepository = ContatoRepository();
  EnderecoRepository _enderecoRepository = EnderecoRepository();

  List<Contato> contatos = List();
  List<Endereco> enderecos = List();
  Aluno aluno = Aluno();

  Future<Null> loadAluno() async {
    Aluno alunoDB = await _alunoRepository.getAluno();
    List<Contato> contatosdb = await _contatoRepository.getContatos(alunoDB.id);
    List<Endereco> enderecosdb = await _enderecoRepository.getContatos(alunoDB.id);
   // await new Future.delayed(const Duration(seconds: 1));
    if (mounted) {
      setState(() {
        aluno = alunoDB;
        contatos = contatosdb;
        enderecos = enderecosdb;
      });
    }
  }

  @override
  void initState() {
    super.initState();
    loadAluno();
  }

  @override
  Widget build(BuildContext context) {

    final contato_items = <Widget>[];
    for (var i = 0; i < contatos.length; i++) {
        contato_items.add(new _ContactItem(
                        lines: <String>[
                          contatos[i].valor,
                          contatos[i].tipo,
                        ],
        ));
    }
        final endereco_items = <Widget>[];
    for (var i = 0; i < enderecos.length; i++) {
        endereco_items.add(new _ContactItem(
                        lines: <String>[
                          enderecos[i].cidade,
                          "Cidade",
                        ],
        ));
    }
    return Theme(
      data: ThemeData(
        brightness: Brightness.light,
        primarySwatch: Colors.blue,
        platform: Theme.of(context).platform,
      ),
      child: Scaffold(
        key: _scaffoldKey,
        body: CustomScrollView(
          slivers: <Widget>[
            SliverAppBar(
              expandedHeight: _appBarHeight,
              pinned: _appBarBehavior == AppBarBehavior.pinned,
              floating: _appBarBehavior == AppBarBehavior.floating ||
                  _appBarBehavior == AppBarBehavior.snapping,
              snap: _appBarBehavior == AppBarBehavior.snapping,
              actions: <Widget>[
                IconButton(
                  icon: const Icon(Icons.create),
                  tooltip: 'Edit',
                  onPressed: () {
                    _scaffoldKey.currentState.showSnackBar(const SnackBar(
                      content: Text("Funcionalidade não suportada."),
                    ));
                  },
                ),
                PopupMenuButton<AppBarBehavior>(
                  onSelected: (AppBarBehavior value) {
                    setState(() {
                      _appBarBehavior = value;
                    });
                  },
                  itemBuilder: (BuildContext context) =>
                      <PopupMenuItem<AppBarBehavior>>[
                    const PopupMenuItem<AppBarBehavior>(
                      value: AppBarBehavior.normal,
                      child: Text('Sair'),
                    ),
                    /*const PopupMenuItem<AppBarBehavior>(
                      value: AppBarBehavior.pinned,
                      child: Text('App bar stays put'),
                    ),
                    const PopupMenuItem<AppBarBehavior>(
                      value: AppBarBehavior.floating,
                      child: Text('App bar floats'),
                    ),
                    const PopupMenuItem<AppBarBehavior>(
                      value: AppBarBehavior.snapping,
                      child: Text('App bar snaps'),
                    ),*/
                  ],
                ),
              ],
              flexibleSpace: FlexibleSpaceBar(
                title: Text(
                  aluno.nome != null ? aluno.nome : "",
                  style: new TextStyle(
                    fontSize: 14.0,
                  ),
                ),
                background: Stack(
                  fit: StackFit.expand,
                  children: <Widget>[
                    // Image.asset(
                    //   'people/ali_landscape.png',
                    //   package: 'flutter_gallery_assets',
                    //   fit: BoxFit.cover,
                    //   height: _appBarHeight,
                    // ),
                    // This gradient ensures that the toolbar icons are distinct
                    // against the background image.
                    const DecoratedBox(
                      decoration: BoxDecoration(
                        gradient: LinearGradient(
                          begin: Alignment(0.0, -1.0),
                          end: Alignment(0.0, -0.4),
                          colors: <Color>[Color(0x60000000), Color(0x00000000)],
                        ),
                      ),
                    ),
                  ],
                ),
              ),
            ),
            SliverList(
              delegate: SliverChildListDelegate(
              <Widget>[
                   AnnotatedRegion<SystemUiOverlayStyle>(
                       value: SystemUiOverlayStyle.dark,
                       child: _ContactCategory(
                              icon: Icons.send,
                              children: contato_items,
                       )
                   ),
                   _ContactCategory(
                              icon: Icons.location_on,
                              children: endereco_items,
                       ),
                   _ContactCategory(
                              icon: Icons.today,
                              children: <Widget>[ _ContactItem(
                        lines: <String>[
                          aluno.dataNascimento != null? DateFormat.yMMMd("pt_BR")
                                  .format(DateTime.parse(aluno.dataNascimento)):"Indefinido",
                          "Aniversário",
                        ],
        ),]
                   )
              ]),
            ),
          ],
        ),
      ),
    );
  }
}
