package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const int_size int = 4
const max_data_file_bytes int = 64000

func usar_banco(qual string) string {
	db_path := filepath.Join(default_db_path, qual)

	if !path_exists(db_path, true) {
		error_signal = true
		return "banco não existe\n"
	}

	using_bd_path = db_path
	rules_path := filepath.Join(using_bd_path, "regras.ns")

	if !path_exists(rules_path, false) {
		error_signal = true
		return "banco sem regras?\n"
	}

	regras, err := os.ReadFile(rules_path)
	err_hand(err, "abrir regras usar banco")

	parse_rules(string(regras))
	if error_signal {
		return "regras malformadas\n"
	}
	show_tabelas()

	return "usando\n"
}

func get_numero_bytes(numero int) string {
	//return strconv.Itoa(numero)

	bytes := make([]byte, int_size)
	binary.BigEndian.PutUint32(bytes, uint32(numero))

	return string(bytes)
}

func process_data_tipo(valor string, tipo regra_tipo) (string, error) {
	switch tipo {
	case tipo_texto:
		return valor, nil
	case tipo_numero:
		numero, err := strconv.Atoi(valor)
		if err != nil {
			return "", err
		}
		return get_numero_bytes(numero), nil
	default:
		return "", errors.New("sigma")
	}
}

func prepare_data_files(table_path string) (string, int, error) {
	dataindex_path := filepath.Join(table_path, "dataindex.ns")
	content, err := os.ReadFile(dataindex_path)
	if err != nil {
		return "", 0, err
	}

	datafile_path := filepath.Join(table_path, "data", fmt.Sprintf("data%v.nsd", string(content)))
	fileinfo, err := os.Stat(datafile_path)
	if err != nil {
		return "", 0, err
	}

	file_size := fileinfo.Size()

	return datafile_path, int(file_size), nil
}

func update_autoindex(auto_index_path string, index int) (string, error) {
	conteudo, err := os.ReadFile(auto_index_path)
	if err != nil {
		error_signal = true
		return "erro abrindo auto index\n", err
	}

	linhas := strings.Split(string(conteudo), "\n")
	index_autual, err := strconv.Atoi(linhas[0])
	if err != nil {
		return "erro com numero auto index\n", err
	}
	println("index:", index_autual)

	return "ok\n", nil
}

func get_tabela_index(nome string) (int, error) {
	var tabela_index int = -1
	for i := 1; i < len(tabelas); i++ {
		if tabelas[i].name == nome {
			tabela_index = i
			break
		}
	}
	if tabela_index == -1 {
		return tabela_index, errors.New("tabela nao existe")
	}

	return tabela_index, nil
}

func inserir_banco(tabela string) (string, error) {
	tabela_index, err := get_tabela_index(tabela)
	if err != nil {
		return "tabela nao existe\n", err
	}

	tabela_path := filepath.Join(using_bd_path, tabela)
	auto_index_path := filepath.Join(tabela_path, "autoindex.ns")

	datafile_path, datafile_length, err := prepare_data_files(tabela_path)
	if err != nil {
		return "erro com datafile", err
	}

	//PLACEHOLDER!!
	dados := [...]string{"pedro", "6767676767676767"}
	//PLACEHOLDER!!

	qtd_rules := len(tabelas[tabela_index].rules)
	if len(dados) != qtd_rules {
		return "tamanho de dados nao condiz com tabelas\n", errors.New("erro tamanho")
	}

	for i, d := range dados {
		store, err := process_data_tipo(d, tabelas[tabela_index].rules[i].tipo)
		if err != nil {
			return "erro com dados para inserir", errors.New("sigma")
		}
		dados[i] = store
	}

	var bloco_data string
	var bloco_headers string
	var bloco_content string

	tamanho_total := int_size + int_size*qtd_rules
	var cur_offset = tamanho_total + datafile_length

	for i, d := range dados {
		tamanho_total += len(d)

		if i > 0 {
			bloco_headers += get_numero_bytes(cur_offset)
			cur_offset += len(dados[i-1])
		}

		bloco_content += dados[i]
	}
	bloco_headers += get_numero_bytes(cur_offset)
	bloco_data = get_numero_bytes(tamanho_total) + bloco_headers + bloco_content

	file, err := os.OpenFile(datafile_path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return "erro abrindo datafile", err
	}
	defer file.Close()

	if _, err := file.WriteString(bloco_data); err != nil {
		return "erro escrevendo em datafile", err
	}

	fmt.Println(bloco_data, len(bloco_data))

	resp, err := update_autoindex(auto_index_path, datafile_length)
	if err != nil {
		return resp, err
	}

	return "feito\n", nil
}
