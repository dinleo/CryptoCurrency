package explorer

import (
	"CryptoCurrency/blockchain"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
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
	Blocks     []*blockchain.Block
	TxOuts     []*blockchain.UTxOut
}

// Handle func
// home
func home(w http.ResponseWriter, r *http.Request) {
	data := pageData{PageTitle: "HomePage", PageHeader: "HomePage", Blocks: blockchain.Blocks(blockchain.Blockchain())}
	tmpl.ExecuteTemplate(w, "home", data)
}

func exit(w http.ResponseWriter, r *http.Request) {
	os.Exit(1)
}

//add block
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

func wallet(w http.ResponseWriter, r *http.Request) {
	data := pageData{PageTitle: "Wallet", PageHeader: "Find Wallet"}
	switch r.Method {
	case "GET":
		tmpl.ExecuteTemplate(w, "wallet", data)
	case "POST":
		r.ParseForm()
		walletName := r.Form.Get("walletName")
		address := "/balance/" + walletName
		http.Redirect(w, r, address, http.StatusPermanentRedirect)
	}
}

func balance(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	address := vars["address"]
	data := pageData{PageTitle: "Balance", PageHeader: "View Balance", TxOuts: blockchain.UTxOutsByAddress(address, blockchain.Blockchain())}
	tmpl.ExecuteTemplate(w, "balance", data)
}

// Start Server
func Start(aPort int) {
	port = fmt.Sprintf(":%d", aPort)
	tmpl = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))
	tmpl = template.Must(tmpl.ParseGlob(templateDir + "partials/*.gohtml"))

	htmlMux := mux.NewRouter()
	htmlMux.HandleFunc("/", home)
	htmlMux.HandleFunc("/exit", exit)
	htmlMux.HandleFunc("/add", add)
	htmlMux.HandleFunc("/wallet", wallet)
	htmlMux.HandleFunc("/balance/{address}", balance)

	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, htmlMux))
	fmt.Println("A")
}
