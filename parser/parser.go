package parser

import (
	"encoding/csv"
	"os"
	"strings"
)

// ParseCSVFile reads a CSV file and structures the data into a nested map.
func ParseCSVFile(csvFilePath string) (map[string]map[string]map[string]bool, error) {
	// Open the CSV file for reading.
	localFilePath, err := os.Open(csvFilePath)
	if err != nil {
		return nil, err // Return error if file cannot be opened.
	}
	defer localFilePath.Close() // Ensure the file is closed after function execution.

	// Create a CSV reader to read the file contents.
	reader := csv.NewReader(localFilePath)

	// Read all records from the CSV file.
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err // Return error if reading fails.
	}

	// Initialize a nested map to store country, state, and city information.
	countries := make(map[string]map[string]map[string]bool)

	// Iterate over each record in the CSV file.
	for i, record := range records {
		if i == 0 {
			// Skip the header row in the CSV.
			continue
		}

		// Extract country, state, and city names from the respective columns.
		// Convert to uppercase for uniformity.
		countryName := strings.ToUpper(record[5]) // Country name (column index 5).
		stateName := strings.ToUpper(record[4])   // State name (column index 4).
		cityName := strings.ToUpper(record[3])    // City name (column index 3).

		// Check if the country exists in the map; if not, create an entry.
		if _, countryExists := countries[countryName]; !countryExists {
			countries[countryName] = make(map[string]map[string]bool)
		}

		// Check if the state exists under the country; if not, create an entry.
		if _, stateExists := countries[countryName][stateName]; !stateExists {
			countries[countryName][stateName] = make(map[string]bool)
		}

		// Check if the city exists under the respective state and country.
		if _, cityExists := countries[countryName][stateName][cityName]; !cityExists {
			// Add the city to the map and mark it as `true`.
			countries[countryName][stateName][cityName] = true
		}
	}

	// Return the structured data.
	return countries, nil
}
