package ponte

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"nosei/data"
	"nosei/shared"
	"os"
	"strings"
)

type hand_type int

const (
	hand_term hand_type = iota
	hand_web
)

var current_hand_type hand_type

func handle_input(input string) string {
	var echo string
	var err error

	separado := strings.Split(input, " ")

	if len(separado) == 1 {
		if input == "PAROU" {
			log.Fatalf("PAROU\n")
		}
	}

	switch separado[0] {
	case "novo":
		echo, err := Input_novo_banco(separado)
		if err != nil {
			return err.Error()
		}
		return echo

	case "usar":
		nome_banco := separado[1]
		echo, err = data.Usar_banco(nome_banco)
		if err != nil {
			return err.Error()
		}
		return echo

	case "put":
		echo, err := Input_nova_entrada(separado)
		if err != nil {
			return err.Error()
		}
		return echo

	case "get":
		echo, err := Input_get_tabela(separado)
		if err != nil {
			return err.Error()
		}
		return echo

	case "upd":
		echo, err := Input_update_entrada(separado)
		if err != nil {
			return err.Error()
		}
		return echo

	case "del":
		echo, err := Input_deletar_entrada(separado)
		if err != nil {
			return err.Error()
		}
		return echo
	}

	return ""
}

func format_output(input string) string {
	if input == ""{
		return "\n"
	}

	if input[len(input)-1] != '\n' {
		return input + "\n"
	}

	return input
}

func web_talk(c net.Conn) {
	defer c.Close()
	reader := bufio.NewReader(c)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		clean := strings.TrimSpace(msg)
		answer := handle_input(clean)
		answer = format_output(answer)

		_, err = c.Write([]byte(answer))
		if err != nil {
			return
		}
	}
}

func Web_request_hand(port string) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	shared.Confirmar(err, "init")
	defer listener.Close()

	fmt.Printf("rodando: localhost:6767\n")
	current_hand_type = hand_web

	for {
		c, err := listener.Accept()
		shared.Confirmar(err, "accept")
		web_talk(c)
	}
}

func Terminal_request_hand() {
	current_hand_type = hand_term
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		resposta := scanner.Text()

		clean := strings.TrimSpace(resposta)
		answer := handle_input(clean)
		answer = format_output(answer)

		fmt.Printf("%v", answer)
	}

	if scanner.Err() != nil {
		return
	}
}
