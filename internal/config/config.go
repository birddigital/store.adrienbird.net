package config

import (
	"os"
	"strconv"
)

type Config struct {
	Server      ServerConfig      `json:"server"`
	Squarespace SquarespaceConfig `json:"squarespace"`
}

type ServerConfig struct {
	Port            int  `json:"port"`
	Mode            string `json:"mode"`
	EnableSwagger   bool `json:"enable_swagger"`
	EnableHealth    bool `json:"enable_health"`
	EnableMetrics   bool `json:"enable_metrics"`
	EnableTracing   bool `json:"enable_tracing"`
}

type SquarespaceConfig struct {
	BaseURL     string `json:"base_url"`
	SiteID      string `json:"site_id"`
	APIKey      string `json:"api_key"`
	AccessToken string `json:"access_token"`
	Environment string `json:"environment"`
}

func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port:            getEnvAsInt("PORT", 8080),
			Mode:            getEnv("GIN_MODE", "debug"),
			EnableSwagger:   getEnvAsBool("ENABLE_SWAGGER", true),
			EnableHealth:    getEnvAsBool("ENABLE_HEALTH", true),
			EnableMetrics:   getEnvAsBool("ENABLE_METRICS", false),
			EnableTracing:   getEnvAsBool("ENABLE_TRACING", false),
		},
		Squarespace: SquarespaceConfig{
			BaseURL:     getEnv("SQUARESPACE_BASE_URL", "https://api.squarespace.com"),
			SiteID:      os.Getenv("SQUARESPACE_SITE_ID"),
			APIKey:      os.Getenv("SQUARESPACE_API_KEY"),
			AccessToken: os.Getenv("SQUARESPACE_ACCESS_TOKEN"),
			Environment: getEnv("NODE_ENV", "development"),
		},
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}