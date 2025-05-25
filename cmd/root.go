/*
Copyright © 2025 HARSH DAWAR <harsh160102@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/HARSH16DAWAR/falcon-cli/cmd/config"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "falcon-cli",
	Short: "Falcon CLI is the command line interface to interact with CrowdStrike Falcon",
	Long: `Falcon CLI is the command line interface to interact with CrowdStrike Falcon. 
Now you might think doesn't CrowdStrike have a falcon-cli itself. Well yes but they arent doing anything with it. So we are going to make our own.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	addSubCommands()

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.falcon-cli/config.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Create .falcon-cli directory if it doesn't exist
		falconDir := filepath.Join(home, ".falcon-cli")
		err = os.MkdirAll(falconDir, 0755)
		cobra.CheckErr(err)

		// Search config in .falcon-cli directory
		viper.AddConfigPath(falconDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		// It's okay if there's no config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Fprintln(os.Stderr, "Error reading config file:", err)
		}
	} else {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func addSubCommands() {
	rootCmd.AddCommand(config.InitCmd)
	rootCmd.AddCommand(hostsCmd)
}
