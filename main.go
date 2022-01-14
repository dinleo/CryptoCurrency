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

type homeData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

func home(rw http.ResponseWriter, rq *http.Request) {
	data := homeData{"Homepage", blockchain.GetBlockchain().AllBlocks()}
	tmpl.ExecuteTemplate(rw, "home", data)
}

func add(rw http.ResponseWriter, rq *http.Request) {
	tmpl.ExecuteTemplate(rw, "add", nil)
}

func main() {
	tmpl = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))
	tmpl = template.Must(tmpl.ParseGlob(templateDir + "partials/*.gohtml"))
	http.HandleFunc("/", home)
	http.HandleFunc("/add", add)
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
