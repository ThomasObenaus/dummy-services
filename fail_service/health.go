package main

import (
	"log"
	"net/http"
	"time"
)

type FailService interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	Start()
	Stop()
}

type failServiceImpl struct {
	healthyFor     int64
	healthyIn      int64
	unHealthyFor   int64
	ticker         *time.Ticker
	healthy        bool
	changeStateAt  int64
	wasHealthyOnce bool
}

func (fs *failServiceImpl) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if fs.healthy {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusGatewayTimeout)
	}
}

func (fs *failServiceImpl) Stop() {
	fs.ticker.Stop()
}

func (fs *failServiceImpl) Start() {

	currentTime := time.Now().Unix()
	fs.changeStateAt = fs.nextEvalStateChange(currentTime)

	fs.ticker = time.NewTicker(time.Millisecond * 1000)
	go func() {
		for _ = range fs.ticker.C {

			currentTime := time.Now().Unix()
			if fs.isChangeState(currentTime) {
				fs.switchHealthy()
				fs.changeStateAt = fs.nextEvalStateChange(currentTime)

				log.Printf("State changed")
				log.Printf("Next state change at %s", time.Unix(fs.changeStateAt, 0).String())
			}

			state := "healthy"
			if !fs.healthy {
				state = "unhealthy"
			}
			log.Printf("State %s", state)
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

	// currently not healthy + were never healthy before ... initially get healthy
	if !fs.wasHealthyOnce {
		return currentTime + fs.healthyIn
	}

	// currently not healthy ... stay unhealthy for ...
	return currentTime + fs.unHealthyFor
}

func NewFailService(healthyIn int64, healthyFor int64, unHealthyFor int64) FailService {

	healthy := false
	// immedately start healthy
	if healthyIn == 0 {
		healthy = true
	}

	result := &failServiceImpl{
		healthyIn:      healthyIn,
		healthyFor:     healthyFor,
		unHealthyFor:   unHealthyFor,
		healthy:        healthy,
		wasHealthyOnce: healthy,
		changeStateAt:  0,
	}

	return result
}
