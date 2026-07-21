package data

import (
	"errors"
	"fmt"
	"io"
	"nosei/shared"
	"os"
	"path/filepath"
	"strconv"
)

// carregras as regras do banco e seu caminho
// (shared.Bd_carregado_path,
// shared.Tabelas)
func Usar_banco(qual string) (string, error) {
	db_path := filepath.Join(shared.Default_db_path, qual)

	if !path_exists(db_path, true) {
		return "", errors.New("banco não existe\n")
	}

	shared.Bd_carregado_path = db_path
	rules_path := filepath.Join(shared.Bd_carregado_path, "regras.ns")

	if !path_exists(rules_path, false) {
		return "", errors.New("banco sem regras\n")
	}

	regras, err := os.ReadFile(rules_path)
	shared.Err_hand(err, "abrir regras usar banco")

	_, err = Parse_rules(string(regras))
	if err != nil {
		return "", err
	}
	//Show_tabelas()

	return "usando\n", nil
}

func process_data_tipo_store(valor string, regra shared.Rule) (string, error) {
	switch regra.Tipo {
	case shared.Tipo_texto:
		res := make([]byte, regra.Bytes)
		copy(res, valor)

		return string(res), nil

	case shared.Tipo_numero:
		numero, err := strconv.Atoi(valor)
		if err != nil {
			return "", err
		}
		return shared.Get_numero_bytes(numero), nil

	default:
		return "", errors.New("sigma")
	}
}

func process_data_tipo_read(valor string, tipo shared.Regra_tipo) (string, error) {
	switch tipo {
	case shared.Tipo_texto:
		return valor, nil
	case shared.Tipo_numero:
		numero := strconv.Itoa(shared.Get_bytes_numero(valor))

		return numero, nil
	default:
		return "", errors.New("sigma")
	}
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

func overwrite_file_from_pos(file_path string, content string, pos int) error {
	f, err := os.OpenFile(file_path, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteAt([]byte(content), int64(pos))
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

	data_index := shared.Get_bytes_numero(string(content))

	datafile_path := filepath.Join(tabela_path, "data", fmt.Sprintf("data%v.nsd", data_index))
	fileinfo, err := os.Stat(datafile_path)
	if err != nil {
		return "", 0, 0, err
	}

	filesize := int(fileinfo.Size())

	if filesize < shared.Max_data_file_bytes {
		return datafile_path, filesize, data_index, nil
	}

	data_index++
	datafile_path = filepath.Join(tabela_path, "data", fmt.Sprintf("data%v.nsd", data_index))

	create_file(datafile_path, "")
	err = overwrite_file_from_pos(dataindex_path, shared.Get_numero_bytes(data_index), 0)
	if err != nil {
		return "", 0, 0, err
	}

	return datafile_path, 0, data_index, nil
}

// atualizar autoindex
// primeira linha (header) incrementaa
// adiciona: [arquivo][index da entrada][pos no arquivo]
func update_autoindex(auto_index_path string, datafile_index int, block_pos int) (string, error) {
	file, err := os.OpenFile(auto_index_path, os.O_RDWR, 0666)
	if err != nil {
		return "", err
	}
	defer file.Close()

	bytes_index_atual := make([]byte, shared.Default_num_size)
	_, err = io.ReadFull(file, bytes_index_atual)
	if err != nil {
		return "", err
	}

	index_atual := shared.Get_bytes_numero(string(bytes_index_atual))
	new_header := shared.Get_numero_bytes(index_atual + 1)

	_, err = file.WriteAt([]byte(new_header), 0)
	if err != nil {
		return "", err
	}

	new_index := shared.Get_numero_bytes(datafile_index) + string(bytes_index_atual) + shared.Get_numero_bytes(block_pos)
	stat, err := file.Stat()
	if err != nil {
		return "", err
	}

	_, err = file.WriteAt([]byte(new_index), stat.Size())
	if err != nil {
		return "", err
	}

	return "ok\n", nil
}

// retorna novo index
func update_index_autoindex(auto_index_path string) (int, error) {
	file, err := os.OpenFile(auto_index_path, os.O_RDWR, 0666)
	if err != nil {
		return -1, err
	}
	defer file.Close()

	bytes_index_atual := make([]byte, shared.Default_num_size)
	_, err = io.ReadFull(file, bytes_index_atual)
	if err != nil {
		return -1, err
	}

	index_atual := shared.Get_bytes_numero(string(bytes_index_atual))
	new_header := shared.Get_numero_bytes(index_atual + 1)

	_, err = file.WriteAt([]byte(new_header), 0)
	if err != nil {
		return -1, err
	}

	return index_atual, nil
}

// retorna o index da tabela "nome" dentro de shared.Tabelas, se existir
func Get_tabela_index(nome string) (int, error) {
	var tabela_index int = -1
	for i := 1; i < len(shared.Tabelas); i++ {
		if shared.Tabelas[i].Name == nome {
			tabela_index = i
			break
		}
	}
	if tabela_index == -1 {
		return tabela_index, errors.New("tabela nao existe\n")
	}

	return tabela_index, nil
}

func preparar_bloco_entrada(tabela_index int, dados []string, file_pos int) (string, error) {
	qtd_rules := len(shared.Tabelas[tabela_index].Rules)
	if len(dados) != qtd_rules {
		return "", errors.New("tamanho de dados nao condiz com shared.Tabelas\n")
	}

	for i := range dados {
		store, err := process_data_tipo_store(dados[i], shared.Tabelas[tabela_index].Rules[i])
		if err != nil {
			return "", err
		}
		dados[i] = store
	}

	var bloco_data string
	var bloco_headers string
	var bloco_content string

	tamanho_total := shared.Default_num_size + shared.Default_num_size*qtd_rules
	var cur_offset = tamanho_total + file_pos

	for i := range dados {
		tamanho_total += len(dados[i])

		if i > 0 {
			bloco_headers += shared.Get_numero_bytes(cur_offset)
			cur_offset += len(dados[i-1])
		}

		bloco_content += dados[i]
	}

	bloco_headers += shared.Get_numero_bytes(cur_offset)
	bloco_data = shared.Get_numero_bytes(tamanho_total) + bloco_headers + bloco_content

	return bloco_data, nil
}

// retorna uma posicao em autoindex que pode ser reocupada
func pop_ifcan_vazioindex(path string) int {
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		return -1
	}
	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		return -1
	}
	len := fileinfo.Size()

	if len <= 0 {
		return -1
	}

	pos := len - int64(shared.Default_num_size)

	ptr := make([]byte, shared.Default_num_size)
	_, err = file.ReadAt(ptr, pos)
	if err != nil && err != io.EOF {
		return -1
	}

	file.Truncate(pos)
	return shared.Get_bytes_numerob(ptr)
}

