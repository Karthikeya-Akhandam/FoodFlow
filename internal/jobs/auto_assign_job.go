package jobs

import (
	"context"

	"foodflow/internal/core/services"

	"github.com/rs/zerolog/log"
)

type AutoAssignJob struct {
	offerService    *services.OfferService
	claimService    *services.ClaimService
	matchingService *services.MatchingService
}

func NewAutoAssignJob(
	offerService *services.OfferService,
	claimService *services.ClaimService,
	matchingService *services.MatchingService,
) *AutoAssignJob {
	return &AutoAssignJob{
		offerService:    offerService,
		claimService:    claimService,
		matchingService: matchingService,
	}
}

func (j *AutoAssignJob) Run(ctx context.Context) error {
	// Get offers that are ready for auto-assignment
	offers, err := j.offerService.GetOffersForAutoAssignment(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get offers for auto-assignment")
		return err
	}

	if len(offers) == 0 {
		return nil // No offers to process
	}

	log.Debug().Int("count", len(offers)).Msg("Processing offers for auto-assignment")

	assignedCount := 0
	for _, offer := range offers {
		// Get claims for this offer
		claims, err := j.claimService.GetClaimsForAutoAssignment(ctx)
		if err != nil {
			log.Error().Err(err).Str("offer_id", offer.ID.String()).Msg("Failed to get claims for auto-assignment")
			continue
		}

		if len(claims) == 0 {
			continue // No claims to process
		}

		// Auto-assign the highest scoring claim
		err = j.offerService.AutoAssignOffer(ctx, offer.ID.String())
		if err != nil {
			log.Error().Err(err).
				Str("offer_id", offer.ID.String()).
				Str("claim_id", claims[0].ID.String()).
				Msg("Failed to auto-assign offer")
			continue
		}

		assignedCount++
		log.Debug().
			Str("offer_id", offer.ID.String()).
			Str("claim_id", claims[0].ID.String()).
			Float64("score", claims[0].PriorityScore).
			Msg("Auto-assigned offer")
	}

	if assignedCount > 0 {
		log.Info().Int("assigned_count", assignedCount).Msg("Auto-assignment job completed")
	}

	return nil
}
