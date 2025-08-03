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
		ID          int    `json:"id"`
	} `json:"weather"`
}

// getWeatherIcon returns an emoji representing the weather condition for a given OpenWeatherMap weather code.
func getWeatherIcon(code int) string {
	switch {
	case code >= 200 && code < 300:
		return "⛈️" // Thunderstorm
	case code >= 300 && code < 400:
		return "💧" // Drizzle
	case code >= 500 && code < 600:
		return "🌧️" // Rain
	case code >= 600 && code < 700:
		return "❄️" // Snow
	case code >= 700 && code < 800:
		return "🌫️" // Atmosphere
	case code == 800:
		return "☀️" // Clear
	case code > 800 && code < 900:
		return "☁️" // Clouds
	default:
		return ""
	}
}

// Make this variable so it can be mocked in tests
// getWeather retrieves the temperature and weather description for a city
// using the configured baseURL and httpClient.
var getWeather = func(apiKey, city, units string) (float64, string, string, error) {
	url := fmt.Sprintf("%s?q=%s&appid=%s&units=%s", baseURL, city, apiKey, units)
	resp, err := httpClient.Get(url)
	if err != nil {
		return 0, "", "", fmt.Errorf("failed to fetch weather: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, "", "", fmt.Errorf("api error: %s", resp.Status)
	}

	var data WeatherData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, "", "", fmt.Errorf("invalid response: %w", err)
	}

	status := ""
	icon := ""
	if len(data.Weather) > 0 {
		status = data.Weather[0].Description
		icon = getWeatherIcon(data.Weather[0].ID)
	}

	return data.Main.Temp, status, icon, nil
}
