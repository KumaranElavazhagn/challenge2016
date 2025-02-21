package parser

import (
	"encoding/csv"
	"os"
	"strings"
)

func ParseCSVFile(csvFilePath string) (map[string]map[string]map[string]bool, error) {
	localFilePath, err := os.Open(csvFilePath)
	if err != nil {
		return nil, err
	}
	defer localFilePath.Close()

	reader := csv.NewReader(localFilePath)

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	countries := make(map[string]map[string]map[string]bool)

	for i, record := range records {
		if i == 0 {
			//skip header row
			continue
		}
		countryName := strings.ToUpper(record[5])
		stateName := strings.ToUpper(record[4])
		cityName := strings.ToUpper(record[3])
		if _, countryExists := countries[countryName]; !countryExists {
			countries[countryName] = make(map[string]map[string]bool)
		}
		if _, stateExists := countries[countryName][stateName]; !stateExists {
			countries[countryName][stateName] = make(map[string]bool)
		}
		if _, cityExists := countries[countryName][stateName][cityName]; !cityExists {
			countries[countryName][stateName][cityName] = true
		}
	}
	return countries, nil
}
