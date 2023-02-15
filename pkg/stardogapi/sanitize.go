package stardogapi

import "strings"

// Sanitizes path values by removing slashes
func sanitizePathValue(value string) string {
	replacer := strings.NewReplacer("%2F", "", "/", "")

	return replacer.Replace(value)
}
