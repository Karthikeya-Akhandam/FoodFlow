package services

import (
	"context"
	"math"
	"strings"
	"time"

	"foodflow/internal/core"
	"foodflow/internal/core/repos"
)

type MatchingService struct {
	offerRepo   repos.OfferRepo
	claimRepo   repos.ClaimRepo
	creditRepo  repos.CreditRepo
	orgRepo     repos.OrgRepo
	profileRepo repos.ProfileRepo
	remoteRepo  repos.RemoteRepo
}

func NewMatchingService(
	offerRepo repos.OfferRepo,
	claimRepo repos.ClaimRepo,
	creditRepo repos.CreditRepo,
	orgRepo repos.OrgRepo,
	profileRepo repos.ProfileRepo,
	remoteRepo repos.RemoteRepo,
) core.MatchingService {
	return &MatchingService{
		offerRepo:   offerRepo,
		claimRepo:   claimRepo,
		creditRepo:  creditRepo,
		orgRepo:     orgRepo,
		profileRepo: profileRepo,
		remoteRepo:  remoteRepo,
	}
}

func (s *MatchingService) CalculateMatchingScore(ctx context.Context, claim *core.DonationClaim, offer *core.DonationOffer) (float64, error) {
	// Get organization and profile for proximity calculation
	org, err := s.orgRepo.GetOrganizationByID(ctx, claim.OrgID.String())
	if err != nil {
		return 0, err
	}

	// Get profile for organization's location
	profile, err := s.profileRepo.GetProfileByUserID(ctx, org.UserID.String())
	if err != nil {
		return 0, err
	}

	// Calculate component scores
	proximityScore := s.calculateProximityScore(offer, profile)
	purposeScore := s.calculatePurposeScore(org, offer)
	creditPressureScore := s.calculateCreditPressureScore(ctx, org)
	reliabilityScore := s.calculateReliabilityScore(ctx, org)
	remoteProxyScore := s.calculateRemoteProxyScore(ctx, org)
	urgencyScore := s.calculateUrgencyScore(offer)
	servingScore := s.calculateServingScore(claim, offer)

	// Weighted combination based on core constants
	finalScore := (proximityScore * core.ScoreWeightProximity +
		purposeScore * core.ScoreWeightPurposeMatch +
		creditPressureScore * core.ScoreWeightCreditPressure +
		reliabilityScore * core.ScoreWeightReliability +
		remoteProxyScore * core.ScoreWeightRemoteProxy) / 100.0

	// Apply urgency and serving efficiency multipliers
	finalScore *= urgencyScore * servingScore

	return math.Max(0.0, math.Min(1.0, finalScore)), nil
}

func (s *MatchingService) FindBestMatches(ctx context.Context, claim *core.DonationClaim, limit int) ([]*core.DonationOffer, error) {
	// Get available offers
	offers, err := s.offerRepo.GetOffersByStatus(ctx, core.DonationStatusOpen, limit*3, 0) // Get more to filter
	if err != nil {
		return nil, err
	}

	// Score and filter offers
	type scoredOffer struct {
		offer *core.DonationOffer
		score float64
	}

	var scoredOffers []scoredOffer
	for _, offer := range offers {
		// Skip expired offers
		if offer.ExpiresAt.Before(time.Now()) {
			continue
		}

		// Calculate matching score
		score, err := s.CalculateMatchingScore(ctx, claim, offer)
		if err != nil {
			continue // Skip on error
		}

		// Only include offers with reasonable scores
		if score > 0.1 {
			scoredOffers = append(scoredOffers, scoredOffer{offer, score})
		}
	}

	// Sort by score (highest first)
	for i := 0; i < len(scoredOffers); i++ {
		for j := i + 1; j < len(scoredOffers); j++ {
			if scoredOffers[j].score > scoredOffers[i].score {
				scoredOffers[i], scoredOffers[j] = scoredOffers[j], scoredOffers[i]
			}
		}
	}

	// Extract top offers
	var result []*core.DonationOffer
	for i, scoredOffer := range scoredOffers {
		if i >= limit {
			break
		}
		result = append(result, scoredOffer.offer)
	}

	return result, nil
}

