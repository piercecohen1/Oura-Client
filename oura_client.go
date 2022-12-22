/*

Author: Pierce Cohen

I wrote this program to view my sleep data for any range of dates, using the Oura API.
It was originally written in Python, but I wanted to learn Go, so I rewrote it in Go.

In the future, I plan to run statistical analysis on the data.

*/

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// Create a struct that matches the structure of the JSON data in the response
type OuraData struct {
	Data []struct {
		Contributors struct {
			DeepSleep   int `json:"deep_sleep"`
			Efficiency  int `json:"efficiency"`
			Latency     int `json:"latency"`
			RemSleep    int `json:"rem_sleep"`
			Restfulness int `json:"restfulness"`
			Timing      int `json:"timing"`
			TotalSleep  int `json:"total_sleep"`
		} `json:"contributors"`
		Day       string `json:"day"`
		Score     int    `json:"score"`
		Timestamp string `json:"timestamp"`
	} `json:"data"`
	NextToken string `json:"next_token"`
}

// This function will process the HTTP GET response from the API
func processResponse(response []byte) {
	var ouraData OuraData

	// Unmarshal the JSON data into the OuraData struct
	err := json.Unmarshal(response, &ouraData)
	// If there is an error, panic
	if err != nil {
		panic(err)
	}

	/*
	   After unmarshalling the JSON data, access the individual fields of the data
	   using the dot notation on the variable
	   For example data.Score would give the score for the first day in the data
	*/
	for _, data := range ouraData.Data {
		// Convert data.Day to a time.Time object
		date, err := time.Parse("2006-01-02", data.Day)
		// If there is an error, panic
		if err != nil {
			panic(err)
		}

		//Print the date for each day in a more readable format
		fmt.Print("Day: ")
		fmt.Println(date.Format("Monday, January 2, 2006"))

		fmt.Print("Sleep Score: ")
		fmt.Println(data.Score)
		fmt.Println("Contributors: ")
		fmt.Print("\tDeep Sleep: ")
		fmt.Println(data.Contributors.DeepSleep)
		fmt.Print("\tEfficiency: ")
		fmt.Println(data.Contributors.Efficiency)
		fmt.Print("\tLatency: ")
		fmt.Println(data.Contributors.Latency)
		fmt.Print("\tREM Sleep: ")
		fmt.Println(data.Contributors.RemSleep)
		fmt.Print("\tRestfulness: ")
		fmt.Println(data.Contributors.Restfulness)
		fmt.Print("\tTiming: ")
		fmt.Println(data.Contributors.Timing)
		fmt.Print("\tTotal Sleep: ")
		fmt.Println(data.Contributors.TotalSleep)
		fmt.Println()
	}
}

func main() {
	// Get environment variables
	token := os.Getenv("OURA_TOKEN")

	// Set the request URL
	url := "https://api.ouraring.com/v2/usercollection/daily_sleep"

	// Specify the request parameters (the start and end dates)
	params := map[string]string{
		"start_date": "2022-12-01",
		"end_date":   "2022-12-10",
	}

	// Create a new http request
	req, err := http.NewRequest("GET", url, nil)
	// If there is an error, panic
	if err != nil {
		panic(err)
	}

	// Set the request headers
	req.Header.Set("Authorization", "Bearer "+token)

	// Create a new URL query
	q := req.URL.Query()

	// For each key, value pair in the params map, add the key, value pair to the query
	for key, value := range params {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

	// Create a new http client
	client := http.Client{}

	// Send the request
	resp, err := client.Do(req)

	// If there is an error, panic
	if err != nil {
		panic(err)
	}

	// Defer the closing of the response body
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)

	// If there is an error, panic
	if err != nil {
		panic(err)
	}

	// Process the response body
	processResponse(body)
}
