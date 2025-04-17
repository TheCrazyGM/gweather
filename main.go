package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

func main() {
   // Load .env file but ignore if it does not exist; warn on other errors
   if err := godotenv.Load(); err != nil {
       if !os.IsNotExist(err) {
           fmt.Fprintln(os.Stderr, "Warning: error loading .env file:", err)
       }
   }

   var metric bool
   cmd := &cobra.Command{
       Use:   "weather [city]",
       Short: "Get the current weather for a specified city",
       Args:  cobra.ExactArgs(1),
       RunE: func(cmd *cobra.Command, args []string) error {
           apiKey := os.Getenv("OPENWEATHER_API_KEY")
           if apiKey == "" {
               return fmt.Errorf("missing API key")
           }

           city := url.QueryEscape(args[0])
           units := "imperial"
           if metric {
               units = "metric"
           }

           temp, status, err := getWeather(apiKey, city, units)
           if err != nil {
               return err
           }

           unit := "F"
           if metric {
               unit = "C"
           }
           fmt.Printf("%.1f°%s [ %s ]\n", temp, unit, status)
           return nil
       },
   }

	cmd.Flags().BoolVarP(&metric, "metric", "m", false, "Use metric units (default is imperial)")
	if err := cmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
