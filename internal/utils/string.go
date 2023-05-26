package utils

func RemoveDuplicates(arr []string) []string {
	uniqueArr := map[string]struct{}{}
	result := []string{}

	for _, str := range arr {
		if _, exists := uniqueArr[str]; !exists {
			result = append(result, str)
			uniqueArr[str] = struct{}{}
		}
	}

	return result
}
