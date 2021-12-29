import 'package:flutter/cupertino.dart';
import 'package:frontend/providers/sensors.dart';
import 'package:provider/provider.dart';

class AddSensorDialog extends StatefulWidget {
  const AddSensorDialog({Key? key}) : super(key: key);

  @override
  _AddSensorDialogState createState() => _AddSensorDialogState();
}

class _AddSensorDialogState extends State<AddSensorDialog> {
  void submit() {
    if (_locationController.text.isEmpty) {
      return;
    }
    Provider.of<SensorProvider>(context, listen: false)
        .addSensor(_locationController.text, initialState);
    Navigator.of(context).pop();
  }

  var initialState = false;

  final _locationController = TextEditingController();
  @override
  Widget build(BuildContext context) {
    return StatefulBuilder(
      builder: (context, setState) {
        return CupertinoAlertDialog(
          title: const Text("Add sensor"),
          content: Column(
            children: [
              CupertinoTextField(
                placeholder: "location",
                keyboardType: TextInputType.name,
                controller: _locationController,
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
                    onChanged: (val) => setState(
                      () {
                        initialState = !initialState;
                      },
                    ),
                  ),
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
              onPressed: submit,
            ),
          ],
        );
      },
    );
  }
}
