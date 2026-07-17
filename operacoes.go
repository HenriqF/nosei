package main

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

type token struct {
	valor     string
	valor_num int
	tipo      rune_tipo
}

type rune_tipo int

const (
	rune_nome rune_tipo = iota
	rune_numero
	rune_op
	rune_delim
	rune_default
)

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
	if t.valor == "u-" {
		return 3
	}
	if t.valor == "!" {
		return 4
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
				tokens = append(tokens, token{current_token, 0, previous_rune})
			}
			current_token = ""
		}
		current_token += string(c)
		previous_rune = ct
	}
	if previous_rune != rune_default {
		tokens = append(tokens, token{current_token, 0, previous_rune})
	}

	for i := range tokens {
		if tokens[i].tipo == rune_numero {
			num, err := strconv.Atoi(tokens[i].valor)
			if err != nil {
				return nil, err
			}

			tokens[i].valor_num = num
			continue
		}

		if tokens[i].tipo == rune_op && tokens[i].valor == "-" {
			if i == 0 {
				tokens[i].valor = "u-"
				continue
			}

			if tokens[i-1].tipo != rune_nome && tokens[i-1].tipo != rune_numero {
				tokens[i].valor = "u-"
				continue
			}
		}

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
				if len(stack) == 0 {
					return final, errors.New("parenteses desbalanceados")
				}

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
			}

			if len(stack) == 0 || p_self > p_past {
				stack = append(stack, t)
				break
			}
		}

	}

	for {
		if len(stack) == 0 {
			break
		}

		if stack[len(stack)-1].valor == "(" {
			return final, errors.New("parenteses desbalanceados")
		}

		final = append(final, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	return final, nil
}

func eval_operation(input []token) (int, error) {
	var stack []int

	for _, d := range input {
		if d.tipo != rune_op {
			stack = append(stack, d.valor_num)
			continue
		}

		if len(stack) == 0 {
			return 0, errors.New("operacao malformada")
		}
		b := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		//UNARIOS
		switch d.valor {
		case "u-":
			stack = append(stack, b*-1)
			continue
		}

		if len(stack) == 0 {
			return 0, errors.New("operacao malformada")
		}
		a := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		//BINARIOS
		switch d.valor {
		case "+":
			stack = append(stack, a+b)
			continue

		case "-":
			stack = append(stack, a-b)
			continue

		case "*":
			stack = append(stack, a*b)
			continue

		case "/":
			stack = append(stack, a/b)
			continue

		}

	}

	return stack[0], nil
}

func process_operation(input string) (string, error) {

	toks, err := get_tokens(input)
	if err != nil {
		return err.Error(), err
	}

	toks_p, err := ordenar_tokens(toks)
	if err != nil {
		return err.Error(), err
	}

	res, err := eval_operation(toks_p)
	if err != nil {
		return err.Error(), err
	}
	fmt.Println(res)

	return input, nil
}
