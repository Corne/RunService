package main

import (
	"net/http"
)

func postHandler(w http.ResponseWriter, r *http.Request) {
	runExample()
}

func main() {
	http.HandleFunc("/post/", postHandler)
	http.ListenAndServe(":8080", nil)
}
