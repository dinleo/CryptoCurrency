package main

import (
	"CryptoCurrency/explorer"
	"CryptoCurrency/rest"
)

func main() {
	go explorer.Start(3000)
	rest.Start(4000)
}
