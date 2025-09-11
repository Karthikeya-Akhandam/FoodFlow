package services

import (
	"context"
	"strconv"
	"time"

	"foodflow/internal/core"
	"foodflow/internal/core/repos"
)

type AdminService struct {
	userRepo       repos.UserRepo
	orgRepo        repos.OrgRepo
	collabRepo     repos.CollabRepo
	offerRepo      repos.OfferRepo
	claimRepo      repos.ClaimRepo
	redemptionRepo repos.RedemptionRepo
}

func NewAdminService(
	userRepo repos.UserRepo,
	orgRepo repos.OrgRepo,
	collabRepo repos.CollabRepo,
	offerRepo repos.OfferRepo,
	claimRepo repos.ClaimRepo,
	redemptionRepo repos.RedemptionRepo,
) core.AdminService {
	return &AdminService{
		userRepo:       userRepo,
		orgRepo:        orgRepo,
		collabRepo:     collabRepo,
		offerRepo:      offerRepo,
		claimRepo:      claimRepo,
		redemptionRepo: redemptionRepo,
	}
}

func (s *AdminService) GetSystemStats(ctx context.Context) (*core.SystemStatsResponse, error) {
	// Get all users
	allUsers, err := s.userRepo.GetAllUsers(ctx, 10000, 0) // Get a large number to count all
	if err != nil {
		return nil, core.NewInternalError("Failed to get users", err)
	}

	// Count active users
	activeUsers := 0
	for _, user := range allUsers {
		if user.Status == core.UserStatusActive {
			activeUsers++
		}
	}

	// Get organizations count
	orgs, err := s.orgRepo.GetAllOrganizations(ctx, 10000, 0)
	if err != nil {
		return nil, core.NewInternalError("Failed to get organizations", err)
	}

	// Get collaborators count
	collabs, err := s.collabRepo.GetAllCollaborators(ctx, 10000, 0)
	if err != nil {
		return nil, core.NewInternalError("Failed to get collaborators", err)
	}

	// Get offers count
	allOffers, err := s.offerRepo.GetOffersByStatus(ctx, "", 10000, 0) // Get all statuses
	if err != nil {
		return nil, core.NewInternalError("Failed to get offers", err)
	}

	// Count active offers
	activeOffers := 0
	for _, offer := range allOffers {
		if offer.Status == core.DonationStatusOpen {
			activeOffers++
		}
	}

	// Get claims count
	allClaims, err := s.claimRepo.GetAllClaims(ctx, 10000, 0)
	if err != nil {
		return nil, core.NewInternalError("Failed to get claims", err)
	}

	// Get redemptions count
	allRedemptions, err := s.redemptionRepo.GetAllRedemptions(ctx, 10000, 0)
	if err != nil {
		return nil, core.NewInternalError("Failed to get redemptions", err)
	}

	return &core.SystemStatsResponse{
		TotalUsers:         len(allUsers),
		ActiveUsers:        activeUsers,
		TotalOrganizations: len(orgs),
		TotalCollaborators: len(collabs),
		TotalOffers:        len(allOffers),
		ActiveOffers:       activeOffers,
		TotalClaims:        len(allClaims),
		TotalRedemptions:   len(allRedemptions),
	}, nil
}

func (s *AdminService) GetUserStats(ctx context.Context) (*core.UserStatsResponse, error) {
	now := time.Now()
	today := now.Format("2006-01-02")
	weekStart := now.AddDate(0, 0, -int(now.Weekday())).Format("2006-01-02")
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02")

	// Get all users for analysis
	allUsers, err := s.userRepo.GetAllUsers(ctx, 10000, 0)
	if err != nil {
		return nil, core.NewInternalError("Failed to get users", err)
	}

	var newUsersToday, newUsersThisWeek, newUsersThisMonth, activeUsersToday, suspendedUsers int

	for _, user := range allUsers {
		userDate := user.CreatedAt.Format("2006-01-02")
		
		// Count new users by period
		if userDate == today {
			newUsersToday++
		}
		if userDate >= weekStart {
			newUsersThisWeek++
		}
		if userDate >= monthStart {
			newUsersThisMonth++
		}

		// Count by status
		if user.Status == core.UserStatusActive {
			// Assume active users today are those who are active status (in real app, would check last login)
			activeUsersToday++
		} else if user.Status == core.UserStatusSuspended {
			suspendedUsers++
		}
	}

	return &core.UserStatsResponse{
		NewUsersToday:     newUsersToday,
		NewUsersThisWeek:  newUsersThisWeek,
		NewUsersThisMonth: newUsersThisMonth,
		ActiveUsersToday:  activeUsersToday,
		SuspendedUsers:    suspendedUsers,
	}, nil
}

