package utils

import "strings"

func ContainsAll(str string, substr ...string) bool {
	for _, s := range substr {
		if !strings.Contains(str, s) {
			return false
		}
	}
	return true
}
