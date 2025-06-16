package filter

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Filter represents a saved filter configuration
type Filter struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Filter      string `json:"filter"`
}

// filterCmd represents the base filter command
var filterCmd = &cobra.Command{
	Use:   "filter",
	Short: "Manage saved filters for Falcon CLI commands",
	Long:  `Save, list, and delete filters that can be reused across different Falcon CLI commands.`,
}

// saveCmd represents the save filter command
var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "Save a filter for later use",
	Long:  `Save a filter with a name and description for later use in various Falcon CLI commands.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		filterType, _ := cmd.Flags().GetString("type")
		description, _ := cmd.Flags().GetString("description")
		filterValue, _ := cmd.Flags().GetString("filter")

		// Create new filter
		newFilter := Filter{
			Name:        name,
			Type:        filterType,
			Description: description,
			Filter:      filterValue,
		}

		// Get existing filters
		var filters []Filter
		if err := viper.UnmarshalKey("filters", &filters); err != nil {
			filters = []Filter{}
		}

		// Check if filter with same name and type already exists
		for i, f := range filters {
			if f.Name == name && f.Type == filterType {
				// Update existing filter
				filters[i] = newFilter
				viper.Set("filters", filters)
				if err := viper.WriteConfig(); err != nil {
					return fmt.Errorf("error saving filter: %v", err)
				}
				fmt.Printf("Updated filter '%s' for type '%s'\n", name, filterType)
				return nil
			}
		}

		// Add new filter
		filters = append(filters, newFilter)
		viper.Set("filters", filters)
		if err := viper.WriteConfig(); err != nil {
			return fmt.Errorf("error saving filter: %v", err)
		}

		fmt.Printf("Saved filter '%s' for type '%s'\n", name, filterType)
		return nil
	},
}

// listCmd represents the list filters command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List saved filters",
	Long:  `List all saved filters, optionally filtered by type.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		filterType, _ := cmd.Flags().GetString("type")

		var filters []Filter
		if err := viper.UnmarshalKey("filters", &filters); err != nil {
			return fmt.Errorf("error reading filters: %v", err)
		}

		if len(filters) == 0 {
			fmt.Println("No filters saved")
			return nil
		}

		fmt.Println("Saved filters:")
		for _, f := range filters {
			if filterType == "" || f.Type == filterType {
				fmt.Printf("\nName: %s\nType: %s\nDescription: %s\nFilter: %s\n",
					f.Name, f.Type, f.Description, f.Filter)
			}
		}
		return nil
	},
}

// deleteCmd represents the delete filter command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a saved filter",
	Long:  `Delete a saved filter by name and type.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		filterType, _ := cmd.Flags().GetString("type")

		var filters []Filter
		if err := viper.UnmarshalKey("filters", &filters); err != nil {
			return fmt.Errorf("error reading filters: %v", err)
		}

		found := false
		var newFilters []Filter
		for _, f := range filters {
			if f.Name == name && f.Type == filterType {
				found = true
				continue
			}
			newFilters = append(newFilters, f)
		}

		if !found {
			return fmt.Errorf("filter '%s' of type '%s' not found", name, filterType)
		}

		viper.Set("filters", newFilters)
		if err := viper.WriteConfig(); err != nil {
			return fmt.Errorf("error deleting filter: %v", err)
		}

		fmt.Printf("Deleted filter '%s' of type '%s'\n", name, filterType)
		return nil
	},
}

// GetCommand returns the filter command
func GetCommand() *cobra.Command {
	// Add flags to save command
	saveCmd.Flags().String("name", "", "Name of the filter")
	saveCmd.Flags().String("type", "", "Type of the filter (e.g., hosts, detections)")
	saveCmd.Flags().String("description", "", "Description of the filter")
	saveCmd.Flags().String("filter", "", "Filter expression to save")
	saveCmd.MarkFlagRequired("name")
	saveCmd.MarkFlagRequired("type")
	saveCmd.MarkFlagRequired("filter")

	// Add flags to list command
	listCmd.Flags().String("type", "", "Filter type to list (optional)")

	// Add flags to delete command
	deleteCmd.Flags().String("name", "", "Name of the filter to delete")
	deleteCmd.Flags().String("type", "", "Type of the filter to delete")
	deleteCmd.MarkFlagRequired("name")
	deleteCmd.MarkFlagRequired("type")

	// Add subcommands
	filterCmd.AddCommand(saveCmd)
	filterCmd.AddCommand(listCmd)
	filterCmd.AddCommand(deleteCmd)

	return filterCmd
}