// tenta inserir dentro de nome_tabela os dados de uma nova entrada
func Nova_entrada(nome_tabela string, dados []string) (string, error) {
	tabela_index, err := Get_tabela_index(nome_tabela)
	if err != nil {
		return "", err
	}

	tabela_path := filepath.Join(shared.Bd_carregado_path, nome_tabela)
	autoindex_path := filepath.Join(tabela_path, "autoindex.ns")
	vazioindex_path := filepath.Join(tabela_path, "vazioindex.ns")

	empty_pos := pop_ifcan_vazioindex(vazioindex_path)
	if empty_pos >= 0 {
		header, err := get_pos_from_autoindex(autoindex_path, empty_pos)
		if err != nil {
			return "", err
		}
		data := shared.Split_fixed(header, shared.Default_num_size)

		entrada_index, err := update_index_autoindex(autoindex_path)
		if err != nil {
			return "", err
		}

		datafile_index := shared.Get_bytes_numero(data[0]) * -1

		new_file_header := shared.Get_numero_bytes(datafile_index) + shared.Get_numero_bytes(entrada_index)
		err = overwrite_file_from_pos(autoindex_path, new_file_header, empty_pos)
		if err != nil {
			return "", err
		}

		bloco_pos := shared.Get_bytes_numero(data[2])
		bloco_data, err := preparar_bloco_entrada(tabela_index, dados, bloco_pos)
		if err != nil {
			return "", err
		}

		datafile_path := filepath.Join(tabela_path, "data", fmt.Sprintf("data%v.nsd", datafile_index))
		err = overwrite_file_from_pos(datafile_path, bloco_data, bloco_pos)
		if err != nil {
			return "", err
		}
		return "", nil
	}

	datafile_path, datafile_length, datafile_index, err := prepare_data_files(tabela_path)
	if err != nil {
		return "", err
	}

	bloco_data, err := preparar_bloco_entrada(tabela_index, dados, datafile_length)
	if err != nil {
		return "", err
	}

	append_to_file(datafile_path, bloco_data)
	resp, err := update_autoindex(autoindex_path, datafile_index, datafile_length)
	if err != nil {
		return resp, err
	}

	return "feito\n", nil
}

// return:
// indexes, [[], [indexes data1], [indexes data2...]
func get_indexes_from_autoindex(autoindex_path string) ([][]string, error) {
	content, err := os.ReadFile(autoindex_path)
	if err != nil {
		return nil, err
	}

	headers := shared.Split_fixed(string(content[shared.Default_num_size:]), shared.Default_num_size*3)
	if headers == nil {
		return nil, errors.New("tabela vazia\n")
	}

	var indexes [][]string

	for i := range headers {
		current_file := shared.Get_bytes_numero(headers[i][:shared.Default_num_size])
		if current_file <= 0 {
			continue
		}

		for k := len(indexes) - 1; k < current_file; k++ {
			var novo []string
			indexes = append(indexes, novo)
		}

		indexes[current_file] = append(indexes[current_file], headers[i][shared.Default_num_size:])
		//[arquivo][index da entrada][pos no arquivo]
	}

	return indexes, nil
}

