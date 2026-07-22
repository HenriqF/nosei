package main

import (
	"fmt"
	"math/rand/v2"
	"nosei/data"
	"nosei/ops"
	"nosei/ponte"
	"nosei/shared"
	"os"
	"time"
)

func main() {
	mode := "web"

	if len(os.Args) > 1 {
		mode = os.Args[1]
	}

	switch mode {
	case "web":
		ponte.Web_request_hand("6767")

	case "terminal":
		ponte.Terminal_request_hand()

	case "testepopular":
		data.Usar_banco("beta")

		init := time.Now()
		for i := 0; i < 1000; i++ {
			_, err := data.Nova_entrada("tabela", []string{
				fmt.Sprintf("dados#%v", i),
				fmt.Sprintf("%v", rand.IntN(1000)),
			})

			shared.Confirmar(err, "vai saber")
		}
		fmt.Println(time.Since(init))

	case "testeops":
		data.Usar_banco("beta")
		_, err := data.Carregar_tabela("tabela")
		shared.Confirmar(err, "deu ruim")

		op, err := ops.Preparar_operacao(" 'nome#1' == dados ")

		shared.Confirmar(err, "deu ruim")

		err = ops.Processar_tabela_carregada("tabela", op)

		shared.Confirmar(err, "deu merda")
	}
}
