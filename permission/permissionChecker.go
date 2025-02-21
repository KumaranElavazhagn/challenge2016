package permission

import (
	"qube-challenge-2016/dto"
	"strings"
)

// CheckPermission verifies if a distributor has access to a given set of regions based on their
// inclusion and exclusion lists.
func CheckPermission(distributorName string, inputData map[string]bool, distributorInformation map[string]dto.Distributor, origin string) []string {
	var validationResult []string
	var errorMsg []string

	// Retrieve distributor data
	distributorData, found := distributorInformation[distributorName]
	if !found {
		return []string{"Distributor " + distributorName + " not found"}
	}

	// Helper function to check access
	checkAccess := func(region string) bool {
		if distributorData.Exclude[region] {
			return false // Excluded regions override inclusion
		}
		return distributorData.Include[region]
	}

	for region := range inputData {
		regionParts := strings.Split(region, "-")
		regionLevel := len(regionParts)

		hasAccess := false

		switch regionLevel {
		case 1: // Country-level access check
			hasAccess = checkAccess(region)

		case 2: // State-Country level check
			countryRegion := regionParts[1]
			hasAccess = checkAccess(region) || (checkAccess(countryRegion) && !distributorData.Exclude[region])

		case 3: // City-State-Country level check
			countryRegion := regionParts[2]
			stateRegion := regionParts[1] + "-" + regionParts[2]

			// If country is accessible
			if checkAccess(countryRegion) {
				if checkAccess(stateRegion) {
					hasAccess = checkAccess(region) || !distributorData.Exclude[region]
				} else {
					hasAccess = !distributorData.Exclude[stateRegion] && !distributorData.Exclude[region]
				}
			} else {
				hasAccess = checkAccess(stateRegion) && !distributorData.Exclude[region] || checkAccess(region)
			}
		}

		// Construct access message
		message := distributorData.Name + " does not have access to " + region
		if hasAccess {
			message = distributorData.Name + " has access to " + region
		}

		validationResult = append(validationResult, message)
		if !hasAccess {
			errorMsg = append(errorMsg, message)
		}
	}

	// Return errors if this check is for sub-distribution creation
	if origin == "subDistributionCreation" {
		return errorMsg
	}
	return validationResult
}
