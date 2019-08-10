import 'dart:ui';

import 'package:flutter/cupertino.dart';

import 'package:flutter/material.dart';

class Colors {

  const Colors();

  static const Color loginGradientStart = const Color(0xFF067BA5);
  static const Color loginGradientEnd = const Color(0xFF258535);//#4e2a8e

  static const primaryGradient = const LinearGradient(
    colors: const [loginGradientStart, loginGradientEnd],
    stops: const [0.0, 1.0],
    begin: Alignment.topCenter,
    end: Alignment.bottomCenter,
  );
}

class CustomIcons {
  static const IconData twitter = IconData(0xe900, fontFamily: "CustomIcons");
  static const IconData facebook = IconData(0xe901, fontFamily: "CustomIcons");
  static const IconData googlePlus =
      IconData(0xe902, fontFamily: "CustomIcons");
  static const IconData linkedin = IconData(0xe903, fontFamily: "CustomIcons");
}