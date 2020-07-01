package utils

func ValueOrDefaultString(val string, defaultValue string) string {
	if len(val) == 0 {
		return defaultValue
	}

	return val
}

func ValueOrDefaultInt(val int, defaultValue int) int {
	if val == 0 {
		return defaultValue
	}

	return val
}
