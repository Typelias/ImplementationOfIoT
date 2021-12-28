// ignore_for_file: non_constant_identifier_names

import 'dart:convert';

import 'package:flutter/material.dart';

class Sensor {
  final bool Status;
  final int Temperature;
  final String Location;

  Sensor(this.Status, this.Temperature, this.Location);

  Sensor.fromJson(Map<String, dynamic> json)
      : Status = json['Status'],
        Temperature = json['Temperature'],
        Location = json['Location'];
}

class SensorProvider with ChangeNotifier {
  static const jsonString =
      '{"bedroom":{"Status":true,"Temperature":21,"Location":"bedroom"},"kitchen":{"Status":true,"Temperature":14,"Location":"kitchen"}}';

  final List<Sensor> _sensorList = [];

  List<Sensor> get sensors => _sensorList;

  void populateList(String jsonData) {
    _sensorList.clear();
    final data = json.decode(jsonString) as Map<String, dynamic>;
    data.forEach((_, value) {
      final s = Sensor.fromJson(value);
      _sensorList.add(s);
    });
  }
}
