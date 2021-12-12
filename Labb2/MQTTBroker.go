package main

import (
	"encoding/binary"
	"fmt"
	"net"
)

var subscriptions map[string][]*net.Conn
var lastValue map[string][]byte
var DEVIDER string

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

type BroadCastMessage struct {
	Topic   string
	Message string
	Packet  []byte
}

func broadcaster(ch chan BroadCastMessage) {
	for {
		message := <-ch
		connections := subscriptions[message.Topic]
		if len(connections) <= 0 {
			continue
		}
		fmt.Println(DEVIDER)
		fmt.Println("Broadcasting message to: ",
			message.Topic, "\ncontaining ", message.Message)
		fmt.Println(DEVIDER)
		fmt.Println()
		for _, c := range subscriptions[message.Topic] {
			(*c).Write(message.Packet)
		}
	}
}

func addSubscription(c *net.Conn, filter string) {
	if subscriptions[filter] == nil {
		subscriptions[filter] = make([]*net.Conn, 0)
	}
	subscriptions[filter] = append(subscriptions[filter], c)
}

func removeSubscription(c *net.Conn, filter string) {
	conns := subscriptions[filter]
	var newConns []*net.Conn

	for _, conn := range conns {
		if conn != c {
			newConns = append(newConns, conn)
		}
	}
	subscriptions[filter] = newConns
}

func removeAllSubs(c *net.Conn) {
	for topic, conns := range subscriptions {
		var newConns []*net.Conn
		for _, conn := range conns {
			if conn != c {
				newConns = append(newConns, conn)
			}
		}
		subscriptions[topic] = newConns
	}
}

func createConnectionAccept() []byte {
	var message []byte
	message = append(message, byte(0b00100000))
	message = append(message, byte(0b00000010))
	message = append(message, byte(0x0))
	message = append(message, byte(0x0))
	return message
}
func createPingAck() []byte {
	var message []byte
	message = append(message, byte(0b11010000))
	message = append(message, byte(0x0))
	return message
}

func createSubAck(identifier []byte, success bool) []byte {
	var message []byte
	message = append(message, byte(0b10010000))
	message = append(message, byte(0b00000011))
	message = append(message, identifier...)
	if success {
		message = append(message, byte(0x00))
	} else {
		message = append(message, byte(0x80))
	}

	return message
}

func createUnsubAck(identifier []byte) []byte {
	var message []byte
	message = append(message, byte(0b10110000))
	message = append(message, byte(0x2))
	message = append(message, identifier...)
	return message
}

func parseSubscribe(message []byte, c *net.Conn, id string) (bool, []byte, []string) {
	body := message[2:]
	var l []byte
	var subscribe string
	var qos int
	var subs []string

	for len(body) > 0 {
		l, body = body[:2], body[2:]
		subLen := binary.BigEndian.Uint16(l)
		if len(body) < int(subLen) {
			return false, message[:2], nil
		}
		subscribe, body = string(body[:subLen]), body[subLen:]
		if len(body) == 0 {
			return false, message[:2], nil
		}
		addSubscription(c, subscribe)
		subs = append(subs, subscribe)
		qos, body = int((body[0] & 0b00000011)), body[1:]
		fmt.Println(DEVIDER)
		fmt.Println(id+" subscribed to: "+subscribe+" \nWith QOS of:", qos)
		fmt.Println(DEVIDER)
		fmt.Println()
	}
	return true, message[:2], subs
}

func parseUnsubscribe(message []byte, c *net.Conn, id string) []byte {
	body := message[2:]
	var l []byte
	var unSub string
	for len(body) > 0 {
		l, body = body[:2], body[2:]
		subLen := binary.BigEndian.Uint16(l)
		unSub, body = string(body[:subLen]), body[subLen:]
		removeSubscription(c, unSub)
		fmt.Println(DEVIDER)
		fmt.Println(id + " unsubscribed from: " + unSub)
		fmt.Println(DEVIDER)
		fmt.Println()
	}
	return message[:2]
}

func parsePublish(message []byte, retain bool, id string) (string, string) {
	data := ""
	topic := ""
	var topicLenBytes []byte

	topicLenBytes, message = message[:2], message[2:]
	topicLength := binary.BigEndian.Uint16(topicLenBytes)
	topic, message = string(message[:topicLength]), message[topicLength:]

	data = string(message)
	fmt.Println(DEVIDER)
	fmt.Println(id + " publshed: ")
	fmt.Println("Topic: ", topic)
	fmt.Println("Data: ", data)
	fmt.Println("Will retain: ", retain)
	fmt.Println(DEVIDER)
	fmt.Println()

	return data, topic
}

