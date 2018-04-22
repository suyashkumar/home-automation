package config

import "os"

func Get(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	defValue, ok := defaults[key]
	if !ok {
		return ""
	}

	return defValue
}

