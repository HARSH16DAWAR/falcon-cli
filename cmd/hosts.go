package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/HARSH16DAWAR/falcon-cli/cmd/filter"
	"github.com/HARSH16DAWAR/falcon-cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

// getFilterValue returns the filter value, either from the --filter flag or from a saved filter
func getFilterValue(cmd *cobra.Command) (string, error) {
	filterValue, _ := cmd.Flags().GetString("filter")
	filterName, _ := cmd.Flags().GetString("filter-name")

	if filterValue != "" && filterName != "" {
		return "", fmt.Errorf("cannot use both --filter and --filter-name")
	}

	if filterName != "" {
		var filters []filter.Filter
		if err := viper.UnmarshalKey("filters", &filters); err != nil {
			return "", fmt.Errorf("error reading filters: %v", err)
		}

		for _, f := range filters {
			if f.Name == filterName && f.Type == "hosts" {
				return f.Filter, nil
			}
		}
		return "", fmt.Errorf("filter '%s' not found for type 'hosts'", filterName)
	}

	return filterValue, nil
}

// hostsCmd represents the hosts command
var hostsCmd = &cobra.Command{
	Use:   "hosts",
	Short: "List hosts in your Falcon environment",
	Long:  `List all hosts in your Falcon environment with their details. You can filter hosts using the --filter flag or a saved filter using --filter-name.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get filter value
		filterValue, err := getFilterValue(cmd)
		if err != nil {
			return err
		}

		// Create Falcon client
		client, err := utils.NewFalconClient()
		if err != nil {
			return fmt.Errorf("error creating Falcon client: %v", err)
		}

		// Prepare query parameters
		params := make(map[string]string)
		if filterValue != "" {
			params["filter"] = filterValue
		}

		// Make request to get hosts
		resp, err := client.Get("/devices/queries/devices/v1", params)
		if err != nil {
			return fmt.Errorf("error getting hosts: %v", err)
		}
		defer resp.Body.Close()

		// Read the raw response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error reading response body: %v", err)
		}

		// Pretty print the raw JSON
		var prettyJSON bytes.Buffer
		if err := json.Indent(&prettyJSON, body, "", "  "); err != nil {
			return fmt.Errorf("error formatting JSON: %v", err)
		}

		// Print the raw JSON response
		fmt.Println(prettyJSON.String())

		return nil
	},
}

func init() {
	hostsCmd.Flags().String("filter", "", "Filter hosts (e.g., platform_name:'Windows')")
	hostsCmd.Flags().String("filter-name", "", "Use a saved filter by name")
	RootCmd.AddCommand(hostsCmd)
}
