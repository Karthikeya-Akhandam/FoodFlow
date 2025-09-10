package services

import (
	"context"
	"strconv"
	"time"

	"foodflow/internal/core"
	"foodflow/internal/core/repos"
)

type ClaimService struct {
	claimRepo   repos.ClaimRepo
	offerRepo   repos.OfferRepo
	orgRepo     repos.OrgRepo
	creditRepo  repos.CreditRepo
	profileRepo repos.ProfileRepo
}

func NewClaimService(
	claimRepo repos.ClaimRepo,
	offerRepo repos.OfferRepo,
	orgRepo repos.OrgRepo,
	creditRepo repos.CreditRepo,
	profileRepo repos.ProfileRepo,
) core.ClaimService {
	return &ClaimService{
		claimRepo:   claimRepo,
		offerRepo:   offerRepo,
		orgRepo:     orgRepo,
		creditRepo:  creditRepo,
		profileRepo: profileRepo,
	}
}

func (s *ClaimService) CreateClaim(ctx context.Context, userID, offerID string, req *core.CreateClaimRequest) (*core.ClaimResponse, error) {
	// Get organization by user ID
	org, err := s.orgRepo.GetOrganizationByUserID(ctx, userID)
	if err != nil {
		return nil, core.NewNotFoundError("Organization not found for user")
	}

	// Get offer to validate it exists and is claimable
	offer, err := s.offerRepo.GetOfferByID(ctx, offerID)
	if err != nil {
		return nil, err
	}

	// Validate offer is open
	if offer.Status != core.DonationStatusOpen {
		return nil, core.NewValidationError("Offer is not available for claims", nil)
	}

	// Validate offer hasn't expired
	if offer.ExpiresAt.Before(time.Now()) {
		return nil, core.NewValidationError(core.ErrMsgOfferExpired, nil)
	}

	// Check if organization has sufficient credits (if specified)
	if req.RequestedServings > 0 {
		currentMonth := time.Now().Format("200601")
		credit, err := s.creditRepo.GetCreditByOrgAndMonth(ctx, org.ID.String(), currentMonth)
		if err != nil {
			return nil, core.NewValidationError("No credits available for this month", nil)
		}

		if credit.CreditsRemaining < req.RequestedServings {
			return nil, core.NewValidationError(core.ErrMsgInsufficientCredits, map[string]interface{}{
				"requested": req.RequestedServings,
				"available": credit.CreditsRemaining,
			})
		}
	}

	// Create claim with initial priority score (will be calculated by matching service)
	claim, err := s.claimRepo.CreateClaim(
		ctx,
		org.ID.String(),
		offerID,
		req.RequestedServings,
		core.ClaimStatusRequested,
	)
	if err != nil {
		return nil, core.NewInternalError("Failed to create claim", err)
	}

	return &core.ClaimResponse{
		ClaimID:           claim.ID,
		OfferID:           claim.OfferID,
		Status:            claim.Status,
		PriorityScore:     claim.PriorityScore,
		RequestedServings: claim.RequestedServings,
		CreatedAt:         claim.CreatedAt,
	}, nil
}

func (s *ClaimService) GetClaim(ctx context.Context, claimID string) (*core.ClaimResponse, error) {
	claim, err := s.claimRepo.GetClaimByID(ctx, claimID)
	if err != nil {
		return nil, err
	}

	return &core.ClaimResponse{
		ClaimID:           claim.ID,
		OfferID:           claim.OfferID,
		Status:            claim.Status,
		PriorityScore:     claim.PriorityScore,
		RequestedServings: claim.RequestedServings,
		CreatedAt:         claim.CreatedAt,
	}, nil
}

