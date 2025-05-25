/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// initCmd represents the init command
var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Falcon CLI configuration",
	Long: `This command will guide you through setting up your Falcon CLI configuration.
It will prompt you for your Falcon API credentials and cloud region.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Questions for the user
		qs := []*survey.Question{
			{
				Name: "client_id",
				Prompt: &survey.Input{
					Message: "Enter your Falcon Client ID:",
				},
				Validate: func(val interface{}) error {
					if str, ok := val.(string); !ok || len(str) == 0 {
						return fmt.Errorf("Client ID cannot be empty")
					}
					return nil
				},
			},
			{
				Name: "client_secret",
				Prompt: &survey.Password{
					Message: "Enter your Falcon Client Secret:",
				},
				Validate: func(val interface{}) error {
					if str, ok := val.(string); !ok || len(str) == 0 {
						return fmt.Errorf("Client Secret cannot be empty")
					}
					return nil
				},
			},
			{
				Name: "cloud_region",
				Prompt: &survey.Select{
					Message: "Select your Falcon Cloud Region:",
					Options: []string{"us-1", "us-2", "eu-1", "us-gov-1", "us-gov-2"},
					Default: "us-1",
				},
			},
		}

		// Answers will be stored here
		answers := struct {
			ClientID     string `survey:"client_id"`
			ClientSecret string `survey:"client_secret"`
			CloudRegion  string `survey:"cloud_region"`
		}{}

		// Perform the questions
		err := survey.Ask(qs, &answers)
		if err != nil {
			fmt.Printf("Failed to get answers: %v\n", err)
			return
		}

		// Get home directory
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Failed to get home directory: %v\n", err)
			return
		}

		// Create .falcon-cli directory if it doesn't exist
		falconDir := filepath.Join(home, ".falcon-cli")
		err = os.MkdirAll(falconDir, 0755)
		if err != nil {
			fmt.Printf("Failed to create config directory: %v\n", err)
			return
		}

		// Set up viper
		viper.Set("falcon.client_id", answers.ClientID)
		viper.Set("falcon.client_secret", answers.ClientSecret)
		viper.Set("falcon.cloud_region", answers.CloudRegion)

		// Save the config
		configPath := filepath.Join(falconDir, "config.yaml")
		err = viper.WriteConfigAs(configPath)
		if err != nil {
			fmt.Printf("Failed to save config: %v\n", err)
			return
		}

		fmt.Println("\nConfiguration saved successfully!")
		fmt.Printf("Config file location: %s\n", configPath)
	},
}

func init() {
	// No flags needed for init command
}
