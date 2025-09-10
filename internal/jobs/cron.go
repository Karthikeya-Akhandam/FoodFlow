package jobs

import (
	"context"
	"time"

	"foodflow/internal/core/services"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
)

type CronManager struct {
	cron              *cron.Cron
	monthlyCreditsJob *MonthlyCreditsJob
	offerCleanupJob   *OfferCleanupJob
	autoAssignJob     *AutoAssignJob
}

func NewCronManager(
	monthlyCreditsService *services.MonthlyCreditsService,
	offerService *services.OfferService,
	claimService *services.ClaimService,
	matchingService *services.MatchingService,
) *CronManager {
	return &CronManager{
		cron:              cron.New(cron.WithLocation(time.UTC)),
		monthlyCreditsJob: NewMonthlyCreditsJob(monthlyCreditsService),
		offerCleanupJob:   NewOfferCleanupJob(offerService),
		autoAssignJob:     NewAutoAssignJob(offerService, claimService, matchingService),
	}
}

func (cm *CronManager) Start() error {
	// Monthly credits job - runs at 00:05 on the 1st of every month
	_, err := cm.cron.AddFunc("5 0 1 * *", func() {
		ctx := context.Background()
		log.Info().Msg("Starting monthly credits job")
		if err := cm.monthlyCreditsJob.Run(ctx); err != nil {
			log.Error().Err(err).Msg("Monthly credits job failed")
		} else {
			log.Info().Msg("Monthly credits job completed successfully")
		}
	})
	if err != nil {
		return err
	}

	// Offer cleanup job - runs every hour
	_, err = cm.cron.AddFunc("0 * * * *", func() {
		ctx := context.Background()
		log.Info().Msg("Starting offer cleanup job")
		if err := cm.offerCleanupJob.Run(ctx); err != nil {
			log.Error().Err(err).Msg("Offer cleanup job failed")
		} else {
			log.Info().Msg("Offer cleanup job completed successfully")
		}
	})
	if err != nil {
		return err
	}

	// Auto assign job - runs every 2 minutes
	_, err = cm.cron.AddFunc("*/2 * * * *", func() {
		ctx := context.Background()
		log.Debug().Msg("Starting auto assign job")
		if err := cm.autoAssignJob.Run(ctx); err != nil {
			log.Error().Err(err).Msg("Auto assign job failed")
		}
	})
	if err != nil {
		return err
	}

	cm.cron.Start()
	log.Info().Msg("Cron jobs started successfully")
	return nil
}

func (cm *CronManager) Stop() {
	cm.cron.Stop()
	log.Info().Msg("Cron jobs stopped")
}
