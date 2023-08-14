package main

import (
	"context"
	"encoding/json"
	"github.com/alitto/pond"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"shopifyspider/pkg"
	"time"
)

func main() {
	// Logger
	logger, _ := zap.NewProduction()

	// database
	db, dbCleanup, err := pkg.NewMongoDb()
	if err != nil {
		logger.Fatal("failed to initialize database", zap.Error(err))
	}
	defer dbCleanup()

	// fiber
	fiberApp := fiber.New(fiber.Config{})
	defer func() {
		logger.Info("stopping fiber")
		if err := fiberApp.Shutdown(); err != nil {
			logger.Fatal("failed to shut down fiber", zap.Error(err))
		}
		logger.Info("fiber is stopped")
	}()

	fiberApp.Get("/healthz", func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})

	go func() {
		logger.Info("fiber is running")
		if err := fiberApp.Listen(":" + os.Getenv("PORT")); err != nil {
			logger.Fatal("failed to start fiber", zap.Error(err))
		}
		logger.Info("fiber is stopped")
	}()

	logger.Info("pubsub is running")

	redisClient, cleanUpRedisClient, err := pkg.NewRedisClient(logger)
	if err != nil {
		logger.Fatal("failed to initialize redis client", zap.Error(err))
	}
	defer cleanUpRedisClient()

	subscriber, cleanupSubscriber, err := pkg.NewSubscriber(redisClient, logger)
	if err != nil {
		logger.Fatal("failed to initialize subscriber", zap.Error(err))
	}
	defer cleanupSubscriber()

	if messages, err := subscriber.Subscribe(context.Background(), pkg.EvtSiteMapAppFound); err != nil {
		logger.Fatal("failed to subscribe to apps. Updated", zap.Error(err))
	} else {
		go func() {
			appCollection := db.Database("shopify").Collection("apps")
			for msg := range messages {
				logger.Info("apps.updated message received", zap.String("message", string(msg.Payload)))

				entry := pkg.SitemapEntry{}

				if err := json.Unmarshal(msg.Payload, &entry); err != nil {
					logger.Error("failed to unmarshal app", zap.Error(err))
					msg.Nack()
					continue
				}

				if entry.ParsedLastModified.IsZero() {
					lastCrawl := entry.ParsedLastModified.AddDate(-1, 0, 0)
					entry.ParsedLastModified = &lastCrawl
				}

				app, err := pkg.GetAppByUrl(appCollection, entry.Location)
				if err != nil {
					logger.Error("failed to get app by url", zap.Error(err))
					msg.Nack()
					continue
				}

				if app == nil {
					lastCrawl := entry.ParsedLastModified.AddDate(-1, 0, 0)
					app = &pkg.App{
						Url:         entry.Location,
						LastUpdated: entry.ParsedLastModified,
						LastCrawl:   &lastCrawl,
					}
					app.ClientId, err = pkg.GetAppClientID(entry.Location)
					if err != nil {
						logger.Error("failed to get app client id", zap.Error(err))
						msg.Nack()
						continue
					}

					if err := pkg.CreateApplication(appCollection, app); err != nil {
						logger.Error("failed to create app", zap.Error(err))
						msg.Nack()
						continue
					}

					logger.Info("app created", zap.String("url", app.Url))
					continue
				}
				time.Sleep(5 * time.Second)
				msg.Ack()
			}
		}()
	}

	// scheduler
	scheduler, cleanupScheduler, err := pkg.NewScheduler(logger)
	if err != nil {
		logger.Fatal("failed to initialize scheduler", zap.Error(err))
	}
	defer cleanupScheduler()
	scheduler.StartAsync()

	publisher, cleanupPublisher, err := pkg.NewPublisher(redisClient, logger)
	if err != nil {
		logger.Fatal("failed to initialize publisher", zap.Error(err))
	}
	defer cleanupPublisher()

	scheduler.Every(10).Minute().SingletonMode().Do(func() {
		pkg.WalkSitemap(publisher)
	})
	// Worker pool
	spiderPool := pond.New(1, 100)
	defer spiderPool.StopAndWait()

	signChan := make(chan os.Signal, 1)
	signal.Notify(signChan, os.Interrupt, os.Kill)

	<-signChan
	logger.Info("shutting down")
}
