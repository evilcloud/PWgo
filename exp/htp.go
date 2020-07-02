package main

import (
	"fmt"
	"log"
	"net/http"
)

func httpCheck(dataString string) {
	fmt.Println("https://" + dataString)

	resp, err := http.Get(dataString)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(resp.StatusCode, http.StatusText(resp.StatusCode))
}

func main() {
	// website := "www.google.com"
	// valid := "kgro"
	// invalid := "34343253535345345gdgfffgdsfsdfdsf"

	httpCheck("www.google.com/")
	httpCheck("twitter.com/kgro")
	httpCheck("twitter.com/112234555fffdd")
}
