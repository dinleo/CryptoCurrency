package rest

import (
	"CryptoCurrency/blockchain"
	"CryptoCurrency/utils"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var port string

// MarshalText Handling
type url string

func (u url) MarshalText() ([]byte, error) {
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
}

// Struct
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
// Home
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
			URL:         url("/blocks/{hash}"),
			Method:      "GET",
			Description: "See A Block",
		},
	}

	// Encoding
	json.NewEncoder(w).Encode(pageData)
}

// View specific block
func block(w http.ResponseWriter, r *http.Request) {
	// Get hash from URL
	vars := mux.Vars(r)

	// Get block from hash
	hash := vars["hash"]

	block, err := blockchain.FindBlock(hash)

	// Encoding
	encoder := json.NewEncoder(w)
	if err == blockchain.ErrNotFound {
		encoder.Encode(errorResponse{fmt.Sprint(err)})
	} else {
		encoder.Encode(block)
	}
}

// View all blocks
func blocks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// Encoding
		json.NewEncoder(w).Encode(blockchain.Blockchain().Blocks())
	case "POST":
		var a addBlockBody

		// Decoding
		utils.HandleErr(json.NewDecoder(r.Body).Decode(&a))

		// Add block
		blockchain.Blockchain().AddBlock(a.Message)

		// Encoding
		w.WriteHeader(http.StatusCreated)
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
	NewMux.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET")

	fmt.Printf("Listening on http://localhost%s\n", port)

	// Run server
	log.Fatal(http.ListenAndServe(port, NewMux))
}
