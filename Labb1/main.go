package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
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
}

func parseOptionsHeader(header byte) (int, int) {
	delta := int((header & 0b11110000) >> 4)
	len := int(header & 0b00001111)
	return delta, len

}

func parseOptions(arr []byte) []COAPOption {
	var options []COAPOption
	fmt.Println("Number of option bytes: ", len(arr))
	cursor := 0
	for cursor < len(arr) {
		delta, len := parseOptionsHeader(arr[cursor])
		cursor++
		fmt.Println("Delta: ", delta, " len: ", len)
		var val []byte
		if len != 0 {
			val = arr[cursor : cursor+len]
		}
		temp := COAPOption{
			Delta:   delta,
			Leangth: len,
			Value:   val,
		}
		options = append(options, temp)
		cursor += len

	}
	return options
}

func main() {
	conn, err := net.Dial("udp", "coap.me:5683")
	if err != nil {
		panic(err)
	}

	msg := []byte{
		0x40,
		0x01,
		0x04,
		0xd2,
		0xb4,
		0x74,
		0x65,
		0x73,
		0x74,
	}

	n, err := conn.Write(msg)
	if err != nil {
		panic(err)
	}

	fmt.Println(n, "bytes written")

	response := make([]byte, 1024)
	n, err = conn.Read(response)
	if err != nil {
		panic(err)
	}
	response = response[:n]

	message := COAPMessage{
		Version:   int((response[0] & 0b11000000) >> 6),
		T:         int((response[0] & 0b00110000) >> 4),
		TKL:       int(response[0] & 0b00001111),
		Code:      int(response[1]),
		MessageID: int(binary.BigEndian.Uint16([]byte{response[2], response[3]})),
	}
	var byteStringArr []string
	var index int
	for i, b := range response {
		if b == 0xFF {
			index = i
		}
		byteStringArr = append(byteStringArr, hex.EncodeToString([]byte{b}))
	}

	optionStart := 4

	if message.TKL != 0 {
		message.Token = string(response[4 : 4+message.TKL])
		optionStart = 5 + message.TKL
	}

	optionByte := response[optionStart:index]
	fmt.Println(byteStringArr[optionStart:index])
	message.Options = parseOptions(optionByte)

	message.Payload = response[index+1 : n]

	fmt.Printf("%+v\n", message)
	fmt.Println(byteStringArr)
	fmt.Println(string(response[index+1 : n]))
	fmt.Println(len(optionByte))
}
