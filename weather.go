package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Make this a variable instead of constant for testing purposes
var baseURL = "http://api.openweathermap.org/data/2.5/weather"

// httpClient allows setting a timeout and is overridable in tests
var httpClient = &http.Client{Timeout: 10 * time.Second}

type WeatherData struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
}

// Make this variable so it can be mocked in tests
// getWeather retrieves the temperature and weather description for a city
// using the configured baseURL and httpClient.
var getWeather = func(apiKey, city, units string) (float64, string, error) {
	url := fmt.Sprintf("%s?q=%s&appid=%s&units=%s", baseURL, city, apiKey, units)
	resp, err := httpClient.Get(url)
	if err != nil {
		return 0, "", fmt.Errorf("failed to fetch weather: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, "", fmt.Errorf("api error: %s", resp.Status)
	}

	var data WeatherData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, "", fmt.Errorf("invalid response: %w", err)
	}

	status := ""
	if len(data.Weather) > 0 {
		status = data.Weather[0].Description
	}

	return data.Main.Temp, status, nil
}