func handleConnection(c *net.Conn, ch chan BroadCastMessage, id string) {
	fmt.Println(DEVIDER)
	fmt.Println("Accepting connection from: " + id)
	fmt.Println(DEVIDER)
	fmt.Println()
	conAck := createConnectionAccept()
	(*c).Write(conAck)
	for {
		constHEAD := make([]byte, 2)
		(*c).Read(constHEAD)
		messageType := int((constHEAD[0] & 0b11110000) >> 4)
		switch messageType {
		case 12: //Ping
			fmt.Println(DEVIDER)
			fmt.Println("Answering Ping Request from: " + id)
			fmt.Println(DEVIDER)
			fmt.Println()
			pingAck := createPingAck()
			(*c).Write(pingAck)
			break
		case 14: // Disconect
			fmt.Println(DEVIDER)
			fmt.Println("Disconnecting: " + id)
			fmt.Println(DEVIDER)
			fmt.Println()
			removeAllSubs(c)
			(*c).Close()
			return
		case 8: // Subscribe
			remainder := int(constHEAD[1])
			messageArr := make([]byte, remainder)
			(*c).Read(messageArr)
			success, code, subs := parseSubscribe(messageArr, c, id)
			subAck := createSubAck(code, success)
			for _, s := range subs {
				if len(lastValue[s]) == 0 {
					continue
				}
				(*c).Write(lastValue[s])
			}
			(*c).Write(subAck)
			break
		case 10: // Unsubscribe
			remainder := int(constHEAD[1])
			messageArr := make([]byte, remainder)
			(*c).Read(messageArr)
			(*c).Write(createUnsubAck(parseUnsubscribe(messageArr, c, id)))
			break
		case 3: // Publish
			remainder := int(constHEAD[1])
			messageArr := make([]byte, remainder)
			retain := int(constHEAD[0] & 0b00000001)
			qos := int((constHEAD[0] & 0b00000110) >> 1)
			if qos != 0 {
				fmt.Println("Unsuported QoS type")
				removeAllSubs(c)
				(*c).Close()
				return
			}
			(*c).Read(messageArr)
			data, topic := parsePublish(messageArr, retain != 0, id)
			ch <- BroadCastMessage{
				Message: data,
				Topic:   topic,
				Packet:  append(constHEAD, messageArr...),
			}
			if retain != 0 {
				lastValue[topic] = append(lastValue[topic],
					append(constHEAD, messageArr...)...)
			}
			break
		default:
			removeAllSubs(c)
			(*c).Close()
			panic("Unsuported packet type")
		}
	}
}

func acceptMessage(c *net.Conn, ch chan BroadCastMessage) {
	constHEAD := make([]byte, 2)
	(*c).Read(constHEAD)
	packetType := int((constHEAD[0] & 0b11110000) >> 4)
	remainder := int(constHEAD[1])
	if packetType != 1 {
		panic("Not a connection")
	}
	message := make([]byte, remainder)
	(*c).Read(message)
	mqttStringLen := binary.BigEndian.Uint16(message[:2])
	if mqttStringLen != 4 {
		(*c).Close()
		panic("Wrong protocoll")
	}
	mqttString := string(message[2:6])
	if mqttString != "MQTT" {
		(*c).Close()
		panic("Wring protocoll")
	}

	version := message[6]
	if version != 4 {
		panic("Protocoll not supported")
	}
	var payload []byte
	if len(message) > 10 {
		payload = message[10:]
	}
	id := "Null Identifier"
	if len(payload) > 0 {
		id = string(payload)
	}
	// fmt.Println(string(payload))
	handleConnection(c, ch, id)
	// Nothing we care about
	// fmt.Printf("%08b\n", message[7])
	// fmt.Println(message[8:9])
	// keepAlive := binary.BigEndian.Uint16(message[8:10])
	// fmt.Println(keepAlive)
	// var payload []byte
	// if len(message) > 10 {
	// 	payload = message[10:]
	// }
	// fmt.Println(string(payload))

}

func main() {
	deviderLen := 50
	for i := 0; i < deviderLen; i++ {
		DEVIDER += "-"
	}
	subscriptions = make(map[string][]*net.Conn)
	lastValue = make(map[string][]byte)
	ch := make(chan BroadCastMessage)
	go broadcaster(ch)
	port := ":1883"
	s, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}
	defer s.Close()

	for {
		c, err := s.Accept()
		if err != nil {
			panic(err)
		}
		go acceptMessage(&c, ch)
	}
}
