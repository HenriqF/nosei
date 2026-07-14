package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type table struct {
	name  string
	rules []string
}

var (
	default_db_path string = "C:/Users/henri/Documentos/nosei"

	recieving_rules bool = false
	rule_input      string
	nome_input      string
	tabelas         []table

	error_signal bool = false
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

func rule_builder(part string) {
	if part == "end" {
		recieving_rules = false
		return
	}

	rule_input += (part + "\n")
}

func request_hand(c net.Conn) {
	defer c.Close()
	reader := bufio.NewReader(c)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		var echo string
		msg_clear := strings.TrimSpace(msg)

		if recieving_rules {
			rule_builder(msg_clear)

			if !recieving_rules && !error_signal {
				echo += parse_rules(rule_input)

				for i := 0; i < len(tabelas); i++ {
					fmt.Printf("TABELA| %v\n", tabelas[i].name)

					for j := 0; j < len(tabelas[i].rules); j++ {
						fmt.Printf("    REGRA| %v\n", tabelas[i].rules[j])
					}
					fmt.Printf("\n")
				}

				echo += create_db()
			}
		} else if msg_clear == "PAROU" {
			return
		} else if strings.HasPrefix(msg_clear, "NEW ") && len(msg_clear) >= 7 {
			recieving_rules = true
			error_signal = false
			rule_input = ""
			nome_input = msg_clear[4:]

			echo = fmt.Sprintf("Envie regras para [%v]...\n", nome_input)
		}

		_, err = c.Write([]byte(echo))
		err_hand(err, "echo")
	}

}

func main() {
	listener, err := net.Listen("tcp", ":6767")
	err_hand(err, "init")
	defer listener.Close()

	fmt.Printf("rodando: localhost:6767\n")

	for {
		c, err := listener.Accept()
		err_hand(err, "accept")
		request_hand(c)
	}
}
