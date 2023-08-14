package pkg

import (
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/garsue/watermillzap"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func NewSubscriber(redisClient *redis.Client, logger *zap.Logger) (message.Subscriber, func(), error) {
	consumerID := uuid.New().String()
	subscriber, err := redisstream.NewSubscriber(
		redisstream.SubscriberConfig{
			Client:        redisClient,
			Unmarshaller:  redisstream.DefaultMarshallerUnmarshaller{},
			ConsumerGroup: "shopify-spider",
			Consumer:      consumerID,
		},
		watermillzap.NewLogger(logger.Named("subscriber")),
	)

	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {}

	return subscriber, cleanup, nil
}
