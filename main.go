package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func err_hand(err error, msg string) {
	if err != nil {
		log.Fatalf("%v\n\ndeu merda:\n\n%v\n", msg, err)
	}
}

func req_hand(c net.Conn) {
	defer c.Close()

	fmt.Printf("Conex: %v", c)

	red := bufio.NewReader(c)
	for {
		msg, err := red.ReadString('\n')
		if err == io.EOF {
			return
		}
		err_hand(err, "reader")

		fmt.Printf("receba: %v\n", msg)
		_, err = c.Write([]byte("echo:" + msg))
		err_hand(err, "echo")

		if d := strings.TrimSpace(msg); d == "PAROOU" {
			return
		}
	}

}

func main() {
	listener, err := net.Listen("tcp", ":6767")
	err_hand(err, "init")
	defer listener.Close()

	for {
		c, err := listener.Accept()
		err_hand(err, "accept")
		go req_hand(c)
	}

}
