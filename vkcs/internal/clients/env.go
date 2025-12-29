package clients

import (
	"os"
	"strings"
)

func getEnv(prefix, key string) string {
	key = strings.ToUpper(key)
	key = strings.ReplaceAll(key, "-", "_")
	return os.Getenv(prefix + key)
}
