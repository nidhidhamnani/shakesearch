package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	indexedData, indexKeys, docs, err := GenerateIndexFromGivenData("completeworks.txt")
	if err != nil {
		log.Fatal(err)
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/search", HandleSearch(indexedData, indexKeys, docs))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	fmt.Printf("Listening on port %s...", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
