package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

func main() {
	godotenv.Load()

	var metric bool
	cmd := &cobra.Command{
		Use:   "weather [city]",
		Short: "Get the current weather for a specified city",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			apiKey := os.Getenv("OPENWEATHER_API_KEY")
			if apiKey == "" {
				fmt.Println("Error: missing API key")
				os.Exit(1)
			}

			city := url.QueryEscape(args[0])
			units := "imperial"
			if metric {
				units = "metric"
			}

			temp, status, err := getWeather(apiKey, city, units)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}

			unit := "F"
			if metric {
				unit = "C"
			}
			fmt.Printf("%.1f°%s (%s)\n", temp, unit, status)
		},
	}

	cmd.Flags().BoolVarP(&metric, "metric", "m", false, "Use metric units (default is imperial)")
	if err := cmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
