package validator

import (
	"qube-challenge-2016/dto"
	"qube-challenge-2016/permission"
	"strings"
)

// ValidateDistributorData validates a distributor's data, checking for empty fields, invalid regions,
// duplicate names, and verifying the existence of a parent distributor if it's a sub-distributor.
func ValidateDistributorData(data dto.Distributor, groupedData map[string]map[string]map[string]bool, distributorInformation map[string]dto.Distributor, subDistributor bool) []string {
	var errorMsg []string

	// Validate distributor name
	if strings.TrimSpace(data.Name) == "" {
		errorMsg = append(errorMsg, "Distributor Name must not be empty, please enter a valid distributor name")
	}

	// Validate include regions
	if len(data.Include) == 0 {
		errorMsg = append(errorMsg, "Include Regions must not be empty, please enter valid regions")
	} else {
		for region := range data.Include {
			if !ValidateRegion(region, groupedData) {
				errorMsg = append(errorMsg, "Include Region '"+region+"' is not present in csv, please enter a valid region")
			}
		}
	}

	// Validate exclude regions
	for region := range data.Exclude {
		if !ValidateRegion(region, groupedData) {
			errorMsg = append(errorMsg, "Exclude Region '"+region+"' is not present in csv, please enter a valid region")
		}
		// Ensure exclude region is not also in include regions
		if _, exists := data.Include[region]; exists {
			errorMsg = append(errorMsg, "Exclude Region '"+region+"' should not be the same as Include Region, please enter a valid region")
		}
	}

	// If the distributor is a sub-distributor, validate parent distributor
	if subDistributor {
		if strings.TrimSpace(data.Parent) == "" {
			errorMsg = append(errorMsg, "Parent distributor Name must not be empty, please enter a valid parent distributor name")
		} else if !ValidateDistributorName(strings.ToUpper(strings.TrimSpace(data.Parent)), distributorInformation) {
			errorMsg = append(errorMsg, "Parent distributor Name does not exist, please enter an existing parent distributor name")
		}

		if len(errorMsg) == 0 {
			// Combine include and exclude regions for permission check
			InputData := make(map[string]bool)

			// Add all included regions
			for key, value := range data.Include {
				InputData[key] = value
			}

			// Add exclude regions only if they are not already in include regions
			for key, value := range data.Exclude {
				if _, exists := InputData[key]; !exists {
					InputData[key] = value
				}
			}

			// Check if the parent distributor allows sub-distribution in these regions
			checkPermissionWithParent := permission.CheckPermission(strings.TrimSpace(data.Parent), InputData, distributorInformation, "subDistributionCreation")
			if len(checkPermissionWithParent) > 0 {
				errorMsg = append(errorMsg, checkPermissionWithParent...)
			}
		}
	}

	return errorMsg
}

// ValidateDistributorName checks if a given distributor name exists in the distributor information map.
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

	// Validate region level (1 = Country, 2 = State-Country, 3 = City-State-Country)
	if regionLevel == 0 || regionLevel > 3 {
		return false
	}

	switch regionLevel {
	case 1: // Country level
		_, countryExists := groupedData[regionParts[0]]
		return countryExists

	case 2: // State level (State-Country)
		countryName := regionParts[1]
		stateName := regionParts[0]

		if states, countryExists := groupedData[countryName]; countryExists {
			_, stateExists := states[stateName]
			return stateExists
		}

	case 3: // City level (City-State-Country)
		countryName := regionParts[2]
		stateName := regionParts[1]
		cityName := regionParts[0]

		if states, countryExists := groupedData[countryName]; countryExists {
			if cities, stateExists := states[stateName]; stateExists {
				_, cityExists := cities[cityName]
				return cityExists
			}
		}
		return false
	}

	return true
}

// ValidateCheckPermissionData validates the distributor name and regions in a CheckPermissionData object.
func ValidateCheckPermissionData(data dto.CheckPermissionData, groupedData map[string]map[string]map[string]bool, distributorInformation map[string]dto.Distributor) []string {
	var errorMsg []string

	// Validate distributor name
	if strings.TrimSpace(data.DistributorName) == "" {
		errorMsg = append(errorMsg, "Distributor Name must not be empty, please enter a valid distributor name")
	} else if !ValidateDistributorName(strings.ToUpper(strings.TrimSpace(data.DistributorName)), distributorInformation) {
		errorMsg = append(errorMsg, "Distributor name does not exist")
	}

	// Validate regions
	for region := range data.Regions {
		if !ValidateRegion(region, groupedData) {
			errorMsg = append(errorMsg, strings.ToUpper(region)+" does not exist in the csv file, please enter a valid region")
		}
	}

	return errorMsg
}
