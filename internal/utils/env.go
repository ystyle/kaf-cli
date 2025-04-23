package utils

import (
	"os"
)

func GetEnv(key, defaultvalue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultvalue
}
