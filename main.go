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

func main() {
	groupedData, err := parser.ParseCSVFile("cities.csv") // Parsing the CSV file containing city data
	if err != nil {
		log.Fatalf("Error parsing CSV file: %v", err) // Exiting if there's an error in parsing CSV
	}

	distributorInformation := make(map[string]dto.Distributor)
	for {
		choice := input.PromptMenu()
		switch choice {
		case "Create a new distributor":
			distributorData := input.PromptDistributorData(false) // Getting data for a new distributor and send false for this the distributor Data

			// If the distributor already exists, notify the user
			if _, exists := distributorInformation[distributorData.Name]; exists {
				fmt.Println("Distributor with this name already exists. Please use a different name.")
				continue
			}

			errorRes := validator.ValidateDistributorData(distributorData, groupedData, distributorInformation, false) // Validating distributor data
			if len(errorRes) > 0 {
				fmt.Println(strings.Join(errorRes, "\n"))
				continue
			}

			distributorInformation[distributorData.Name] = distributorData // Add distributor to the map
		case "Create a sub distributor":
			subDistributorData := input.PromptDistributorData(true) // Getting data for a new sub-distributor and send true for this the sub-distributor Data

			// If the distributor already exists, notify the user
			if _, exists := distributorInformation[subDistributorData.Name]; exists {
				fmt.Println("Distributor with this name already exists. Please use a different name.")
				continue
			}

			errorRes := validator.ValidateDistributorData(subDistributorData, groupedData, distributorInformation, true) // Validating sub-distributor data
			if len(errorRes) > 0 {
				fmt.Println(strings.Join(errorRes, "\n"))
				continue
			}

			parentDistributor := distributorInformation[subDistributorData.Parent]

			// First, add all include values
			maps.Copy(subDistributorData.Exclude, parentDistributor.Exclude)

			distributorInformation[subDistributorData.Name] = subDistributorData // Add sub-distributor to the map
		case "Check permission for a distributor":
			checkPermissionData := input.PromptCheckPermissionData()                                                    // Getting data to check permission
			errorRes := validator.ValidateCheckPermissionData(checkPermissionData, groupedData, distributorInformation) // Validating permission check data
			if len(errorRes) > 0 {
				fmt.Println(strings.Join(errorRes, "\n"))
				continue
			}
			checkPermissionResult := permission.CheckPermission(checkPermissionData.DistributorName, checkPermissionData.Regions, distributorInformation, "Check Permission") // Checking permission
			fmt.Println("Check Permission Result:\n", strings.Join(checkPermissionResult, "\n"))
		case "View Distributors information":
			ViewDistributorsInfo(distributorInformation)
		case "Exit":
			fmt.Println("Exiting the program")
			return // Exiting the program
		}
	}
}

func ViewDistributorsInfo(distributorInformation map[string]dto.Distributor) {
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

func GetKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
