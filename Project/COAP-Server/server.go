package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	linuxproc "github.com/c9s/goprocinfo/linux"
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
			r := rand.Intn(3-1) + 1
			change := rand.Intn(4-1) + 1
			if r == 1 {
				sensor.Temperature = sensor.Temperature - change
			} else {
				sensor.Temperature = sensor.Temperature + change
			}
			//sensor.Temperature = rand.Intn(max-min) + min
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

var CPU string
var MEM string

func handleCpu(w mux.ResponseWriter, r *mux.Message) {
	code := r.Code.String()

	if code != "GET" {
		return
	}
	w.SetResponse(codes.Content, message.TextPlain, bytes.NewReader([]byte(CPU+":"+MEM)))
}

func updateCPUAndMEM() {
	var prevIdle, prevTot uint64
	first := false
	for {
		mem, _ := linuxproc.ReadMemInfo("/proc/meminfo")
		totMemMeg := float64(mem.MemTotal) * 0.0009765625
		totalMem := strconv.Itoa(int(math.Round(totMemMeg)))
		usedMemMeg := float64(mem.MemTotal-mem.MemFree) * 0.0009765625
		usedMeme := strconv.Itoa(int(math.Round(usedMemMeg)))
		MEM = usedMeme + "/" + totalMem
		stat, _ := linuxproc.ReadStat("/proc/stat")
		cpu := stat.CPUStatAll
		tot := cpu.User + cpu.Nice + cpu.System + cpu.Idle + cpu.IOWait + cpu.IRQ + cpu.SoftIRQ + cpu.Steal + cpu.Guest + cpu.GuestNice
		if first {
			deltaIdle := cpu.Idle - prevIdle
			deltaTot := tot - prevTot
			cpuUsage := (1.0 - float64(deltaIdle)/float64(deltaTot)) * 100.0
			proc := fmt.Sprintf("%6.3f", cpuUsage)
			CPU = proc + "%"
		} else {
			first = true
		}
		prevIdle = cpu.Idle
		prevTot = tot
		time.Sleep(time.Second)

	}
}

func main() {
	CPU = ""
	MEM = ""
	go updateCPUAndMEM()
	rand.Seed(time.Now().UnixNano())
	thermostats = make(map[string]TemperatureSensor)
	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	r.Handle("/temp", mux.HandlerFunc(handleTemp))
	r.Handle("/all", mux.HandlerFunc(func(w mux.ResponseWriter, r *mux.Message) {
		data, _ := json.Marshal(thermostats)
		w.SetResponse(codes.Content, message.TextPlain, bytes.NewReader(data))
	}))
	r.Handle("/pi", mux.HandlerFunc(handleCpu))
	fmt.Println("started server on port 5683")
	log.Fatal(coap.ListenAndServe("udp", ":5683", r))

}
