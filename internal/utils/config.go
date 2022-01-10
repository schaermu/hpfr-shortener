package utils

type Config struct {
	MongoDSN string `mapstructure:"MONGO_DSN"`
	MongoDB  string `mapstructure:"MONGO_DB"`
	BaseURL  string `mapstructure:"BASE_URL"`
}