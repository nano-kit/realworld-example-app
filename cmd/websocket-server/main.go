package main

import "net/http"

func main() {
	panic(http.ListenAndServe(":9100", http.FileServer(http.Dir("./html"))))
}
