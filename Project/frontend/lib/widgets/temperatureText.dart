import 'package:flutter/material.dart';

class TemperatureText extends StatelessWidget {
  const TemperatureText(this.temp);
  final int temp;

  Color getColor() {
    if (temp < -10) {
      return Colors.purple;
    } else if (temp <= 0) {
      return Colors.blue;
    } else if (temp <= 10) {
      return Colors.lightBlue;
    } else if (temp <= 20) {
      return Colors.yellow;
    }
    return Colors.orange;
  }

  @override
  Widget build(BuildContext context) {
    return Text(
      temp.toString(),
      textAlign: TextAlign.center,
      textScaleFactor: 1.2,
      style: TextStyle(color: getColor()),
    );
  }
}
