package ponte

//orquestrar comunicacoes mais complexas

import (
	"encoding/hex"
	"errors"
	"nosei/data"
	"nosei/ops"
	"nosei/shared"
	"strconv"
)

func Input_novo_banco(args []string) (string, error) {
	args = args[1:]
	if len(args) != 2 {
		return "", errors.New("quantia de argumentos indevidos para criação de banco\n")
	}
	shared.Nome_input = args[0]
	regras, err := hex.DecodeString(args[1])
	if err != nil {
		return "", errors.New("erro com expressão hex\n")
	}

	// fmt.Println(string(regras))

	_, err = data.Parse_rules(string(regras))
	if err != nil {
		return "", err
	}

	echo, err := data.Create_db()
	if err != nil {
		return "", err
	}

	echo_b, err := data.Usar_banco(args[0])

	return echo + echo_b, nil
}

func Input_nova_entrada(args []string) (string, error) {
	if len(args) < 2 {
		return "", errors.New("comando de inserir incompleto\n")
	}

	r, _ := data.Get_tabela_rules(args[1])
	esperado := len(r)

	if len(args)-2 != esperado {
		return "", errors.New("quantia de argumentos indevidos para tabela\n")
	}

	for i := 2; i < len(args); i++ {
		r, err := hex.DecodeString(args[i])
		if err != nil {
			return "", errors.New("erro com expressão hex\n")
		}

		args[i] = string(r)
	}

	echo, err := data.Nova_entrada(args[1], args[2:])

	if err != nil {
		return "", err
	}

	return echo, nil
}

func Input_update_entrada(args []string) (string, error) {
	if len(args) < 3 {
		return "", errors.New("comando de update incompleto\n")
	}

	r, _ := data.Get_tabela_rules(args[1])
	esperado := len(r)

	if len(args)-3 != esperado {
		return "", errors.New("quantia de argumentos indevidos para tabela\n")
	}

	for i := 2; i < len(args); i++ {
		r, err := hex.DecodeString(args[i])
		if err != nil {
			return "", errors.New("erro com expressão hex\n")
		}

		args[i] = string(r)
	}

	index, err := strconv.Atoi(args[2])
	if err != nil {
		return "", errors.New("Numero de index malformado\n")
	}

	echo, err := data.Update_entrada(args[1], index, args[3:])

	if err != nil {
		return "", err
	}

	return echo, nil
}

func Input_get_tabela(args []string) (string, error) {
	args = args[1:]
	if len(args) != 1 && len(args) != 2 {
		return "", errors.New("É necessár io um nome de tabela\n")
	}

	_, err := data.Carregar_tabela(args[0])
	if err != nil {
		return "", err
	}

	busca := false
	if len(args) == 2 {
		busca = true

		expressao, err := hex.DecodeString(args[1])
		if err != nil {
			return "", errors.New("erro com expressão hex\n")
		}

		op, err := ops.Preparar_operacao(string(expressao))
		if err != nil {
			return "", err
		}

		err = ops.Processar_tabela_carregada(args[0], op)

		if err != nil {
			return "", err
		}
	}

	switch current_hand_type {
	case hand_term:
		data.Show_tabela_carregada(args[0], busca, false)
	case hand_web:
		res := data.Show_tabela_carregada(args[0], busca, true)
		return res, nil
	}

	return "mostrando\n", nil
}

func Input_deletar_entrada(args []string) (string, error) {
	args = args[1:]
	if len(args) != 2 {
		return "", errors.New("É necessário um nome de tabela e index\n")
	}

	numero_hex, err := hex.DecodeString(args[1])
	if err != nil {
		return "", errors.New("Numero de index malformado\n")
	}
	numero, err := strconv.Atoi(string(numero_hex))
	if err != nil {
		return "", errors.New("Numero de index malformado\n")
	}

	echo, err := data.Deletar_entrada(args[0], numero)

	if err != nil {
		return "", err
	}

	return echo, nil
}
