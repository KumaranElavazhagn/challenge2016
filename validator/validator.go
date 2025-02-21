package validator

import (
	"qube-challenge-2016/dto"
	"qube-challenge-2016/permission"
	"strings"
)

// The function `ValidateDistributorData` validates the data of a sub-distributor, checking for
// errors such as empty fields, duplicate names, invalid regions, and incorrect parent distributor
// name.
func ValidateDistributorData(data dto.Distributor, groupedData map[string]map[string]map[string]bool, distributorInformation map[string]dto.Distributor, subDistributor bool) []string {
	var errorMsg []string

	if strings.TrimSpace(data.Name) == "" {
		errorMsg = append(errorMsg, "Distributor Name must not be empty, please enter a valid distributor name")
	}

	if len(data.Include) == 0 {
		errorMsg = append(errorMsg, "Include Regions must not be empty, please enter valid regions")
	} else {
		for region := range data.Include {
			if !ValidateRegion(region, groupedData) {
				errorMsg = append(errorMsg, "Include Region '"+region+"' is not present in csv, please enter a valid region")
			}
		}
	}

	for region := range data.Exclude {
		if !ValidateRegion(region, groupedData) {
			errorMsg = append(errorMsg, "Exclude Region '"+region+"' is not present in csv, please enter a valid region")
		}
		if _, exists := data.Include[region]; exists {
			errorMsg = append(errorMsg, "Exclude Region '"+region+"' should not be the same as Include Region, please enter a valid region")
		}
	}

	if subDistributor {
		if strings.TrimSpace(data.Parent) == "" {
			errorMsg = append(errorMsg, "Parent distributor Name must not be empty, please enter a valid parent distributor name")
		} else if !ValidateDistributorName(strings.ToUpper(strings.TrimSpace(data.Parent)), distributorInformation) {
			errorMsg = append(errorMsg, "Parent distributor Name does not exist, please enter an existing parent distributor name")
		}

		if len(errorMsg) == 0 {
			// InputData := append(GetKeys(data.Include), GetKeys(data.Exclude)...)
			InputData := make(map[string]bool)

			// First, add all include values
			for key, value := range data.Include {
				InputData[key] = value
			}

			// Only add exclude values if the key is not already in InputData
			for key, value := range data.Exclude {
				if _, exists := InputData[key]; !exists {
					InputData[key] = value
				}
			}

			checkPermissionWithParent := permission.CheckPermission(strings.TrimSpace(data.Parent), InputData, distributorInformation, "subDistributionCreation")
			if len(checkPermissionWithParent) > 0 {
				errorMsg = append(errorMsg, checkPermissionWithParent...)
			}
		}
	}

	return errorMsg
}

// The function "ValidateDistributorName" checks if a given distributor name exists in a list of
// distributor information.
func ValidateDistributorName(distributorName string, distributorInformation map[string]dto.Distributor) bool {
	for _, distributor := range distributorInformation {
		if strings.EqualFold(distributor.Name, distributorName) {
			return true
		}
	}
	return false
}

// ValidateRegion checks if a given region exists in the structured map data.
func ValidateRegion(region string, groupedData map[string]map[string]map[string]bool) bool {

	regionParts := strings.Split(strings.ToUpper(strings.TrimSpace(region)), "-")
	regionLevel := len(regionParts)

	if regionLevel == 0 || regionLevel > 3 {
		return false
	}

	switch regionLevel {
	case 1: // Country level
		if _, countryExists := groupedData[regionParts[0]]; !countryExists {
			return false
		}

	case 2: // State level (State-Country)
		countryName := regionParts[1]
		stateName := regionParts[0]

		if states, countryExists := groupedData[countryName]; countryExists {
			if _, stateExists := states[stateName]; !stateExists {
				return false
			}
		} else {
			return false
		}

	case 3: // City level (City-State-Country)
		countryName := regionParts[2]
		stateName := regionParts[1]
		cityName := regionParts[0]

		if states, countryExists := groupedData[countryName]; countryExists {
			if cities, stateExists := states[stateName]; stateExists {
				if _, cityExists := cities[cityName]; cityExists {
					return true
				}
			}
		}
		return false // If city not found
	}

	return true // All regions exist in the map
}

// The function `ValidateCheckPermissionData` validates the `CheckPermissionData` object by checking if
// the distributor name is not empty and exists in the distributor information, and if all the regions
// in the data exist in the grouped data.
func ValidateCheckPermissionData(data dto.CheckPermissionData, groupedData map[string]map[string]map[string]bool, distributorInformation map[string]dto.Distributor) []string {
	var errorMsg []string

	if strings.TrimSpace(data.DistributorName) == "" {
		errorMsg = append(errorMsg, "Distributor Name must not be empty, please enter a valid distributor name")
	} else if !ValidateDistributorName(strings.ToUpper(strings.TrimSpace(data.DistributorName)), distributorInformation) {
		errorMsg = append(errorMsg, "Distributor name does not exist")
	}

	for region := range data.Regions {
		if !ValidateRegion(region, groupedData) {
			errorMsg = append(errorMsg, strings.ToUpper(region)+" does not exist in the csv file, please enter a valid region")
		}
	}

	return errorMsg
}
