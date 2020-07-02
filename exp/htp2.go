package main

import (
	"fmt"
	"log"
	"net/http"
)

func thp(stringData string) {
	fmt.Println(stringData)
	resp, err := http.Get("https://" + stringData)
	if err != nil {
		log.Fatal(err)
	}

	// Print the HTTP Status Code and Status Name
	fmt.Println("HTTP Response Status:", resp.StatusCode, http.StatusText(resp.StatusCode))

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		fmt.Println("HTTP Status is in the 2xx range")
	} else {
		fmt.Println("Argh! Broken")
	}
}

func main() {
	thp("golangcode.com")
	thp("twitter.com")
	thp("twitter.com/kgro")
	thp("twitter.com/fjdlfdjflklr54er4r4ef")
}
