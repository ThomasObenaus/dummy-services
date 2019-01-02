package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/namsral/flag"
)

func main() {

	var localPort = flag.Int("p", 8080, "The port where the application instance listens to. Defaults to 8080.")
	var healthyFor = flag.Int("healthy-for", 0, "Number of seconds the health end-point will return a 200. 0 means forever.")
	var healthyIn = flag.Int("healthy-in", 0, "Number of seconds the health end-point will start returning a 200.")
	var unhealthyFor = flag.Int("unhealthy-for", 0, "Number of seconds the health end-point will keep returning a !200.")
	var restartsToGetHealthy = flag.Int("restarts-to-get-healthy", 0, "The number of restarts the service waits until he starts with the procedure and gets healthy after 'healthy-in'. For this feature consul KV is needed as persistence component. The location of consul can be specified using 'consul-server-addr'.")
	var addrOfConsul = flag.String("consul-server-addr", "", "The addr of the consul-server. Defaults to \"\". If not given the provider is searched using DNS.")
	flag.Parse()

	log.Println("Cfg:")
	log.Printf("\thealthyIn: %d", *healthyIn)
	log.Printf("\thealthyFor: %d", *healthyFor)
	log.Printf("\tunhealthyFor: %d", *unhealthyFor)
	log.Printf("\trestartsToGetHealthy: %d", *restartsToGetHealthy)

	addrOfConsulMsg := "N/A"
	if len(*addrOfConsul) > 0 {
		addrOfConsulMsg = *addrOfConsul
	}

	log.Printf("\taddrOfConsul: %s", addrOfConsulMsg)

	failService := NewFailService(int64(*healthyIn), int64(*healthyFor), int64(*unhealthyFor))
	http.Handle("/health", failService)
	failService.Start()

	//start the web server
	log.Printf("Starts listening at %d.\n", *localPort)

	if err := http.ListenAndServe(":"+strconv.Itoa(*localPort), nil); err != nil {
		log.Fatal("ListenAndServer:", err)
	}

	failService.Stop()
	log.Println("Exiting")
}
