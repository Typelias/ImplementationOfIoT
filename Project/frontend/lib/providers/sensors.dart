// ignore_for_file: non_constant_identifier_names, avoid_print

import 'dart:convert';
import 'dart:math';

import 'package:flutter/foundation.dart';
import 'package:mqtt_client/mqtt_client.dart';
import 'package:mqtt_client/mqtt_server_client.dart';
import 'package:mutex/mutex.dart';
import 'package:uuid/uuid.dart';

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
  /* MQTT client part */
  late MqttClient client;
  //MqttServerClient.withPort("typelias.se", "Flutter", 1883) as MQTTClient;
  final m = Mutex();
  late final Uuid uuid;
  late final String devId;

  SensorProvider() {
    uuid = Uuid();
    devId = uuid.v4();
    client = MqttServerClient.withPort("typelias.se", "Flutter-$devId", 1883);
  }

  Future<bool> connect() async {
    //client.logging(on: true);
    client.onConnected = onConnected;
    client.onDisconnected = onDisconnected;
    client.onSubscribed = onSubscribed;
    client.onSubscribeFail = onSubscribeFail;
    client.onUnsubscribed = (s) => onUnsubscribed(s ?? "");
    client.pongCallback = pong;

    var msg = MqttConnectMessage().withProtocolName("MQTT");
    client.connectionMessage = msg;

    try {
      await client.connect();
      subscribe("all/$devId");
      subscribe("home/change");
      subscribe("home/delete");
      subscribe("home/add");
      subscribe("pi/cpu");
      subscribe("pi/mem");
      return true;
    } catch (e) {
      print("Conection failed");
      return false;
    }
  }

  void subscribe(String topic) {
    client.subscribe(topic, MqttQos.atMostOnce);
  }

  void publish(String message, String topic) {
    final builder = MqttClientPayloadBuilder();
    builder.addString(message);
    client.publishMessage(topic, MqttQos.atMostOnce, builder.payload!);
  }

  Future<void> init() async {
    client.updates?.listen(subscribHandler);
    final builder = MqttClientPayloadBuilder();
    builder.addString("GET:$devId");
    client.publishMessage("all", MqttQos.atMostOnce, builder.payload!);
  }

  final List<Duration> b = [];
  final Mutex bench = Mutex();
  bool runningBenchmark = false;
  DateTime lastStartTime = DateTime.now();
  int totalTime = 1;

  Future<List<Duration>> benchmark(int runs) async {
    runningBenchmark = true;
    if (bench.isLocked) bench.release();
    b.clear();
    final benchStart = DateTime.now();
    for (var i = 0; i < runs; i++) {
      await bench.acquire();
      print("run: $i");
      lastStartTime = DateTime.now();
      final builder = MqttClientPayloadBuilder();
      builder.addString("GET:$devId");
      client.publishMessage("all", MqttQos.atMostOnce, builder.payload!);
    }
    await bench.acquire();
    runningBenchmark = false;
    totalTime = DateTime.now().difference(benchStart).inMilliseconds;
    bench.release();

    return b;
  }

  void subscribHandler(List<MqttReceivedMessage<MqttMessage>> list) {
    for (var element in list) {
      if (element.topic == "all/$devId") {
        handleAll(element);
      } else if (element.topic.contains("home/")) {
        final split = element.topic.split("/");
        final splitLen = split.length;
        if (splitLen != 2) {
          return;
        }
        if (split.last == "add" ||
            split.last == "change" ||
            split.last == "delete") {
          reload();
          return;
        }
        handleTemperatureChange(element);
      } else if (element.topic.contains("pi/")) {
        final splitLen = element.topic.split("/").length;
        if (splitLen != 2) {
          return;
        }
        handlePi(element);
      }
    }
  }

  void reload() {
    final builder = MqttClientPayloadBuilder();
    builder.addString("GET:$devId");
    client.publishMessage("all", MqttQos.atMostOnce, builder.payload!);
  }

  void handlePi(MqttReceivedMessage<MqttMessage> msg) {
    final info = msg.topic.split("/")[1];
    switch (info) {
      case "cpu":
        final MqttPublishMessage conv = msg.payload as MqttPublishMessage;
        final payload =
            MqttPublishPayload.bytesToStringAsString(conv.payload.message);
        _cpuList.add(payload);
        break;
      case "mem":
        final MqttPublishMessage conv = msg.payload as MqttPublishMessage;
        final payload =
            MqttPublishPayload.bytesToStringAsString(conv.payload.message);
        _memUsage = payload;
        break;
    }
  }

