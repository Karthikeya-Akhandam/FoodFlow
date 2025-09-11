package services

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"foodflow/internal/core"
	"foodflow/internal/core/repos"
)

type CreditService struct {
	creditRepo repos.CreditRepo
	orgRepo    repos.OrgRepo
}

func NewCreditService(creditRepo repos.CreditRepo, orgRepo repos.OrgRepo) core.CreditService {
	return &CreditService{
		creditRepo: creditRepo,
		orgRepo:    orgRepo,
	}
}

func (s *CreditService) GetUserCredits(ctx context.Context, userID string) (*core.CreditsResponse, error) {
	// Get organization for this user
	org, err := s.orgRepo.GetOrganizationByUserID(ctx, userID)
	if err != nil {
		return nil, core.NewNotFoundError("Organization not found for user")
	}

	// Get all credits for the organization
	allCredits, err := s.creditRepo.GetCreditsByOrganization(ctx, org.ID.String())
	if err != nil {
		return nil, core.NewInternalError("Failed to retrieve credits", err)
	}

	// Calculate total credits
	totalCredits := 0
	monthlyCredits := make([]core.CreditResponse, len(allCredits))
	
	for i, credit := range allCredits {
		totalCredits += credit.CreditsRemaining
		
		monthlyCredits[i] = core.CreditResponse{
			Month:            credit.Month,
			CreditsIssued:    credit.CreditsIssued,
			CreditsRemaining: credit.CreditsRemaining,
			ExpiresAt:        credit.ExpiresAt.Format(time.RFC3339),
		}
	}

	return &core.CreditsResponse{
		TotalCredits:   totalCredits,
		MonthlyCredits: monthlyCredits,
	}, nil
}

func (s *CreditService) GetCreditHistory(ctx context.Context, userID, page, limit string) (*core.CreditHistoryResponse, error) {
	// Parse pagination parameters
	pageNum := 1
	limitNum := 20
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

	// Get organization for this user
	org, err := s.orgRepo.GetOrganizationByUserID(ctx, userID)
	if err != nil {
		return nil, core.NewNotFoundError("Organization not found for user")
	}

	// Get credit history for the organization
	creditHistory, err := s.creditRepo.GetCreditHistory(ctx, org.ID.String(), limitNum, (pageNum-1)*limitNum)
	if err != nil {
		return nil, core.NewInternalError("Failed to retrieve credit history", err)
	}

	// Convert to response format
	transactions := make([]core.CreditTransactionResponse, len(creditHistory))
	for i, tx := range creditHistory {
		transactions[i] = core.CreditTransactionResponse{
			TransactionID: tx.ID,
			Amount:        tx.CreditsIssued,
			Type:          "ISSUED",
			Description:   fmt.Sprintf("Monthly credits for %s", tx.Month),
			CreatedAt:     tx.CreatedAt,
		}
	}

	return &core.CreditHistoryResponse{
		Items: transactions,
		Page:  pageNum,
		Limit: limitNum,
		Total: len(transactions), // In real app, would be separate count query
	}, nil
}

func (s *CreditService) SpendCredits(ctx context.Context, userID string, req *core.SpendCreditsRequest) (*core.CreditTransactionResponse, error) {
	// Get organization for this user
	org, err := s.orgRepo.GetOrganizationByUserID(ctx, userID)
	if err != nil {
		return nil, core.NewNotFoundError("Organization not found for user")
	}

	// Get current month's credits
	currentMonth := time.Now().Format("200601")
	credit, err := s.creditRepo.GetCreditByOrgAndMonth(ctx, org.ID.String(), currentMonth)
	if err != nil {
		return nil, core.NewNotFoundError("No credits available for current month")
	}

	// Check if sufficient credits available
	if credit.CreditsRemaining < req.Amount {
		return nil, core.NewValidationError("Insufficient credits", map[string]interface{}{
			"requested": req.Amount,
			"available": credit.CreditsRemaining,
		})
	}

	// Spend the credits
	transaction, err := s.creditRepo.SpendCredits(ctx, org.ID.String(), req.Amount, req.Reason)
	if err != nil {
		return nil, core.NewInternalError("Failed to spend credits", err)
	}

	return &core.CreditTransactionResponse{
		TransactionID: transaction.ID,
		Amount:        transaction.CreditsRemaining, // This represents spending amount
		Type:          "SPENT",
		Description:   fmt.Sprintf("Credits spent: %s", req.Reason),
		CreatedAt:     transaction.CreatedAt,
	}, nil
}

func (s *CreditService) IssueMonthlyCredits(ctx context.Context) error {
	// Get all organizations for monthly credit issuance
	orgs, err := s.orgRepo.GetAllOrganizations(ctx, 10000, 0)
	if err != nil {
		return core.NewInternalError("Failed to get organizations", err)
	}

	currentMonth := time.Now().Format("200601")
	
	for _, org := range orgs {
		// Calculate credits based on last month's people fed
		baseCredits := core.DefaultMonthlyCreditsBase
		peopleFed := org.LastMonthPeopleFed
		
		// Formula: base + (people_fed * multiplier)
		creditsToIssue := baseCredits + int(float64(peopleFed)*core.DefaultMonthlyCreditsK)
		
		// Issue credits for the month
		err = s.creditRepo.IssueCredits(ctx, org.ID.String(), currentMonth, creditsToIssue, false)
		if err != nil {
			// Log error but continue with other organizations
			continue
		}
	}

	return nil
}
