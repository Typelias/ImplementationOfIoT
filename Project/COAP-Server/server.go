package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	coap "github.com/plgd-dev/go-coap/v2"
	"github.com/plgd-dev/go-coap/v2/message"
	"github.com/plgd-dev/go-coap/v2/message/codes"
	"github.com/plgd-dev/go-coap/v2/mux"
)

/*
coap "github.com/plgd-dev/go-coap/v2"
	"github.com/plgd-dev/go-coap/v2/message"
	"github.com/plgd-dev/go-coap/v2/message/codes"
	"github.com/plgd-dev/go-coap/v2/mux"
*/

type TemperatureSensor struct {
	Status      bool
	Temperature int
	Location    string
}

var thermostats map[string]TemperatureSensor
var max = 30
var min = -20

func loggingMiddleware(next mux.Handler) mux.Handler {
	return mux.HandlerFunc(func(w mux.ResponseWriter, r *mux.Message) {
		log.Printf("ClientAddress %v, %v\n", w.Client().RemoteAddr(), r.String())
		next.ServeCOAP(w, r)
	})
}

func parseInitState(state string) (bool, error) {
	if state == "ON" {
		return true, nil
	} else if state == "OFF" {
		return false, nil
	}
	return false, errors.New("invalid input")

}

func parseSplit(body string) (string, bool, error) {
	split := strings.Split(body, ":")
	if len(split) != 2 {
		return "", false, errors.New("invalid request")
	}
	state, err := parseInitState(split[1])
	if err != nil {
		return "", false, err
	}
	return split[0], state, nil
}

func handleTemp(w mux.ResponseWriter, r *mux.Message) {
	method := r.Code.String()
	body := make([]byte, 128)
	if r.Body == nil {
		w.SetResponse(codes.BadRequest, message.TextPlain, bytes.NewReader([]byte(string("No body included"))))
		return
	}
	n, _ := r.Body.Read(body)
	body = body[:n]
	fmt.Println(string(body))
	switch method {
	case "GET":
		sensor, ok := thermostats[string(body)]
		if !ok {
			w.SetResponse(codes.NotFound, message.TextPlain, bytes.NewReader([]byte(string(body)+" was not found")))
			test, _ := json.Marshal(thermostats)
			fmt.Println(string(test))
			break
		}
		if sensor.Status {
			w.SetResponse(codes.Content, message.TextPlain, bytes.NewReader([]byte(strconv.Itoa(sensor.Temperature))))
			sensor.Temperature = rand.Intn(max-min) + min
			thermostats[string(body)] = sensor
		} else {
			w.SetResponse(codes.Content, message.TextPlain, bytes.NewReader([]byte("Sensor is offline")))
		}

		test, _ := json.Marshal(thermostats)
		fmt.Println(string(test))
		break
	case "POST":
		location, initState, err := parseSplit(string(body))
		if err != nil {
			w.SetResponse(codes.BadRequest, message.TextPlain, bytes.NewReader([]byte(err.Error())))
			break
		}
		newThermo := TemperatureSensor{
			Status:      initState,
			Location:    location,
			Temperature: rand.Intn(max-min) + min,
		}
		thermostats[location] = newThermo
		test, _ := json.Marshal(thermostats)
		fmt.Println(string(test))
		w.SetResponse(codes.Created, message.TextPlain, bytes.NewReader([]byte("creted new device")))
		break
	case "PUT":
		location, newStatus, err := parseSplit(string(body))
		if err != nil {
			w.SetResponse(codes.BadRequest, message.TextPlain, bytes.NewReader([]byte(err.Error())))
			break
		}
		sensor, ok := thermostats[location]
		if !ok {
			w.SetResponse(codes.NotFound, message.TextPlain, bytes.NewReader([]byte(string(body)+" was not found")))
		}
		sensor.Status = newStatus
		thermostats[location] = sensor
		test, _ := json.Marshal(thermostats)
		print(string(test))
		w.SetResponse(codes.Changed, message.TextPlain, bytes.NewReader([]byte("state changed")))
		break
	case "DELETE":
		_, ok := thermostats[string(body)]
		if !ok {
			w.SetResponse(codes.NotFound, message.TextPlain, bytes.NewReader([]byte(string(body)+" was not found")))
			test, _ := json.Marshal(thermostats)
			fmt.Println(string(test))
			break
		}
		delete(thermostats, string(body))
		w.SetResponse(codes.Deleted, message.TextPlain, bytes.NewReader([]byte("deleted: "+string(body))))
	default:
		w.SetResponse(codes.BadRequest, message.TextPlain, bytes.NewReader([]byte("Method not supported")))
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	thermostats = make(map[string]TemperatureSensor)
	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	r.Handle("/temp", mux.HandlerFunc(handleTemp))
	r.Handle("/all", mux.HandlerFunc(func(w mux.ResponseWriter, r *mux.Message) {
		data, _ := json.Marshal(thermostats)
		w.SetResponse(codes.Content, message.TextPlain, bytes.NewReader(data))
	}))

	log.Fatal(coap.ListenAndServe("udp", ":5683", r))

}
