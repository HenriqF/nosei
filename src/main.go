package main

import (
	"nosei/ponte"
	"os"
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
	}
}
