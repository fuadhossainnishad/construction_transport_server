package utils

import "strconv"

func StringToInt(s string, fallback int) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		return fallback
	}
	return val
}
