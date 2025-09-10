package services

import (
	"context"
	"database/sql"
	"math"
	"strconv"
	"time"

	"foodflow/config"
	"foodflow/internal/core"
	"foodflow/internal/core/repos"
	"foodflow/internal/db"

	"github.com/jackc/pgx/v5"
)

type RedemptionService struct {
	db             *db.DB
	offerRepo      repos.OfferRepo
	claimRepo      repos.ClaimRepo
	creditRepo     repos.CreditRepo
	tokenRepo      repos.TokenRepo
	redemptionRepo repos.RedemptionRepo
	config         config.AppConfig
}

func NewRedemptionService(
	db *db.DB,
	offerRepo repos.OfferRepo,
	claimRepo repos.ClaimRepo,
	creditRepo repos.CreditRepo,
	tokenRepo repos.TokenRepo,
	redemptionRepo repos.RedemptionRepo,
	config config.AppConfig,
) *RedemptionService {
	return &RedemptionService{
		db:             db,
		offerRepo:      offerRepo,
		claimRepo:      claimRepo,
		creditRepo:     creditRepo,
		tokenRepo:      tokenRepo,
		redemptionRepo: redemptionRepo,
		config:         config,
	}
}

// RedeemOffer processes a redemption transaction
func (s *RedemptionService) RedeemOffer(ctx context.Context, offerID, orgID string, servingsAccepted int) (*core.RedemptionResponse, error) {
	var result *core.RedemptionResponse

	err := s.db.WithTransaction(ctx, func(tx pgx.Tx) error {
		// Get offer details
		offer, err := s.offerRepo.GetOfferByID(ctx, offerID)
		if err != nil {
			return err
		}

		// Validate offer status
		if offer.Status != core.DonationStatusPendingConfirm {
			return core.NewValidationError(core.ErrMsgInvalidOfferStatus, map[string]interface{}{
				"current_status":  offer.Status,
				"expected_status": core.DonationStatusPendingConfirm,
			})
		}

		// Validate servings
		if servingsAccepted <= 0 || servingsAccepted > offer.EstimatedServings {
			return core.NewValidationError(core.ErrMsgInvalidServings, map[string]interface{}{
				"requested": servingsAccepted,
				"available": offer.EstimatedServings,
			})
		}

		// Calculate credits and tokens
		creditsSpent := int(math.Ceil(float64(servingsAccepted) / float64(s.config.ServingsPerCredit)))
		tokensAwarded := creditsSpent * s.config.TokenMultiplier

		// Get current month
		currentMonth := time.Now().Format("200601")

		// Lock and update credits (SELECT FOR UPDATE)
		credit, err := s.creditRepo.GetCreditByOrgAndMonth(ctx, orgID, currentMonth)
		if err != nil {
			return err
		}

		if credit.CreditsRemaining < creditsSpent {
			return core.NewValidationError(core.ErrMsgInsufficientCredits, map[string]interface{}{
				"requested": creditsSpent,
				"available": credit.CreditsRemaining,
			})
		}

		// Update credits
		err = s.creditRepo.UpdateCreditRemaining(ctx, credit.ID.String(), credit.CreditsRemaining-creditsSpent)
		if err != nil {
			return err
		}

		// Update or create tokens for collaborator
		token, err := s.tokenRepo.GetTokenByCollabAndMonth(ctx, offer.CollabID.String(), currentMonth)
		if err != nil {
			if err == sql.ErrNoRows {
				// Create new token record
				_, err = s.tokenRepo.CreateToken(ctx, offer.CollabID.String(), currentMonth, tokensAwarded)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		} else {
			// Update existing token record
			err = s.tokenRepo.UpdateTokenEarned(ctx, token.ID.String(), token.TokensEarned+tokensAwarded)
			if err != nil {
				return err
			}
		}

		// Create redemption record
		redemption, err := s.redemptionRepo.CreateRedemption(ctx, offerID, orgID, offer.CollabID.String(), servingsAccepted, creditsSpent, tokensAwarded)
		if err != nil {
			return err
		}

		// Update offer status to CLAIMED
		_, err = s.offerRepo.UpdateOfferStatus(ctx, offerID, core.DonationStatusClaimed)
		if err != nil {
			return err
		}

		result = &core.RedemptionResponse{
			RedemptionID:  redemption.ID,
			CreditsSpent:  creditsSpent,
			TokensAwarded: tokensAwarded,
			OfferStatus:   string(core.DonationStatusClaimed),
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetRedemptionsByOrganization returns redemptions for an organization
func (s *RedemptionService) GetRedemptionsByOrganization(ctx context.Context, orgID string, limit, offset int) ([]*core.Redemption, error) {
	return s.redemptionRepo.GetRedemptionsByOrganization(ctx, orgID, limit, offset)
}

// GetRedemptionsByCollaborator returns redemptions for a collaborator
func (s *RedemptionService) GetRedemptionsByCollaborator(ctx context.Context, collabID string, limit, offset int) ([]*core.Redemption, error) {
	return s.redemptionRepo.GetRedemptionsByCollaborator(ctx, collabID, limit, offset)
}

// CreateRedemption creates a redemption (interface implementation)
func (s *RedemptionService) CreateRedemption(ctx context.Context, userID string, req *core.CreateRedemptionRequest) (*core.RedemptionResponse, error) {
	// TODO: Get organization ID from userID - for now we'll use userID as orgID
	// In a full implementation, you'd look up the organization associated with this user
	orgID := userID
	
	return s.RedeemOffer(ctx, req.OfferID.String(), orgID, req.ServingsAccepted)
}

// GetRedemption gets a redemption by ID
func (s *RedemptionService) GetRedemption(ctx context.Context, redemptionID string) (*core.RedemptionResponse, error) {
	redemption, err := s.redemptionRepo.GetRedemptionByID(ctx, redemptionID)
	if err != nil {
		return nil, err
	}

	return &core.RedemptionResponse{
		RedemptionID:  redemption.ID,
		CreditsSpent:  redemption.CreditsSpent,
		TokensAwarded: redemption.TokensAwarded,
		OfferStatus:   "CLAIMED", // Status would be derived from offer status
	}, nil
}

// UpdateRedemptionStatus updates redemption status
func (s *RedemptionService) UpdateRedemptionStatus(ctx context.Context, redemptionID string, req *core.UpdateRedemptionStatusRequest) (*core.RedemptionResponse, error) {
	// Get existing redemption
	redemption, err := s.redemptionRepo.GetRedemptionByID(ctx, redemptionID)
	if err != nil {
		return nil, err
	}

	// For now, redemptions don't have a status field in the entity
	// In a full implementation, you might add status tracking
	// This is a placeholder implementation

	return &core.RedemptionResponse{
		RedemptionID:  redemption.ID,
		CreditsSpent:  redemption.CreditsSpent,
		TokensAwarded: redemption.TokensAwarded,
		OfferStatus:   req.Status,
	}, nil
}

// ListRedemptions lists redemptions with pagination
func (s *RedemptionService) ListRedemptions(ctx context.Context, page, limit, status, orgID string) (*core.RedemptionListResponse, error) {
	// Parse pagination parameters
	pageNum := 1
	limitNum := 10
	if page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			pageNum = p
		}
	}
	if limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 && l <= 100 {
			limitNum = l
		}
	}
	offset := (pageNum - 1) * limitNum

	// Get redemptions based on filters
	var redemptions []*core.Redemption
	var err error

	if orgID != "" {
		redemptions, err = s.redemptionRepo.GetRedemptionsByOrganization(ctx, orgID, limitNum, offset)
	} else {
		redemptions, err = s.redemptionRepo.GetAllRedemptions(ctx, limitNum, offset)
	}

	if err != nil {
		return nil, core.NewInternalError("Failed to retrieve redemptions", err)
	}

	// Convert to response format
	redemptionResponses := make([]core.RedemptionResponse, len(redemptions))
	for i, redemption := range redemptions {
		redemptionResponses[i] = core.RedemptionResponse{
			RedemptionID:  redemption.ID,
			CreditsSpent:  redemption.CreditsSpent,
			TokensAwarded: redemption.TokensAwarded,
			OfferStatus:   "CLAIMED", // Default status
		}
	}

	return &core.RedemptionListResponse{
		Items: redemptionResponses,
		Page:  pageNum,
		Limit: limitNum,
		Total: len(redemptions), // In a real app, this would be a separate count query
	}, nil
}
