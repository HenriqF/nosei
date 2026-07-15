package main

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

func is_nome_valido(nome string) bool {
	for i := 0; i < len(nome); i++ {
		char := nome[i]

		if !unicode.IsLetter(rune(char)) && !unicode.IsDigit(rune(char)) && char != '_' {
			return false
		}
	}

	return true
}

func show_tabelas() {
	for i := 1; i < len(tabelas); i++ {
		fmt.Printf("TABELA| %v\n", tabelas[i].name)

		for j := 0; j < len(tabelas[i].rules); j++ {
			fmt.Printf("    REGRA| %v\n", tabelas[i].rules[j].comando)
		}
		fmt.Printf("\n")
	}
}

func get_rule(conteudo string) (rule, error) {
	nova_regra := rule{conteudo, "", tipo_nada}

	separado := strings.Split(conteudo, "->")
	if len(separado) < 2 {
		return nova_regra, errors.New("regra mal construida")
	}
	tipo := strings.TrimSpace(separado[0])

	switch tipo {
	case "numero":
		nova_regra.tipo = tipo_numero
	case "letras":
		nova_regra.tipo = tipo_texto
	case "refere":
		nova_regra.tipo = tipo_referencia
	case "binary":
		nova_regra.tipo = tipo_binary
	default:
		fmt.Printf("%v", tipo)

		return nova_regra, errors.New("regra sem tipo especifico")
	}

	valor := strings.TrimSpace(separado[1])
	if !is_nome_valido(valor) {
		return nova_regra, errors.New("regra com nome de valor indevido")
	}

	nova_regra.valor = valor

	return nova_regra, nil
}

func parse_rules(conteudo string) string {
	tabelas = tabelas[:0]
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
				error_signal = true
				return "nome de tabela repetido\n"
			}

			if !is_nome_valido(nome) {
				error_signal = true
				return "nome com caracteres invalidos\n"
			}

			nomes_usados[nome] = true
			nova_tab := table{nome, nil}
			tabelas = append(tabelas, nova_tab)
			continue
		}

		nova_regra, err := get_rule(clear_line)
		if err != nil {
			error_signal = true
			return "erro com uma regra\n"
		}

		tabelas[len(tabelas)-1].rules = append(tabelas[len(tabelas)-1].rules, nova_regra)
	}

	for i := 0; i < len(tabelas); i++ {

		if len(tabelas[i].rules) == 0 && i != 0 {
			error_signal = true
			return "sem regras\n"
		}

	}

	return "regras ok\n"
}