func (s *ClaimService) UpdateClaimStatus(ctx context.Context, claimID string, req *core.UpdateClaimStatusRequest) (*core.ClaimResponse, error) {
	// Get existing claim to validate
	claim, err := s.claimRepo.GetClaimByID(ctx, claimID)
	if err != nil {
		return nil, err
	}

	// Validate status transition
	validTransitions := map[core.ClaimStatus][]core.ClaimStatus{
		core.ClaimStatusRequested: {core.ClaimStatusWon, core.ClaimStatusLost, core.ClaimStatusCancelled},
		core.ClaimStatusWon:       {core.ClaimStatusCancelled}, // Won claims can only be cancelled
		core.ClaimStatusLost:      {},                         // No transitions allowed
		core.ClaimStatusCancelled: {},                         // No transitions allowed
	}

	validNextStatuses, exists := validTransitions[claim.Status]
	if !exists {
		return nil, core.NewValidationError("Invalid current status", nil)
	}

	isValidTransition := false
	for _, validStatus := range validNextStatuses {
		if validStatus == req.Status {
			isValidTransition = true
			break
		}
	}

	if !isValidTransition {
		return nil, core.NewValidationError("Invalid status transition", map[string]interface{}{
			"current": claim.Status,
			"requested": req.Status,
		})
	}

	// Update claim status
	updatedClaim, err := s.claimRepo.UpdateClaimStatus(ctx, claimID, req.Status)
	if err != nil {
		return nil, core.NewInternalError("Failed to update claim status", err)
	}

	// If claim is won, award credits for next month
	if req.Status == core.ClaimStatusWon {
		// Get organization for credit award
		org, err := s.orgRepo.GetOrganizationByID(ctx, claim.OrgID.String())
		if err == nil {
			// Award credits for next month based on servings claimed
			nextMonth := time.Now().AddDate(0, 1, 0).Format("200601")
			// This would need an AwardCredits method in the credit repo
			// For now, we'll leave this as a placeholder
			_ = nextMonth
			_ = org
		}
	}

	return &core.ClaimResponse{
		ClaimID:           updatedClaim.ID,
		OfferID:           updatedClaim.OfferID,
		Status:            updatedClaim.Status,
		PriorityScore:     updatedClaim.PriorityScore,
		RequestedServings: updatedClaim.RequestedServings,
		CreatedAt:         updatedClaim.CreatedAt,
	}, nil
}

func (s *ClaimService) ListClaims(ctx context.Context, page, limit, status, orgID string) (*core.ClaimListResponse, error) {
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

	// Get claims based on filters
	var claims []*core.DonationClaim
	var err error

	if orgID != "" {
		claims, err = s.claimRepo.GetClaimsByOrganization(ctx, orgID, limitNum, offset)
	} else {
		// Parse status string to ClaimStatus
		var claimStatus core.ClaimStatus
		if status != "" {
			claimStatus = core.ClaimStatus(status)
		} else {
			claimStatus = core.ClaimStatusRequested // Default
		}
		claims, err = s.claimRepo.GetClaimsByStatus(ctx, claimStatus, limitNum, offset)
	}
	
	total := len(claims) // For simplicity, using slice length as total

	if err != nil {
		return nil, core.NewInternalError("Failed to retrieve claims", err)
	}

	// Convert to response format
	claimResponses := make([]core.ClaimResponse, len(claims))
	for i, claim := range claims {
		claimResponses[i] = core.ClaimResponse{
			ClaimID:           claim.ID,
			OfferID:           claim.OfferID,
			Status:            claim.Status,
			PriorityScore:     claim.PriorityScore,
			RequestedServings: claim.RequestedServings,
			CreatedAt:         claim.CreatedAt,
		}
	}

	return &core.ClaimListResponse{
		Items: claimResponses,
		Page:  pageNum,
		Limit: limitNum,
		Total: total,
	}, nil
}

// Additional method for jobs
func (s *ClaimService) GetClaimsForAutoAssignment(ctx context.Context) ([]*core.DonationClaim, error) {
	// Get all requested claims that need auto-assignment
	return s.claimRepo.GetClaimsByStatus(ctx, core.ClaimStatusRequested, 100, 0)
}
