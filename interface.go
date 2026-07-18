package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

func input_inserir_banco(args []string) (string, error) {
	if len(args) < 2 {
		return "", errors.New("comando de inserir incompleto\n")
	}

	_, err := carregar_tabela(args[1])
	if err != nil {
		return "", err
	}

	r, _ := get_tabela_rules(args[1])
	esperado := len(r)

	if len(args)-2 != esperado {
		return "", errors.New("quantia de argumentos indevidos para tabela\n")
	}

	for i := 2; i < len(args); i++ {
		r, err := hex.DecodeString(args[i])
		if err != nil {
			return "", err
		}

		args[i] = string(r)
	}

	t := time.Now()
	echo, err := inserir_banco(args[1], args[2:])
	d := time.Since(t)

	if err != nil {
		return "", err
	}

	return echo + fmt.Sprintf("%v\n", d), nil
}
