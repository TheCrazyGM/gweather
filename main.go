package main

import (
	"encoding/json"
	"fmt"
	"io"
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
)
type WeatherData struct {
	Cod     int           `json:"cod"`
	Main    MainData      `json:"main"`
	Weather []WeatherInfo `json:"weather"`
}

type MainData struct {
	Temp float64 `json:"temp"`
}

type WeatherInfo struct {
	MainStatus Status `json:"main"`
}

type Status string

func getWeatherData(apiKey, city, units string) (*WeatherData, error) {
	url := fmt.Sprintf("%s?q=%s&appid=%s&units=%s", baseURL, city, apiKey,
		units)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weather data: %w", err)
	}
	defer resp.Body.Close()

	var weatherData WeatherData
	dec := json.NewDecoder(resp.Body)

	// Use the Decode method to iterate over the JSON objects until an error occurs.
	for {
		if err := dec.Decode(&weatherData); err == io.EOF {
			break // Break the loop if there is no more data.
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
