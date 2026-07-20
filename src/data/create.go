package data

import (
	"errors"
	"fmt"
	"log"
	"nosei/shared"
	"os"
	"path/filepath"
)

func path_exists(path string, is_folder bool) bool {
	info, err := os.Stat(path)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false
		}

		log.Printf("ERRO ABRINDO DIRETORIO\n")
		return true
	}
	if info.IsDir() {
		return is_folder
	}

	return !is_folder
}

func create_folder(path string) {
	err := os.MkdirAll(path, 0755)
	shared.Err_hand(err, "criar BD")
}

func create_file(path string, content string) {
	if path_exists(path, false) {
		return
	}

	file, err := os.Create(path)
	shared.Err_hand(err, "create_file")
	defer file.Close()

	file.Write([]byte(content))
}

func Create_db() (string, error) {
	db_path := filepath.Join(shared.Default_db_path, shared.Nome_input)
	rule_path := filepath.Join(db_path, "regras.ns")

	if path_exists(db_path, true) {
		return "", fmt.Errorf("banco já existe [%v]\n", db_path)
	}

	create_folder(db_path)

	var regras_fix string
	for i := 1; i < len(shared.Tabelas); i++ {
		regras_fix += fmt.Sprintf("[%v]\n", shared.Tabelas[i].Name)
		for j := 0; j < len(shared.Tabelas[i].Rules); j++ {
			regras_fix += fmt.Sprintf("%v\n", shared.Tabelas[i].Rules[j].Comando)
		}

		tabela_path := filepath.Join(db_path, shared.Tabelas[i].Name)
		create_folder(tabela_path)

		data_path := filepath.Join(tabela_path, "data")
		create_folder(data_path)

		index_path := filepath.Join(tabela_path, "autoindex.ns")
		free_path := filepath.Join(tabela_path, "vazioindex.ns")
		data_index_path := filepath.Join(tabela_path, "dataindex.ns")
		data_zero_path := filepath.Join(data_path, "data1.nsd")

		create_file(index_path, shared.Get_numero_bytes(0))
		create_file(free_path, "")
		create_file(data_index_path, shared.Get_numero_bytes(1))
		create_file(data_zero_path, "")
	}

	create_file(rule_path, regras_fix)

	return "Criado\n", nil
}
