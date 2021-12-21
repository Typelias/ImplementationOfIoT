package main

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"
)

//40 01 04 d2 b4 74 65 73 74

type COAPOption struct {
	Delta   int
	Leangth int
	Value   []byte
}

type COAPMessage struct {
	Version   int
	T         int
	TKL       int
	Code      int
	MessageID int
	Token     string
	Options   []COAPOption
	Payload   []byte
	Format    string
}

func parseMethodCode(c int) string {
	switch c {
	case 1:
		return "GET"
	case 2:
		return "POST"
	case 3:
		return "PUT"
	case 4:
		return "DELETE"
	case 65:
		return "Created"
	case 66:
		return "Deleted"
	case 67:
		return "Valid"
	case 68:
		return "Changed"
	case 69:
		return "Continue"
	case 132:
		return "Not Found"
	}

	return "not in list: " + strconv.Itoa(c)
}

func parseOptionCode(c int) (string, string) {
	switch c {
	case 4:
		return "Etag", "opaque"
	case 8:
		return "Location Path", "string"
	case 11:
		return "Uri-path", "string"
	case 12:
		return "Content Format", "string"
	}

	return "Number " + strconv.Itoa(c) + " is not in list", "null"
}

func createOption(name string, data []byte, lastDelta int) ([]byte, int) {
	switch name {
	case "uri":
		delta := (11 - lastDelta) << 4
		lastDelta = 11 - lastDelta
		l := len(data)
		header := delta + l
		return append([]byte{byte(header)}, data...), lastDelta
	case "contentType":
		delta := (12 - lastDelta) << 4
		lastDelta = 12 - lastDelta
		l := len(data)
		header := delta + l
		return append([]byte{byte(header)}, data...), lastDelta
	}

	return nil, 0
}

func parseOptionsHeader(header byte) (int, int) {
	delta := int((header & 0b11110000) >> 4)
	len := int(header & 0b00001111)
	return delta, len
}

func parseOptions(arr []byte) []COAPOption {
	var options []COAPOption
	cursor := 0
	lastDelta := 0
	for cursor < len(arr) {
		delta, len := parseOptionsHeader(arr[cursor])
		lastDelta += delta
		cursor++
		var val []byte
		if len != 0 {
			val = arr[cursor : cursor+len]
		}
		temp := COAPOption{
			Delta:   lastDelta,
			Leangth: len,
			Value:   val,
		}
		options = append(options, temp)
		cursor += len

	}
	return options
}

func parseMessage(arr []byte, n int) COAPMessage {
	message := COAPMessage{
		Version:   int((arr[0] & 0b11000000) >> 6),
		T:         int((arr[0] & 0b00110000) >> 4),
		TKL:       int(arr[0] & 0b00001111),
		Code:      int(arr[1]),
		MessageID: int(binary.BigEndian.Uint16([]byte{arr[2], arr[3]})),
	}
	var index int
	for i, b := range arr {
		if b == 0xFF {
			index = i
		}
	}
	if index == 0 {
		index = n
	}

	optionStart := 4

	if message.TKL != 0 {
		message.Token = string(arr[4 : 4+message.TKL])
		optionStart = 5 + message.TKL
	}

	optionByte := arr[optionStart:index]
	message.Options = parseOptions(optionByte)

	if index != n {
		message.Payload = arr[index+1 : n]
	}

	return message
}

func createGet(path string) []byte {
	var ret []byte
	firstByte := byte(0b01010000)
	ret = append(ret, firstByte)
	code := byte(0x01)
	ret = append(ret, code)
	r := rand.Uint32()
	id := make([]byte, 4)
	binary.BigEndian.PutUint32(id, r)
	id = id[:2]
	ret = append(ret, id...)
	lastDelta := 0
	option, lastDelta := createOption("uri", []byte(path), lastDelta)
	ret = append(ret, option...)
	option, _ = createOption("contentType", make([]byte, 0), lastDelta)
	ret = append(ret, option...)

	return ret
}

func createPost(path string, payload string) []byte {
	ret := make([]byte, 0)
	firstByte := byte(0b01010000)
	ret = append(ret, firstByte)
	code := byte(0x02)
	ret = append(ret, code)
	r := rand.Uint32()
	id := make([]byte, 4)
	binary.BigEndian.PutUint32(id, r)
	id = id[:2]
	ret = append(ret, id...)
	lastDelta := 0
	option, lastDelta := createOption("uri", []byte(path), lastDelta)
	ret = append(ret, option...)
	option, _ = createOption("contentType", make([]byte, 0), lastDelta)
	ret = append(ret, option...)
	ret = append(ret, byte(0xFF))
	ret = append(ret, []byte(payload)...)

	return ret
}

