package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

// main is the entry point for the weather CLI tool, handling argument parsing, environment variable loading, and output formatting to display the current weather for a specified city.
func main() {
	// Load .env file but ignore if it does not exist; warn on other errors
	if err := godotenv.Load(); err != nil {
		if !os.IsNotExist(err) {
			fmt.Fprintln(os.Stderr, "Warning: error loading .env file:", err)
		}
	}

	var metric bool
	var apiKeyFlag string
	var noEmoji bool
	cmd := &cobra.Command{
		Use:   "weather [city]",
		Short: "Get the current weather for a specified city",
		Args:  cobra.MinimumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// If help flag is set, skip running main logic
			help, _ := cmd.Flags().GetBool("help")
			if help || cmd.CalledAs() == "help" {
				return nil
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// If help flag is set, skip running main logic
			help, _ := cmd.Flags().GetBool("help")
			if help || cmd.CalledAs() == "help" {
				return nil
			}

			apiKey := apiKeyFlag
			if apiKey == "" {
				apiKey = os.Getenv("OPENWEATHER_API_KEY")
				if apiKey == "" {
					return fmt.Errorf("missing API key")
				}
			}

			city := url.QueryEscape(args[0])
			units := "imperial"
			if metric {
				units = "metric"
			}

			temp, status, icon, err := getWeather(apiKey, city, units)
			if err != nil {
				return err
			}

			unit := "F"
			if metric {
				unit = "C"
			}
			if noEmoji {
				fmt.Printf("%.1f°%s [ %s ]\n", temp, unit, status)
			} else {
				fmt.Printf("%.1f°%s [ %s %s ]\n", temp, unit, icon, status)
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&metric, "metric", "m", false, "Use metric units (default is imperial)")
	cmd.Flags().StringVarP(&apiKeyFlag, "api-key", "k", "", "OpenWeather API key")
	cmd.Flags().BoolVarP(&noEmoji, "no-emoji", "n", false, "Disable emoji icons in output")
	if err := cmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
