package main

import (
	"context"
	"log"

	"github.com/osamikoyo/dark-fantasy-land/internal/config"
	"github.com/osamikoyo/dark-fantasy-land/pkg/logger"
	"github.com/osamikoyo/dark-fantasy-land/pkg/retrier"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func main() {
	cfg := config.NewConfig()

	logCfg := logger.Config{
		LogFile:   "app.log",
		LogLevel:  "debug",
		AppName:   "dark-fantasy-land",
		AddCaller: true,
	}

	if err := logger.Init(logCfg); err != nil {
		log.Fatal(err)

		return
	}

	logger := logger.Get()

	logger.Info("starting dark-fantasy land...", zap.Any("cfg", cfg))

	db, err := retrier.Connect(3, 5, func() (*mongo.Database, error) {
		client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.MongoUrl))
		if err != nil {
			return nil, err
		}

		return client.Database("dark-fantasy"), nil
	})
	if err != nil {
		logger.Error("failed connecto to mongodb", zap.Error(err))

		return
	}

	logger.Info("successfully connected to mongo db")

	redisDB, err := retrier.Connect(3, 5, func() (*redis.Client, error) {
		client := redis.NewClient(&redis.Options{
			Addr:     cfg.RedisUrl,
			Password: "",
			DB:       0,
		})

		return client, client.Ping(context.Background()).Err()
	})
	if err != nil {
		logger.Error("failed connect to redis", zap.Error(err))

		return
	}

	logger.Info("successfully connected to redis")
}
