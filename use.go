package main

import (
	"os"
	"path/filepath"
)

func usar_banco(qual string) string {
	db_path := filepath.Join(default_db_path, qual)

	if !path_exists(db_path, true) {
		error_signal = true
		return "banco não existe\n"
	}

	using_path = db_path
	rules_path := filepath.Join(using_path, "regras.ns")

	if !path_exists(rules_path, false) {
		error_signal = true
		return "banco sem regras?\n"
	}

	regras, err := os.ReadFile(rules_path)
	err_hand(err, "abrir regras usar banco")

	parse_rules(string(regras))
	if error_signal {
		return "regras malformadas"
	}
	show_tabelas()

	return "usando\n"
}
