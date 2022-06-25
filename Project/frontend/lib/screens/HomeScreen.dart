import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:frontend/providers/sensors.dart';
import 'package:frontend/widgets/addSensorDialog.dart';
import 'package:frontend/widgets/sensorGrid.dart';
import 'package:provider/provider.dart';
import 'package:stats/stats.dart';

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
    final stat = Stats.fromData(l.map((e) => e.inMilliseconds).toList())
        .withPrecision(3);
    final reqPerSec = l.length /
        (Provider.of<SensorProvider>(context, listen: false).totalTime / 1000);
    const style = TextStyle(fontSize: 14);
    return Column(
      children: [
        Text("Result for ${l.length} requests:"),
        SizedBox(
          width: double.infinity,
          child: Table(
            border: TableBorder.all(width: 3),
            columnWidths: const <int, TableColumnWidth>{
              0: FixedColumnWidth(100),
              1: FixedColumnWidth(20),
            },
            defaultVerticalAlignment: TableCellVerticalAlignment.middle,
            children: [
              TableRow(
                children: [
                  const Text("Request/Second:", style: style),
                  Text(reqPerSec.toStringAsFixed(3) + "st", style: style)
                ],
              ),
              TableRow(children: [
                const TableCell(
                    child: Text("Standard Deviation:", style: style)),
                TableCell(
                    child: Text(stat.standardDeviation.toString() + "ms",
                        style: style))
              ]),
              TableRow(children: [
                const TableCell(child: Text("Mean:", style: style)),
                TableCell(
                    child: Text(stat.average.toString() + "ms", style: style))
              ]),
              TableRow(children: [
                const TableCell(child: Text("Min:", style: style)),
                Text(stat.min.toString() + "ms", style: style)
              ]),
              TableRow(children: [
                const TableCell(child: Text("Max:", style: style)),
                TableCell(child: Text(stat.max.toString() + "ms", style: style))
              ]),
            ],
          ),
        ),
      ],
    );
  }

  void runBench() {
    showCupertinoDialog(
        context: context,
        builder: (ctx) => CupertinoAlertDialog(
              title: const Text("Benchmark"),
              content: FutureBuilder(
                //BENCHMARK
                future: Provider.of<SensorProvider>(ctx, listen: false)
                    .benchmark(1000),
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
