package explorer

import (
	"CryptoCurrency/blockchain"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
)

const (
	templateDir string = "explorer/templates/"
)

var (
	tmpl *template.Template
	port string
)

type pageData struct {
	PageTitle  string
	PageHeader string
	Wallets    []Wallets
	Blocks     []*blockchain.Block
	UTxOuts    []*blockchain.UTxOut
	MemTx      []*blockchain.Tx
}

type Wallets struct {
	WalletName string
	Balance    int
}

type addTxPayload struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount int    `json:"amount"`
}

// Handle func
func home(w http.ResponseWriter, r *http.Request) {
	data := pageData{PageTitle: "HomePage", PageHeader: "HomePage", Blocks: blockchain.Blocks(blockchain.Blockchain())}
	tmpl.ExecuteTemplate(w, "home", data)
}

func wallets(w http.ResponseWriter, r *http.Request) {
	var wallets []Wallets
	for _, u := range blockchain.Wallets(blockchain.Blockchain()) {
		wallet := Wallets{WalletName: u, Balance: blockchain.BalanceByAddress(u, blockchain.Blockchain())}
		wallets = append(wallets, wallet)
	}
	data := pageData{PageTitle: "Wallets", PageHeader: "Wallet", Wallets: wallets}
	switch r.Method {
	case "GET":
		tmpl.ExecuteTemplate(w, "wallets", data)
	case "POST":
		var payload addTxPayload
		r.ParseForm()
		payload.From = r.Form.Get("from")
		payload.To = r.Form.Get("to")
		payload.Amount, _ = strconv.Atoi(r.Form.Get("amount"))
		err := blockchain.Mempool.AddTx(payload.From, payload.To, payload.Amount)
		if err == nil {
			http.Redirect(w, r, "/mempool", http.StatusPermanentRedirect)
			fmt.Println("Make transaction", payload)
		}
		fmt.Fprintf(w, "err: %s", err)
	}
}

func mempool(w http.ResponseWriter, r *http.Request) {
	Tx := blockchain.Mempool.Txs
	data := pageData{PageTitle: "Mempool", PageHeader: "Memory Pool", MemTx: Tx}
	tmpl.ExecuteTemplate(w, "mempool", data)
}

func add(w http.ResponseWriter, r *http.Request) {
	data := pageData{PageTitle: "AddPage", PageHeader: "Add Block"}
	switch r.Method {
	case "GET":
		tmpl.ExecuteTemplate(w, "add", data)
	case "POST":
		r.ParseForm()
		name := r.Form.Get("blockName")
		blockchain.Blockchain().AddBlock(name)
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		fmt.Println("Add Block", name)
	}
}

func exit(w http.ResponseWriter, r *http.Request) {
	os.Exit(1)
}

// Start Server
func Start(aPort int) {
	port = fmt.Sprintf(":%d", aPort)
	tmpl = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))
	tmpl = template.Must(tmpl.ParseGlob(templateDir + "partials/*.gohtml"))

	htmlMux := mux.NewRouter()
	htmlMux.HandleFunc("/", home)
	htmlMux.HandleFunc("/wallets", wallets)
	htmlMux.HandleFunc("/mempool", mempool)
	htmlMux.HandleFunc("/add", add)
	htmlMux.HandleFunc("/exit", exit)

	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, htmlMux))
	fmt.Println("A")
}
