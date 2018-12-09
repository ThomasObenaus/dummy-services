package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/namsral/flag"
)

const version = "v1"

func main() {

	var localPort = flag.Int("p", 8080, "The port where the application instance listens to. Defaults to 8080.")
	flag.Parse()

	http.HandleFunc("/health", healthHandler)

	//start the web server
	log.Printf("Starts listening at %d.\n", *localPort)

	if err := http.ListenAndServe(":"+strconv.Itoa(*localPort), nil); err != nil {
		log.Fatal("ListenAndServer:", err)
	}

	log.Println("Exiting")
}
