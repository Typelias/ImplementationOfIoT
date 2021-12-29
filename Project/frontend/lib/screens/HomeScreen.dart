import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:frontend/providers/sensors.dart';
import 'package:frontend/widgets/addSensorDialog.dart';
import 'package:frontend/widgets/sensorGrid.dart';
import 'package:provider/provider.dart';

class HomeScreen extends StatefulWidget {
  const HomeScreen({Key? key}) : super(key: key);

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

  Widget getResult(List<Duration> l) {
    int totalTime = 0;
    for (var element in l) {
      totalTime += element.inMilliseconds;
    }
    final avg = (totalTime / l.length) / 1000;
    l.sort((el1, el2) => el1.inMilliseconds.compareTo(el2.inMilliseconds));
    final min = l.first.inMilliseconds / 1000;
    final max = l.last.inMilliseconds / 1000;
    // print(l.first.inMilliseconds / 1000);
    // print(l.last.inMilliseconds / 1000);
    // print(avg);
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Column(
          children: [const Text("Average"), Text(avg.toStringAsFixed(3) + "s")],
        ),
        Column(
          children: [const Text("Max"), Text(max.toStringAsFixed(3) + "s")],
        ),
        Column(
          children: [const Text("Min"), Text(min.toStringAsFixed(3) + "s")],
        )
      ],
    );
  }

  void runBench() {
    showCupertinoDialog(
        context: context,
        builder: (ctx) => CupertinoAlertDialog(
              title: const Text("Benchmark"),
              content: FutureBuilder(
                future: Provider.of<SensorProvider>(ctx, listen: false)
                    .benchmark(100),
                builder: (context, snapshot) {
                  return snapshot.connectionState == ConnectionState.waiting
                      ? const Text("Running benchmark")
                      : getResult(snapshot.data as List<Duration>);
                },
              ),
              actions: [
                CupertinoDialogAction(
                  child: const Text("Close"),
                  onPressed: () {
                    if (Provider.of<SensorProvider>(ctx, listen: false)
                        .runningBenchmark) return;
                    Navigator.of(ctx).pop();
                  },
                )
              ],
            ));
  }

  @override
  Widget build(BuildContext context) {
    final mqtt = Provider.of<SensorProvider>(context, listen: false);
    return CupertinoPageScaffold(
        navigationBar: CupertinoNavigationBar(
          leading: CupertinoButton(
            child: const Icon(Icons.run_circle),
            onPressed: runBench,
          ),
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
