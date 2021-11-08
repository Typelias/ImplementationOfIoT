package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

type HTTPRequest struct {
	Method      string
	Endpoint    string
	Body        string
	ContentType string
}

type Sensor struct {
	ID     int   `json:"id"`
	Values []int `json:"values"`
}

type Sensors struct {
	Sensors []Sensor `json:"sensors"`
}

func getSensors() Sensors {
	data, _ := os.ReadFile("database.json")
	//fmt.Println(string(data))
	var s Sensors
	json.Unmarshal(data, &s)
	return s
}

func saveSensors(s Sensors) {
	se := Sensor{
		ID:     3,
		Values: []int{2, 3},
	}
	s.Sensors = append(s.Sensors, se)
	file, _ := json.MarshalIndent(s, "", " ")
	ioutil.WriteFile("database.json", file, 0644)

}

func sensorToJSON(s Sensor) string {
	b, _ := json.Marshal(s)
	return string(b)
}

func parseRequestString(req string) HTTPRequest {
	fmt.Println(req)
	split := strings.Split(req, "\r\n")
	//fmt.Println(split[len(split)-1])
	requestInfo := strings.Split(split[0], " ")
	//fmt.Println(requestInfo)
	cont := ""
	for _, s := range split {
		if strings.Contains(s, "Content-Type:") {
			sp := strings.Split(s, ": ")
			cont = sp[len(sp)-1]
		}
	}
	httpReq := HTTPRequest{
		Method:      requestInfo[0],
		Endpoint:    requestInfo[1],
		Body:        split[len(split)-1],
		ContentType: cont,
	}

	return httpReq
}

func checkID(id int) bool {
	s := getSensors()
	for _, se := range s.Sensors {
		if se.ID == id {
			return true
		}
	}
	return false
}

func getIDIndex(s Sensors, id int) int {
	for i, se := range (s).Sensors {
		if se.ID == id {
			return i
		}
	}
	return -1
}

func parseGET(req HTTPRequest) string {
	fmt.Println("Parsing GET")
	fmt.Println(req.Endpoint)
	if req.Endpoint == "/" {
		resp := "HTTP/1.1 200 OK\r\n"
		resp += "Content-Type: text/html\r\n"
		resp += "\r\n"
		resp += "<h1>Hello </h1>"
		return resp
	} else if req.Endpoint == "/sensor" {
		resp := "HTTP/1.1 200 OK\r\n"
		resp += "Content-Type: application/json\r\n"
		resp += "\r\n"
		s, _ := json.Marshal(getSensors())
		resp += string(s)

		return resp
	} else if strings.Contains(req.Endpoint, "?") {
		moses := strings.Split(req.Endpoint, "?")
		if moses[0] != "/sensor" {
			return "HTTP/1.1 400 Bad Request\r\n\r\n"
		}
		gud := strings.Split(moses[1], "=")
		if gud[0] != "id" {
			return "HTTP/1.1 400 Bad Request\r\n\r\n"
		}
		id, _ := strconv.Atoi(gud[1])
		se := getSensors()
		index := getIDIndex(se, id)
		if index == -1 {
			//HEJ ADIL
			return "HTTP/1.1 404 Deez Nutts\r\n\r\n"
		}
		s := se.Sensors[index]
		resp := "HTTP/1.1 200 OK\r\n"
		resp += "Content-Type: application/json\r\n"
		resp += "\r\n"
		resp += sensorToJSON(s)
		return resp
	}

	return "HTTP/1.1 404 Not Found\r\n\r\n"
}

func parsePOST(req HTTPRequest) string {
	fmt.Println("Parsing POST")
	if req.ContentType != "application/json" {
		return "HTTP/1.1 400 Bad Request\r\n\r\n"
	}
	if req.Endpoint == "/sensor" {
		data := gjson.Get(req.Body, "id")
		if !data.Exists() {
			return "HTTP/1.1 400 Bad Request\r\n\r\n"
		}
		id := int(data.Num)
		if checkID(id) {
			return "HTTP/1.1 409 Conflict\r\n\r\n"
		}
		se := getSensors()
		s := Sensor{
			ID:     id,
			Values: []int{},
		}
		se.Sensors = append(se.Sensors, s)
		saveSensors(se)
		return "HTTP/1.1 201 Created\r\n\r\n"
	}

	return "HTTP/1.1 404 Not Found\r\n\r\n"
}

func parsePUT(req HTTPRequest) string {
	if req.ContentType != "application/json" {
		return "HTTP/1.1 400 Bad Request\r\n\r\n"
	}
	if req.Endpoint == "/sensor" {
		data := gjson.Get(req.Body, "id")
		if !data.Exists() {
			return "HTTP/1.1 400 Bad Request\r\n\r\n"
		}
		id := int(data.Num)
		se := getSensors()
		index := getIDIndex(se, id)
		if index == -1 {
			return "HTTP/1.1 404 Not Found\r\n\r\n"
		}
		data = gjson.Get(req.Body, "value")
		if !data.Exists() {
			return "HTTP/1.1 400 Bad Request\r\n\r\n"
		}
		value := int(data.Num)
		se.Sensors[index].Values = append(se.Sensors[index].Values, value)
		saveSensors(se)
		return "HTTP/1.1 200 OK\r\n\r\n"
	}
	return "HTTP/1.1 404 Not Found\r\n\r\n"
}

func parseDELETE(req HTTPRequest) string {
	if req.ContentType != "application/json" {
		return "HTTP/1.1 400 Bad Request\r\n\r\n"
	}
	if req.Endpoint == "/sensor" {
		data := gjson.Get(req.Body, "id")
		if !data.Exists() {
			return "HTTP/1.1 400 Bad Request\r\n\r\n"
		}
		id := int(data.Num)
		fmt.Println(id)
		if !checkID(id) {
			return "HTTP/1.1 404 Not Found\r\n\r\n"
		}
		se := getSensors()
		var newArr []Sensor
		for _, s := range se.Sensors {
			if s.ID != id {
				newArr = append(newArr, s)
			}
		}
		se.Sensors = newArr
		saveSensors(se)
		return "HTTP/1.1 200 Stonks\r\n\r\n"
	}
	return "HTTP/1.1 404 Deez Nutts\r\n\r\n"

}

func parseRequest(req HTTPRequest) string {
	switch req.Method {
	case "GET":
		return parseGET(req)
	case "POST":
		return parsePOST(req)
	case "PUT":
		return parsePUT(req)
	case "DELETE":
		return parseDELETE(req)
	default:
		return ""

	}
	return ""
}

func handleConnection(c net.Conn) {
	fmt.Printf("Serving %s \n", c.RemoteAddr().String())
	requestString := ""
	bufsize := 1024
	buf := make([]byte, bufsize)
	//reader := bufio.NewReader(c)
	n, err := c.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(n)
	requestString += string(buf)
	if n == bufsize {
		for n == bufsize {
			buf = make([]byte, bufsize)
			n, err = c.Read(buf)
			if err != nil {
				fmt.Println(err)
				return
			}
			requestString += string(buf)
		}
	}
	req := parseRequestString(requestString)

	fmt.Println("Parsing request")
	returnString := parseRequest(req)

	c.Write([]byte(returnString))

	c.Close()
}

func main() {
	PORT := ":80"
	l, err := net.Listen("tcp4", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(c)
	}

}
