package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

const max_data_file_bytes int = 256000
const default_num_size int = 4

type regra_tipo int

const (
	tipo_texto regra_tipo = iota
	tipo_numero
	tipo_referencia
	tipo_nada
)

type rule struct {
	comando   string
	descritor string
	bytes     int
	tipo      regra_tipo
}

type tabela struct {
	name  string
	rules []rule
}

var (
	default_db_path string = "C:/Users/henri/Documentos/nosei"

	// tabelas
	tabelas []tabela

	//criacao de banco
	nomes_usados    map[string]bool
	recieving_rules bool = false
	rule_input      string
	nome_input      string

	//uso de banco
	bd_carregado_path     string
	tabela_carregada      [][][]byte
	nome_tabela_carregada string
)

// funcs ajuda
func split_fixed(s string, n int) []string {
	if len(s)%n != 0 {
		return nil
	}

	var sigma []string
	for i := 0; i < len(s); i += n {
		sigma = append(sigma, s[i:i+n])
	}

	return sigma
}

func get_numero_bytes(numero int) string {
	//return strconv.Itoa(numero)

	bytes := make([]byte, default_num_size)
	binary.BigEndian.PutUint32(bytes, uint32(numero))

	return string(bytes)
}

func get_bytes_numero(str string) int {

	n := binary.BigEndian.Uint32([]byte(str))

	return int(n)
}

func get_bytes_numerob(bytes []byte) int {
	n := binary.BigEndian.Uint32(bytes)

	return int(n)
}

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

//-----------

func rule_builder(part string) {
	if part == "end" {
		recieving_rules = false
		return
	}

	rule_input += (part + "\n")
}

func finalize_rules() string {
	echo, err := parse_rules(rule_input)

	if err == nil {
		show_tabelas()
		echo += create_db()
	}

	return echo
}

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
		echo, err = usar_banco(nome_banco)
		if err != nil {
			return err.Error()
		}
		return echo

	case "novo":
		recieving_rules = true
		rule_input = ""

		nome_input = separado[1]
		echo = fmt.Sprintf("Envie regras para [%v]...\n", nome_input)

		return echo

	case "put":
		echo, err := input_inserir_banco(separado)
		if err != nil {
			return err.Error()
		}
		return echo

	}

	return ""
}

// conversariada
func web_request_hand(c net.Conn) {
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

func terminal_request_hand() {
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

//----------

func main() {
	mode := "web"

	if len(os.Args) > 1 {
		mode = os.Args[1]
	}

	switch mode {
	case "web":
		listener, err := net.Listen("tcp", ":6767")
		err_hand(err, "init")
		defer listener.Close()

		fmt.Printf("rodando: localhost:6767\n")

		for {
			c, err := listener.Accept()
			err_hand(err, "accept")
			web_request_hand(c)
		}

	case "usar":
		terminal_request_hand()

	case "terminal":
		fmt.Printf("rodando: terminal\n")

		usar_banco("beta")
		echo, _ := carregar_tabela("tabela")
		fmt.Println(echo)

		show_tabela_carregada("tabela")

	case "pop":
		fmt.Println("popin")
		fmt.Println(usar_banco("beta"))

		init := time.Now()
		for i := 0; i < 5; i++ {
			_, err := inserir_banco("tabela", []string{
				fmt.Sprintf("pedrinho#%v", i),
				fmt.Sprintf("Eu gosto muito de %v", i+3),
				fmt.Sprintf("A#%v", i),
				fmt.Sprintf("%v", i),
				fmt.Sprintf("%v", i),
				fmt.Sprintf("%v", i+420),
			})

			err_hand(err, "vai saber")
		}
		fmt.Println(time.Since(init))

		fmt.Println("ok!")

	case "op":
		usar_banco("beta")
		_, err := carregar_tabela("tabela")
		err_hand(err, "deu ruim")

		op, err := preparar_operacao("nome == descricao")
		err_hand(err, "deu ruim")

		err = processar_tabela_carregada("tabela", op)

		err_hand(err, "deu merda")

	}

}
