package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/namsral/flag"
)

func main() {

	var localPort = flag.Int("p", 8080, "The port where the application instance listens to. Defaults to 8080.")
	var healtyFor = flag.Int("healthy-for", 0, "Number of seconds the health end-point will return a 200. 0 means forever.")
	//var healtyIn = flag.Int("healthy-in", 0, "Number of seconds the health end-point will return a 200. 0 means forever.")
	//var unhealtyFor = flag.Int("unhealthy-for", 0, "Number of seconds the health end-point will return a 200. 0 means forever.")
	flag.Parse()

	log.Println("Cfg:")
	log.Printf("\thealtyFor: %d", *healtyFor)

	failService := NewFailService()

	http.Handle("/health", &failService)

	//start the web server
	log.Printf("Starts listening at %d.\n", *localPort)

	if err := http.ListenAndServe(":"+strconv.Itoa(*localPort), nil); err != nil {
		log.Fatal("ListenAndServer:", err)
	}

	log.Println("Exiting")
}