func createPut(path string, payload string) []byte {
	ret := make([]byte, 0)
	firstByte := byte(0b01010000)
	ret = append(ret, firstByte)
	code := byte(0x03)
	ret = append(ret, code)
	r := rand.Uint32()
	id := make([]byte, 4)
	binary.BigEndian.PutUint32(id, r)
	id = id[:2]
	ret = append(ret, id...)
	lastDelta := 0
	option, lastDelta := createOption("uri", []byte(path), lastDelta)
	ret = append(ret, option...)
	option, _ = createOption("contentType", make([]byte, 0), lastDelta)
	ret = append(ret, option...)
	ret = append(ret, byte(0xFF))
	ret = append(ret, []byte(payload)...)

	return ret
}

func createDelete(path string) []byte {
	ret := make([]byte, 0)
	firstByte := byte(0b01010000)
	ret = append(ret, firstByte)
	code := byte(0x04)
	ret = append(ret, code)
	r := rand.Uint32()
	id := make([]byte, 4)
	binary.BigEndian.PutUint32(id, r)
	id = id[:2]
	ret = append(ret, id...)
	lastDelta := 0
	option, lastDelta := createOption("uri", []byte(path), lastDelta)
	ret = append(ret, option...)
	option, _ = createOption("contentType", make([]byte, 0), lastDelta)
	ret = append(ret, option...)

	return ret
}

func printOptions(options []COAPOption) string {
	format := ""
	for _, option := range options {
		optionName, optionFormat := parseOptionCode(option.Delta)
		if optionName == "Content Format" {
			if len(option.Value) == 0 {
				format = "text/plain"
				optionFormat = "CF"
			}
		}
		fmt.Print("\t" + optionName + ": ")
		switch optionFormat {
		case "string":
			fmt.Println(string(option.Value))
		case "CF":
			fmt.Println(format)
		case "opaque":
			fmt.Println(option.Value)
		}

	}
	return format
}

func printCOAP(c COAPMessage) {
	fmt.Println("Version:", c.Version, "Message Type:", c.T, "Token leangth:", c.TKL)
	fmt.Println("Metod/Response:", "\""+parseMethodCode(c.Code)+"\"", "Message id:", c.MessageID)
	fmt.Println("Token:", c.Token)
	fmt.Println("Options:")
	c.Format = printOptions(c.Options)
	if c.Format != "text/plain" {
		fmt.Println("FIX FORMAT")
	} else {
		fmt.Println("Payload:\n", "\t"+string(c.Payload))
	}
}

func sendCOAP(method string) {
	conn, err := net.Dial("udp", "coap.me:5683")
	if err != nil {
		panic(err)
	}
	var msg []byte
	var uri string
	fmt.Println("Type endpoint")
	fmt.Scanln(&uri)
	switch method {
	case "GET":
		msg = createGet(uri)
		break
	case "POST":
		var payload string
		fmt.Println("Type payload")
		fmt.Scanln(&payload)
		msg = createPost(uri, payload)
		break
	case "PUT":
		var payload string
		fmt.Println("Type payload")
		fmt.Scanln(&payload)
		msg = createPut(uri, payload)
		break
	case "DELETE":
		msg = createDelete(uri)
		break
	default:
		msg = make([]byte, 0)
		break
	}

	fmt.Println("Created following Message:")
	fmt.Println("------------------------------------------------")
	printCOAP(parseMessage(msg, len(msg)))
	fmt.Println("------------------------------------------------")
	_, err = conn.Write(msg)
	if err != nil {
		panic(err)
	}
	response := make([]byte, 1024)
	n, err := conn.Read(response)
	if err != nil {
		panic(err)
	}
	response = response[:n]
	fmt.Print("\n")
	fmt.Println("Got following response:")
	fmt.Println("------------------------------------------------")
	printCOAP(parseMessage(response, n))
	fmt.Println("------------------------------------------------")
	fmt.Println("Press any key to continue")
	fmt.Scanln()
	cmd := exec.Command("clear") //Linux example, its tested
	cmd.Stdout = os.Stdout
	cmd.Run()

}

func main() {
	rand.Seed(time.Now().UnixMicro())
	var option string
	run := true
	/* msg := createGet()c
	printCOAP(parseMessage(msg, len(msg))) */
	for run {
		fmt.Println("Select message to send :")
		fmt.Println("1: GET")
		fmt.Println("2: POST")
		fmt.Println("3: PUT")
		fmt.Println("4: DELETE")
		fmt.Println("Any other key to exit")
		fmt.Scanln(&option)

		switch option {
		case "1":
			sendCOAP("GET")
		case "2":
			sendCOAP("POST")
		case "3":
			sendCOAP("PUT")
		case "4":
			sendCOAP("DELETE")
		default:
			run = false
		}
		option = ""
	}
}
