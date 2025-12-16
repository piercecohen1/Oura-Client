/*

Author: Pierce Cohen

I wrote this program to view my sleep data for any range of dates, using the Oura API.
It was originally written in Python, but I wanted to learn Go, so I rewrote it in Go.
Now transformed into a fully interactive web application using the Oura API v2!

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

const baseURL = "https://api.ouraring.com/v2/usercollection"

// ============================================
// Data Models for Oura API v2
// ============================================

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

// DailySleep represents a single day's sleep score data
type DailySleep struct {
	ID           string            `json:"id"`
	Contributors SleepContributors `json:"contributors"`
	Day          string            `json:"day"`
	Score        int               `json:"score"`
	Timestamp    string            `json:"timestamp"`
}

// SleepPeriod represents detailed sleep session data
type SleepPeriod struct {
	ID                   string  `json:"id"`
	AverageBreath        float64 `json:"average_breath"`
	AverageHeartRate     float64 `json:"average_heart_rate"`
	AverageHRV           int     `json:"average_hrv"`
	AwakeTime            int     `json:"awake_time"`
	BedtimeEnd           string  `json:"bedtime_end"`
	BedtimeStart         string  `json:"bedtime_start"`
	Day                  string  `json:"day"`
	DeepSleepDuration    int     `json:"deep_sleep_duration"`
	Efficiency           int     `json:"efficiency"`
	Latency              int     `json:"latency"`
	LightSleepDuration   int     `json:"light_sleep_duration"`
	LowBatteryAlert      bool    `json:"low_battery_alert"`
	LowestHeartRate      int     `json:"lowest_heart_rate"`
	REMSleepDuration     int     `json:"rem_sleep_duration"`
	Restless             int     `json:"restless_periods"`
	TimeInBed            int     `json:"time_in_bed"`
	TotalSleepDuration   int     `json:"total_sleep_duration"`
	Type                 string  `json:"type"`
}

// ActivityContributors represents activity score contributors
type ActivityContributors struct {
	MeetDailyTargets  int `json:"meet_daily_targets"`
	MoveEveryHour     int `json:"move_every_hour"`
	RecoveryTime      int `json:"recovery_time"`
	StayActive        int `json:"stay_active"`
	TrainingFrequency int `json:"training_frequency"`
	TrainingVolume    int `json:"training_volume"`
}

// DailyActivity represents daily activity data
type DailyActivity struct {
	ID                       string               `json:"id"`
	Class5Min                string               `json:"class_5_min"`
	Score                    int                  `json:"score"`
	ActiveCalories           int                  `json:"active_calories"`
	AverageMetMinutes        float64              `json:"average_met_minutes"`
	Contributors             ActivityContributors `json:"contributors"`
	EquivalentWalkingDistance int                 `json:"equivalent_walking_distance"`
	HighActivityMetMinutes   int                  `json:"high_activity_met_minutes"`
	HighActivityTime         int                  `json:"high_activity_time"`
	InactivityAlerts         int                  `json:"inactivity_alerts"`
	LowActivityMetMinutes    int                  `json:"low_activity_met_minutes"`
	LowActivityTime          int                  `json:"low_activity_time"`
	MediumActivityMetMinutes int                  `json:"medium_activity_met_minutes"`
	MediumActivityTime       int                  `json:"medium_activity_time"`
	MetersToTarget           int                  `json:"meters_to_target"`
	NonWearTime              int                  `json:"non_wear_time"`
	RestingTime              int                  `json:"resting_time"`
	SedentaryMetMinutes      int                  `json:"sedentary_met_minutes"`
	SedentaryTime            int                  `json:"sedentary_time"`
	Steps                    int                  `json:"steps"`
	TargetCalories           int                  `json:"target_calories"`
	TargetMeters             int                  `json:"target_meters"`
	TotalCalories            int                  `json:"total_calories"`
	Day                      string               `json:"day"`
	Timestamp                string               `json:"timestamp"`
}

// ReadinessContributors represents readiness score contributors
type ReadinessContributors struct {
	ActivityBalance     int `json:"activity_balance"`
	BodyTemperature     int `json:"body_temperature"`
	HRVBalance          int `json:"hrv_balance"`
	PreviousDayActivity int `json:"previous_day_activity"`
	PreviousNight       int `json:"previous_night"`
	RecoveryIndex       int `json:"recovery_index"`
	RestingHeartRate    int `json:"resting_heart_rate"`
	SleepBalance        int `json:"sleep_balance"`
}

// DailyReadiness represents daily readiness data
type DailyReadiness struct {
	ID                   string                `json:"id"`
	Contributors         ReadinessContributors `json:"contributors"`
	Day                  string                `json:"day"`
	Score                int                   `json:"score"`
	TemperatureDeviation float64               `json:"temperature_deviation"`
	TemperatureTrendDeviation float64          `json:"temperature_trend_deviation"`
	Timestamp            string                `json:"timestamp"`
}

// HeartRate represents heart rate data point
type HeartRate struct {
	Bpm       int    `json:"bpm"`
	Source    string `json:"source"`
	Timestamp string `json:"timestamp"`
}

// PersonalInfo represents user personal information
type PersonalInfo struct {
	ID          string  `json:"id"`
	Age         int     `json:"age"`
	Weight      float64 `json:"weight"`
	Height      float64 `json:"height"`
	BiologicalSex string `json:"biological_sex"`
	Email       string  `json:"email"`
}

// DailyStress represents daily stress data
type DailyStress struct {
	ID              string `json:"id"`
	Day             string `json:"day"`
	StressHigh      int    `json:"stress_high"`
	RecoveryHigh    int    `json:"recovery_high"`
	DaySummary      string `json:"day_summary"`
}

// ============================================
// API Response Wrappers
// ============================================

type DailySleepResponse struct {
	Data      []DailySleep `json:"data"`
	NextToken string       `json:"next_token,omitempty"`
}

type SleepPeriodResponse struct {
	Data      []SleepPeriod `json:"data"`
	NextToken string        `json:"next_token,omitempty"`
}

type DailyActivityResponse struct {
	Data      []DailyActivity `json:"data"`
	NextToken string          `json:"next_token,omitempty"`
}

type DailyReadinessResponse struct {
	Data      []DailyReadiness `json:"data"`
	NextToken string           `json:"next_token,omitempty"`
}

type HeartRateResponse struct {
	Data      []HeartRate `json:"data"`
	NextToken string      `json:"next_token,omitempty"`
}

type DailyStressResponse struct {
	Data      []DailyStress `json:"data"`
	NextToken string        `json:"next_token,omitempty"`
}

// APIResponse is the generic response sent to the frontend
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// DashboardData combines all data for the dashboard
type DashboardData struct {
	Sleep      []DailySleep     `json:"sleep"`
	SleepPeriods []SleepPeriod  `json:"sleep_periods"`
	Activity   []DailyActivity  `json:"activity"`
	Readiness  []DailyReadiness `json:"readiness"`
	HeartRate  []HeartRate      `json:"heart_rate"`
	Stress     []DailyStress    `json:"stress"`
}

// ============================================
// HTTP Client Functions
// ============================================

func getToken() (string, error) {
	token := os.Getenv("OURA_TOKEN")
	if token == "" {
		return "", fmt.Errorf("OURA_TOKEN environment variable is not set")
	}
	return token, nil
}

func fetchFromOura(endpoint, startDate, endDate string) ([]byte, error) {
	token, err := getToken()
	if err != nil {
		return nil, err
	}

	url := baseURL + "/" + endpoint

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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// ============================================
// Data Fetch Functions
// ============================================

func fetchDailySleep(startDate, endDate string) ([]DailySleep, error) {
	body, err := fetchFromOura("daily_sleep", startDate, endDate)
	if err != nil {
		return nil, err
	}

	var response DailySleepResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse sleep response: %v", err)
	}

	return response.Data, nil
}

func fetchSleepPeriods(startDate, endDate string) ([]SleepPeriod, error) {
	body, err := fetchFromOura("sleep", startDate, endDate)
	if err != nil {
		return nil, err
	}

	var response SleepPeriodResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse sleep periods response: %v", err)
	}

	return response.Data, nil
}

func fetchDailyActivity(startDate, endDate string) ([]DailyActivity, error) {
	body, err := fetchFromOura("daily_activity", startDate, endDate)
	if err != nil {
		return nil, err
	}

	var response DailyActivityResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse activity response: %v", err)
	}

	return response.Data, nil
}

func fetchDailyReadiness(startDate, endDate string) ([]DailyReadiness, error) {
	body, err := fetchFromOura("daily_readiness", startDate, endDate)
	if err != nil {
		return nil, err
	}

	var response DailyReadinessResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse readiness response: %v", err)
	}

	return response.Data, nil
}

func fetchHeartRate(startDate, endDate string) ([]HeartRate, error) {
	body, err := fetchFromOura("heartrate", startDate, endDate)
	if err != nil {
		return nil, err
	}

	var response HeartRateResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse heart rate response: %v", err)
	}

	return response.Data, nil
}

func fetchDailyStress(startDate, endDate string) ([]DailyStress, error) {
	body, err := fetchFromOura("daily_stress", startDate, endDate)
	if err != nil {
		return nil, err
	}

	var response DailyStressResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse stress response: %v", err)
	}

	return response.Data, nil
}

func fetchPersonalInfo() (*PersonalInfo, error) {
	token, err := getToken()
	if err != nil {
		return nil, err
	}

	url := baseURL + "/personal_info"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var info PersonalInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, fmt.Errorf("failed to parse personal info response: %v", err)
	}

	return &info, nil
}

// ============================================
// HTTP Handlers
// ============================================

func parseDateParams(r *http.Request) (string, string) {
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	if startDate == "" || endDate == "" {
		endDate = time.Now().Format("2006-01-02")
		startDate = time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	}

	return startDate, endDate
}

func validateDates(startDate, endDate string) error {
	_, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return fmt.Errorf("invalid start_date format. Use YYYY-MM-DD")
	}

	_, err = time.Parse("2006-01-02", endDate)
	if err != nil {
		return fmt.Errorf("invalid end_date format. Use YYYY-MM-DD")
	}

	return nil
}

func handleDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Success: false, Error: "Method not allowed"})
		return
	}

	startDate, endDate := parseDateParams(r)
	if err := validateDates(startDate, endDate); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{Success: false, Error: err.Error()})
		return
	}

	dashboard := DashboardData{}
	var fetchErr error

	// Fetch all data types (continue even if some fail)
	if data, err := fetchDailySleep(startDate, endDate); err == nil {
		dashboard.Sleep = data
	} else {
		log.Printf("Warning: failed to fetch sleep data: %v", err)
		fetchErr = err
	}

	if data, err := fetchSleepPeriods(startDate, endDate); err == nil {
		dashboard.SleepPeriods = data
	} else {
		log.Printf("Warning: failed to fetch sleep periods: %v", err)
	}

	if data, err := fetchDailyActivity(startDate, endDate); err == nil {
		dashboard.Activity = data
	} else {
		log.Printf("Warning: failed to fetch activity data: %v", err)
	}

	if data, err := fetchDailyReadiness(startDate, endDate); err == nil {
		dashboard.Readiness = data
	} else {
		log.Printf("Warning: failed to fetch readiness data: %v", err)
	}

	if data, err := fetchHeartRate(startDate, endDate); err == nil {
		dashboard.HeartRate = data
	} else {
		log.Printf("Warning: failed to fetch heart rate data: %v", err)
	}

	if data, err := fetchDailyStress(startDate, endDate); err == nil {
		dashboard.Stress = data
	} else {
		log.Printf("Warning: failed to fetch stress data: %v", err)
	}

	// If all fetches failed, return error
	if dashboard.Sleep == nil && dashboard.Activity == nil && dashboard.Readiness == nil && fetchErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{Success: false, Error: fetchErr.Error()})
		return
	}

	json.NewEncoder(w).Encode(APIResponse{Success: true, Data: dashboard})
}

func handleSleep(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	startDate, endDate := parseDateParams(r)
	if err := validateDates(startDate, endDate); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{Success: false, Error: err.Error()})
		return
	}

	data, err := fetchDailySleep(startDate, endDate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{Success: false, Error: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(APIResponse{Success: true, Data: data})
}

func handleSleepPeriods(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	startDate, endDate := parseDateParams(r)
	if err := validateDates(startDate, endDate); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{Success: false, Error: err.Error()})
		return
	}

	data, err := fetchSleepPeriods(startDate, endDate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{Success: false, Error: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(APIResponse{Success: true, Data: data})
}

func handleActivity(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	startDate, endDate := parseDateParams(r)
	if err := validateDates(startDate, endDate); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{Success: false, Error: err.Error()})
		return
	}

	data, err := fetchDailyActivity(startDate, endDate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{Success: false, Error: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(APIResponse{Success: true, Data: data})
}

func handleReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	startDate, endDate := parseDateParams(r)
	if err := validateDates(startDate, endDate); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{Success: false, Error: err.Error()})
		return
	}

	data, err := fetchDailyReadiness(startDate, endDate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{Success: false, Error: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(APIResponse{Success: true, Data: data})
}

func handleHeartRate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	startDate, endDate := parseDateParams(r)
	if err := validateDates(startDate, endDate); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{Success: false, Error: err.Error()})
		return
	}

	data, err := fetchHeartRate(startDate, endDate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{Success: false, Error: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(APIResponse{Success: true, Data: data})
}

func handlePersonalInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data, err := fetchPersonalInfo()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{Success: false, Error: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(APIResponse{Success: true, Data: data})
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "api_version": "v2"})
}

func main() {
	// API v2 routes
	http.HandleFunc("/api/dashboard", handleDashboard)
	http.HandleFunc("/api/sleep", handleSleep)
	http.HandleFunc("/api/sleep/periods", handleSleepPeriods)
	http.HandleFunc("/api/activity", handleActivity)
	http.HandleFunc("/api/readiness", handleReadiness)
	http.HandleFunc("/api/heartrate", handleHeartRate)
	http.HandleFunc("/api/personal", handlePersonalInfo)
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

	log.Printf("Oura Sleep Dashboard (API v2) starting on http://localhost:%s", port)
	log.Printf("Make sure OURA_TOKEN environment variable is set")

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
