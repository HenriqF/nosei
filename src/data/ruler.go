package data

import (
	"errors"
	"fmt"
	"nosei/shared"
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

func Show_tabelas() {
	for i := 1; i < len(shared.Tabelas); i++ {
		fmt.Printf("TABELA| %v\n", shared.Tabelas[i].Name)

		for j := 0; j < len(shared.Tabelas[i].Rules); j++ {
			fmt.Printf("    REGRA[%v]b| %v\n", shared.Tabelas[i].Rules[j].Bytes, shared.Tabelas[i].Rules[j].Comando)
		}
		fmt.Printf("\n")
	}
}

func get_rule(conteudo string) (shared.Rule, error) {
	nova_regra := shared.Rule{
		Comando:   conteudo,
		Descritor: "",
		Bytes:     0,
		Tipo:      shared.Tipo_nada,
	}

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
		nova_regra.Tipo = shared.Tipo_numero
		nova_regra.Bytes = shared.Default_num_size

	case "letras":
		nova_regra.Tipo = shared.Tipo_texto

		num, err := strconv.Atoi(prefixo)
		if err != nil {
			return nova_regra, errors.New("numero de bytes invalido\n")
		}

		nova_regra.Bytes = num

	case "refere":
		nova_regra.Tipo = shared.Tipo_referencia
		nova_regra.Bytes = shared.Default_num_size

	default:
		fmt.Printf("%v", tipo)
		return nova_regra, errors.New("regra sem tipo especifico\n")
	}

	descritor := strings.TrimSpace(partes[1])
	if !is_nome_regra_valido(descritor) {
		return nova_regra, errors.New("regra com nome de descritor indevido\n")
	}

	nova_regra.Descritor = descritor
	return nova_regra, nil
}

func get_tabela_name(line string) (string, error) {
	l := len(line)

	if line[0] == '[' && line[l-1] == ']' {
		nome := line[1 : l-1]

		if shared.Nomes_usados[nome] {
			return "", errors.New("nome repetido de tabela")
		}

		if !is_nome_tabela_valido(nome) {
			return "", errors.New("nome de tabela indevido")
		}

		return nome, nil
	}

	return "", nil
}

// interpreta de forma simples as tabelas e suas regras e coloca em shared.Tabelas
func Parse_rules(conteudo string) (string, error) {
	shared.Tabelas = shared.Tabelas[:0]
	shared.Nomes_usados = make(map[string]bool)

	nova_tab := shared.Tabela{
		Name:  "BDINFO",
		Rules: nil,
	}

	shared.Tabelas = append(shared.Tabelas, nova_tab)
	shared.Nomes_usados["BDINFO"] = true

	lines := strings.Split(string(conteudo), "\n")

	for i := 0; i < len(lines); i++ {
		clear_line := strings.TrimSpace(lines[i])
		if len(clear_line) <= 0 {
			continue
		}

		tabela_name, err := get_tabela_name(clear_line)
		if err != nil {
			return "", err
		}
		if tabela_name != "" {
			shared.Nomes_usados[tabela_name] = true
			shared.Tabelas = append(shared.Tabelas, shared.Tabela{Name: tabela_name, Rules: nil})
			continue
		}

		nova_regra, err := get_rule(clear_line)
		if err != nil {
			return "", err
		}

		shared.Tabelas[len(shared.Tabelas)-1].Rules = append(shared.Tabelas[len(shared.Tabelas)-1].Rules, nova_regra)
	}

	for i := 0; i < len(shared.Tabelas); i++ {
		if len(shared.Tabelas[i].Rules) == 0 && i != 0 {
			return "", errors.New("sem regras\n")
		}
	}

	return "regras ok\n", nil
}
