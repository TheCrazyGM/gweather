package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const baseURL = "http://api.openweathermap.org/data/2.5/weather"

type WeatherData struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
}

func getWeather(apiKey, city, units string) (float64, string, string, error) {
	url := fmt.Sprintf("%s?q=%s&appid=%s&units=%s", baseURL, city, apiKey, units)
	resp, err := http.Get(url)
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
	if len(data.Weather) > 0 {
		status = data.Weather[0].Description
	}

	icon := ""
	if len(data.Weather) > 0 {
		icon = data.Weather[0].Icon
	}
	// Set the icon based on the weather condition
	var ICON string
	switch icon {
	case "01d":
		ICON = "☀️"
	case "02d":
		ICON = "🌤"
	case "03d", "04d":
		ICON = "☁️"
	case "09d":
		ICON = "🌧"
	case "10d":
		ICON = "🌦"
	case "11d":
		ICON = "⛈"
	case "13d":
		ICON = "❄️"
	case "50d":
		ICON = "🌫"
	case "01n", "02n":
		ICON = "🌜"
	case "03n", "04n":
		ICON = "☁️"
	case "09n":
		ICON = "🌧"
	case "10n":
		ICON = "🌦"
	case "11n":
		ICON = "⛈"
	case "13n":
		ICON = "❄️"
	case "50n":
		ICON = "🌫"
	default:
		ICON = "❓"
	}

	return data.Main.Temp, status, ICON, nil
}
