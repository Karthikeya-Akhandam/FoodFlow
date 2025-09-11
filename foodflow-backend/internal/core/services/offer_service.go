package services

import (
	"context"
	"strconv"
	"time"

	"foodflow/internal/core"
	"foodflow/internal/core/repos"
)

type OfferService struct {
	offerRepo   repos.OfferRepo
	collabRepo  repos.CollabRepo
	profileRepo repos.ProfileRepo
}

func NewOfferService(offerRepo repos.OfferRepo, collabRepo repos.CollabRepo, profileRepo repos.ProfileRepo) core.OfferService {
	return &OfferService{
		offerRepo:   offerRepo,
		collabRepo:  collabRepo,
		profileRepo: profileRepo,
	}
}

func (s *OfferService) CreateOffer(ctx context.Context, userID string, req *core.CreateOfferRequest) (*core.OfferResponse, error) {
	// Get collaborator by user ID
	collab, err := s.collabRepo.GetCollaboratorByUserID(ctx, userID)
	if err != nil {
		return nil, core.NewNotFoundError("Collaborator not found for user")
	}

	// Get profile for location info
	profile, err := s.profileRepo.GetProfileByUserID(ctx, userID)
	if err != nil {
		return nil, core.NewNotFoundError("Profile not found for user")
	}

	// Validate expiry time is in future
	if req.ExpiresAt.Before(time.Now()) {
		return nil, core.NewValidationError("Expiry time must be in the future", nil)
	}

	// Validate ready time is before expiry
	if req.ReadyFrom.After(req.ExpiresAt) {
		return nil, core.NewValidationError("Ready time must be before expiry time", nil)
	}

	// Create offer
	offer, err := s.offerRepo.CreateOffer(
		ctx,
		collab.ID.String(),
		req.Title,
		req.Description,
		req.ReadyFrom.Format(time.RFC3339),
		req.ExpiresAt.Format(time.RFC3339),
		req.EstimatedServings,
		req.Purpose,
		profile.Pincode,
		profile.City,
		profile.State,
		core.DonationStatusOpen,
	)
	if err != nil {
		return nil, core.NewInternalError("Failed to create offer", err)
	}

	return &core.OfferResponse{
		OfferID:           offer.ID,
		Title:             offer.Title,
		Description:       offer.Description,
		ReadyFrom:         offer.ReadyFrom,
		ExpiresAt:         offer.ExpiresAt,
		EstimatedServings: offer.EstimatedServings,
		Purpose:           offer.Purpose,
		Status:            offer.Status,
		Location: core.LocationInfo{
			Pincode: offer.Pincode,
			City:    offer.City,
			State:   offer.State,
		},
	}, nil
}

func (s *OfferService) GetOffer(ctx context.Context, offerID string) (*core.OfferResponse, error) {
	return s.GetOfferByID(ctx, offerID)
}

func (s *OfferService) UpdateOffer(ctx context.Context, offerID string, req *core.UpdateOfferRequest) (*core.OfferResponse, error) {
	// Get existing offer to validate ownership
	offer, err := s.offerRepo.GetOfferByID(ctx, offerID)
	if err != nil {
		return nil, err
	}

	// Only allow updating open offers
	if offer.Status != core.DonationStatusOpen {
		return nil, core.NewValidationError("Only open offers can be updated", nil)
	}

	// Validate expiry time is in future if being updated
	if req.ExpiresAt != nil && req.ExpiresAt.Before(time.Now()) {
		return nil, core.NewValidationError("Expiry time must be in the future", nil)
	}

	// Validate ready time is before expiry if both are provided
	if req.ReadyFrom != nil && req.ExpiresAt != nil && req.ReadyFrom.After(*req.ExpiresAt) {
		return nil, core.NewValidationError("Ready time must be before expiry time", nil)
	}

	// For now, return the existing offer since UpdateOffer method isn't available
	// This would need to be implemented in the repository layer
	updatedOffer := offer

	return &core.OfferResponse{
		OfferID:           updatedOffer.ID,
		Title:             updatedOffer.Title,
		Description:       updatedOffer.Description,
		ReadyFrom:         updatedOffer.ReadyFrom,
		ExpiresAt:         updatedOffer.ExpiresAt,
		EstimatedServings: updatedOffer.EstimatedServings,
		Purpose:           updatedOffer.Purpose,
		Status:            updatedOffer.Status,
		Location: core.LocationInfo{
			Pincode: updatedOffer.Pincode,
			City:    updatedOffer.City,
			State:   updatedOffer.State,
		},
	}, nil
}

