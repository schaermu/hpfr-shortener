package utils

import "os"

type Config struct {
	MongoDSN string
	MongoDB  string
}

func NewConfigFromEnv() Config {
	return Config{
		MongoDSN: os.Getenv("MONGO_DSN"),
		MongoDB:  os.Getenv("MONGO_DB"),
	}
}
