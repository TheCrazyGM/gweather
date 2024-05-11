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
	// baseURL is the base URL for the OpenWeatherMap API.
	baseURL = "http://api.openweathermap.org/data/2.5/weather"
)

// Errors
var (
	// ErrMissingAPIKey is returned when the API key is not found.
	ErrMissingAPIKey = fmt.Errorf("API key not found")
  ErrMissingCity = fmt.Errorf("Please provided a city")
)

// City represents the city name.
type City string

// Temperature represents the temperature in Fahrenheit.
type Temperature float64

// getWeatherData retrieves the weather data from the OpenWeatherMap API.
// It takes the API key, city name, and units as input, and returns the
// weather data as a map[string]interface{} and any error that occurred.
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

	return weatherData, nil
}

// getCurrentTemperature extracts the current temperature from the weather data.
// It takes the weather data as a map[string]interface{} and returns the
// current temperature as a Temperature and any error that occurred.
func getCurrentTemperature(weatherData map[string]interface{}) (Temperature, error) {
	if cod, ok := weatherData["cod"].(float64); !ok || cod != 200 {
		return 0, fmt.Errorf("failed to retrieve temperature: invalid response code")
	}

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

// displayTemperature prints the temperature.
func displayTemperature(temperature Temperature) {
	fmt.Printf("%.2f°F\n", temperature)
}

// getAPIKey retrieves the API key from the environment.
// It returns the API key as a string and any error that occurred.
func getAPIKey() (string, error) {
	apiKey, ok := os.LookupEnv("OPENWEATHER_API_KEY")
	if !ok {
		return "", fmt.Errorf("%w", ErrMissingAPIKey)
	}
	return apiKey, nil
}

// getCity retrieves the city name from the command line arguments.
// It returns the city name as a City and any error that occurred.
func getCity() (City, error) {
	if len(os.Args) < 2 {
		return "", fmt.Errorf("%w", ErrMissingCity)
	}
  city := url.QueryEscape(os.Args[1])
	return City(city), nil
}

func main() {
	// Load the .env file
	godotenv.Load()

	// Retrieve the API key from the environment
	apiKey, err := getAPIKey()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Retrieve the city name from the command line arguments
	city, err := getCity()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	units := "imperial"

	// Retrieve the weather data from the OpenWeatherMap API
	weatherData, err := getWeatherData(apiKey, string(city), units)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Extract the current temperature from the weather data
	temperature, err := getCurrentTemperature(weatherData)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Display the current temperature
	displayTemperature(temperature)
}
