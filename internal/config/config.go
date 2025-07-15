package config

import "os"

type Config struct {
	Port     string
	Host     string
	MongoUrl string
	RedisUrl string
	MinioUrl string
	MinioAccessKey string
	MinioSecretKey string
	MinioBuckets [2]string
	MinioSSL bool
}

func NewConfig() *Config {
	url := os.Getenv("MONGO_URI")
	if url == "" {
		url = "mongodb://mongo:27017"
	}

	redisUrl := os.Getenv("REDIS_URI")
	if redisUrl == "" {
		redisUrl = "redis:6379"
	}

	return &Config{
		Port:     "8080",
		Host:     "localhost",
		MongoUrl: url,
		RedisUrl: redisUrl,
		MinioUrl: "localhost:9000",
		MinioAccessKey: "minioadmin",
		MinioSecretKey: "minioadmin",
		MinioBuckets: [2]string{
			"wallpapers",
			"mems",
		},
		MinioSSL: false,
	}
}
