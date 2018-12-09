package main

import (
	"net/http"
)

type FailService struct {
	HealthyFor   int
	HealthyIn    int
	UnHealthyFor int
}

func (p *FailService) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
}

func NewFailService() FailService {

	return FailService{}
}
