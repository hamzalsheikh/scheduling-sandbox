package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"sandbox/client"
)

func main() {

	client.RegisterHandlers()

	// choose a port randomly between 1024 to 49151
	port := fmt.Sprint(rand.Intn(49151-1024) + 1024)
	fmt.Printf("Client port is %v\n", port)
	var server http.Server
	server.Addr = ":" + port

	server.ListenAndServe()
}
