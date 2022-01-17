package rest

import (
	"CryptoCurrency/blockchain"
	"CryptoCurrency/utils"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

var port string

// Url handling
type url string

func (u url) MarshalText() ([]byte, error) {
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
}

// struct
type urlDescription struct {
	URL         url    `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"`
}

type addBlockBody struct {
	Message string
}

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

// Handle func
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
			URL:         url("/blocks/{height}"),
			Method:      "GET",
			Description: "See A Block",
		},
	}

	// Encoding
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pageData)

}

func blocks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// Encoding
		json.NewEncoder(w).Encode(blockchain.GetBlockchain().AllBlocks())
	case "POST":
		var a addBlockBody

		// Decoding
		utils.HandleErr(json.NewDecoder(r.Body).Decode(&a))

		// Add block
		blockchain.GetBlockchain().AddBlock(a.Message)

		// Encoding
		w.WriteHeader(http.StatusCreated)
	}
}

func block(w http.ResponseWriter, r *http.Request) {
	// Get height from URL
	vars := mux.Vars(r)

	// Get block from height
	blockNum, err := strconv.Atoi(vars["height"])
	utils.HandleErr(err)
	block, err := blockchain.GetBlockchain().GetBlock(blockNum)

	// Encoding
	encoder := json.NewEncoder(w)
	if err == blockchain.ErrNotFound {
		encoder.Encode(errorResponse{fmt.Sprint(err)})
	} else {
		encoder.Encode(block)
	}
}

// Middleware
func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// Start Server
func Start(aPort int) {
	port = fmt.Sprintf(":%d", aPort)

	// Mux handling
	NewMux := mux.NewRouter()
	NewMux.Use(jsonContentTypeMiddleware)
	NewMux.HandleFunc("/", documentation).Methods("GET")
	NewMux.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	NewMux.HandleFunc("/blocks/{height:[0-9]+}", block).Methods("GET")

	fmt.Printf("Listening on http://localhost%s\n", port)

	// Run server
	log.Fatal(http.ListenAndServe(port, NewMux))
}