// #region Handle temp call
  Future<void> handleTemperatureChange(
      MqttReceivedMessage<MqttMessage> msg) async {
    await Future.delayed(const Duration(seconds: 1));
    await m.acquire();
    final sensorLocation = msg.topic.split("/")[1];
    final MqttPublishMessage conv = msg.payload as MqttPublishMessage;
    final payload =
        MqttPublishPayload.bytesToStringAsString(conv.payload.message);
    final index =
        _sensorList.indexWhere((element) => element.Location == sensorLocation);
    if (index == -1) {
      m.release();
      notifyListeners();
      return;
    }
    int p = 0;
    try {
      p = int.parse(payload);
    } catch (e) {
      m.release();
      return;
    }
    _sensorList[index] = Sensor(_sensorList[index].Status, p, sensorLocation);
    m.release();
    notifyListeners();
  }
// #endregion

// #region Handle All sub
  void handleAll(MqttReceivedMessage<MqttMessage> input) {
    if (runningBenchmark) {
      final stopTime = DateTime.now();
      final dur = stopTime.difference(lastStartTime).abs();
      b.add(dur);
      bench.release();
      return;
    }
    final MqttPublishMessage msg = input.payload as MqttPublishMessage;
    final x = MqttPublishPayload.bytesToStringAsString(msg.payload.message);
    if (x == "GET") {
      return;
    }
    populateList(x);
  }
// #endregion

// #region
  void onConnected() {
    print('Connected');
  }

  void onDisconnected() {
    print('Disconnected');
  }

  void onSubscribed(String topic) {
    print('Subscribed topic: $topic');
  }

  void onSubscribeFail(String topic) {
    print('Failed to subscribe $topic');
  }

  void onUnsubscribed(String topic) {
    print('Unsubscribed topic: $topic');
  }

  void pong() {
    print('Ping response client callback invoked');
  }

  void disconect() {
    client.disconnect();
  }

// #endregion
  bool isConnected() {
    return client.connectionStatus!.state == MqttConnectionState.connected;
  }

  /*Provider part*/

  final List<Sensor> _sensorList = [];
  final List<String> _cpuList = [];
  String _memUsage = "";

  List<Sensor> get sensors => _sensorList;
  List<String> get cpu => _cpuList;
  String get mem => _memUsage;

  void populateList(String jsonData) {
    _sensorList.clear();
    final data = json.decode(jsonData) as Map<String, dynamic>;
    data.forEach((_, value) {
      final s = Sensor.fromJson(value);
      subscribe("home/${s.Location}");
      _sensorList.add(s);
    });
    notifyListeners();
  }

  Future<void> addSensor(String name, bool state) async {
    await m.acquire();
    const max = 100;
    const min = -100;
    final rn = Random(DateTime.now().microsecondsSinceEpoch);
    final num = min + rn.nextInt(max - min);
    _sensorList.add(Sensor(state, num, name));
    final s = state ? "ON" : "OFF";
    publish("$name:$s", "home/add");
    subscribe("home/$name");
    m.release();
    notifyListeners();
  }

  Future<void> changePowerStatus(bool newState, int index) async {
    await m.acquire();
    _sensorList[index] = Sensor(
        newState, _sensorList[index].Temperature, _sensorList[index].Location);
    final s = newState ? "ON" : "OFF";
    publish(_sensorList[index].Location + ":$s", "home/change");
    m.release();
    notifyListeners();
  }

  Future<void> removeSensor(int index) async {
    await m.acquire();
    final s = _sensorList.removeAt(index);
    publish(s.Location, "home/delete");
    m.release();
    notifyListeners();
  }
}
