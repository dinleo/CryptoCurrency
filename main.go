package main

import (
	"CryptoCurrency/blockchain"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

const (
	port        string = ":4000"
	templateDir string = "templates/"
)

var tmpl *template.Template

type pageData struct {
	PageTitle  string
	PageHeader string
	Blocks     []*blockchain.Block
}

func home(rw http.ResponseWriter, rq *http.Request) {
	data := pageData{"HomePage", "HomePage", blockchain.GetBlockchain().AllBlocks()}
	tmpl.ExecuteTemplate(rw, "home", data)
}

func add(w http.ResponseWriter, r *http.Request) {
	data := pageData{"AddPage", "Add Block", nil}
	switch r.Method {
	case "GET":
		tmpl.ExecuteTemplate(w, "add", data)
	case "POST":
		r.ParseForm()
		blockData := r.Form.Get("blockData")
		blockchain.GetBlockchain().AddBlock(blockData)
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		fmt.Println("Add Block", blockData)
	}
}

func main() {
	tmpl = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))
	tmpl = template.Must(tmpl.ParseGlob(templateDir + "partials/*.gohtml"))

	http.HandleFunc("/", home)
	http.HandleFunc("/add", add)

	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))

}
