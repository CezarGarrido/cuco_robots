// Copyright 2018 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

import 'package:flutter/gestures.dart';
import 'package:flutter/foundation.dart' show defaultTargetPlatform;
import 'package:flutter/material.dart';

import 'package:url_launcher/url_launcher.dart';

class _LinkTextSpan extends TextSpan {

  // Beware!
  //
  // This class is only safe because the TapGestureRecognizer is not
  // given a deadline and therefore never allocates any resources.
  //
  // In any other situation -- setting a deadline, using any of the less trivial
  // recognizers, etc -- you would have to manage the gesture recognizer's
  // lifetime and call dispose() when the TextSpan was no longer being rendered.
  //
  // Since TextSpan itself is @immutable, this means that you would have to
  // manage the recognizer from outside the TextSpan, e.g. in the State of a
  // stateful widget that then hands the recognizer to the TextSpan.

  _LinkTextSpan({ TextStyle style, String url, String text }) : super(
    style: style,
    text: text ?? url,
    recognizer: TapGestureRecognizer()..onTap = () {
      launch(url, forceSafariVC: false);
    }
  );
}

void showGalleryAboutDialog(BuildContext context) {
  final ThemeData themeData = Theme.of(context);
  final TextStyle aboutTextStyle = themeData.textTheme.body2;
  final TextStyle linkStyle = themeData.textTheme.body2.copyWith(color: themeData.accentColor);

  showAboutDialog(
    context: context,
    applicationVersion: 'Janeiro 2019',
    applicationIcon: const FlutterLogo(),
    applicationLegalese: '© 2019 The Cuco Authors',
    children: <Widget>[
      Padding(
        padding: const EdgeInsets.only(top: 24.0),
        child: RichText(
          text: TextSpan(
            children: <TextSpan>[
              TextSpan(
                style: aboutTextStyle,
                text: 'Cuco é um projeto de código aberto feito para ajudar estudantes da Uems, '
                      'trazendo dinamismo e controle sobre o dia a dia na faculdade. '
                      'Esta versão esta em fase beta trazendo algumas funcionalidades em sua codebase. '
                      "O aplicativo foi feito com a linguagem de programação Golang e Flutter da Google. "
                      'Saiba mais sobre Flutter em ',
              ),
              _LinkTextSpan(
                style: linkStyle,
                url: 'https://flutter.dev',
              ),
                TextSpan(
                style: aboutTextStyle,
                text: '.\n\nSaiba mais sobre Golang em ',
              ),
              _LinkTextSpan(
                style: linkStyle,
                url: 'https://golang.org/',
                text: 'site golang',
              ),
              TextSpan(
                style: aboutTextStyle,
                text: '.\n\nPara ver o código-fonte deste aplicativo, visite ',
              ),
              _LinkTextSpan(
                style: linkStyle,
                url: 'https://github.com/CezarGarrido/cuco_robots',
                text: 'cuco github repo',
              ),
              TextSpan(
                style: aboutTextStyle,
                text: '.',
              ),
            ],
          ),
        ),
      ),
    ],
  );
}