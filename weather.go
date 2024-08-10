package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	baseURL = "http://api.openweathermap.org/data/2.5/weather"
)

var (
	ErrMissingAPIKey = fmt.Errorf("API key not found")
	ErrMissingCity   = fmt.Errorf("please provide a city")
)

type (
	City        string
	Temperature float64
	Status      string
	WeatherData struct {
		Cod     int           `json:"cod"`
		Main    MainData      `json:"main"`
		Weather []WeatherInfo `json:"weather"`
	}
	MainData struct {
		Temp float64 `json:"temp"`
	}
	WeatherInfo struct {
		MainStatus Status `json:"main"`
	}
)

func getWeatherData(apiKey, city, units string) (*WeatherData, error) {
	url := fmt.Sprintf("%s?q=%s&appid=%s&units=%s", baseURL, city, apiKey, units)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weather data: %w", err)
	}
	defer resp.Body.Close()

	var weatherData WeatherData
	dec := json.NewDecoder(resp.Body)

	for {
		if err := dec.Decode(&weatherData); err == io.EOF {
			break
		} else if err != nil {
			return nil, fmt.Errorf("failed to decode weather data: %w", err)
		}
	}

	if weatherData.Cod != 200 {
		return nil, fmt.Errorf("failed to retrieve weather data: invalid response code")
	}

	return &weatherData, nil
}

func getCurrentTemperature(weatherData *WeatherData) (Temperature, error) {
	temp := weatherData.Main.Temp
	return Temperature(temp), nil
}

func getWeatherStatus(weatherData *WeatherData) (Status, error) {
	status := ""
	for _, weather := range weatherData.Weather {
		status = string(weather.MainStatus)
	}

	if status == "" {
		return "", fmt.Errorf("failed to retrieve status: invalid status data")
	}

	return Status(status), nil
}

func displayWeather(temperature Temperature, weatherStatus Status) {
	fmt.Printf("%.1f°F (%s)\n", temperature, weatherStatus)
}
