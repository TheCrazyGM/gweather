package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetWeather(t *testing.T) {
	// Mock server to simulate the OpenWeather API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		
		// Check for required parameters
		if q.Get("q") == "" || q.Get("appid") == "" || q.Get("units") == "" {
			http.Error(w, "Missing required parameters", http.StatusBadRequest)
			return
		}

		// Return different responses based on city parameter
		switch q.Get("q") {
		case "London":
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{
				"main": {"temp": 15.5},
				"weather": [{"description": "cloudy"}]
			}`))
		case "InvalidCity":
			http.Error(w, "City not found", http.StatusNotFound)
		default:
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{
				"main": {"temp": 25.0},
				"weather": [{"description": "sunny"}]
			}`))
		}
	}))
	defer server.Close()

	// Save the original base URL and restore it after the test
	originalBaseURL := baseURL
	baseURL = server.URL
	defer func() { baseURL = originalBaseURL }()

	tests := []struct {
		name        string
		city        string
		units       string
		wantTemp    float64
		wantStatus  string
		wantErr     bool
		description string
	}{
		{
			name:        "Valid city London",
			city:        "London",
			units:       "metric",
			wantTemp:    15.5,
			wantStatus:  "cloudy",
			wantErr:     false,
			description: "Should return temperature and weather status for London",
		},
		{
			name:        "Invalid city",
			city:        "InvalidCity",
			units:       "metric",
			wantTemp:    0,
			wantStatus:  "",
			wantErr:     true,
			description: "Should return an error for a non-existent city",
		},
		{
			name:        "Default city",
			city:        "DefaultCity",
			units:       "imperial",
			wantTemp:    25.0,
			wantStatus:  "sunny",
			wantErr:     false,
			description: "Should return temperature and weather status for default case",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTemp, gotStatus, err := getWeather("test-api-key", tt.city, tt.units)
			
			// Check error expectations
			if (err != nil) != tt.wantErr {
				t.Errorf("getWeather() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			// Skip further checks if we expected an error
			if tt.wantErr {
				return
			}
			
			// Check temperature
			if gotTemp != tt.wantTemp {
				t.Errorf("getWeather() gotTemp = %v, want %v", gotTemp, tt.wantTemp)
			}
			
			// Check weather status
			if gotStatus != tt.wantStatus {
				t.Errorf("getWeather() gotStatus = %v, want %v", gotStatus, tt.wantStatus)
			}
		})
	}
}

func TestGetWeatherNetworkError(t *testing.T) {
	// Set baseURL to a non-existent server to simulate network error
	originalBaseURL := baseURL
	baseURL = "http://non-existent-server.example"
	defer func() { baseURL = originalBaseURL }()

	_, _, err := getWeather("test-api-key", "London", "metric")
	if err == nil {
		t.Error("getWeather() expected error for network failure, got nil")
	}
}