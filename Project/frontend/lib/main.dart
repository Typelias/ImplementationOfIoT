import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:frontend/screens/HomeScreen.dart';
import 'package:frontend/providers/sensors.dart';
import 'package:provider/provider.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({Key? key}) : super(key: key);

  // This widget is the root of your application.
  @override
  Widget build(BuildContext context) {
    return ChangeNotifierProvider<SensorProvider>(
      create: (_) => SensorProvider(),
      child: CupertinoApp(
        title: 'Flutter Demo',
        theme: const CupertinoThemeData(
          primaryColor: Colors.purple,
          brightness: Brightness.dark,
        ),
        home: HomeScreen(),
      ),
    );
  }
}
