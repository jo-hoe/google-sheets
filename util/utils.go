package util

import "fmt"

func CSVSlicesToMap(csvData [][]string) map[string][]string {
	result := make(map[string][]string)
	if len(csvData) == 0 {
		return result
	}

	header := csvData[0]
	for i := 1; i < len(csvData); i++ {
		value := make([]string, 0)
		value = append(value, csvData[i]...)
		result[header[i-1]] = value
	}
	return result
}

func GetMaxLengthOfSlices(item map[string][]string) int {
	max := 0

	for _, v := range item {
		length := len(v)
		if max < length {
			max = length
		}
	}

	return max
}

// Checks if key are in map. If all key are found the returned error will be nil
// otherwise an error will be returned
func ValidatedKeys(items map[string][]string, keysToValidate ...string) error {
	missingKeys := make([]string, 0)

	for _, key := range keysToValidate {
		if _, ok := items[key]; !ok {
			missingKeys = append(missingKeys, key)
		}
	}

	if len(missingKeys) != 0 {
		return fmt.Errorf("missing the following keys '%v'", missingKeys)
	}
	return nil
}
