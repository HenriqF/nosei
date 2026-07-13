package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

const (
	data_name string = "data.txt"
)

func err_hand(err error, msg string) {
	if err != nil {
		log.Fatalf("%v\n\ndeu merda:\n\n%v\n", msg, err)
	}
}

func nc_err_hand(err error, msg string) bool {
	if err != nil {
		fmt.Printf("%v\n\ndeu merda:\n\n%v\n", msg, err)
		return true
	}
	return false
}

func req_hand(c net.Conn, input_chain chan string) {
	defer c.Close()
	fmt.Printf("Conex: %v\n", c)

	red := bufio.NewReader(c)
	for {
		msg, err := red.ReadString('\n')

		if err == io.EOF {
			return
		}
		err_hand(err, "reader")

		fmt.Printf("%v: %v\n", c, msg)

		var echo string

		d := strings.TrimSpace(msg)
		if d == "PAROU" {
			return
		} else if strings.HasPrefix(d, "WRITE ") {
			resto := d[6:]
			input_chain <- resto

			echo = fmt.Sprintf("WROTE: [%v]\n", resto)
		}

		_, err = c.Write([]byte(echo))
		err_hand(err, "echo")
	}

}

func linewriter_init(file *os.File, ic chan string, ok chan bool) {
	for line := range ic {
		file.WriteString(line + "\n")
	}

	ok <- true
}

func main() {
	file, _ := os.OpenFile(data_name, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	defer file.Close()

	input_chain := make(chan string, 50)
	lw_done := make(chan bool)
	go linewriter_init(file, input_chain, lw_done)

	listener, err := net.Listen("tcp", ":6767")
	err_hand(err, "init")
	defer listener.Close()

	for {
		c, err := listener.Accept()
		err_hand(err, "accept")
		go req_hand(c, input_chain)
	}
}
