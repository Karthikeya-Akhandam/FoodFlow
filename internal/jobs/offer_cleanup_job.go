package jobs

import (
	"context"
	"time"

	"foodflow/internal/core/services"

	"github.com/rs/zerolog/log"
)

type OfferCleanupJob struct {
	offerService *services.OfferService
}

func NewOfferCleanupJob(offerService *services.OfferService) *OfferCleanupJob {
	return &OfferCleanupJob{
		offerService: offerService,
	}
}

func (j *OfferCleanupJob) Run(ctx context.Context) error {
	startTime := time.Now()
	log.Info().Msg("Starting offer cleanup job")

	// Clean up expired offers
	err := j.offerService.CleanupExpiredOffers(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to cleanup expired offers")
		return err
	}

	// Clean up stale claims
	err = j.offerService.CleanupStaleClaims(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to cleanup stale claims")
		return err
	}

	duration := time.Since(startTime)
	log.Info().
		Dur("duration", duration).
		Msg("Offer cleanup job completed successfully")

	return nil
}
