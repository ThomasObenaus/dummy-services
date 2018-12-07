package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

type pingResponse struct {
	Message string `json:"message"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

type PingService struct {
	Name         string
	Version      string
	ProviderName string
	ProviderAddr string
	ConsulClient Client
}

const maxHop = 10

func (s *PingService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var message string

	hopStr := r.URL.Query().Get("hop")
	hop, err := strconv.Atoi(hopStr)
	if err != nil {
		hop = 0
	}
	hop++

	hasNoProvider := (len(s.ProviderAddr) == 0) && (len(s.ProviderName) == 0)

	if hasNoProvider || hop > maxHop {
		message = "(PONG)"
	} else {
		var err error
		addrOfProvider := s.ProviderAddr

		// do service-discovery using consul-http api
		if s.ConsulClient != nil {
			addr, err := s.ConsulClient.FindProvider(s.ProviderName)
			if err != nil {
				log.Println("[" + s.Name + "," + s.Version + "] Error: Failed to discover service ... " + err.Error())
			} else {
				addrOfProvider = addr
			}
		}

		message, err = s.getMessage(addrOfProvider, hop)
		if err != nil {
			log.Println("[" + s.Name + "," + s.Version + "] Error: " + err.Error())
		}
	}

	message = "/[" + s.Name + "," + s.Version + "]" + message
	response := &pingResponse{
		Message: message,
		Name:    s.Name,
		Version: s.Version,
	}

	log.Printf("[%s,%s] responds: %s\n", s.Name, s.Version, message)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *PingService) getMessage(providerAddr string, hop int) (string, error) {

	url := "http://" + providerAddr + "/ping?hop=" + strconv.Itoa(hop)
	log.Println("Call: ", url)

	client := &http.Client{
		Timeout: time.Second * 1,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "(BOING)", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "(BOING)", err
	}
	defer resp.Body.Close()

	response := &pingResponse{}
	err = json.NewDecoder(resp.Body).Decode(response)
	if err != nil {
		return "(BOING)", err
	}

	message := response.Message
	return message, nil
}
