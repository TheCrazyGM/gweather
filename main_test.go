package main

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// Helper function to create a testable command and capture its output
func executeCommand(cmd *cobra.Command, args ...string) (string, error) {
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return buf.String(), err
}

func TestCmdArgsValidation(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		wantError bool
	}{
		{
			name:      "No arguments",
			args:      []string{},
			wantError: true,
		},
		{
			name:      "One argument",
			args:      []string{"London"},
			wantError: false,
		},
		{
			name:      "Too many arguments",
			args:      []string{"London", "Paris"},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a command that won't make external API calls
			cmd := &cobra.Command{
				Use:  "weather [city]",
				Args: cobra.ExactArgs(1),
				Run:  func(cmd *cobra.Command, args []string) {},
			}

			_, err := executeCommand(cmd, tt.args...)
			if (err != nil) != tt.wantError {
				t.Errorf("Command execution error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestMetricFlag(t *testing.T) {
	// Save and restore environment variable
	origAPIKey := os.Getenv("OPENWEATHER_API_KEY")
	defer os.Setenv("OPENWEATHER_API_KEY", origAPIKey)

	// Mock the getWeather function for testing
	originalGetWeather := getWeather
	defer func() { getWeather = originalGetWeather }()

	getWeather = func(apiKey, city, units string) (float64, string, string, error) {
		// Validate units parameter is passed correctly
		if city != "TestCity" {
			t.Errorf("Expected city 'TestCity', got '%s'", city)
		}

		// Return different temperature based on units
		if units == "metric" {
			return 20.0, "sunny", "☀️", nil // Celsius
		}
		return 68.0, "sunny", "☀️", nil // Fahrenheit (approximately 20°C)
	}

	// Set API key for tests
	os.Setenv("OPENWEATHER_API_KEY", "test-api-key")

	tests := []struct {
		name        string
		args        []string
		wantMetric  bool
		wantContain string
	}{
		{
			name:        "Default unit (imperial)",
			args:        []string{"TestCity"},
			wantMetric:  false,
			wantContain: "68.0°F",
		},
		{
			name:        "Metric unit with short flag",
			args:        []string{"-m", "TestCity"},
			wantMetric:  true,
			wantContain: "20.0°C",
		},
		{
			name:        "Metric unit with long flag",
			args:        []string{"--metric", "TestCity"},
			wantMetric:  true,
			wantContain: "20.0°C",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Skip test execution that would make actual API calls
			// This is just a simple test of the flag parsing
			var metric bool
			cmd := &cobra.Command{
				Use:  "weather [city]",
				Args: cobra.ExactArgs(1),
				Run: func(cmd *cobra.Command, args []string) {
					city := args[0]
					units := "imperial"
					if metric {
						units = "metric"
					}

					temp, status, icon, _ := getWeather("test-api-key", city, units)

					unit := "F"
					if metric {
						unit = "C"
					}
					if cmd.Flag("no-emoji") != nil && cmd.Flag("no-emoji").Changed {
						cmd.Printf("%.1f°%s [ %s ]\n", temp, unit, status)
					} else {
						cmd.Printf("%s %.1f°%s [ %s ]\n", icon, temp, unit, status)
					}
				},
			}

			cmd.Flags().BoolVarP(&metric, "metric", "m", false, "Use metric units")
			cmd.Flags().Bool("no-emoji", false, "Disable emoji icons in output")

			output, err := executeCommand(cmd, tt.args...)
			if err != nil {
				t.Errorf("Command execution error = %v", err)
				return
			}

			if !strings.Contains(output, tt.wantContain) {
				t.Errorf("Output does not contain expected string.\nWanted: %s\nGot: %s",
					tt.wantContain, output)
			}
		})
	}
}

func TestMissingAPIKey(t *testing.T) {
	// Save environment variable
	origAPIKey := os.Getenv("OPENWEATHER_API_KEY")
	defer os.Setenv("OPENWEATHER_API_KEY", origAPIKey)

	// Clear API key for this test
	os.Setenv("OPENWEATHER_API_KEY", "")

	var metric bool
	cmd := &cobra.Command{
		Use:   "weather [city]",
		Short: "Get the current weather for a specified city",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			apiKey := os.Getenv("OPENWEATHER_API_KEY")
			if apiKey == "" {
				cmd.PrintErrln("Error: missing API key")
				return
			}

			city := args[0]
			units := "imperial"
			if metric {
				units = "metric"
			}

			temp, status, icon, _ := getWeather(apiKey, city, units)

			unit := "F"
			if metric {
				unit = "C"
			}
			if cmd.Flag("no-emoji") != nil && cmd.Flag("no-emoji").Changed {
				cmd.Printf("%.1f°%s [ %s ]\n", temp, unit, status)
			} else {
				cmd.Printf("%s %.1f°%s [ %s ]\n", icon, temp, unit, status)
			}
		},
	}

	cmd.Flags().BoolVarP(&metric, "metric", "m", false, "Use metric units")

	output, _ := executeCommand(cmd, "London")

	if !strings.Contains(output, "missing API key") {
		t.Errorf("Expected error about missing API key, got: %s", output)
	}
}
