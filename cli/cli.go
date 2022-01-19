package cli

import (
	"CryptoCurrency/explorer"
	"CryptoCurrency/rest"
	"flag"
	"fmt"
	"os"
)

func usage() {
	fmt.Printf("Welcome to CyptoWorld\n\n")
	fmt.Printf("Please use following flags:\n")
	fmt.Printf("-port=4000:     Set the Port of the server\n")
	fmt.Printf("-mode=rest:     Choose between 'html' and 'rest\n\n")
	os.Exit(1)
}

func Start() {
	if len(os.Args) < 2 {
		usage()
	}

	port := flag.Int("port", 4000, "Set port of the server")
	mode := flag.String("mode", "rest", "Choose between 'html' and 'rest")

	flag.Parse()
	switch *mode {
	case "rest":
		//start rest api
		rest.Start(*port)
	case "html":
		//start html explorer
		explorer.Start(*port)
	default:
		usage()
	}
}
