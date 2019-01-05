package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/namsral/flag"
)

func main() {

	var localPort = flag.Int("p", 8080, "The port where the application instance listens to. Defaults to 8080.")
	var healthyFor = flag.Int("healthy-for", 0, "Number of seconds the health end-point will return a 200. A -1 will result in the service staying healthy forever (prio 1, if multiple flags are -1).")
	var healthyIn = flag.Int("healthy-in", 0, "Number of seconds the health end-point will start returning a 200. A -1 will result in the service NEVER getting healthy (prio 2, if multiple flags are -1).")
	var unhealthyFor = flag.Int("unhealthy-for", 0, "Number of seconds the health end-point will keep returning a !200. A -1 will result in the service staying unhealthy forever (prio 3, if multiple flags are -1).")
	flag.Parse()

	healthyForConverted := *healthyFor
	healthyForStr := strconv.Itoa(*healthyFor)
	if healthyForConverted < 0 {
		healthyForStr = "FOREVER (value <0 was specified)"
		healthyForConverted = 60 * 60 * 24 * 30 * 12 * 50 // 50 years unhealthy
	}

	healthyInConverted := *healthyIn
	healthyInStr := strconv.Itoa(*healthyIn)
	if healthyInConverted < 0 {
		healthyInStr = "NEVER (value <0 was specified)"
		healthyInConverted = 60 * 60 * 24 * 30 * 12 * 50 // 50 years unhealthy
	}

	unhealthyForConverted := *unhealthyFor
	unhealthyForStr := strconv.Itoa(*unhealthyFor)
	if unhealthyForConverted < 0 {
		unhealthyForStr = "FOREVER, but only if it got unhealthy once (value <0 was specified)"
		unhealthyForConverted = 60 * 60 * 24 * 30 * 12 * 50 // 50 years unhealthy
	}

	log.Println("Cfg:")
	log.Printf("\thealthyIn: %s", healthyInStr)
	log.Printf("\thealthyFor: %s", healthyForStr)
	log.Printf("\tunhealthyFor: %s", unhealthyForStr)

	failService := NewFailService(int64(healthyInConverted), int64(healthyForConverted), int64(unhealthyForConverted))
	http.HandleFunc("/health", failService.HealthEndpointHandler)
	http.HandleFunc("/sethealthy", failService.SetHealthyEndpointHandler)
	http.HandleFunc("/setunhealthy", failService.SetUnHealthyEndpointHandler)
	failService.Start()

	//start the web server
	log.Printf("Starts listening at %d.\n", *localPort)

	if err := http.ListenAndServe(":"+strconv.Itoa(*localPort), nil); err != nil {
		log.Fatal("ListenAndServer:", err)
	}

	failService.Stop()
	log.Println("Exiting")
}
