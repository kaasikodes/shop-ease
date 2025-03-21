package env

import (
	"log"
	"os"
	"strconv"
)

func GetString(key string, fallback string) string {
    value, ok := os.LookupEnv(key)
    if !ok {
        return fallback
    }
    
    // Remove any surrounding quotes if present
    if len(value) >= 2 && value[0] == '"' && value[len(value)-1] == '"' {
        value = value[1 : len(value)-1]
    }
    
    log.Println("Clean value:", value)
    return value
}

func GetBool(key string, fallback bool) bool {
	_value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	value, err := strconv.ParseBool(_value)
	if err != nil {
		return fallback
	}
	return value
	
}
func GetInt(key string, fallback int) int {
	_value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	value, err := strconv.Atoi(_value)
	if err != nil {
		return fallback
	}
	return value
	
}