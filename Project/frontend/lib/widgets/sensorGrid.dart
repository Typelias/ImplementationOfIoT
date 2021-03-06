import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:frontend/providers/sensors.dart';
import 'package:frontend/widgets/temperatureText.dart';
import 'package:provider/provider.dart';

class SensorGrid extends StatelessWidget {
  const SensorGrid({Key? key}) : super(key: key);

  String parseMem(String m) {
    if (m.isEmpty) return "Waiting for data";
    final split = m.split("/");
    return split[0] + "mb/" + split[1] + "mb";
  }

  @override
  Widget build(BuildContext context) {
    final sensorData = Provider.of<SensorProvider>(context);
    final mediaData = MediaQuery.of(context);

    return Column(
      children: [
        SizedBox(
          height: 10 + mediaData.padding.top,
        ),
        Flexible(
            child: GridView(
          padding: EdgeInsets.zero,
          shrinkWrap: true,
          gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
            crossAxisCount: 2,
            childAspectRatio: 3 / 2,
            crossAxisSpacing: 10,
            mainAxisSpacing: 10,
          ),
          children: [
            Card(
              color: Colors.white10,
              child: Padding(
                padding: const EdgeInsets.all(8.0),
                child: GridTile(
                  header: const Text(
                    "CPU usage",
                    style: TextStyle(color: Colors.white60),
                    textScaleFactor: 1.5,
                    textAlign: TextAlign.center,
                  ),
                  child: GridTileBar(
                    subtitle: Text(
                      sensorData.cpu.isEmpty ? "0%" : sensorData.cpu.last,
                      textAlign: TextAlign.center,
                    ),
                  ),
                ),
              ),
            ),
            Card(
              color: Colors.white10,
              child: GridTile(
                header: const Text(
                  "RAM usage",
                  style: TextStyle(color: Colors.white60),
                  textScaleFactor: 1.5,
                  textAlign: TextAlign.center,
                ),
                child: GridTileBar(
                  subtitle: Text(
                    parseMem(sensorData.mem),
                    textAlign: TextAlign.center,
                  ),
                ),
              ),
            ),
          ],
        )),
        const SizedBox(
          height: 10,
        ),
        Flexible(
          child: GridView.builder(
              padding: EdgeInsets.zero,
              gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                crossAxisCount: 2,
                childAspectRatio: 3 / 2,
                crossAxisSpacing: 10,
                mainAxisSpacing: 10,
              ),
              itemCount: sensorData.sensors.length,
              itemBuilder: (ctx, i) {
                return GestureDetector(
                  onLongPress: () {
                    openDialog(context, i);
                  },
                  child: Card(
                    color: Colors.white10,
                    child: Padding(
                      padding: const EdgeInsets.all(8.0),
                      child: GridTile(
                        header: Text(
                          sensorData.sensors[i].Location.toUpperCase(),
                          style: const TextStyle(color: Colors.white60),
                          textScaleFactor: 1.5,
                          textAlign: TextAlign.center,
                        ),
                        child: GridTileBar(
                          subtitle: sensorData.sensors[i].Status
                              ? TemperatureText(
                                  sensorData.sensors[i].Temperature)
                              : const Text(
                                  "Offline",
                                  textAlign: TextAlign.center,
                                ),
                        ),
                        footer: Row(
                          mainAxisAlignment: MainAxisAlignment.spaceBetween,
                          children: [
                            const Text(
                              "Power status: ",
                              textScaleFactor: 1,
                              style: TextStyle(color: Colors.white),
                            ),
                            CupertinoSwitch(
                              value: sensorData.sensors[i].Status,
                              onChanged: (newVal) =>
                                  sensorData.changePowerStatus(newVal, i),
                            ),
                          ],
                        ),
                      ),
                    ),
                  ),
                );
              }),
        ),
      ],
    );
  }

  void openDialog(BuildContext context, int index) {
    final prov = Provider.of<SensorProvider>(context, listen: false);
    showCupertinoDialog(
        context: context,
        builder: (ctx) => CupertinoAlertDialog(
              title: const Text("Remove Sensor?"),
              actions: [
                CupertinoDialogAction(
                  child: const Text("No"),
                  onPressed: () => Navigator.of(ctx).pop(),
                ),
                CupertinoDialogAction(
                    child: const Text("Yes"),
                    isDestructiveAction: true,
                    onPressed: () {
                      prov.removeSensor(index);
                      Navigator.of(context).pop();
                    }),
              ],
            ));
  }
}
