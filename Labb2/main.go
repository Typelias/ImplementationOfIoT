package main

import (
	"fmt"
	"net"
)

var subscriptions map[string][]net.Conn

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
}

func broadcaster(ch chan BroadCastMessage) {
	for {
		message := <-ch
		connections := subscriptions[message.Topic]
		if len(connections) <= 0 {
			continue
		}
		for _, c := range subscriptions[message.Topic] {
			//Implement send
			fmt.Println(c)
		}
	}
}

func main() {
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
		bufsize := 1024
		buf := make([]byte, bufsize)
		n, _ := c.Read(buf)
		buf = buf[:n]
		fmt.Println("Bytes read: ", n)
		fmt.Println(buf)
	}
}