func (s *MatchingService) calculateUrgencyScore(offer *core.DonationOffer) float64 {
	// Higher score for offers expiring sooner
	now := time.Now()
	timeToExpiry := offer.ExpiresAt.Sub(now)
	
	// If already expired, score is 0
	if timeToExpiry <= 0 {
		return 0.0
	}
	
	// Normalize to hours
	hoursToExpiry := timeToExpiry.Hours()
	
	// Higher urgency for offers expiring within 24 hours
	if hoursToExpiry <= 24 {
		return 1.0
	} else if hoursToExpiry <= 48 {
		return 0.8
	} else if hoursToExpiry <= 72 {
		return 0.6
	} else {
		return 0.4
	}
}

func (s *MatchingService) calculateServingScore(claim *core.DonationClaim, offer *core.DonationOffer) float64 {
	// Score based on how well the serving amounts match
	requestedServings := float64(claim.RequestedServings)
	availableServings := float64(offer.EstimatedServings)
	
	if requestedServings <= availableServings {
		// Perfect or under-utilization
		ratio := requestedServings / availableServings
		return math.Min(1.0, ratio + 0.1) // Small bonus for efficient use
	} else {
		// Over-requested, penalty
		ratio := availableServings / requestedServings
		return ratio * 0.8 // Penalty for over-requesting
	}
}

func (s *MatchingService) calculateProximityScore(offer *core.DonationOffer, profile *core.Profile) float64 {
	// Exact pincode match gets highest score
	if offer.Pincode != nil && profile.Pincode != nil && *offer.Pincode == *profile.Pincode {
		return 1.0
	}
	
	// Same city gets good score
	if offer.City != nil && profile.City != nil && strings.EqualFold(*offer.City, *profile.City) {
		return 0.7
	}
	
	// Same state gets moderate score
	if offer.State != nil && profile.State != nil && strings.EqualFold(*offer.State, *profile.State) {
		return 0.4
	}
	
	// Different state gets low score
	return 0.1
}

func (s *MatchingService) calculatePurposeScore(org *core.Organization, offer *core.DonationOffer) float64 {
	// If organization has purpose focus and offer has purpose
	if len(org.PurposeFocus) > 0 && offer.Purpose != nil {
		// Check if offer purpose matches any of org's focus areas
		for _, orgPurpose := range org.PurposeFocus {
			if orgPurpose == *offer.Purpose {
				// Perfect purpose match
				return 1.0
			}
		}
		
		// Check for purpose compatibility
		for _, orgPurpose := range org.PurposeFocus {
			switch orgPurpose {
			case "CHILDREN", "ELDERLY", "WOMEN":
				// Specialized NGOs can accept general food too
				if *offer.Purpose == "GENERAL" {
					return 0.6
				}
			case "GENERAL":
				// General NGOs can accept any purpose
				return 0.8
			}
		}
		
		// No compatibility
		return 0.3
	}
	
	// If either has no specific purpose, neutral score
	return 0.5
}

func (s *MatchingService) calculateCreditPressureScore(ctx context.Context, org *core.Organization) float64 {
	// Get current month's credit info
	currentMonth := time.Now().Format("200601")
	credit, err := s.creditRepo.GetCreditByOrgAndMonth(ctx, org.ID.String(), currentMonth)
	if err != nil {
		return 0.5 // Neutral if can't determine
	}
	
	if credit.CreditsIssued <= 0 {
		return 0.5 // Neutral if no credits issued
	}
	
	// Calculate usage ratio
	usageRatio := float64(credit.CreditsIssued - credit.CreditsRemaining) / float64(credit.CreditsIssued)
	
	// Higher pressure (more credits used) gets higher score
	if usageRatio >= 0.8 {
		return 1.0 // High pressure
	} else if usageRatio >= 0.6 {
		return 0.8
	} else if usageRatio >= 0.4 {
		return 0.6
	} else {
		return 0.4 // Low pressure
	}
}

func (s *MatchingService) calculateReliabilityScore(ctx context.Context, org *core.Organization) float64 {
	// This would typically look at historical performance
	// For now, return neutral score - can be enhanced later
	return 0.7
}

func (s *MatchingService) calculateRemoteProxyScore(ctx context.Context, org *core.Organization) float64 {
	// For now, return a neutral score for all organizations
	// This would be enhanced when remote org assignment logic is implemented
	return 0.5
}