# Implementation of IoT protocolls
## Labb0
Go implementation of HTTP server for IoT device.
- POST. Used to add new sensor id. 201 success. 409 if sensor already exists.
  Endpoint: 
  ```js
  "/sensor"
  ```
  Expects JSON input:
  ```json
  {
      "id": ID_OF_DEVICE
  }
  ```
- PUT. Used to add sensor data to device id. 200 for success. 404 if sensor is not found.
  Endpoint: 
  ```js
  "/sensor"
  ```
  Expects JSON input:
  ```json
  {
      "id": ID_OF_DEVICE,
      "value": NEW_SENSOR_VALUE
  }
  ```
- Get. Used to get values for a sensor.
  Endpoint(gets all sensors):
  ```js
  "/sensor"
  ```
  Endpoint(gets specific sensor):
  ```js
  "/sensor?id=ID_OF_DEVICE"
  ```
- DELETE. Used to remove sensor device.
Endpoint: 
  ```js
  "/sensor"
  ```
  Expects JSON input:
  ```json
  {
      "id": ID_OF_DEVICE
  }
  ```
## TCPSocketC++
TCP socket connection was planed to be HTTP server but moved to Golang. It is mostly C code