func get_pos_from_autoindex(autoindex_path string, pos int) (string, error) {
	cont, err := os.ReadFile(autoindex_path)
	if err != nil {
		return "", err
	}

	bloco := cont[pos : pos+(3*shared.Default_num_size)]
	return string(bloco), nil
}

// return:
// posicao do bloco que contem os metadados da entrada index dentro do arquivo autoindex
// o bloco em si
func find_entrada_em_autoindex(autoindex_path string, index int) (int, string, error) {
	content, err := os.ReadFile(autoindex_path)
	if err != nil {
		return -1, "", err
	}

	if len(content)-shared.Default_num_size <= 0 {
		return -1, "", errors.New("tabela vazia\n")
	}
	header_size := shared.Default_num_size * 3

	for i := shared.Default_num_size; i < len(content); i += (header_size) {
		header := string(content[i : i+header_size])
		res := shared.Split_fixed(header, shared.Default_num_size)
		if shared.Get_bytes_numero(res[1]) == index && shared.Get_bytes_numero(res[0]) > 0 {
			return i, header, nil
		}

	}

	return -1, "", errors.New("index nao existe\n")
}

// procura por todos "indexes" dentro de datafile_path e retorna as informacoes em [][][]byte (lista de tabelas -> tabela -> regra)
func get_entradas_from_datafile(datafile_path string, indexes []string, tabela_index int) ([][][]byte, error) {
	qtd_regras := len(shared.Tabelas[tabela_index].Rules)

	cont, err := os.ReadFile(datafile_path)
	if err != nil {
		return nil, err
	}

	var res [][][]byte
	block_offset := 0
	for i := range indexes {
		var linha [][]byte

		pos := shared.Get_bytes_numero(indexes[i][shared.Default_num_size:])
		total_block_size := shared.Get_bytes_numerob(cont[pos : pos+shared.Default_num_size])

		bloco := string(cont[pos+shared.Default_num_size : pos+total_block_size])

		offset := 0
		total_size_read := 0
		block_offset += total_block_size

		//fmt.Printf("index [%v]:\n", shared.Get_bytes_numero(indexes[i][:shared.Default_num_size]))

		linha = append(linha, []byte(indexes[i][:shared.Default_num_size]))
		for j := 0; j < qtd_regras; j++ {
			read_from := shared.Get_bytes_numero(bloco[offset : offset+shared.Default_num_size])
			to_read := read_from + shared.Tabelas[tabela_index].Rules[j].Bytes

			//col_rule := shared.Tabelas[tabela_index].Rules[j]
			raw_data := cont[read_from:to_read]
			linha = append(linha, raw_data)
			//data, err := process_data_tipo_read(string(cont[read_from:to_read]), col_rule.tipo)
			// if err != nil {
			// 	return nil, err
			// }

			//fmt.Printf("%v -> [%v] | %v(%v)\n", col_rule.descritor, 10, raw_data, len(raw_data))

			total_size_read += to_read
			offset += shared.Default_num_size
		}

		res = append(res, linha)
	}

	return res, nil
}

// mostra o conteudo dentro de shared.Tabela_carregada com as regras de nome_tabela
func Show_tabela_carregada(nome_tabela string) {
	tabela_index, err := Get_tabela_index(nome_tabela)
	if err != nil {
		return
	}

	regras := shared.Tabelas[tabela_index].Rules

	for i := range shared.Tabela_carregada {
		for j, d := range shared.Tabela_carregada[i] {
			if j > 0 {
				processado, err := process_data_tipo_read(string(d), regras[j-1].Tipo)
				if err != nil {
					return
				}

				fmt.Printf("%v -> %v\n", regras[j-1].Descritor, processado)
			} else {
				fmt.Printf("\nindex -> %v\n", shared.Get_bytes_numero(string(d)))
			}
		}
	}
}

// mostra o conteudo dentro de shared.Tabela_carregada com as regras de nome_tabela onde index está dentro de index_s
func Show_tabela_carregada_indexes(nome_tabela string, index_s map[int]bool) {
	tabela_index, err := Get_tabela_index(nome_tabela)
	if err != nil {
		return
	}
	regras := shared.Tabelas[tabela_index].Rules

	for i := range shared.Tabela_carregada {
		for j, d := range shared.Tabela_carregada[i] {
			if j > 0 {
				processado, err := process_data_tipo_read(string(d), regras[j-1].Tipo)
				if err != nil {
					return
				}
				fmt.Printf("%v -> %v\n", regras[j-1].Descritor, processado)
			} else {
				idx := shared.Get_bytes_numero(string(d))
				if !index_s[idx] {
					break
				}
				fmt.Printf("\nindex -> %v\n", idx)
			}
		}
	}
}

