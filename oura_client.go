/*

Author: Pierce Cohen

I wrote this program to view my sleep data for any range of dates, using the Oura API.
It was originally written in Python, but I wanted to learn Go, so I rewrote it in Go.
Now transformed into a fully interactive web application!

*/

package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

//go:embed static/*
var staticFiles embed.FS

// SleepContributors represents the breakdown of sleep score factors
type SleepContributors struct {
	DeepSleep   int `json:"deep_sleep"`
	Efficiency  int `json:"efficiency"`
	Latency     int `json:"latency"`
	RemSleep    int `json:"rem_sleep"`
	Restfulness int `json:"restfulness"`
	Timing      int `json:"timing"`
	TotalSleep  int `json:"total_sleep"`
}

// SleepData represents a single day's sleep data
type SleepData struct {
	Contributors SleepContributors `json:"contributors"`
	Day          string            `json:"day"`
	Score        int               `json:"score"`
	Timestamp    string            `json:"timestamp"`
}

// OuraResponse represents the API response from Oura
type OuraResponse struct {
	Data      []SleepData `json:"data"`
	NextToken string      `json:"next_token,omitempty"`
}

// APIResponse is the response sent to the frontend
type APIResponse struct {
	Success bool        `json:"success"`
	Data    []SleepData `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// fetchOuraData fetches sleep data from the Oura API
func fetchOuraData(startDate, endDate string) ([]SleepData, error) {
	token := os.Getenv("OURA_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("OURA_TOKEN environment variable is not set")
	}

	url := "https://api.ouraring.com/v2/usercollection/daily_sleep"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	q := req.URL.Query()
	q.Add("start_date", startDate)
	q.Add("end_date", endDate)
	req.URL.RawQuery = q.Encode()

	client := http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var ouraData OuraResponse
	err = json.Unmarshal(body, &ouraData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return ouraData.Data, nil
}

// handleSleepData handles API requests for sleep data
func handleSleepData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{
			Success: false,
			Error:   "Method not allowed",
		})
		return
	}

	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	if startDate == "" || endDate == "" {
		// Default to last 30 days
		endDate = time.Now().Format("2006-01-02")
		startDate = time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	}

	// Validate date format
	_, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{
			Success: false,
			Error:   "Invalid start_date format. Use YYYY-MM-DD",
		})
		return
	}

	_, err = time.Parse("2006-01-02", endDate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{
			Success: false,
			Error:   "Invalid end_date format. Use YYYY-MM-DD",
		})
		return
	}

	data, err := fetchOuraData(startDate, endDate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(APIResponse{
		Success: true,
		Data:    data,
	})
}

// handleHealth handles health check requests
func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func main() {
	// API routes
	http.HandleFunc("/api/sleep", handleSleepData)
	http.HandleFunc("/api/health", handleHealth)

	// Serve static files from embedded filesystem
	staticHandler := http.FileServer(http.FS(staticFiles))
	http.Handle("/static/", staticHandler)

	// Serve index.html for root path
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		content, err := staticFiles.ReadFile("static/index.html")
		if err != nil {
			http.Error(w, "Could not load page", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write(content)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🌙 Oura Sleep Dashboard starting on http://localhost:%s", port)
	log.Printf("📊 Make sure OURA_TOKEN environment variable is set")

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
