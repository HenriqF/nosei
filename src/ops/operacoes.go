package ops

import (
	"errors"
	"nosei/data"
	"nosei/shared"
	"strconv"
	"strings"
	"unicode"
)

type token struct {
	valor        string
	valor_num    int
	valor_string string
	tipo         rune_tipo
	tipo_og      rune_tipo
}

type rune_tipo int

const (
	rune_nome rune_tipo = iota
	rune_numero
	rune_op
	rune_delim
	rune_string
	rune_default
)

func bool_to_int(b bool) int {
	if b {
		return 1
	}
	return 0
}

func int_to_bool(i int) bool {
	if i == 0 {
		return false
	}
	return true
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

	if input == '(' || input == ')' || input == '\'' || input == '"' {
		return rune_delim
	}

	return rune_default
}

func is_rune_operator(r rune) bool {
	if r == '+' || r == '-' || r == '*' || r == '/' {
		return true
	}
	if r == '=' || r == '>' || r == '<' || r == '!' {
		return true
	}
	if r == '|' || r == '&' || r == '!' {
		return true
	}
	return false
}

func get_token_precedence(t token) int {
	if t.valor == "||" {
		return -2
	}
	if t.valor == "&&" {
		return -1
	}
	if t.valor == "==" || t.valor == ">" || t.valor == "<" || t.valor == ">=" || t.valor == "<=" || t.valor == "!=" {
		return 0
	}
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

func find_string_end(op string, i int) int {
	pair := rune(op[i])
	size := len(op) - 1

	i++
	for {
		c := rune(op[i])
		ct := get_char_tipo(c)

		if ct == rune_delim && c == pair {
			break
		}

		if i == size {
			return -1
		}

		i++
	}

	return i
}

func get_tokens(input string) ([]token, error) {
	var tokens []token

	var current_token string
	previous_rune := rune_default

	for i := 0; i < len(input); i++ {
		c := rune(input[i])
		ct := get_char_tipo(c)

		if c == '\'' || c == '"' {
			end := find_string_end(input, i)
			if end < 0 {
				return nil, errors.New("string nao terminada")
			}

			tokens = append(tokens, token{input[i+1 : end], 0, input[i+1 : end], rune_string, rune_string})

			i = end
			current_token = ""
			previous_rune = rune_default

			continue
		}

		if ct != previous_rune || ct == rune_delim {
			if previous_rune != rune_default {
				tokens = append(tokens, token{current_token, 0, "", previous_rune, previous_rune})
			}
			current_token = ""
		}

		current_token += string(c)
		previous_rune = ct
	}
	if previous_rune != rune_default {
		tokens = append(tokens, token{current_token, 0, "", previous_rune, previous_rune})
	}

	for i := range tokens {
		if tokens[i].tipo == rune_nome {
			tokens[i].valor_num = -1
		}

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
		if t.tipo == rune_nome || t.tipo == rune_numero || t.tipo == rune_string {
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

func eval_operacao(input []token) (int, error) {
	var stack []token

	for _, d := range input {
		if d.tipo != rune_op {
			stack = append(stack, d)
			continue
		}

		if len(stack) == 0 {
			return 0, errors.New("operacao malformada")
		}
		b := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		bn, bs := b.valor_num, b.valor_string

		if b.tipo != rune_string {
			switch d.valor {
			case "u-":
				stack = append(stack, token{"", bn * -1, bs, b.tipo, b.tipo_og})
				continue

			case "!":
				stack = append(stack, token{"", bool_to_int(!int_to_bool(bn)), bs, b.tipo, b.tipo_og})
				continue
			}
		}

		if len(stack) == 0 {
			return 0, errors.New("operacao malfarmada")
		}
		a := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		an, as := a.valor_num, a.valor_string

		if a.tipo != b.tipo {
			return 0, errors.New("operacao com tipos diferentes")
		}

		if a.tipo == rune_numero {
			switch d.valor {
			case "+":
				stack = append(stack, token{"", an + bn, "", b.tipo, b.tipo_og})
				continue

			case "-":
				stack = append(stack, token{"", an - bn, "", b.tipo, b.tipo_og})
				continue

			case "*":
				stack = append(stack, token{"", an * bn, "", b.tipo, b.tipo_og})
				continue

			case "/":
				stack = append(stack, token{"", an / bn, "", b.tipo, b.tipo_og})
				continue

			case "==":
				stack = append(stack, token{"", bool_to_int(an == bn), "", b.tipo, b.tipo_og})
				continue

			case ">":
				stack = append(stack, token{"", bool_to_int(an > bn), "", b.tipo, b.tipo_og})
				continue

			case "<":
				stack = append(stack, token{"", bool_to_int(an < bn), "", b.tipo, b.tipo_og})
				continue

			case "<=":
				stack = append(stack, token{"", bool_to_int(an <= bn), "", b.tipo, b.tipo_og})
				continue

			case ">=":
				stack = append(stack, token{"", bool_to_int(an >= bn), "", b.tipo, b.tipo_og})
				continue

			case "!=":
				stack = append(stack, token{"", bool_to_int(an != bn), "", b.tipo, b.tipo_og})
				continue

			case "&&":
				stack = append(stack, token{"", bool_to_int(int_to_bool(an) && int_to_bool(bn)), "", b.tipo, b.tipo_og})
				continue

			case "||":
				stack = append(stack, token{"", bool_to_int(int_to_bool(an) || int_to_bool(bn)), "", b.tipo, b.tipo_og})
				continue

			default:
				return 0, errors.New("operador inexistente")
			}
		}

		if a.tipo == rune_string {
			switch d.valor {
			case "+":
				stack = append(stack, token{"", 0, as + bs, rune_numero, b.tipo_og})
				continue

			case "==":
				stack = append(stack, token{"", bool_to_int(as == bs), "", rune_numero, b.tipo_og})
				continue

			case "!=":
				stack = append(stack, token{"", bool_to_int(as != bs), "", rune_numero, b.tipo_og})
				continue

			default:
				return 0, errors.New("operador invalido com strings")
			}
		}

	}

	return stack[0].valor_num, nil
}

func Preparar_operacao(input string) ([]token, error) {
	if input == "" {
		return nil, errors.New("operacao vazia")
	}

	toks, err := get_tokens(input)
	if err != nil {
		return nil, err
	}

	toks_p, err := ordenar_tokens(toks)
	if err != nil {
		return nil, err
	}

	// _, err = eval_operacao(toks_p)
	// if err != nil {
	// 	return nil, err
	// }

	return toks_p, nil
}

func trocar_nome_por_valor(operacao []token, regras []shared.Rule, tabela [][]byte) ([]token, error) {

	for k, t := range operacao {
		if t.tipo_og != rune_nome {
			continue
		}

		if t.valor == "index" {
			operacao[k].valor_num = shared.Get_bytes_numero(string(tabela[0]))
			operacao[k].tipo = rune_numero
			continue
		}

		for i := range regras {
			if regras[i].Descritor != t.valor {
				continue
			}

			if regras[i].Tipo == shared.Tipo_numero {
				operacao[k].valor_num = shared.Get_bytes_numero(string(tabela[i+1]))
				operacao[k].tipo = rune_numero
				break
			}
			if regras[i].Tipo == shared.Tipo_texto {
				operacao[k].valor_num = 0
				operacao[k].valor_string = strings.TrimRight(string(tabela[i+1]), "\x00")
				operacao[k].tipo = rune_string
				break
			}
		}

		if operacao[k].valor_num == -1 {
			return nil, errors.New("nome de regra inexistente na tabela")
		}

	}

	return operacao, nil
}

func Processar_tabela_carregada(nome_tabela string, operacao []token) (map[int]bool, error) {
	tabela_index, err := data.Get_tabela_index(nome_tabela)
	if err != nil {
		return nil, err
	}

	regras := shared.Tabelas[tabela_index].Rules

	corretos := make(map[int]bool)
	for _, t := range shared.Tabela_carregada {
		op, err := trocar_nome_por_valor(operacao, regras, t)

		if err != nil {
			return nil, err
		}

		res, err := eval_operacao(op)
		if err != nil {
			return nil, err
		}
		if res != 0 {
			corretos[shared.Get_bytes_numero(string(t[0]))] = true
			//fmt.Printf("index: %v\n", shared.Get_bytes_numero(string(t[0])))
		}
	}

	return corretos, nil
}
