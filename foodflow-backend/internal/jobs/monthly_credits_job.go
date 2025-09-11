package jobs

import (
	"context"
	"time"

	"foodflow/internal/core/services"

	"github.com/rs/zerolog/log"
)

type MonthlyCreditsJob struct {
	monthlyCreditsService *services.MonthlyCreditsService
}

func NewMonthlyCreditsJob(monthlyCreditsService *services.MonthlyCreditsService) *MonthlyCreditsJob {
	return &MonthlyCreditsJob{
		monthlyCreditsService: monthlyCreditsService,
	}
}

func (j *MonthlyCreditsJob) Run(ctx context.Context) error {
	startTime := time.Now()
	log.Info().Msg("Starting monthly credits issuance job")

	// Run the monthly credits service
	err := j.monthlyCreditsService.IssueMonthlyCredits(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to issue monthly credits")
		return err
	}

	duration := time.Since(startTime)
	log.Info().
		Dur("duration", duration).
		Msg("Monthly credits issuance job completed successfully")

	return nil
}