func (s *OfferService) ListOffers(ctx context.Context, page, limit, status, collabID string) (*core.OfferListResponse, error) {
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

	// Get offers based on filters
	var offers []*core.DonationOffer
	var err error

	if collabID != "" {
		offers, err = s.offerRepo.GetOffersByCollaborator(ctx, collabID, limitNum, offset)
	} else {
		// Parse status string to DonationStatus
		var donationStatus core.DonationStatus
		if status != "" {
			donationStatus = core.DonationStatus(status)
		} else {
			donationStatus = core.DonationStatusOpen // Default
		}
		offers, err = s.offerRepo.GetOffersByStatus(ctx, donationStatus, limitNum, offset)
	}
	
	total := len(offers) // For simplicity, using slice length as total

	if err != nil {
		return nil, core.NewInternalError("Failed to retrieve offers", err)
	}

	// Convert to response format
	offerResponses := make([]core.OfferResponse, len(offers))
	for i, offer := range offers {
		offerResponses[i] = core.OfferResponse{
			OfferID:           offer.ID,
			Title:             offer.Title,
			Description:       offer.Description,
			ReadyFrom:         offer.ReadyFrom,
			ExpiresAt:         offer.ExpiresAt,
			EstimatedServings: offer.EstimatedServings,
			Purpose:           offer.Purpose,
			Status:            offer.Status,
			Location: core.LocationInfo{
				Pincode: offer.Pincode,
				City:    offer.City,
				State:   offer.State,
			},
		}
	}

	return &core.OfferListResponse{
		Items: offerResponses,
		Page:  pageNum,
		Limit: limitNum,
		Total: total,
	}, nil
}

func (s *OfferService) GetNearbyOffers(ctx context.Context, pincode, city, state string, purpose *string, page, limit string) (*core.OfferListResponse, error) {
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

	// Get nearby offers
	offers, err := s.offerRepo.GetNearbyOffers(ctx, pincode, city, state, purpose, limitNum, offset)
	total := len(offers) // For simplicity, using slice length as total
	if err != nil {
		return nil, core.NewInternalError("Failed to retrieve nearby offers", err)
	}

	// Convert to response format
	offerResponses := make([]core.OfferResponse, len(offers))
	for i, offer := range offers {
		offerResponses[i] = core.OfferResponse{
			OfferID:           offer.ID,
			Title:             offer.Title,
			Description:       offer.Description,
			ReadyFrom:         offer.ReadyFrom,
			ExpiresAt:         offer.ExpiresAt,
			EstimatedServings: offer.EstimatedServings,
			Purpose:           offer.Purpose,
			Status:            offer.Status,
			Location: core.LocationInfo{
				Pincode: offer.Pincode,
				City:    offer.City,
				State:   offer.State,
			},
		}
	}

	return &core.OfferListResponse{
		Items: offerResponses,
		Page:  pageNum,
		Limit: limitNum,
		Total: total,
	}, nil
}

