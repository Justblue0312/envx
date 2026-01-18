//go:build !windows

package envx

import (
	"os"
	"strings"
)

func lookupEnv(key string) (string, bool) {
	return os.LookupEnv(key)
}

func tryNestedKeys(key string) string {
	if !strings.Contains(key, "_") {
		return ""
	}

	parts := strings.Split(key, "_")
	for i := len(parts); i > 1; i-- {
		testKey := strings.Join(parts[:i], "_")
		if value, ok := lookupEnv(testKey); ok {
			return value
		}
	}

	return ""
}
