package cmd

import (
	"fmt"

	"github.com/HARSH16DAWAR/falcon-cli/utils"
	"github.com/spf13/cobra"
)

// HostsResponse represents the response from the hosts API
type HostsResponse struct {
	Resources []string `json:"resources"`
	Errors    []struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`
	Meta struct {
		QueryTime float64 `json:"query_time"`
		PoweredBy string  `json:"powered_by"`
		TraceID   string  `json:"trace_id"`
	} `json:"meta"`
}

// hostsCmd represents the hosts command
var hostsCmd = &cobra.Command{
	Use:   "hosts",
	Short: "List hosts in your Falcon environment",
	Long:  `List all hosts in your Falcon environment with their details.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create Falcon client
		client, err := utils.NewFalconClient()
		if err != nil {
			return fmt.Errorf("error creating Falcon client: %v", err)
		}

		// Make request to get hosts
		resp, err := client.Get("/devices/queries/devices/v1", nil)
		if err != nil {
			return fmt.Errorf("error getting hosts: %v", err)
		}

		// Parse response
		var result HostsResponse
		if err := client.ParseResponse(resp, &result); err != nil {
			return fmt.Errorf("error parsing response: %v", err)
		}

		// Print results
		fmt.Printf("Found %d hosts\n", len(result.Resources))
		fmt.Printf("Query time: %.2f seconds\n", result.Meta.QueryTime)
		fmt.Printf("Trace ID: %s\n", result.Meta.TraceID)

		// Print any errors
		if len(result.Errors) > 0 {
			fmt.Println("\nErrors:")
			for _, err := range result.Errors {
				fmt.Printf("- %s (code: %d)\n", err.Message, err.Code)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(hostsCmd)
}
