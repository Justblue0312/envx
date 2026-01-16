//go:build windows

package envx

import (
	"os"
	"strings"
)

func lookupEnv(key string) (string, bool) {
	for _, env := range os.Environ() {
		kv := strings.SplitN(env, "=", 2)
		if len(kv) == 2 && strings.EqualFold(kv[0], key) {
			return kv[1], true
		}
	}
	return "", false
}

func tryNestedKeys(key string) string {
	if !strings.Contains(key, "__") {
		return ""
	}

	parts := strings.Split(key, "__")
	for i := len(parts); i > 1; i-- {
		testKey := strings.Join(parts[:i], "__")
		if value := tryWindowsLookup(testKey); value != "" {
			return value
		}
	}

	return ""
}

func tryWindowsLookup(key string) string {
	for _, env := range os.Environ() {
		kv := strings.SplitN(env, "=", 2)
		if len(kv) == 2 && strings.EqualFold(kv[0], key) {
			return kv[1]
		}
	}
	return ""
}
