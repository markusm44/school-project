package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	setUpHandlers()

	fmt.Println("Serving frontend at: http://localhost:8081")
	listenServeErr := (http.ListenAndServe(":8081", nil))
	if listenServeErr != nil {
		log.Fatal("ListenAndServe: ", listenServeErr)
	}
}

func setUpHandlers() {
	//html handler
	http.Handle("/", http.FileServer(http.Dir("../frontend/dist/")))
}
