package redis

import "strings"

func extractIDFromKey(key string) string {
	split := strings.Split(key, "::")
	if len(split) > 1 {
		return split[0]
	}

	return ""
}
