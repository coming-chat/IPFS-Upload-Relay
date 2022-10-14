package inits

import (
	"context"
	"fmt"
	"github.com/NaturalSelectionLabs/IPFS-Upload-Relay/global"
	"github.com/go-redis/redis/v9"
	"os"
	"time"
)

func Redis() error {
	// Get conn string
	redisConnString, exist := os.LookupEnv("REDIS_CONNECTION_STRING")
	if !exist {
		return fmt.Errorf("env virable REDIS_CONNECTION_STRING not found")
	}

	// Parse connect string
	redisConfig, err := redis.ParseURL(redisConnString)
	if err != nil {
		return fmt.Errorf("failed to parse redis connection string: %v", err)
	}

	// Connect to server
	global.Redis = redis.NewClient(redisConfig)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Try connection
	err = global.Redis.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to connect to redis: %v", err)
	}

	return nil
}
