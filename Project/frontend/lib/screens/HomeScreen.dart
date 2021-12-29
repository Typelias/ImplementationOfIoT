import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:frontend/helpers/mqtt.dart';
import 'package:frontend/providers/sensors.dart';
import 'package:frontend/widgets/addSensorDialog.dart';
import 'package:frontend/widgets/sensorGrid.dart';
import 'package:provider/provider.dart';

class HomeScreen extends StatefulWidget {
  HomeScreen({Key? key}) : super(key: key);

  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  bool initialState = true;

  Widget init(BuildContext context) {
    final mqtt = Provider.of<SensorProvider>(context, listen: false);
    return FutureBuilder<void>(
      builder: (ctx, snapshot) {
        if (snapshot.connectionState == ConnectionState.waiting) {
          return const Center(
            child: CircularProgressIndicator(),
          );
        }
        //Provider.of<SensorProvider>(context, listen: false).populateList("");
        return const SensorGrid();
      },
      future: mqtt.init(),
    );
  }

  void openAdd() {
    showCupertinoDialog(
        context: context,
        builder: (ctx) {
          return const AddSensorDialog();
        });
  }

  @override
  Widget build(BuildContext context) {
    final mqtt = Provider.of<SensorProvider>(context, listen: false);
    return CupertinoPageScaffold(
        navigationBar: CupertinoNavigationBar(
          middle: const Text("Home Temperatures"),
          trailing: CupertinoButton(
            child: const Icon(CupertinoIcons.add),
            onPressed: openAdd,
          ),
        ),
        child: FutureBuilder<bool>(
          future: mqtt.connect(),
          builder: (ctx, snapshot) {
            return snapshot.connectionState == ConnectionState.waiting
                ? const Center(
                    child: CircularProgressIndicator(),
                  )
                : init(ctx);
          },
        ));
  }
}
