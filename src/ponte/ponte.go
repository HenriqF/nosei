package ponte

import (
	"encoding/hex"
	"errors"
	"fmt"
	"nosei/data"
	"nosei/ops"
	"nosei/shared"
	"time"
)

//orquestrar comunicacoes mais complexas

func Input_inserir_banco(args []string) (string, error) {
	if len(args) < 2 {
		return "", errors.New("comando de inserir incompleto\n")
	}

	// _, err := data.Carregar_tabela(args[1])
	// if err != nil {
	// 	return "", err
	// }

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

	t := time.Now()
	echo, err := data.Inserir_banco(args[1], args[2:])
	d := time.Since(t)

	if err != nil {
		return "", err
	}

	return echo + fmt.Sprintf("%v\n", d), nil
}

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

	fmt.Println(string(regras))

	_, err = data.Parse_rules(string(regras))
	if err != nil {
		return "", err
	}

	echo, err := data.Create_db()
	if err != nil {
		return "", err
	}

	return echo, nil
}

func Input_ver_tabela(args []string) (string, error) {
	args = args[1:]
	if len(args) <= 0 {
		return "", errors.New("É necessário um nome de tabela\n")
	}
	_, err := data.Carregar_tabela(args[0])
	if err != nil {
		return "", err
	}

	if len(args) != 2 {
		data.Show_tabela_carregada(args[0])
		return "mostrado\n", nil
	}

	expressao, err := hex.DecodeString(args[1])
	if err != nil {
		return "", errors.New("erro com expressão hex\n")
	}

	op, err := ops.Preparar_operacao(string(expressao))
	if err != nil {
		return "", err
	}
	res, err := ops.Processar_tabela_carregada(args[0], op)
	if err != nil {
		return "", err
	}

	data.Show_tabela_carregada_indexes(args[0], res)

	return "mostrado com operacao\n", nil
}
