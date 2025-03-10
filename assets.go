package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Audience struct {
	Gender                string `json:"gender"`
	AgeFrom               int    `json:"ageFrom"`
	AgeTo                 int    `json:"ageTo"`
	Country               string `json:"country"`
	DailyHoursOnSocial    int    `json:"dailyHoursOnSocialMedia"`
	NumPurchasesLastMonth int    `json:"numPurchasesLastMonth"`
}

type Insight struct {
	Text string `json:"text"`
}

type Chart struct {
	Title string    `json:"title"`
	Axes  []string  `json:"axes"`
	Data  []float32 `json:"data"`
}

type Assets struct {
	Audiences []Audience `json:"audiences"`
	Insights  []Insight  `json:"insights"`
	Charts    []Chart    `json:"charts"`
}

// AssetsServer is an HTTP server for in memory assets.
type AssetsServer struct {
	assets Assets
	delay  time.Duration // for testing. A handler delays before answering a request
}

func (as *AssetsServer) insightsHandler(w http.ResponseWriter, r *http.Request) {
	genericHandler(as.assets.Insights, as.delay, w, r)
}

func (as *AssetsServer) chartsHandler(w http.ResponseWriter, r *http.Request) {
	genericHandler(as.assets.Charts, as.delay, w, r)
}

func (as *AssetsServer) audiencesHandler(w http.ResponseWriter, r *http.Request) {
	genericHandler(as.assets.Audiences, as.delay, w, r)
}

// genericHandler extract the slice index from the request path and responds with the item
func genericHandler[T any, S ~[]T](data S, delay time.Duration, w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("key"))
	if err != nil || id < 0 || id >= len(data) {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	log.Printf("request from %s for %T/%d", r.RemoteAddr, data[id], id)

	if err := json.NewEncoder(w).Encode(data[id]); err != nil {
		log.Printf("failed to decode %T %d: %v", data[id], id, err)
	}

	if delay > 0 {
		time.Sleep(delay)
	}
}
