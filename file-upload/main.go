package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

var tmp *template.Template

func init() {
	tmp = template.Must(template.ParseFiles("index.gohtml"))
}

func main() {
	var port string
	var arg string
	if len(os.Args) > 1 {
		arg = os.Args[1]
	}

	if arg != "" {
		port = arg
	} else {
		port = ":8080"
	}

	fmt.Println("Server started at port ", port)

	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.HandleFunc("/", defaultpagehandler)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalln(err)
	}
}

func defaultpagehandler(rew http.ResponseWriter, req *http.Request) {

	var s string

	if req.Method == http.MethodPost {
		f, fh, err := req.FormFile("ufile")

		if handleerror(rew, err) {
			return
		}
		defer f.Close()

		bs, err := ioutil.ReadAll(f)

		if handleerror(rew, err) {
			return
		}

		ptf, err := os.Create(filepath.Join("./uploads/", fh.Filename))

		if handleerror(rew, err) {
			return
		}

		defer ptf.Close()

		_, err = ptf.Write(bs)

		if handleerror(rew, err) {
			return
		}

		s = string(bs)
	}

	rew.Header().Set("Content-Type", "text/html; charset ")
	tmp.Execute(rew, s)
}

func handleerror(rew http.ResponseWriter, err error) bool {
	if err != nil {
		http.Error(rew, err.Error(), http.StatusInternalServerError)
		return true
	}

	return false
}
