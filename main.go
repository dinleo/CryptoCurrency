package main

import (
	"CryptoCurrency/db"
	"CryptoCurrency/explorer"
)

func main() {
	defer db.Close()
	explorer.Start(4000)
}
