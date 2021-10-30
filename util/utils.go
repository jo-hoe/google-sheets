package util

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
