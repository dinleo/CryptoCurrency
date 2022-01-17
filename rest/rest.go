package rest

import (
	"CryptoCurrency/blockchain"
	"CryptoCurrency/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var port string

type url string

func (u url) MarshalText() ([]byte, error) {
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
}

type urlDescription struct {
	URL         url    `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"`
}

type addBlockBody struct {
	Message string
}

func documentation(w http.ResponseWriter, r *http.Request) {
	pageData := []urlDescription{
		{
			URL:         url("/"),
			Method:      "GET",
			Description: "See Documentation",
		},
		{
			URL:         url("/blocks"),
			Method:      "POST",
			Description: "Add A Block",
			Payload:     "data:string",
		},
		{
			URL:         url("/blocks/{id}"),
			Method:      "GET",
			Description: "See A Block",
		},
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pageData)

}

func blocks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(blockchain.GetBlockchain().AllBlocks())
	case "POST":
		var a addBlockBody
		utils.HandleErr(json.NewDecoder(r.Body).Decode(&a))
		blockchain.GetBlockchain().AddBlock(a.Message)
		w.WriteHeader(http.StatusCreated)
	}
}

func Start(aPort int) {
	port = fmt.Sprintf(":%d", aPort)

	mux := http.NewServeMux()
	mux.HandleFunc("/", documentation)
	mux.HandleFunc("/blocks", blocks)

	fmt.Printf("Listening on http://localhost%s\n", port)

	log.Fatal(http.ListenAndServe(port, mux))
}
