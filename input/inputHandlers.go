package input

import (
	"fmt"
	"qube-challenge-2016/dto"
	"strings"

	"github.com/manifoldco/promptui"
)

// PromptMenu displays a menu for selecting different choices and returns the selected option.
func PromptMenu() string {
	fmt.Println("Please specify the regions you wish to include or exclude for this distributor (use hyphens for specifying location hierarchy, e.g., Chennai-Tamil Nadu-India, Karnataka-India)")

	// Define the selection menu with available choices
	prompt := promptui.Select{
		Label: "Select one of the below choices",
		Items: []string{
			"Create a new distributor",
			"Create a sub distributor",
			"Check permission for a distributor",
			"View Distributors information",
			"Exit",
		},
	}

	// Run the prompt and get user selection
	_, result, err := prompt.Run()
	if err != nil {
		fmt.Println("Error in prompt selection. Please try again.")
		return ""
	}

	return result
}

// PromptDistributorData prompts the user for distributor details and stores Include/Exclude regions as maps.
func PromptDistributorData(subDistributor bool) dto.Distributor {
	var distributor dto.Distributor

	// Prompt user to enter distributor name
	promptName := promptui.Prompt{
		Label:       "Enter distributor name:",
		HideEntered: true,
	}
	name, _ := promptName.Run()
	distributor.Name = name
	fmt.Println(promptName.Label, name)

	// Prompt user to enter include regions
	promptInclude := promptui.Prompt{
		Label:       "Enter the regions you want to include for this distributor (comma separated):",
		HideEntered: true,
	}
	includeInput, _ := promptInclude.Run()
	distributor.Include = convertToMap(includeInput)
	fmt.Println(promptInclude.Label, includeInput)

	// Prompt user to enter exclude regions
	promptExclude := promptui.Prompt{
		Label:       "Enter the regions you want to exclude for this distributor (comma separated):",
		HideEntered: true,
	}
	excludeInput, _ := promptExclude.Run()
	distributor.Exclude = convertToMap(excludeInput)
	fmt.Println(promptExclude.Label, excludeInput)

	// If creating a sub-distributor, prompt for parent distributor name
	if subDistributor {
		promptParent := promptui.Prompt{
			Label:       "Enter the name of the parent distributor:",
			HideEntered: true,
		}
		parent, _ := promptParent.Run()
		distributor.Parent = parent
		fmt.Println(promptParent.Label, parent)
	}

	return distributor
}

// PromptCheckPermissionData prompts the user for distributor name and regions.
func PromptCheckPermissionData() dto.CheckPermissionData {
	var data dto.CheckPermissionData

	// Prompt user to enter distributor name to check permission
	promptName := promptui.Prompt{
		Label:       "Enter distributor name that needs to be checked:",
		HideEntered: true,
	}
	data.DistributorName, _ = promptName.Run()
	fmt.Println(promptName.Label, data.DistributorName)

	// Prompt user to enter regions for permission check
	promptRegions := promptui.Prompt{
		Label:       "Enter regions that need to be checked (comma separated):",
		HideEntered: true,
	}
	regionsInput, _ := promptRegions.Run()
	data.Regions = convertToMap(regionsInput)
	fmt.Println(promptRegions.Label, regionsInput)

	return data
}

// convertToMap is a helper function to convert a comma-separated string into a map[string]bool.
func convertToMap(input string) map[string]bool {
	result := make(map[string]bool)

	// Split input string by commas and trim spaces
	for _, region := range strings.Split(input, ",") {
		trimmedRegion := strings.TrimSpace(region)
		if trimmedRegion != "" {
			result[trimmedRegion] = true
		}
	}
	return result
}
