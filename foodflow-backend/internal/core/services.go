package core

import "context"

// Service interfaces for dependency injection

type AuthService interface {
	Signup(ctx context.Context, req *SignupRequest) (*UserResponse, error)
	Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error)
	GetMe(ctx context.Context, userID string) (*UserResponse, error)
}

type ProfileService interface {
	GetProfile(ctx context.Context, userID string) (*ProfileResponse, error)
	UpdateProfile(ctx context.Context, userID string, req *UpdateProfileRequest) (*ProfileResponse, error)
}

type OrgService interface {
	CreateOrganization(ctx context.Context, userID string, req *CreateOrgRequest) (*OrgResponse, error)
	GetOrganization(ctx context.Context, orgID string) (*OrgResponse, error)
	UpdateOrganization(ctx context.Context, orgID string, req *UpdateOrgRequest) (*OrgResponse, error)
	ListOrganizations(ctx context.Context, page, limit string) (*OrgListResponse, error)
}

type CollabService interface {
	CreateCollaborator(ctx context.Context, userID string, req *CreateCollabRequest) (*CollabResponse, error)
	GetCollaborator(ctx context.Context, collabID string) (*CollabResponse, error)
	UpdateCollaborator(ctx context.Context, collabID string, req *UpdateCollabRequest) (*CollabResponse, error)
	ListCollaborators(ctx context.Context, page, limit string) (*CollabListResponse, error)
}

type OfferService interface {
	CreateOffer(ctx context.Context, userID string, req *CreateOfferRequest) (*OfferResponse, error)
	GetOffer(ctx context.Context, offerID string) (*OfferResponse, error)
	GetOfferByID(ctx context.Context, offerID string) (*OfferResponse, error)
	UpdateOffer(ctx context.Context, offerID string, req *UpdateOfferRequest) (*OfferResponse, error)
	ListOffers(ctx context.Context, page, limit, status, collabID string) (*OfferListResponse, error)
	GetOffers(ctx context.Context, userID string, userRole UserRole, status string, limit, offset int) ([]*OfferResponse, int, error)
	GetNearbyOffers(ctx context.Context, pincode, city, state string, purpose *string, page, limit string) (*OfferListResponse, error)
	GetNearbyOffersForUser(ctx context.Context, userID string, purpose string, limit, offset int) ([]*NearbyOfferResponse, int, error)
	CancelOffer(ctx context.Context, offerID, userID string) (*OfferResponse, error)
}

type ClaimService interface {
	CreateClaim(ctx context.Context, userID, offerID string, req *CreateClaimRequest) (*ClaimResponse, error)
	GetClaim(ctx context.Context, claimID string) (*ClaimResponse, error)
	UpdateClaimStatus(ctx context.Context, claimID string, req *UpdateClaimStatusRequest) (*ClaimResponse, error)
	ListClaims(ctx context.Context, page, limit, status, orgID string) (*ClaimListResponse, error)
}

type MatchingService interface {
	CalculateMatchingScore(ctx context.Context, claim *DonationClaim, offer *DonationOffer) (float64, error)
	FindBestMatches(ctx context.Context, claim *DonationClaim, limit int) ([]*DonationOffer, error)
}

type CreditService interface {
	GetUserCredits(ctx context.Context, userID string) (*CreditsResponse, error)
	GetCreditHistory(ctx context.Context, userID, page, limit string) (*CreditHistoryResponse, error)
	SpendCredits(ctx context.Context, userID string, req *SpendCreditsRequest) (*CreditTransactionResponse, error)
	IssueMonthlyCredits(ctx context.Context) error
}

type TokenService interface {
	GetUserTokens(ctx context.Context, userID string) (*TokensResponse, error)
	GetTokenHistory(ctx context.Context, userID, page, limit string) (*TokenHistoryResponse, error)
	RedeemTokens(ctx context.Context, userID string, req *RedeemTokensRequest) (*TokenRedemptionResponse, error)
	AwardTokens(ctx context.Context, collabID string, creditsSpent int) error
}

type RedemptionService interface {
	CreateRedemption(ctx context.Context, userID string, req *CreateRedemptionRequest) (*RedemptionResponse, error)
	GetRedemption(ctx context.Context, redemptionID string) (*RedemptionResponse, error)
	UpdateRedemptionStatus(ctx context.Context, redemptionID string, req *UpdateRedemptionStatusRequest) (*RedemptionResponse, error)
	ListRedemptions(ctx context.Context, page, limit, status, orgID string) (*RedemptionListResponse, error)
}

type RemoteService interface {
	CreateRemoteOrg(ctx context.Context, userID string, req *CreateRemoteOrgRequest) (*RemoteOrgResponse, error)
	GetRemoteOrg(ctx context.Context, remoteOrgID string) (*RemoteOrgResponse, error)
	UpdateRemoteOrg(ctx context.Context, remoteOrgID string, req *UpdateRemoteOrgRequest) (*RemoteOrgResponse, error)
	ListRemoteOrgs(ctx context.Context, page, limit string) (*RemoteOrgListResponse, error)
	AssignRemoteOrg(ctx context.Context, userID string, req *AssignRemoteOrgRequest) (*RemoteOrgAssignmentResponse, error)
}

type AdminService interface {
	GetSystemStats(ctx context.Context) (*SystemStatsResponse, error)
	GetUserStats(ctx context.Context) (*UserStatsResponse, error)
	GetDonationStats(ctx context.Context) (*DonationStatsResponse, error)
	ListUsers(ctx context.Context, page, limit, role, status string) (*UserListResponse, error)
	UpdateUserStatus(ctx context.Context, userID string, req *UpdateUserStatusRequest) (*UserResponse, error)
	GetAuditLogs(ctx context.Context, page, limit, action, userID string) (*AuditLogResponse, error)
}
