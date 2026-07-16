package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

const max_data_file_bytes int = 2000

func usar_banco(qual string) string {
	db_path := filepath.Join(default_db_path, qual)

	if !path_exists(db_path, true) {
		error_signal = true
		return "banco não existe\n"
	}

	bd_carregado_path = db_path
	rules_path := filepath.Join(bd_carregado_path, "regras.ns")

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
	//show_tabelas()

	return "usando\n"
}

func process_data_tipo_store(valor string, regra rule) (string, error) {
	switch regra.tipo {
	case tipo_texto:
		res := make([]byte, regra.bytes)
		copy(res, valor)

		return string(res), nil

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

func process_data_tipo_read(valor string, tipo regra_tipo) (string, error) {
	switch tipo {
	case tipo_texto:
		return valor, nil
	case tipo_numero:
		numero := strconv.Itoa(get_bytes_numero(valor))

		return numero, nil
	default:
		return "", errors.New("sigma")
	}
}

func update_file_index_header(file_path string, new_header string) error {
	f, err := os.OpenFile(file_path, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteAt([]byte(new_header), 0)
	if err != nil {
		return err
	}

	return nil
}

// ret:
// Caminho do arquivo data.nsd
// Tamanho do arquivo data.nsd
// O index do arquivo data
func prepare_data_files(tabela_path string) (string, int, int, error) {

	dataindex_path := filepath.Join(tabela_path, "dataindex.ns")
	content, err := os.ReadFile(dataindex_path)
	if err != nil {
		return "", 0, 0, err
	}

	data_index := get_bytes_numero(string(content))

	datafile_path := filepath.Join(tabela_path, "data", fmt.Sprintf("data%v.nsd", strconv.Itoa(data_index)))
	fileinfo, err := os.Stat(datafile_path)
	if err != nil {
		return "", 0, 0, err
	}

	filesize := int(fileinfo.Size())

	if filesize < max_data_file_bytes {
		return datafile_path, filesize, data_index, nil
	}

	data_index++
	datafile_path = filepath.Join(tabela_path, "data", fmt.Sprintf("data%v.nsd", strconv.Itoa(data_index)))

	create_file(datafile_path, "")
	err = update_file_index_header(dataindex_path, get_numero_bytes(data_index))
	if err != nil {
		return "", 0, 0, err
	}

	return datafile_path, 0, data_index, nil
}

func append_to_file(file_path string, content string) error {
	f, err := os.OpenFile(file_path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write([]byte(content))
	if err != nil {
		return err
	}

	return nil
}

// atualizar autoindex
// primeira linha (header) incrementaa
// adiciona: [arquivo][index da entrada][pos no arquivo]
func update_autoindex(auto_index_path string, datafile_index int, block_pos int) (string, error) {
	file, err := os.Open(auto_index_path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	bytes_index_atual := make([]byte, default_num_size)
	_, err = io.ReadFull(file, bytes_index_atual)
	if err != nil {
		return "", err
	}

	index_atual := get_bytes_numero(string(bytes_index_atual))

	new_header := get_numero_bytes(index_atual + 1)
	err = update_file_index_header(auto_index_path, new_header)
	if err != nil {
		return "", err
	}

	new_index := get_numero_bytes(datafile_index) + string(bytes_index_atual) + get_numero_bytes(block_pos)
	err = append_to_file(auto_index_path, new_index)
	if err != nil {
		return "", err
	}

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

func inserir_banco(nome_tabela string, dados []string) (string, error) {
	tabela_index, err := get_tabela_index(nome_tabela)
	if err != nil {
		return "tabela nao existe\n", err
	}

	tabela_path := filepath.Join(bd_carregado_path, nome_tabela)
	autoindex_path := filepath.Join(tabela_path, "autoindex.ns")

	datafile_path, datafile_length, datafile_index, err := prepare_data_files(tabela_path)
	if err != nil {
		return "erro com datafile\n", err
	}

	qtd_rules := len(tabelas[tabela_index].rules)
	if len(dados) != qtd_rules {
		return "tamanho de dados nao condiz com tabelas\n", errors.New("erro tamanho")
	}

	for i := range dados {
		store, err := process_data_tipo_store(dados[i], tabelas[tabela_index].rules[i])
		if err != nil {
			return "erro com dados para inserir\n", errors.New("sigma")
		}
		dados[i] = store
	}

	var bloco_data string
	var bloco_headers string
	var bloco_content string

	tamanho_total := default_num_size + default_num_size*qtd_rules
	var cur_offset = tamanho_total + datafile_length

	for i := range dados {
		tamanho_total += len(dados[i])

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
		return "erro abrindo datafile\n", err
	}
	defer file.Close()

	if _, err := file.WriteString(bloco_data); err != nil {
		return "erro escrevendo em datafile\n", err
	}

	resp, err := update_autoindex(autoindex_path, datafile_index, datafile_length)
	if err != nil {
		return resp, err
	}

	return "feito\n", nil
}

// return:
// indexes, [[indexes arq.0] [indexes arq.1]...]
func get_indexes_from_autoindex(autoindex_path string) ([][]string, error) {
	content, err := os.ReadFile(autoindex_path)
	if err != nil {
		return nil, err
	}

	headers := split_fixed(string(content[default_num_size:]), default_num_size*3)
	if headers == nil {
		return nil, errors.New("tabela vazia")
	}

	var indexes [][]string

	for i := range headers {
		current_file := get_bytes_numero(headers[i][:default_num_size])

		for k := len(indexes) - 1; k < current_file; k++ {
			var novo []string
			indexes = append(indexes, novo)
		}

		indexes[current_file] = append(indexes[current_file], headers[i][default_num_size:])
		//[arquivo][index da entrada][pos no arquivo]
	}

	return indexes, nil
}

func get_indexes_from_file(datafile_path string, indexes []string, tabela_index int) ([][][]byte, error) {
	qtd_regras := len(tabelas[tabela_index].rules)

	cont, err := os.ReadFile(datafile_path)
	if err != nil {
		return nil, err
	}

	var res [][][]byte
	block_offset := 0
	for i := range indexes {
		var linha [][]byte

		pos := get_bytes_numero(indexes[i][default_num_size:])
		total_block_size := get_bytes_numerob(cont[pos : pos+default_num_size])

		bloco := string(cont[pos+default_num_size : pos+total_block_size])

		offset := 0
		total_size_read := 0
		block_offset += total_block_size

		//fmt.Printf("index [%v]:\n", get_bytes_numero(indexes[i][:default_num_size]))

		linha = append(linha, []byte(indexes[i][:default_num_size]))
		for j := 0; j < qtd_regras; j++ {
			read_from := get_bytes_numero(bloco[offset : offset+default_num_size])
			to_read := read_from + tabelas[tabela_index].rules[j].bytes

			//col_rule := tabelas[tabela_index].rules[j]
			raw_data := cont[read_from:to_read]
			linha = append(linha, raw_data)
			//data, err := process_data_tipo_read(string(cont[read_from:to_read]), col_rule.tipo)
			// if err != nil {
			// 	return nil, err
			// }

			//fmt.Printf("%v -> [%v] | %v(%v)\n", col_rule.descritor, 10, raw_data, len(raw_data))

			total_size_read += to_read
			offset += default_num_size
		}

		res = append(res, linha)
	}

	return res, nil
}

// mostra o conteudo dentro de tabela_carregada com as regras de nome_tabela
func show_tabela_carregada(nome_tabela string) {
	tabela_index, err := get_tabela_index(nome_tabela)
	if err != nil {
		return
	}

	regras := tabelas[tabela_index].rules

	for i := range tabela_carregada {
		for j, d := range tabela_carregada[i] {
			if j > 0 {
				processado, err := process_data_tipo_read(string(d), regras[j-1].tipo)
				if err != nil {
					return
				}

				fmt.Printf("%v -> %v\n", regras[j-1].descritor, processado)
			} else {
				fmt.Printf("index -> %v\n", get_bytes_numero(string(d)))
			}
		}
		fmt.Println()
	}
}

func carregar_tabela(nome_tabela string) (string, error) {
	tabela_carregada = nil

	tabela_path := filepath.Join(bd_carregado_path, nome_tabela)
	autoindex_path := filepath.Join(tabela_path, "autoindex.ns")
	data_path := filepath.Join(tabela_path, "data")

	tabela_index, err := get_tabela_index(nome_tabela)
	if err != nil {
		return err.Error(), err
	}

	indexes, err := get_indexes_from_autoindex(autoindex_path)
	if err != nil {
		return err.Error(), err
	}

	for i := range indexes {
		file_path := filepath.Join(data_path, fmt.Sprintf("data%v.nsd", i))

		res, err := get_indexes_from_file(file_path, indexes[i], tabela_index)
		if err != nil {
			return err.Error(), err
		}

		tabela_carregada = append(tabela_carregada, res...)
	}

	return "ok\n", nil
}
