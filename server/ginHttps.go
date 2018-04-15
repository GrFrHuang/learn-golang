package main

import (
	"net/http"
	"fmt"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w,
		"Hi, GrFrHuang!")
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServeTLS("www.bloghuang.com:8088",
		"server.crt", "serverKey.pem", nil)
}
