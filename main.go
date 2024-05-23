package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

const (
	baseURL = "http://api.openweathermap.org/data/2.5/weather"
)

var (
	// ErrMissingAPIKey is an error indicating that the API key is missing.
	ErrMissingAPIKey = fmt.Errorf("API key not found")
	// ErrMissingCity is an error indicating that the city is missing.
	ErrMissingCity = fmt.Errorf("please provide a city")
)

type (
	// City represents a city name.
	City string
	// Temperature represents a temperature value.
	Temperature float64
	// Status represents the weather status.
	Status string
)

// getWeatherData retrieves weather data from the OpenWeatherMap API.
func getWeatherData(apiKey, city, units string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s?q=%s&appid=%s&units=%s", baseURL, city, apiKey, units)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weather data: %w", err)
	}
	defer resp.Body.Close()

	var weatherData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&weatherData); err != nil {
		return nil, fmt.Errorf("failed to decode weather data: %w", err)
	}

	if cod, ok := weatherData["cod"].(float64); !ok || cod != 200 {
		return nil, fmt.Errorf("failed to retrieve weather data: invalid response code")
	}

	return weatherData, nil
}

// getCurrentTemperature retrieves the current temperature from the weather data.
func getCurrentTemperature(weatherData map[string]interface{}) (Temperature, error) {
	main, ok := weatherData["main"].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("failed to retrieve temperature: invalid main data")
	}

	temp, ok := main["temp"].(float64)
	if !ok {
		return 0, fmt.Errorf("failed to retrieve temperature: invalid temperature data")
	}

	return Temperature(temp), nil
}

// getWeatherStatus retrieves the current weather status from the weather data.
func getWeatherStatus(weatherData map[string]interface{}) (Status, error) {
	weather, ok := weatherData["weather"].([]interface{})
	if !ok || len(weather) == 0 {
		return "", fmt.Errorf("failed to retrieve status: invalid weather data")
	}

	first, ok := weather[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("failed to retrieve current status: invalid status data")
	}

	status, ok := first["main"].(string)
	if !ok {
		return "", fmt.Errorf("failed to retrieve current status: invalid status data")
	}

	return Status(status), nil
}

// displayWeather displays the temperature and weather status.
func displayWeather(temperature Temperature, weatherStatus Status) {
	fmt.Printf("%.2f°F (%s)\n", temperature, weatherStatus)
}

// getAPIKey retrieves the API key from the environment variables.
func getAPIKey() (string, error) {
	apiKey, ok := os.LookupEnv("OPENWEATHER_API_KEY")
	if !ok {
		return "", fmt.Errorf("%w", ErrMissingAPIKey)
	}
	return apiKey, nil
}

// getCity retrieves the city from the command-line arguments.
func getCity() (City, error) {
	if len(os.Args) < 2 {
		return "", fmt.Errorf("%w", ErrMissingCity)
	}
	city := url.QueryEscape(os.Args[1])
	return City(city), nil
}

func main() {
	godotenv.Load()

	apiKey, err := getAPIKey()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	city, err := getCity()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	units := "imperial"

	weatherData, err := getWeatherData(apiKey, string(city), units)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	temperature, err := getCurrentTemperature(weatherData)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	weatherStatus, err := getWeatherStatus(weatherData)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	displayWeather(temperature, weatherStatus)
}
