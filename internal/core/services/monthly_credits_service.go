package services

import (
	"context"
	"time"

	"foodflow/internal/core"
	"foodflow/internal/core/repos"
)

type MonthlyCreditsService struct {
	creditRepo repos.CreditRepo
	orgRepo    repos.OrgRepo
}

func NewMonthlyCreditsService(creditRepo repos.CreditRepo, orgRepo repos.OrgRepo) *MonthlyCreditsService {
	return &MonthlyCreditsService{
		creditRepo: creditRepo,
		orgRepo:    orgRepo,
	}
}

func (s *MonthlyCreditsService) IssueMonthlyCredits(ctx context.Context) error {
	// This method is very similar to CreditService.IssueMonthlyCredits
	// In a real application, you might consolidate these or have this call the credit service
	
	// Get all organizations for monthly credit issuance
	orgs, err := s.orgRepo.GetAllOrganizations(ctx, 10000, 0)
	if err != nil {
		return err
	}

	currentMonth := time.Now().Format("200601")
	
	for _, org := range orgs {
		// Calculate credits based on last month's people fed
		baseCredits := core.DefaultMonthlyCreditsBase
		peopleFed := org.LastMonthPeopleFed
		
		// Formula: base + (people_fed * multiplier)
		creditsToIssue := baseCredits + int(float64(peopleFed)*core.DefaultMonthlyCreditsK)
		
		// Check if organization is remote and eligible for bonus
		isRemote := false // Would check remote assignments in real implementation
		
		// Issue credits for the month
		err = s.creditRepo.IssueCredits(ctx, org.ID.String(), currentMonth, creditsToIssue, isRemote)
		if err != nil {
			// Log error but continue with other organizations
			continue
		}
	}

	return nil
}