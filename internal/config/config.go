package config

import "os"

type (
	Buckets struct {
		WallpaperFull  string
		WallpaperWatch string
		Mems           string
	}

	Config struct {
		Port           string
		Host           string
		MongoUrl       string
		RedisUrl       string
		MinioUrl       string
		MinioAccessKey string
		MinioSecretKey string
		MinioBuckets   Buckets
		MinioSSL       bool
	}
)

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
		Port:           "8080",
		Host:           "localhost",
		MongoUrl:       url,
		RedisUrl:       redisUrl,
		MinioUrl:       "localhost:9000",
		MinioAccessKey: "minioadmin",
		MinioSecretKey: "minioadmin",
		MinioBuckets: Buckets{
			WallpaperFull:  "wallpaper_full",
			WallpaperWatch: "wallpaper_watch",
			Mems:           "mem",
		},
		MinioSSL: false,
	}
}
