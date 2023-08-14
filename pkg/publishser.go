package pkg

import (
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/garsue/watermillzap"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func NewPublisher(redisClient *redis.Client, logger *zap.Logger) (message.Publisher, func(), error) {
	publisher, err := redisstream.NewPublisher(
		redisstream.PublisherConfig{
			Client: redisClient,
		},
		watermillzap.NewLogger(logger.Named("publisher")),
	)

	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		if err := publisher.Close(); err != nil {
			logger.Error("Router: error closing router", zap.Error(err))
			return
		}
	}

	return publisher, cleanup, nil
}