func (s *AdminService) GetDonationStats(ctx context.Context) (*core.DonationStatsResponse, error) {
	now := time.Now()
	today := now.Format("2006-01-02")
	weekStart := now.AddDate(0, 0, -int(now.Weekday())).Format("2006-01-02")
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02")

	// Get all offers for analysis
	allOffers, err := s.offerRepo.GetOffersByStatus(ctx, "", 10000, 0)
	if err != nil {
		return nil, core.NewInternalError("Failed to get offers", err)
	}

	// Get all redemptions for analysis
	allRedemptions, err := s.redemptionRepo.GetAllRedemptions(ctx, 10000, 0)
	if err != nil {
		return nil, core.NewInternalError("Failed to get redemptions", err)
	}

	var offersToday, offersThisWeek, offersThisMonth int
	var totalServingsOffered int

	for _, offer := range allOffers {
		offerDate := offer.CreatedAt.Format("2006-01-02")
		
		// Count offers by period
		if offerDate == today {
			offersToday++
		}
		if offerDate >= weekStart {
			offersThisWeek++
		}
		if offerDate >= monthStart {
			offersThisMonth++
		}

		totalServingsOffered += offer.EstimatedServings
	}

	// Calculate redemption stats
	var totalServingsRedeemed int
	for _, redemption := range allRedemptions {
		totalServingsRedeemed += redemption.ServingsAccepted
	}

	// Calculate rates and averages
	redemptionRate := float64(0)
	if totalServingsOffered > 0 {
		redemptionRate = float64(totalServingsRedeemed) / float64(totalServingsOffered) * 100
	}

	avgServingsPerOffer := float64(0)
	if len(allOffers) > 0 {
		avgServingsPerOffer = float64(totalServingsOffered) / float64(len(allOffers))
	}

	return &core.DonationStatsResponse{
		OffersToday:         offersToday,
		OffersThisWeek:      offersThisWeek,
		OffersThisMonth:     offersThisMonth,
		ServingsOffered:     totalServingsOffered,
		ServingsRedeemed:    totalServingsRedeemed,
		RedemptionRate:      redemptionRate,
		AvgServingsPerOffer: avgServingsPerOffer,
	}, nil
}

func (s *AdminService) ListUsers(ctx context.Context, page, limit, role, status string) (*core.UserListResponse, error) {
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

	// Get users with filters
	allUsers, err := s.userRepo.GetAllUsers(ctx, limitNum*10, 0) // Get more to filter
	if err != nil {
		return nil, core.NewInternalError("Failed to get users", err)
	}

	// Filter users based on role and status
	var filteredUsers []*core.User
	for _, user := range allUsers {
		includeUser := true

		// Filter by role
		if role != "" && string(user.Role) != role {
			includeUser = false
		}

		// Filter by status
		if status != "" && string(user.Status) != status {
			includeUser = false
		}

		if includeUser {
			filteredUsers = append(filteredUsers, user)
		}
	}

	// Apply pagination to filtered results
	total := len(filteredUsers)
	start := offset
	end := offset + limitNum

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	paginatedUsers := filteredUsers[start:end]

	// Convert to response format
	userResponses := make([]core.UserResponse, len(paginatedUsers))
	for i, user := range paginatedUsers {
		userResponses[i] = core.UserResponse{
			UserID: user.ID,
			Email:  user.Email,
			Role:   user.Role,
			Status: user.Status,
		}
	}

	return &core.UserListResponse{
		Items: userResponses,
		Page:  pageNum,
		Limit: limitNum,
		Total: total,
	}, nil
}

func (s *AdminService) UpdateUserStatus(ctx context.Context, userID string, req *core.UpdateUserStatusRequest) (*core.UserResponse, error) {
	// Update user status
	updatedUser, err := s.userRepo.UpdateUserStatus(ctx, userID, req.Status)
	if err != nil {
		return nil, core.NewInternalError("Failed to update user status", err)
	}

	return &core.UserResponse{
		UserID: updatedUser.ID,
		Email:  updatedUser.Email,
		Role:   updatedUser.Role,
		Status: updatedUser.Status,
	}, nil
}

func (s *AdminService) GetAuditLogs(ctx context.Context, page, limit, action, userID string) (*core.AuditLogResponse, error) {
	// Parse pagination parameters
	pageNum := 1
	limitNum := 50
	if page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			pageNum = p
		}
	}
	if limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 && l <= 1000 {
			limitNum = l
		}
	}

	// For now, return empty audit logs as we haven't implemented event tracking yet
	// In a real implementation, this would query an events/audit_logs table
	auditLogs := []core.AuditLogEntry{}

	return &core.AuditLogResponse{
		Items: auditLogs,
		Page:  pageNum,
		Limit: limitNum,
		Total: 0,
	}, nil
}