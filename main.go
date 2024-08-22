package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

func getAPIKey() (string, error) {
	apiKey, ok := os.LookupEnv("OPENWEATHER_API_KEY")
	if !ok {
		return "", fmt.Errorf("%w", ErrMissingAPIKey)
	}
	return apiKey, nil
}

func getCity() (City, error) {
	if len(os.Args) < 2 {
		return "", fmt.Errorf("%w: please provide a city name", ErrMissingCity)
	}
	city := url.QueryEscape(os.Args[1])
	if city == "" {
		return "", fmt.Errorf("%w: city name cannot be empty", ErrMissingCity)
	}
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
