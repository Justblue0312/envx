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
	if !strings.Contains(key, "__") {
		return ""
	}

	parts := strings.Split(key, "__")
	for i := len(parts); i > 1; i-- {
		testKey := strings.Join(parts[:i], "__")
		if value, ok := lookupEnv(testKey); ok {
			return value
		}
	}

	return ""
}
