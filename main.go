package main

import (
	"CryptoCurrency/cli"
	"CryptoCurrency/db"
)

func main() {
	defer db.Close()
	cli.Start()
}
