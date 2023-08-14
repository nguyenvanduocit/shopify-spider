package pkg

import (
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"os"
)

func RedisConfigFromEnv() (*redis.Options, error) {
	addr, ok := os.LookupEnv("REDISHOST")
	if !ok {
		return nil, errors.New("REDISHOST is not set")
	}

	port, ok := os.LookupEnv("REDISPORT")
	if !ok {
		return nil, errors.New("REDISPORT is not set")
	}

	username, ok := os.LookupEnv("REDISUSER")
	if !ok {
		return nil, errors.New("REDISUSER is not set")
	}

	password, ok := os.LookupEnv("REDISPASSWORD")
	if !ok {
		return nil, errors.New("REDISPASSWORD is not set")
	}

	return &redis.Options{
		Addr:     addr + ":" + port,
		Username: username,
		Password: password,
		DB:       0,
	}, nil
}

func NewRedisClient(
	logSvc *zap.Logger,
) (*redis.Client, func(), error) {
	logger := logSvc.With(zap.Strings("tags", []string{"redis-client"}))

	config, err := RedisConfigFromEnv()
	if err != nil {
		return nil, nil, err
	}

	redisClient := redis.NewClient(config)

	cleanup := func() {
		logger.Info("Router: Cleaning up")
		if err := redisClient.Close(); err != nil {
			if errors.Is(err, redis.ErrClosed) {
				logger.Info("Router: router already closed")
				return
			}
			logger.Error("Router: error closing router", zap.Error(err))
			return
		}
		logger.Info("Router: router closed")
	}

	return redisClient, cleanup, nil
}
