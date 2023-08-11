package pkg

import (
	"strings"
)

func GetUrlType(url string) string {
	parts := strings.Split(url, "/")
	if len(parts) == 4 {
		return "apps"
	}

	return parts[3]
}
