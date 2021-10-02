package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w,"Hello World",r.URL.Path)
}
func main(){
	http.HandleFunc("/",handler)
	http.ListenAndServe(":8080",nil)
}