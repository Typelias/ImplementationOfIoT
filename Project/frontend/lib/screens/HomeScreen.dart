import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:frontend/helpers/mqtt.dart';
import 'package:frontend/providers/sensors.dart';
import 'package:frontend/widgets/sensorGrid.dart';
import 'package:provider/provider.dart';

class HomeScreen extends StatefulWidget {
  HomeScreen({Key? key}) : super(key: key);

  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  MQTTClient client = MQTTClient();
  bool initialState = true;

  @override
  void dispose() {
    client.client.disconnect();
    super.dispose();
  }

  Widget init() {
    return FutureBuilder<void>(
      builder: (ctx, snapshot) {
        if (snapshot.connectionState == ConnectionState.waiting) {
          return const Center(
            child: CircularProgressIndicator(),
          );
        }
        Provider.of<SensorProvider>(context, listen: false).populateList("");
        return SensorGrid();
      },
      future: client.init(),
    );
  }

  void openAdd() {
    showCupertinoDialog(
        context: context,
        builder: (ctx) {
          return StatefulBuilder(builder: (context, setState) {
            return CupertinoAlertDialog(
              title: const Text("Add sensor"),
              content: Column(
                children: [
                  const CupertinoTextField(
                    placeholder: "location",
                    keyboardType: TextInputType.name,
                  ),
                  const SizedBox(
                    height: 10,
                  ),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      const Text(
                        "Power status: ",
                        textScaleFactor: 1.5,
                      ),
                      CupertinoSwitch(
                          value: initialState,
                          onChanged: (val) => setState(() {
                                initialState = !initialState;
                              })),
                    ],
                  ),
                ],
              ),
              actions: [
                CupertinoDialogAction(
                  child: const Text("Close"),
                  onPressed: () => Navigator.of(context).pop(),
                  isDestructiveAction: true,
                ),
                CupertinoDialogAction(
                  child: const Text("Add"),
                  onPressed: () {},
                ),
              ],
            );
          });
        });
  }

  @override
  Widget build(BuildContext context) {
    return CupertinoPageScaffold(
        navigationBar: CupertinoNavigationBar(
          middle: const Text("My app"),
          trailing: CupertinoButton(
            child: const Icon(CupertinoIcons.add),
            onPressed: openAdd,
          ),
        ),
        child: FutureBuilder<bool>(
          future: client.connect(),
          builder: (ctx, snapshot) {
            return snapshot.connectionState == ConnectionState.waiting
                ? const Center(
                    child: CircularProgressIndicator(),
                  )
                : init();
          },
        ));
  }
}
