package permission

import (
	"qube-challenge-2016/dto"
	"strings"
)

// Importing DTO package for data transfer objects

// The CheckPermission function checks if a distributor has access to certain test data based on their
// inclusion and exclusion lists.
func CheckPermission(distributorName string, inputData map[string]bool, distributorInformation map[string]dto.Distributor, origin string) []string {
	var validationResult []string
	var errorMsg []string

	// Get distributor data by name
	distributorData, found := distributorInformation[distributorName]
	if !found {
		return []string{"Distributor " + distributorName + " not found"}
	}

	// Helper function to check access
	checkAccess := func(region string) bool {
		if distributorData.Exclude[region] {
			return false
		}
		if distributorData.Include[region] {
			return true
		}
		return false
	}

	for data := range inputData {
		regionParts := strings.Split(data, "-")
		regionLevel := len(regionParts)

		hasAccess := false
		switch regionLevel {
		case 1:
			hasAccess = checkAccess(data)
		case 2:
			countryRegion := regionParts[1]
			hasAccess = (checkAccess(countryRegion) && !distributorData.Exclude[data]) || checkAccess(data)
		case 3:
			countryRegion := regionParts[2]
			stateRegion := regionParts[1] + "-" + regionParts[2]

			if checkAccess(countryRegion) {
				if checkAccess(stateRegion) {
					hasAccess = checkAccess(data) || !distributorData.Exclude[data]
				} else {
					hasAccess = !distributorData.Exclude[stateRegion] && !distributorData.Exclude[data]
				}
			} else {
				hasAccess = checkAccess(stateRegion) && !distributorData.Exclude[data] || checkAccess(data)
			}
		}

		message := distributorData.Name + " does not have access to " + data
		if hasAccess {
			message = distributorData.Name + " has access to " + data
		}

		validationResult = append(validationResult, message)
		if !hasAccess {
			errorMsg = append(errorMsg, message)
		}
	}

	if origin == "subDistributionCreation" {
		return errorMsg
	}
	return validationResult
}
