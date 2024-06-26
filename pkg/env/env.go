package env

import (
	"fmt"
	"os"
)

func MustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("Env variable %s is not set", key))
	}
	return value
}