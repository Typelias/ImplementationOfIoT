// ignore_for_file: avoid_print

import 'package:mqtt_client/mqtt_client.dart';
import 'package:mqtt_client/mqtt_server_client.dart';

class MQTTClient {
  MqttServerClient client =
      MqttServerClient.withPort("localhost", "Flutter", 1883);

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
      subscribe("all");
      subscribe("home/change");
      subscribe("home/delete");
      return true;
    } catch (e) {
      print("Conection failed");
      return false;
    }
  }

  void subscribe(String topic) {
    client.subscribe(topic, MqttQos.atMostOnce);
  }

  void publish(String message) {}

  Future<void> init() async {
    //client.subscribe("all", MqttQos.atMostOnce);
    //await Future<void>.delayed(const Duration(seconds: 2));
    client.updates?.listen(subscribHandler);
    final builder = MqttClientPayloadBuilder();
    builder.addString("GET");
    client.publishMessage("all", MqttQos.atMostOnce, builder.payload!);
  }

  void subscribHandler(List<MqttReceivedMessage<MqttMessage>> list) {
    for (var element in list) {
      switch (element.topic) {
        case "all":
          handleAll(element);
          break;
      }
    }
  }

  void handleAll(MqttReceivedMessage<MqttMessage> input) {
    final MqttPublishMessage msg = input.payload as MqttPublishMessage;
    final x = MqttPublishPayload.bytesToStringAsString(msg.payload.message);
    if (x == "GET") {
      return;
    }
    print(input.topic);
    print(x);
  }

  // connection succeeded
  void onConnected() {
    print('Connected');
  }

// unconnected
  void onDisconnected() {
    print('Disconnected');
  }

// subscribe to topic succeeded
  void onSubscribed(String topic) {
    print('Subscribed topic: $topic');
  }

// subscribe to topic failed
  void onSubscribeFail(String topic) {
    print('Failed to subscribe $topic');
  }

// unsubscribe succeeded
  void onUnsubscribed(String topic) {
    print('Unsubscribed topic: $topic');
  }

// PING response received
  void pong() {
    print('Ping response client callback invoked');
  }

  void disconect() {
    client.disconnect();
  }

  bool isConnected() {
    return client.connectionStatus!.state == MqttConnectionState.connected;
  }
}
