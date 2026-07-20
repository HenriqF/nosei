package main

import (
	"fmt"
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
		data.Usar_banco("beta")
		ponte.Terminal_request_hand()

	case "testeshow":
		data.Usar_banco("beta")
		echo, _ := data.Carregar_tabela("tabela")
		fmt.Println(echo)
		data.Show_tabela_carregada("tabela")

	case "testepopular":
		data.Usar_banco("beta")

		init := time.Now()
		for i := 0; i < 30; i++ {
			_, err := data.Nova_entrada("tabela", []string{
				fmt.Sprintf("nome#%v", i),
				fmt.Sprintf("%v", i),
			})

			shared.Err_hand(err, "vai saber")
		}
		fmt.Println(time.Since(init))

	case "testeops":
		data.Usar_banco("beta")
		_, err := data.Carregar_tabela("tabela")
		shared.Err_hand(err, "deu ruim")

		op, err := ops.Preparar_operacao("index != 200")
		shared.Err_hand(err, "deu ruim")

		_, err = ops.Processar_tabela_carregada("tabela", op)

		shared.Err_hand(err, "deu merda")
	}
}
