package main

import (
	"errors"
	"fmt"
	"log"
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
	err_hand(err, "criar BD")
}

func create_file(path string, content string) {
	if path_exists(path, false) {
		return
	}

	file, err := os.Create(path)
	err_hand(err, "create_file")
	defer file.Close()

	file.Write([]byte(content))
}

func create_db() string {
	db_path := filepath.Join(default_db_path, nome_input)
	rule_path := filepath.Join(db_path, "regras")

	if path_exists(db_path, true) {
		return fmt.Sprintf("banco já existe [%v]\n", db_path)
	}

	create_folder(db_path)

	var regras_fix string
	for i := 1; i < len(tabelas); i++ {
		regras_fix += fmt.Sprintf("[%v]\n", tabelas[i].name)
		for j := 0; j < len(tabelas[i].rules); j++ {
			regras_fix += fmt.Sprintf("%v\n", tabelas[i].rules[j])
		}

		tabela_path := filepath.Join(db_path, tabelas[i].name)
		create_folder(tabela_path)

		data_path := filepath.Join(tabela_path, "data")
		create_folder(data_path)

		index_path := filepath.Join(tabela_path, "autoindex")
		free_path := filepath.Join(tabela_path, "vazioindex")

		create_file(index_path, "")
		create_file(free_path, "")
	}

	create_file(rule_path, regras_fix)

	return "Criado\n"
}
