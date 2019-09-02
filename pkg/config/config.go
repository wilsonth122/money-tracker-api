package config

import (
	"os"
	"strconv"
	"strings"
)

type APIConfig struct {
	Port           string
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

type DatabaseConfig struct {
	Addresses         []string
	Username          string
	Password          string
	AdminDatabase     string
	AppDatabase       string
	UserCollection    string
	ExpenseCollection string
}

type AuthConfig struct {
	TokenPassword string
}

type Config struct {
	API      APIConfig
	Database DatabaseConfig
	Auth     AuthConfig
}

// New returns a new Config struct
func New() *Config {
	return &Config{
		API: APIConfig{
			Port:           getEnv("PORT", "8080"),
			AllowedOrigins: getEnvAsSlice("ALLOWED_ORIGINS", []string{""}, ","),
			AllowedMethods: getEnvAsSlice("ALLOWED_METHODS", []string{""}, ","),
			AllowedHeaders: getEnvAsSlice("ALLOWED_HEADERS", []string{""}, ","),
		},
		Database: DatabaseConfig{
			Addresses:         getEnvAsSlice("DATABASE_ADDRESSES", []string{""}, ","),
			Username:          getEnv("DATABASE_USERNAME", ""),
			Password:          getEnv("DATABASE_PASSWORD", ""),
			AdminDatabase:     getEnv("ADMIN_DATABASE", ""),
			AppDatabase:       getEnv("APP_DATABASE", ""),
			UserCollection:    getEnv("USER_COLLECTION", ""),
			ExpenseCollection: getEnv("EXPENSE_COLLECTION", ""),
		},
		Auth: AuthConfig{
			TokenPassword: getEnv("TOKEN_PASSWORD", ""),
		},
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

// Simple helper function to read an environment variable into integer or return a default value
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

// Helper to read an environment variable into a bool or return default value
func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}

// Helper to read an environment variable into a string slice or return default value
func getEnvAsSlice(name string, defaultVal []string, sep string) []string {
	valStr := getEnv(name, "")

	if valStr == "" {
		return defaultVal
	}

	val := strings.Split(valStr, sep)

	return val
}
