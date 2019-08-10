import 'dart:async';
import 'package:flutter/material.dart';
import 'dart:convert';
import 'package:http/http.dart' as http;

class Notas extends StatefulWidget {
  @override
  _NotasState createState() => _NotasState();
}

class _NotasState extends State<Notas> {
  var list = List();

  _loadList() async {
    final response =
        await http.get("https://jsonplaceholder.typicode.com/posts/");
    if (response.statusCode == 200) {
      await new Future.delayed(const Duration(seconds: 1));
      if (mounted) {
        setState(() {
          list = json.decode(response.body) as List;
        });
      }
    } else {
      throw Exception('Failed to load posts');
    }
  }

  @override
  void initState() {
    _loadList();
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return ListView.builder(
      itemCount: list.length,
      itemBuilder: (BuildContext context, int index) {
        final data = list[index];
        if (index % 2 == 0) {
          return ListTile(
            title: Text(
              "Matemática",
            ),
            trailing: Icon(Icons.keyboard_arrow_right),
          );
        }
        return ListTile(
          leading: CircleAvatar(
            backgroundColor: Colors.redAccent,
            child: Text('5'),
          ),
          // leading: new Icon(Icons.brightness_1, size: 28.0,
          //       color: Colors.redAccent),
          contentPadding: EdgeInsets.all(10.0),
          title: Text(data['title']),
          subtitle: Text(
            data['body'],
            maxLines: 2,
            overflow: TextOverflow.ellipsis,
          ),
        );
      },
    );
  }
}
