package main

import (
	"strings"
)

func parse_rules(conteudo string) string {
	nomes_usados := make(map[string]bool)

	nova_tab := table{"BDINFO", nil}
	tabelas = append(tabelas, nova_tab)
	nomes_usados["BDINFO"] = true

	lines := strings.Split(string(conteudo), "\n")

	for i := 0; i < len(lines); i++ {
		clear_line := strings.TrimSpace(lines[i])
		l := len(clear_line)

		if l <= 0 {
			continue
		}

		if clear_line[0] == '[' && clear_line[l-1] == ']' {
			nome := clear_line[1 : l-1]

			if nomes_usados[nome] {
				tabelas = tabelas[:0]
				error_signal = true
				return "nome de tabela repetido\n"
			}

			nomes_usados[nome] = true
			nova_tab := table{nome, nil}
			tabelas = append(tabelas, nova_tab)
			continue
		}

		tabelas[len(tabelas)-1].rules = append(tabelas[len(tabelas)-1].rules, clear_line)
	}

	for i := 0; i < len(tabelas); i++ {

		if len(tabelas[i].rules) == 0 && i != 0 {
			tabelas = tabelas[:0]
			error_signal = true
			return "sem regras\n"
		}

	}

	return "regras ok\n"
}


//VERIFICAR NOMES!!!