package main

import (
	"net/http"

	"google.golang.org/appengine"
)

func main() {
	http.HandleFunc("/", Redirect)
	appengine.Main()
}
