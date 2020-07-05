package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

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

	cookie, err := req.Cookie("acccount")

	if err == http.ErrNoCookie {
		cookie = &http.Cookie{
			Name:  "acccount",
			Value: "0",
		}
	}

	acccount, err := strconv.Atoi(cookie.Value)
	handleerror(rew, err)
	acccount++
	cookie.Value = strconv.Itoa(acccount)

	http.SetCookie(rew, cookie)
	fmt.Fprintln(rew, cookie.Value)
}

func handleerror(rew http.ResponseWriter, err error) bool {
	if err != nil {
		http.Error(rew, err.Error(), http.StatusInternalServerError)
		return true
	}

	return false
}
