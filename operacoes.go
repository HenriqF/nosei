package main

import (
	"fmt"
	"unicode"
)

type token struct {
	valor string
	tipo  rune_tipo
}

type rune_tipo int

const (
	rune_nome rune_tipo = iota
	rune_numero
	rune_op
	rune_delim
	rune_default
)

/*
1. separar tudo em tokens
2. reverse polish notation

precedencia:
-, + : 1
*, / : 2
p    : -1
se precedencia for menor/= que topo do stack, remova e insira


*/

func is_rune_operator(r rune) bool {
	if r == '+' || r == '-' || r == '*' || r == '/' || r == '=' {
		return true
	}
	return false
}

func get_char_tipo(input rune) rune_tipo {
	if unicode.IsLetter(input) || input == '_' {
		return rune_nome
	}

	if unicode.IsDigit(input) {
		return rune_numero
	}

	if is_rune_operator(input) {
		return rune_op
	}

	if input == '(' || input == ')' {
		return rune_delim
	}

	return rune_default
}

func get_token_precedence(t token) int {
	if t.valor == "+" || t.valor == "-" {
		return 1
	}
	if t.valor == "*" || t.valor == "/" {
		return 2
	}

	return -67
}

func get_tokens(input string) ([]token, error) {
	var tokens []token
	var current_token string

	previous_rune := rune_default

	for _, c := range input {
		ct := get_char_tipo(c)
		if ct != previous_rune || ct == rune_delim {
			if previous_rune != rune_default {
				tokens = append(tokens, token{current_token, previous_rune})
			}
			current_token = ""
		}
		current_token += string(c)
		previous_rune = ct
	}
	if previous_rune != rune_default {
		tokens = append(tokens, token{current_token, previous_rune})
	}

	return tokens, nil
}

func ordenar_tokens(input []token) ([]token, error) {
	var stack []token
	var final []token

	for _, t := range input {
		if t.tipo == rune_nome || t.tipo == rune_numero {
			final = append(final, t)
			continue
		}

		if len(stack) <= 0 {
			stack = append(stack, t)
			continue
		}

		if t.valor == ")" {
			for {
				if stack[len(stack)-1].valor == "(" {
					stack = stack[:len(stack)-1]
					break
				}
				final = append(final, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}

			continue
		}

		if t.valor == "(" {
			stack = append(stack, t)
			continue
		}

		for {
			p_self := get_token_precedence(t)
			p_past := get_token_precedence(stack[len(stack)-1])

			if p_self <= p_past {
				final = append(final, stack[len(stack)-1])
				stack = stack[:len(stack)-1]

			} else {
				stack = append(stack, t)
				break
			}
		}

	}

	for {
		if len(stack) == 0 {
			break
		}
		final = append(final, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	return final, nil
}

func process_operation(input string) (string, error) {

	toks, err := get_tokens(input)
	if err != nil {
		return err.Error(), err
	}

	fmt.Println(toks)

	toks_p, err := ordenar_tokens(toks)
	if err != nil {
		return err.Error(), err
	}

	fmt.Println(toks_p)

	return input, nil
}
