package ponte

import (
	"encoding/hex"
	"errors"
	"fmt"
	"nosei/data"
	"nosei/ops"
	"nosei/shared"
	"strconv"
	"time"
)

//orquestrar comunicacoes mais complexas

// args: put nome_tabela regras...
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

	t := time.Now()
	echo, err := data.Nova_entrada(args[1], args[2:])
	d := time.Since(t)

	if err != nil {
		return "", err
	}

	return echo + fmt.Sprintf("%v\n", d), nil
}

// args: update nome_tabela index regras...
func Input_update_entrada(args []string) (string, error) {
	if len(args) < 3 {
		return "", errors.New("comando de update incompleto\n")
	}

	r, _ := data.Get_tabela_rules(args[1])
	esperado := len(r)

	println(esperado, len(args)-3)
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

	println(args[1], index, args[3:])
	t := time.Now()
	echo, err := data.Update_entrada(args[1], index, args[3:])
	d := time.Since(t)

	if err != nil {
		return "", err
	}

	return echo + fmt.Sprintf("%v\n", d), nil
}

// args: novo nome_banco regras
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

// args: get nome_tabela (operacao)?
// se operacao != 0, tabela é mostrada.
func Input_ver_tabela(args []string) (string, error) {
	args = args[1:]
	if len(args) <= 0 {
		return "", errors.New("É necessár io um nome de tabela\n")
	}

	t := time.Now()
	_, err := data.Carregar_tabela(args[0])
	dt := time.Since(t)
	if err != nil {
		return "", err
	}

	if len(args) != 2 {
		data.Show_tabela_carregada(args[0])
		return fmt.Sprintf("mostrando (%v)\n", dt), nil
	}

	expressao, err := hex.DecodeString(args[1])
	if err != nil {
		return "", errors.New("erro com expressão hex\n")
	}

	op, err := ops.Preparar_operacao(string(expressao))
	if err != nil {
		return "", err
	}

	t = time.Now()
	res, err := ops.Processar_tabela_carregada(args[0], op)
	dt = time.Since(t)

	if err != nil {
		return "", err
	}

	data.Show_tabela_carregada_indexes(args[0], res)

	return fmt.Sprintf("mostrando (%v)\n", dt), nil
}
