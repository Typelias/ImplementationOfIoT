import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:frontend/providers/sensors.dart';
import 'package:frontend/widgets/temperatureText.dart';
import 'package:provider/provider.dart';

class SensorGrid extends StatelessWidget {
  const SensorGrid({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final sensorData = Provider.of<SensorProvider>(context);
    final mediaData = MediaQuery.of(context);

    return GridView.builder(
        padding: EdgeInsets.only(top: 10 + mediaData.padding.top),
        gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
          crossAxisCount: 2,
          childAspectRatio: 3 / 2,
          crossAxisSpacing: 10,
          mainAxisSpacing: 10,
        ),
        itemCount: sensorData.sensors.length,
        itemBuilder: (ctx, i) {
          return Card(
            color: Colors.white10,
            child: Padding(
              padding: const EdgeInsets.all(8.0),
              child: GridTile(
                header: Text(
                  sensorData.sensors[i].Location.toUpperCase(),
                  style: TextStyle(color: Colors.white60),
                  textScaleFactor: 1.5,
                  textAlign: TextAlign.center,
                ),
                child: GridTileBar(
                    subtitle:
                        TemperatureText(sensorData.sensors[i].Temperature)),
                footer: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    const Text(
                      "Power status: ",
                      textScaleFactor: 1,
                      style: TextStyle(color: Colors.white),
                    ),
                    CupertinoSwitch(
                        value: sensorData.sensors[i].Status, onChanged: null),
                  ],
                ),
              ),
            ),
          );
        });
  }
}
