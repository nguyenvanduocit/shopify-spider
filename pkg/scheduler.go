package pkg

import (
	"github.com/go-co-op/gocron"
	"go.uber.org/zap"
	"time"
)

func NewScheduler(logger *zap.Logger) (*gocron.Scheduler, func(), error) {
	scheduler := gocron.NewScheduler(time.UTC)
	cleanup := func() {
		logger.Info("stopping scheduler")
		scheduler.StopBlockingChan()
		logger.Info("scheduler stopped")
	}
	logger.Info("scheduler started")

	return scheduler, cleanup, nil
}
