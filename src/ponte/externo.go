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

	case "usar":
		nome_banco := separado[1]
		echo, err = data.Usar_banco(nome_banco)
		if err != nil {
			return err.Error()
		}
		return echo

	case "novo":
		echo, err := Input_novo_banco(separado)
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
		echo, err := Input_ver_tabela(separado)
		if err != nil {
			return err.Error()
		}
		return echo

	case "update":
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

		_, err = c.Write([]byte(answer))
		if err != nil {
			return
		}
	}
}

func Web_request_hand(port string) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	shared.Err_hand(err, "init")
	defer listener.Close()

	fmt.Printf("rodando: localhost:6767\n")

	for {
		c, err := listener.Accept()
		shared.Err_hand(err, "accept")
		web_talk(c)
	}
}

func Terminal_request_hand() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		resposta := scanner.Text()

		clean := strings.TrimSpace(resposta)
		answer := handle_input(clean)

		fmt.Printf("%v", answer)
	}

	if scanner.Err() != nil {
		return
	}
}
