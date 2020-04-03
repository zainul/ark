package array

// InArrayInt checks if the given int is in the given array
func InArrayInt(arr []int, value int) bool {
	for i := range arr {
		if arr[i] == value {
			return true
		}
	}
	return false
}

// RemoveDuplicate value in array of string
func RemoveDuplicate(arr []string) []string {
	m := make(map[string]bool)

	for _, item := range arr {
		m[item] = true
	}

	result := make([]string, 0)

	for item := range m {
		result = append(result, item)
	}

	return result
}

// InArrayString checks if the given string is in the given array
func InArrayString(arr []string, value string) bool {
	for i := range arr {
		if arr[i] == value {
			return true
		}
	}
	return false
}
