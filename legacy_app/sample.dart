import 'package:flutter/material.dart';
import 'dart:convert';

class LegacyFlutterWidget extends StatelessWidget {
  final String data;

  LegacyFlutterWidget({required this.data});

  void processData() {
    print("Processing $data");
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      child: Text(data),
    );
  }
}
