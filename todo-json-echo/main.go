package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

func handle(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	log.Println(string(bodyBytes))
	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/", handle)
	http.ListenAndServe(":8000", nil)
}
