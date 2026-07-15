package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type regra_tipo int

const (
	tipo_texto regra_tipo = iota
	tipo_numero
	tipo_binary
	tipo_referencia
	tipo_nada
)

type rule struct {
	comando string
	valor   string
	tipo    regra_tipo
}

type table struct {
	name  string
	rules []rule
}

var (
	default_db_path string = "C:/Users/henri/Documentos/nosei"

	// tabelas
	tabelas []table

	//criacao de banco
	recieving_rules bool = false
	rule_input      string
	nome_input      string

	//uso de banco
	using_path string

	//erro
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

func finalize_rules() string {
	echo := parse_rules(rule_input)

	if !error_signal {
		show_tabelas()

		echo += create_db()
	}

	return echo
}

func handle_input(input string) string {
	var echo string
	if recieving_rules {
		rule_builder(input)

		if !recieving_rules && !error_signal {
			echo += finalize_rules()
		}

	} else if input == "PAROU" {
		log.Fatalf("PAROU\n")

	} else if strings.HasPrefix(input, "novo ") && len(input) >= 7 {
		recieving_rules = true
		error_signal = false
		rule_input = ""
		nome_input = input[5:]

		echo = fmt.Sprintf("Envie regras para [%v]...\n", nome_input)
	} else if strings.HasPrefix(input, "usar ") && len(input) >= 7 {
		error_signal = false
		nome_banco := input[5:]

		echo = usar_banco(nome_banco)
	}

	return echo
}

func request_hand(c net.Conn) {
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
		err_hand(err, "echo")
	}

}

func terminal_request_hand() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		resposta := scanner.Text()

		clean := strings.TrimSpace(resposta)
		answer := handle_input(clean)

		fmt.Printf("%v", answer)
	}

	err_hand(scanner.Err(), "scanner temrinal")
}

func main() {
	web := true

	if len(os.Args) > 1 {
		if os.Args[1] == "noweb" {
			web = false
		}
	}

	if web {
		listener, err := net.Listen("tcp", ":6767")
		err_hand(err, "init")
		defer listener.Close()

		fmt.Printf("rodando: localhost:6767\n")

		for {
			c, err := listener.Accept()
			err_hand(err, "accept")
			request_hand(c)
		}
	} else {
		fmt.Printf("rodando: terminal\n")

		terminal_request_hand()
	}

}
