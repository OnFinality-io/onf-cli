package utils

func String(s string) *string {
	return &s
}

func Contains(arr []string, item string) bool {
	for _, a := range arr {
		if a == item {
			return true
		}
	}
	return false
}
