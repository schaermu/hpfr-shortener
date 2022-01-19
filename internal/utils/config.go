package utils

import (
	"fmt"
	"os"
)

type Config struct {
	MongoDSN string
	MongoDB  string
	BaseURL  string
}

var (
	requiredEnvVars = []string{"MONGO_DSN", "MONGO_DB"}
)

func NewConfigFromEnv() (cfg *Config, err error) {
	for _, envVar := range requiredEnvVars {
		if _, found := os.LookupEnv(envVar); !found {
			return nil, fmt.Errorf("env var %q not set", envVar)
		}
	}

	return &Config{
		MongoDSN: os.Getenv("MONGO_DSN"),
		MongoDB:  os.Getenv("MONGO_DB"),
		BaseURL:  os.Getenv("BASE_URL"),
	}, nil
}
