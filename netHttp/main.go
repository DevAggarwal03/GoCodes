package main

import (
	"fmt"
	"net/http"
	"strings"
)
func hello(w http.ResponseWriter, req *http.Request){
	fmt.Fprintf(w, "hello world");
}
func listHeaders(w http.ResponseWriter, req *http.Request){
	token := req.Header["Authorization"][0]
	fmt.Println(strings.Split(token, " ")[1]);
	fmt.Fprintf(w, "TOKEN: " + strings.Split(token, " ")[1])
}

func main () {
	fmt.Println("hit");
	http.HandleFunc("/Hello", hello)
	http.HandleFunc("/headers", listHeaders)


	http.ListenAndServe(":8080", nil);
}