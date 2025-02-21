package main

import (
	"fmt"
	"log"
	"maps"
	"qube-challenge-2016/dto"
	"qube-challenge-2016/input"
	"qube-challenge-2016/parser"
	"qube-challenge-2016/permission"
	"qube-challenge-2016/validator"
	"strings"
)

// main function to handle distributor management operations.
func main() {
	// Parsing the CSV file containing city data.
	groupedData, err := parser.ParseCSVFile("cities.csv")
	if err != nil {
		log.Fatalf("Error parsing CSV file: %v", err) // Exit if there's an error in parsing CSV.
	}

	// Initialize a map to store distributor information.
	distributorInformation := make(map[string]dto.Distributor)

	// Infinite loop to handle user menu choices.
	for {
		// Display menu and get user choice.
		choice := input.PromptMenu()

		switch choice {
		case "Create a new distributor":
			// Prompt user for new distributor data (passing false since it's not a sub-distributor).
			distributorData := input.PromptDistributorData(false)

			// Check if the distributor name already exists.
			if _, exists := distributorInformation[distributorData.Name]; exists {
				fmt.Println("Distributor with this name already exists. Please use a different name.")
				continue
			}

			// Validate distributor data.
			errorRes := validator.ValidateDistributorData(distributorData, groupedData, distributorInformation, false)
			if len(errorRes) > 0 {
				fmt.Println(strings.Join(errorRes, "\n")) // Print validation errors.
				continue
			}

			// Add the new distributor to the map.
			distributorInformation[distributorData.Name] = distributorData

		case "Create a sub distributor":
			// Prompt user for sub-distributor data (passing true since it's a sub-distributor).
			subDistributorData := input.PromptDistributorData(true)

			// Check if the sub-distributor name already exists.
			if _, exists := distributorInformation[subDistributorData.Name]; exists {
				fmt.Println("Distributor with this name already exists. Please use a different name.")
				continue
			}

			// Validate sub-distributor data.
			errorRes := validator.ValidateDistributorData(subDistributorData, groupedData, distributorInformation, true)
			if len(errorRes) > 0 {
				fmt.Println(strings.Join(errorRes, "\n")) // Print validation errors.
				continue
			}

			// Get the parent distributor's data.
			parentDistributor := distributorInformation[subDistributorData.Parent]

			// Copy parent's exclude values to the sub-distributor.
			maps.Copy(subDistributorData.Exclude, parentDistributor.Exclude)

			// Add the sub-distributor to the map.
			distributorInformation[subDistributorData.Name] = subDistributorData

		case "Check permission for a distributor":
			// Prompt user for permission check details.
			checkPermissionData := input.PromptCheckPermissionData()

			// Validate the input data for permission check.
			errorRes := validator.ValidateCheckPermissionData(checkPermissionData, groupedData, distributorInformation)
			if len(errorRes) > 0 {
				fmt.Println(strings.Join(errorRes, "\n")) // Print validation errors.
				continue
			}

			// Perform permission check for the distributor.
			checkPermissionResult := permission.CheckPermission(
				checkPermissionData.DistributorName,
				checkPermissionData.Regions,
				distributorInformation,
				"Check Permission",
			)

			// Display the permission check results.
			fmt.Println("Check Permission Result:\n", strings.Join(checkPermissionResult, "\n"))

		case "View Distributors information":
			// Display distributor information.
			ViewDistributorsInfo(distributorInformation)

		case "Exit":
			// Exit the program.
			fmt.Println("Exiting the program")
			return
		}
	}
}

// ViewDistributorsInfo displays information about distributors from the provided map.
func ViewDistributorsInfo(distributorInformation map[string]dto.Distributor) {
	// Check if there are no distributors in the map.
	if len(distributorInformation) == 0 {
		fmt.Println("No distributors found.")
		return
	}

	fmt.Println("\n--- Distributor Information ---")

	count := 0
	for name, distributor := range distributorInformation {
		count++
		fmt.Printf("%d. Name: %s, Include: %v, Exclude: %v, Parent: %s\n",
			count, name, GetKeys(distributor.Include), GetKeys(distributor.Exclude), distributor.Parent)
	}
}

// GetKeys extracts the keys from a map[string]bool and returns them as a slice of strings.
func GetKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
