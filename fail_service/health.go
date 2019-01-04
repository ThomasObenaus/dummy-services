package main

import (
	"log"
	"net/http"
	"time"
)

// FailService is a simple interface for a service that implements 3 endpoints
// The /health EP responds with the health status of the service
// The /sethealthy EP switches the health state of the service to healthy
// The /setunhealthy EP switches the health state of the service to unhealthy
type FailService interface {
	HealthEndpointHandler(w http.ResponseWriter, r *http.Request)
	SetHealthyEndpointHandler(w http.ResponseWriter, r *http.Request)
	SetUnHealthyEndpointHandler(w http.ResponseWriter, r *http.Request)
	Start()
	Stop()
}

var errorResponseWrongMethod = []byte("{ \"error\": \"Invalid mehtod used. You have to use the PUT mehtod.\" }")

type failServiceImpl struct {
	healthyFor            int64
	healthyIn             int64
	unHealthyFor          int64
	ticker                *time.Ticker
	healthy               bool
	changeStateAt         int64
	wasHealthyOnce        bool
	overwrittenByEndpoint bool
}

func validateHTTPMethod(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodPut {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorResponseWrongMethod)
		return false
	}

	return true
}

func (fs *failServiceImpl) SetHealthyEndpointHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("SetHealthyEndpointHandler called")

	if validateHTTPMethod(w, r) {
		fs.overwrittenByEndpoint = true
		fs.healthy = true
	}
}

func (fs *failServiceImpl) SetUnHealthyEndpointHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("SetUnHealthyEndpointHandler called")

	if validateHTTPMethod(w, r) {
		fs.overwrittenByEndpoint = true
		fs.healthy = false
	}
}

func (fs *failServiceImpl) HealthEndpointHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("HealthEndpointHandler called")

	if fs.healthy {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusGatewayTimeout)
	}
}

func (fs *failServiceImpl) Stop() {
	fs.ticker.Stop()
}

func stateToStr(healthy bool) string {
	state := "healthy"
	if !healthy {
		state = "unhealthy"
	}
	return state
}

func (fs *failServiceImpl) Start() {

	currentTime := time.Now().Unix()
	fs.changeStateAt = fs.nextEvalStateChange(currentTime)

	fs.ticker = time.NewTicker(time.Millisecond * 1000)
	go func() {
		for range fs.ticker.C {

			if !fs.overwrittenByEndpoint {
				currentTime := time.Now().Unix()
				if fs.isChangeState(currentTime) {
					fs.switchHealthy()
					fs.changeStateAt = fs.nextEvalStateChange(currentTime)

					log.Printf("State changed")
					log.Printf("Next state change at %s", time.Unix(fs.changeStateAt, 0).String())
				}
			}

			overwrittenByEndpointStr := ""
			if fs.overwrittenByEndpoint {
				overwrittenByEndpointStr = " - Was set and thus fixed by endpoint."
			}
			log.Printf("State %s %s", stateToStr(fs.healthy), overwrittenByEndpointStr)
		}
	}()
}

func (fs *failServiceImpl) switchHealthy() {
	if !fs.wasHealthyOnce && fs.healthy {
		fs.wasHealthyOnce = true
	}
	fs.healthy = !fs.healthy
}

func (fs *failServiceImpl) isChangeState(currentTime int64) bool {

	if fs.healthy && fs.healthyFor == 0 {
		return false
	}

	if currentTime > fs.changeStateAt {
		return true
	}
	return false
}

func (fs *failServiceImpl) nextEvalStateChange(currentTime int64) int64 {

	// currently healthy ... stay healthy for ...
	if fs.healthy {
		return currentTime + fs.healthyFor
	}

	// in case healthyIn is -1 or negative at all
	healthyIn := fs.healthyIn
	if healthyIn < 0 {
		healthyIn = 0
	}

	// currently not healthy + were never healthy before ... initially get healthy
	if !fs.wasHealthyOnce {
		return currentTime + healthyIn
	}

	// currently not healthy ... stay unhealthy for ...
	return currentTime + fs.unHealthyFor
}

// NewFailService creates a new instance of a FailService implementation
func NewFailService(healthyIn int64, healthyFor int64, unHealthyFor int64) FailService {

	healthy := false
	// immediately start healthy
	if healthyIn == 0 {
		healthy = true
	}

	result := &failServiceImpl{
		healthyIn:             healthyIn,
		healthyFor:            healthyFor,
		unHealthyFor:          unHealthyFor,
		healthy:               healthy,
		wasHealthyOnce:        healthy,
		changeStateAt:         0,
		overwrittenByEndpoint: false,
	}

	return result
}
