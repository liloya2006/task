package main

import (
	"log"
	"net/http"
	"test-task/api"
)

func main() {
	http.HandleFunc("/", api.HomePageHandler)
	http.HandleFunc("/add", api.AddUser)
	err := http.ListenAndServe(":8080", nil)
	log.Fatal(err)
}