func Get_tabela_rules(nome_tabela string) ([]shared.Rule, error) {
	tabela_index, err := Get_tabela_index(nome_tabela)
	if err != nil {
		return nil, err
	}
	return shared.Tabelas[tabela_index].Rules, err
}

// coloca os dados da tabela nome_tabela em
// (shared.Nome_tabela_carregada, shared.Tabela_carregada) para uso posterior
func Carregar_tabela(nome_tabela string) (string, error) {
	// if shared.Nome_tabela_carregada == nome_tabela {
	// 	return "tabela já está carregada\n", nil
	// }

	shared.Tabela_carregada = nil
	shared.Nome_tabela_carregada = ""

	tabela_path := filepath.Join(shared.Bd_carregado_path, nome_tabela)
	autoindex_path := filepath.Join(tabela_path, "autoindex.ns")
	data_path := filepath.Join(tabela_path, "data")

	tabela_index, err := Get_tabela_index(nome_tabela)
	if err != nil {
		return err.Error(), err
	}

	indexes, err := get_indexes_from_autoindex(autoindex_path)
	if err != nil {
		return err.Error(), err
	}

	for i := range indexes {
		if i == 0 {
			continue
		}

		file_path := filepath.Join(data_path, fmt.Sprintf("data%v.nsd", i))

		res, err := get_entradas_from_datafile(file_path, indexes[i], tabela_index)
		if err != nil {
			return err.Error(), err
		}

		shared.Tabela_carregada = append(shared.Tabela_carregada, res...)
	}

	shared.Nome_tabela_carregada = nome_tabela
	return "tabela carregada\n", nil
}

// atualiza na tabela a informacao da entrada index para dados
func Update_entrada(nome_tabela string, index int, dados []string) (string, error) {
	tabela_index, err := Get_tabela_index(nome_tabela)
	if err != nil {
		return "", err
	}

	tabela_path := filepath.Join(shared.Bd_carregado_path, nome_tabela)
	autoindex_path := filepath.Join(tabela_path, "autoindex.ns")

	indexes, err := get_indexes_from_autoindex(autoindex_path)
	if err != nil {
		return "", err
	}

	var file_entrada int = -1
	var pos_entrada int = -1

	brk := false
	for d := range indexes {
		for _, i := range indexes[d] {
			idx, pos := shared.Get_bytes_numero(string(i[:4])), shared.Get_bytes_numero(string(i[4:]))

			if idx == index {
				pos_entrada = pos
				file_entrada = d
				brk = true
				break
			}

		}

		if brk {
			break
		}
	}

	if pos_entrada == -1 || file_entrada == -1 {
		return "", errors.New("index nao existe dentro da tabela.\n")
	}

	datafile_path := filepath.Join(tabela_path, "data", fmt.Sprintf("data%v.nsd", file_entrada))

	bloco, err := preparar_bloco_entrada(tabela_index, dados, pos_entrada)
	if err != nil {
		return "", err
	}

	err = overwrite_file_from_pos(datafile_path, bloco, pos_entrada)
	if err != nil {
		return "", err
	}

	return "atualizado\n", nil
}

// deleta a entrada index de nome_tabela, atualiza seus metadados de arquivo no autoindex e cria um ponteiro do espaco vazio de autoindex em vazioindex
func Deletar_entrada(nome_tabela string, index int) (string, error) {
	_, err := Get_tabela_index(nome_tabela)
	if err != nil {
		return "", err
	}

	tabela_path := filepath.Join(shared.Bd_carregado_path, nome_tabela)
	autoindex_path := filepath.Join(tabela_path, "autoindex.ns")
	vazioindex_path := filepath.Join(tabela_path, "vazioindex.ns")

	pos, bloco, err := find_entrada_em_autoindex(autoindex_path, index)
	if err != nil {
		return "", err
	}

	err = append_to_file(vazioindex_path, shared.Get_numero_bytes(pos))
	if err != nil {
		return "", err
	}

	//merdas quando append vazioindex ok e overwrite falha em seguida, nunca aconteceu
	novo_file_index := shared.Get_numero_bytes(shared.Get_bytes_numero(bloco[:shared.Default_num_size]) * -1)
	novo_bloco := novo_file_index + bloco[shared.Default_num_size:]

	err = overwrite_file_from_pos(autoindex_path, novo_bloco, pos)
	if err != nil {
		return "", err
	}

	return "deletado\n", nil
}
