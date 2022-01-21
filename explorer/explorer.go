package explorer

import (
	"CryptoCurrency/blockchain"
	"fmt"
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
}

// Handle func
// home
func home(rw http.ResponseWriter, rq *http.Request) {
	data := pageData{"HomePage", "HomePage", blockchain.Blockchain().Blocks()}
	tmpl.ExecuteTemplate(rw, "home", data)
}

//add block
func add(w http.ResponseWriter, r *http.Request) {
	data := pageData{"AddPage", "Add Block", nil}
	switch r.Method {
	case "GET":
		tmpl.ExecuteTemplate(w, "add", data)
	case "POST":
		r.ParseForm()
		blockData := r.Form.Get("blockData")
		if blockData == "시스템 종료" {
			os.Exit(1)
		}
		blockchain.Blockchain().AddBlock(blockData)
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		fmt.Println("Add Block", blockData)
	}
}

// Start Server
func Start(aPort int) {
	port = fmt.Sprintf(":%d", aPort)
	tmpl = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))
	tmpl = template.Must(tmpl.ParseGlob(templateDir + "partials/*.gohtml"))

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/add", add)

	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, mux))
	fmt.Println("A")
}
