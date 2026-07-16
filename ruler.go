package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func is_nome_tabela_valido(nome string) bool {
	if len(nome) < 1 {
		return false
	}

	for i := 0; i < len(nome); i++ {
		char := nome[i]

		if !unicode.IsLetter(rune(char)) && !unicode.IsDigit(rune(char)) && char != '_' {
			return false
		}
	}

	return true
}

func is_nome_regra_valido(nome string) bool {
	if len(nome) < 1 || nome == "index" {
		return false
	}

	for i := 0; i < len(nome); i++ {
		char := nome[i]

		if !unicode.IsLetter(rune(char)) && char != '_' {
			return false
		}
	}

	return true
}

func show_tabelas() {
	for i := 1; i < len(tabelas); i++ {
		fmt.Printf("TABELA| %v\n", tabelas[i].name)

		for j := 0; j < len(tabelas[i].rules); j++ {
			fmt.Printf("    REGRA[%v]b| %v\n", tabelas[i].rules[j].bytes, tabelas[i].rules[j].comando)
		}
		fmt.Printf("\n")
	}
}

func get_rule(conteudo string) (rule, error) {
	nova_regra := rule{conteudo, "", 0, tipo_nada}
	partes := strings.Split(conteudo, "->")
	if len(partes) < 2 {
		return nova_regra, errors.New("regra mal construida\n")
	}

	var prefixo string
	tipo := strings.TrimSpace(partes[0])

	ver := strings.Split(tipo, " ")
	if len(ver) > 1 {
		prefixo = ver[0]
		tipo = ver[1]
	}

	switch tipo {
	case "numero":
		nova_regra.tipo = tipo_numero
		nova_regra.bytes = default_num_size

	case "letras":
		nova_regra.tipo = tipo_texto

		num, err := strconv.Atoi(prefixo)
		if err != nil {
			return nova_regra, errors.New("numero de bytes invalido\n")
		}

		nova_regra.bytes = num

	case "refere":
		nova_regra.tipo = tipo_referencia
		nova_regra.bytes = default_num_size

	default:
		fmt.Printf("%v", tipo)
		return nova_regra, errors.New("regra sem tipo especifico\n")
	}

	descritor := strings.TrimSpace(partes[1])
	if !is_nome_regra_valido(descritor) {
		return nova_regra, errors.New("regra com nome de descritor indevido\n")
	}

	nova_regra.descritor = descritor
	return nova_regra, nil
}

func get_tabela_name(line string) (string, error) {
	l := len(line)

	if line[0] == '[' && line[l-1] == ']' {
		nome := line[1 : l-1]

		if nomes_usados[nome] {
			return "", errors.New("nome repetido de tabela")
		}

		if !is_nome_tabela_valido(nome) {
			return "", errors.New("nome de tabela indevido")
		}

		return nome, nil
	}

	return "", nil
}

func parse_rules(conteudo string) (string, error) {
	tabelas = tabelas[:0]
	nomes_usados = make(map[string]bool)

	nova_tab := tabela{"BDINFO", nil}
	tabelas = append(tabelas, nova_tab)
	nomes_usados["BDINFO"] = true

	lines := strings.Split(string(conteudo), "\n")

	for i := 0; i < len(lines); i++ {
		clear_line := strings.TrimSpace(lines[i])
		if len(clear_line) <= 0 {
			continue
		}

		tabela_name, err := get_tabela_name(clear_line)
		if err != nil {
			return err.Error(), err
		}
		if tabela_name != "" {
			nomes_usados[tabela_name] = true
			tabelas = append(tabelas, tabela{tabela_name, nil})
			continue
		}

		nova_regra, err := get_rule(clear_line)
		if err != nil {
			return err.Error(), err
		}

		tabelas[len(tabelas)-1].rules = append(tabelas[len(tabelas)-1].rules, nova_regra)
	}

	for i := 0; i < len(tabelas); i++ {
		if len(tabelas[i].rules) == 0 && i != 0 {
			return "sem regras\n", errors.New("sem regras")
		}
	}

	return "regras ok\n", nil
}
