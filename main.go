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
	var serviceName = flag.String("service-name", "foo", "The name of the consumer service instance (this application instance). Defaults to foo.")
	var nameOfProvider = flag.String("provider", "", "The service_name of the provider (another instance of this application). Defaults to \"\".")
	var addrOfProvider = flag.String("provider-addr", "", "The address of the provider (another instance of this application). Defaults to \"\".")
	var addrOfConsul = flag.String("consul-server-addr", "", "The addr of the consul-server. Defaults to \"\". If not given the provider is searched using DNS.")
	flag.Parse()

	var consulClient Client

	if len(*addrOfConsul) > 0 {
		consul, err := NewConsulClient(*addrOfConsul)
		if err != nil {
			log.Println("Error unable to create consul at: ", *addrOfConsul)
		}
		consulClient = consul
	} else {
		log.Println("Service Discovery over consul is disabled.")
	}

	http.Handle("/ping", &PingService{Name: *serviceName, ProviderAddr: *addrOfProvider, ProviderName: *nameOfProvider, Version: version, ConsulClient: consulClient})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Path '/' is not implemented")
		http.Error(w, "Path '/' is not implemented", http.StatusInternalServerError)
	})

	//start the web server
	log.Printf("%s starts listening at %d.\n", *serviceName, *localPort)

	provider := *nameOfProvider
	if len(provider) == 0 {
		provider = *addrOfProvider
	}

	if len(provider) > 0 {
		log.Printf("The provider at %s is used.\n", provider)
	} else {
		log.Println("No provider is used.")
	}
	if err := http.ListenAndServe(":"+strconv.Itoa(*localPort), nil); err != nil {
		log.Fatal("ListenAndServer:", err)
	}

	log.Println("Exiting")
}
