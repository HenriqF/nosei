package shared

import (
	"encoding/binary"
	"fmt"
	"log"
)

const Max_data_file_bytes int = 256000
// const Max_data_file_bytes int = 4000

const Default_num_size int = 4

type Regra_tipo int

const (
	Tipo_texto Regra_tipo = iota
	Tipo_numero
	Tipo_referencia
	Tipo_nada
)

type Rule struct {
	Comando   string
	Descritor string
	Bytes     int
	Tipo      Regra_tipo
}

type Tabela struct {
	Name  string
	Rules []Rule
}

var (
	Default_db_path string = "C:/Users/henri/Documentos/nosei"

	// tabelas (Parse_rules())
	Tabelas []Tabela

	//criacao de banco
	Nomes_usados map[string]bool
	Nome_input   string

	//uso de banco
	Bd_carregado_path     string
	Tabela_carregada      [][][]byte
	Nome_tabela_carregada string
)

func Split_fixed(s string, n int) []string {
	if len(s)%n != 0 {
		return nil
	}

	var sigma []string
	for i := 0; i < len(s); i += n {
		sigma = append(sigma, s[i:i+n])
	}

	return sigma
}

func Get_numero_bytes(numero int) string {
	//return strconv.Itoa(numero)

	bytes := make([]byte, Default_num_size)
	binary.BigEndian.PutUint32(bytes, uint32(numero))

	return string(bytes)
}

func Get_bytes_numero(str string) int {
	n := binary.BigEndian.Uint32([]byte(str))

	m := int32(n)

	return int(m)
}

func Get_bytes_numerob(bytes []byte) int {
	n := binary.BigEndian.Uint32(bytes)

	m := int32(n)

	return int(m)
}

func Err_hand(err error, msg string) {
	if err != nil {
		log.Fatalf("%v\n\ndeu merda:\n\n%v\n", msg, err)
	}
}

func Nc_err_hand(err error, msg string) bool {
	if err != nil {
		fmt.Printf("%v\n\ndeu merda:\n\n%v\n", msg, err)
		return true
	}
	return false
}