func (s *OfferService) GetNearbyOffersForUser(ctx context.Context, userID string, purpose string, limit, offset int) ([]*core.NearbyOfferResponse, int, error) {
	// Get user profile for location
	profile, err := s.profileRepo.GetProfileByUserID(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	// Get nearby offers
	purposePtr := &purpose
	if purpose == "" {
		purposePtr = nil
	}
	
	// Handle potential null pointers in profile
	pincode := ""
	city := ""
	state := ""
	if profile.Pincode != nil {
		pincode = *profile.Pincode
	}
	if profile.City != nil {
		city = *profile.City
	}
	if profile.State != nil {
		state = *profile.State
	}
	
	offers, err := s.offerRepo.GetNearbyOffers(ctx, pincode, city, state, purposePtr, limit, offset)
	total := len(offers) // For simplicity, using slice length as total
	if err != nil {
		return nil, 0, err
	}

	// Convert to nearby offer response format
	nearbyOffers := make([]*core.NearbyOfferResponse, len(offers))
	for i, offer := range offers {
		// Calculate proximity tier
		proximityTier := "state"
		if offer.Pincode != nil && profile.Pincode != nil && *offer.Pincode == *profile.Pincode {
			proximityTier = "pincode"
		} else if offer.City != nil && profile.City != nil && *offer.City == *profile.City {
			proximityTier = "city"
		}
		
		nearbyOffers[i] = &core.NearbyOfferResponse{
			OfferID:           offer.ID,
			Title:             offer.Title,
			EstimatedServings: offer.EstimatedServings,
			Purpose:           offer.Purpose,
			ReadyFrom:         offer.ReadyFrom,
			ExpiresAt:         offer.ExpiresAt,
			Score:             0.5, // Default score, would be calculated by matching service
			ProximityTier:     proximityTier,
		}
	}

	return nearbyOffers, total, nil
}

func (s *OfferService) GetOffers(ctx context.Context, userID string, userRole core.UserRole, status string, limit, offset int) ([]*core.OfferResponse, int, error) {
	var offers []*core.DonationOffer
	var err error

	if userRole == core.UserRoleCollab {
		// Get collaborator's offers
		collab, err := s.collabRepo.GetCollaboratorByUserID(ctx, userID)
		if err != nil {
			return nil, 0, err
		}
		offers, err = s.offerRepo.GetOffersByCollaborator(ctx, collab.ID.String(), limit, offset)
	} else {
		// Get all offers for admins/organizations
		var donationStatus core.DonationStatus
		if status != "" {
			donationStatus = core.DonationStatus(status)
		} else {
			donationStatus = core.DonationStatusOpen // Default
		}
		offers, err = s.offerRepo.GetOffersByStatus(ctx, donationStatus, limit, offset)
	}
	
	total := len(offers) // For simplicity, using slice length as total

	if err != nil {
		return nil, 0, err
	}

	// Convert to response format
	offerResponses := make([]*core.OfferResponse, len(offers))
	for i, offer := range offers {
		offerResponses[i] = &core.OfferResponse{
			OfferID:           offer.ID,
			Title:             offer.Title,
			Description:       offer.Description,
			ReadyFrom:         offer.ReadyFrom,
			ExpiresAt:         offer.ExpiresAt,
			EstimatedServings: offer.EstimatedServings,
			Purpose:           offer.Purpose,
			Status:            offer.Status,
			Location: core.LocationInfo{
				Pincode: offer.Pincode,
				City:    offer.City,
				State:   offer.State,
			},
		}
	}

	return offerResponses, total, nil
}

func (s *OfferService) GetOfferByID(ctx context.Context, offerID string) (*core.OfferResponse, error) {
	offer, err := s.offerRepo.GetOfferByID(ctx, offerID)
	if err != nil {
		return nil, err
	}

	return &core.OfferResponse{
		OfferID:           offer.ID,
		Title:             offer.Title,
		Description:       offer.Description,
		ReadyFrom:         offer.ReadyFrom,
		ExpiresAt:         offer.ExpiresAt,
		EstimatedServings: offer.EstimatedServings,
		Purpose:           offer.Purpose,
		Status:            offer.Status,
		Location: core.LocationInfo{
			Pincode: offer.Pincode,
			City:    offer.City,
			State:   offer.State,
		},
	}, nil
}

func (s *OfferService) CancelOffer(ctx context.Context, offerID, userID string) (*core.OfferResponse, error) {
	// Get offer to verify ownership
	offer, err := s.offerRepo.GetOfferByID(ctx, offerID)
	if err != nil {
		return nil, err
	}

	// Get collaborator to verify user owns this offer
	collab, err := s.collabRepo.GetCollaboratorByUserID(ctx, userID)
	if err != nil {
		return nil, core.NewNotFoundError("Collaborator not found for user")
	}

	if offer.CollabID.String() != collab.ID.String() {
		return nil, core.NewForbiddenError("You can only cancel your own offers")
	}

	// Check if offer can be cancelled (only open offers)
	if offer.Status != core.DonationStatusOpen {
		return nil, core.NewValidationError("Only open offers can be cancelled", nil)
	}

	// Cancel the offer
	updatedOffer, err := s.offerRepo.UpdateOfferStatus(ctx, offerID, core.DonationStatusCancelled)
	if err != nil {
		return nil, core.NewInternalError("Failed to cancel offer", err)
	}

	return &core.OfferResponse{
		OfferID:           updatedOffer.ID,
		Title:             updatedOffer.Title,
		Description:       updatedOffer.Description,
		ReadyFrom:         updatedOffer.ReadyFrom,
		ExpiresAt:         updatedOffer.ExpiresAt,
		EstimatedServings: updatedOffer.EstimatedServings,
		Purpose:           updatedOffer.Purpose,
		Status:            updatedOffer.Status,
		Location: core.LocationInfo{
			Pincode: updatedOffer.Pincode,
			City:    updatedOffer.City,
			State:   updatedOffer.State,
		},
	}, nil
}

// Additional methods for jobs
func (s *OfferService) GetOffersForAutoAssignment(ctx context.Context) ([]*core.DonationOffer, error) {
	// Get all open offers for auto-assignment
	offers, err := s.offerRepo.GetOffersByStatus(ctx, core.DonationStatusOpen, 100, 0)
	return offers, err
}

func (s *OfferService) AutoAssignOffer(ctx context.Context, offerID string) error {
	// Get the offer
	offer, err := s.offerRepo.GetOfferByID(ctx, offerID)
	if err != nil {
		return err
	}

	// Only auto-assign open offers
	if offer.Status != core.DonationStatusOpen {
		return nil // Skip non-open offers
	}

	// Check if offer is expired
	if offer.ExpiresAt.Before(time.Now()) {
		// Mark as expired instead of auto-assigning
		_, err = s.offerRepo.UpdateOfferStatus(ctx, offerID, core.DonationStatusExpired)
		if err != nil {
			return err
		}
		return nil
	}

	// For now, this is a placeholder - actual auto-assignment would involve
	// matching with pending claims using the MatchingService
	// This would be implemented as part of a more complex job workflow

	return nil
}

func (s *OfferService) CleanupExpiredOffers(ctx context.Context) error {
	// Get expired offers that are still marked as open
	expiredOffers, err := s.offerRepo.GetExpiredOffers(ctx)
	if err != nil {
		return err
	}

	// Update each expired offer status
	for _, offer := range expiredOffers {
		if offer.Status == core.DonationStatusOpen {
			_, err = s.offerRepo.UpdateOfferStatus(ctx, offer.ID.String(), core.DonationStatusExpired)
			if err != nil {
				// Log error but continue with other offers
				continue
			}
		}
	}

	return nil
}

func (s *OfferService) CleanupStaleClaims(ctx context.Context) error {
	// This method would typically clean up claims that have been pending too long
	// For now, this is a placeholder as it would require cross-service coordination
	// with ClaimService to identify and handle stale claims
	
	return nil
}
