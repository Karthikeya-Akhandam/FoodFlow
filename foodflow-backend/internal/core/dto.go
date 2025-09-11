package core

import (
	"time"

	"github.com/google/uuid"
)

// Request DTOs

type SignupRequest struct {
	Email    string   `json:"email" validate:"required,email"`
	Password string   `json:"password" validate:"required,min=8"`
	Role     UserRole `json:"role" validate:"required,oneof=ORG COLLAB"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UpdateProfileRequest struct {
	Name        *string `json:"name" validate:"omitempty,min=2,max=100"`
	Phone       *string `json:"phone" validate:"omitempty,len=10"`
	Pincode     *string `json:"pincode" validate:"omitempty,len=6"`
	City        *string `json:"city" validate:"omitempty,min=2,max=50"`
	District    *string `json:"district" validate:"omitempty,min=2,max=50"`
	State       *string `json:"state" validate:"omitempty,min=2,max=50"`
	Address     *string `json:"address" validate:"omitempty,max=200"`
	IsRemoteOrg *bool   `json:"is_remote_org"`
}

type UpdateOrgRequest struct {
	OrgType            *string  `json:"org_type" validate:"omitempty,min=2,max=50"`
	PurposeFocus       []string `json:"purpose_focus" validate:"omitempty,dive,oneof=CHILDREN ELDERLY WOMEN GENERAL EMERGENCY"`
	LastMonthPeopleFed *int     `json:"last_month_people_fed" validate:"omitempty,min=0"`
}

type UpdateCollabRequest struct {
	CollabType *string `json:"collab_type" validate:"omitempty,min=2,max=50"`
}

type CreateOfferRequest struct {
	Title             string    `json:"title" validate:"required,min=5,max=100"`
	Description       *string   `json:"description" validate:"omitempty,max=500"`
	ReadyFrom         time.Time `json:"ready_from" validate:"required"`
	ExpiresAt         time.Time `json:"expires_at" validate:"required,gtfield=ReadyFrom"`
	EstimatedServings int       `json:"estimated_servings" validate:"required,min=1,max=10000"`
	Purpose           *string   `json:"purpose" validate:"omitempty,oneof=CHILDREN ELDERLY WOMEN GENERAL EMERGENCY"`
	Pincode           *string   `json:"pincode" validate:"omitempty,len=6"`
	City              *string   `json:"city" validate:"omitempty,min=2,max=50"`
	State             *string   `json:"state" validate:"omitempty,min=2,max=50"`
}

type CreateClaimRequest struct {
	RequestedServings   int        `json:"requested_servings" validate:"required,min=1,max=10000"`
	OnBehalfRemoteOrgID *uuid.UUID `json:"on_behalf_remote_org_id" validate:"omitempty"`
}

type AssignOfferRequest struct {
	ClaimID uuid.UUID `json:"claim_id" validate:"required"`
}

type RedeemOfferRequest struct {
	ServingsAccepted int `json:"servings_accepted" validate:"required,min=1,max=10000"`
}

type CreateRemoteAssignmentRequest struct {
	RemoteOrgID uuid.UUID `json:"remote_org_id" validate:"required"`
	ProxyOrgID  uuid.UUID `json:"proxy_org_id" validate:"required"`
	Active      bool      `json:"active"`
	Notes       *string   `json:"notes" validate:"omitempty,max=500"`
}

// Request DTOs continued

type CreateOrgRequest struct {
	OrgType            string   `json:"org_type" validate:"required,min=2,max=50"`
	PurposeFocus       []string `json:"purpose_focus" validate:"required,dive,oneof=CHILDREN ELDERLY WOMEN GENERAL EMERGENCY"`
	LastMonthPeopleFed int      `json:"last_month_people_fed" validate:"required,min=0"`
}

type CreateCollabRequest struct {
	CollabType string `json:"collab_type" validate:"required,min=2,max=50"`
}

type UpdateOfferRequest struct {
	Title             *string    `json:"title" validate:"omitempty,min=5,max=100"`
	Description       *string    `json:"description" validate:"omitempty,max=500"`
	ReadyFrom         *time.Time `json:"ready_from"`
	ExpiresAt         *time.Time `json:"expires_at"`
	EstimatedServings *int       `json:"estimated_servings" validate:"omitempty,min=1,max=10000"`
	Purpose           *string    `json:"purpose" validate:"omitempty,oneof=CHILDREN ELDERLY WOMEN GENERAL EMERGENCY"`
	Pincode           *string    `json:"pincode" validate:"omitempty,len=6"`
	City              *string    `json:"city" validate:"omitempty,min=2,max=50"`
	State             *string    `json:"state" validate:"omitempty,min=2,max=50"`
}

type UpdateClaimStatusRequest struct {
	Status ClaimStatus `json:"status" validate:"required,oneof=REQUESTED WON LOST CANCELLED"`
}

type SpendCreditsRequest struct {
	Amount int    `json:"amount" validate:"required,min=1"`
	Reason string `json:"reason" validate:"required,min=5,max=200"`
}

type RedeemTokensRequest struct {
	Amount int    `json:"amount" validate:"required,min=1"`
	Reason string `json:"reason" validate:"required,min=5,max=200"`
}

type CreateRedemptionRequest struct {
	OfferID          uuid.UUID `json:"offer_id" validate:"required"`
	ServingsAccepted int       `json:"servings_accepted" validate:"required,min=1,max=10000"`
}

type UpdateRedemptionStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=PENDING CONFIRMED CANCELLED"`
}

type CreateRemoteOrgRequest struct {
	Name             string  `json:"name" validate:"required,min=2,max=100"`
	ContactEmail     string  `json:"contact_email" validate:"required,email"`
	ContactPhone     *string `json:"contact_phone" validate:"omitempty,len=10"`
	Address          string  `json:"address" validate:"required,min=10,max=200"`
	Pincode          string  `json:"pincode" validate:"required,len=6"`
	City             string  `json:"city" validate:"required,min=2,max=50"`
	State            string  `json:"state" validate:"required,min=2,max=50"`
	ServingCapacity  int     `json:"serving_capacity" validate:"required,min=1"`
	PurposeFocus     []string `json:"purpose_focus" validate:"required,dive,oneof=CHILDREN ELDERLY WOMEN GENERAL EMERGENCY"`
}

type UpdateRemoteOrgRequest struct {
	Name            *string  `json:"name" validate:"omitempty,min=2,max=100"`
	ContactEmail    *string  `json:"contact_email" validate:"omitempty,email"`
	ContactPhone    *string  `json:"contact_phone" validate:"omitempty,len=10"`
	Address         *string  `json:"address" validate:"omitempty,min=10,max=200"`
	ServingCapacity *int     `json:"serving_capacity" validate:"omitempty,min=1"`
	PurposeFocus    []string `json:"purpose_focus" validate:"omitempty,dive,oneof=CHILDREN ELDERLY WOMEN GENERAL EMERGENCY"`
}

type AssignRemoteOrgRequest struct {
	RemoteOrgID uuid.UUID `json:"remote_org_id" validate:"required"`
	Notes       *string   `json:"notes" validate:"omitempty,max=500"`
}

type UpdateUserStatusRequest struct {
	Status UserStatus `json:"status" validate:"required,oneof=ACTIVE SUSPENDED"`
}

// Response DTOs

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type UserResponse struct {
	UserID uuid.UUID  `json:"user_id"`
	Email  string     `json:"email"`
	Role   UserRole   `json:"role"`
	Status UserStatus `json:"status"`
}

type ProfileResponse struct {
	Name        *string `json:"name"`
	Phone       *string `json:"phone"`
	Pincode     *string `json:"pincode"`
	City        *string `json:"city"`
	District    *string `json:"district"`
	State       *string `json:"state"`
	Address     *string `json:"address"`
	IsRemoteOrg bool    `json:"is_remote_org"`
}

type OrgResponse struct {
	OrgID              uuid.UUID `json:"org_id"`
	OrgType            string    `json:"org_type"`
	PurposeFocus       []string  `json:"purpose_focus"`
	LastMonthPeopleFed int       `json:"last_month_people_fed"`
}

type CollabResponse struct {
	CollabID   uuid.UUID `json:"collab_id"`
	CollabType string    `json:"collab_type"`
}

type OfferResponse struct {
	OfferID           uuid.UUID      `json:"offer_id"`
	CollabID          uuid.UUID      `json:"collab_id"`
	Title             string         `json:"title"`
	Description       *string        `json:"description"`
	ReadyFrom         time.Time      `json:"ready_from"`
	ExpiresAt         time.Time      `json:"expires_at"`
	EstimatedServings int            `json:"estimated_servings"`
	Purpose           *string        `json:"purpose"`
	Location          LocationInfo   `json:"location"`
	Status            DonationStatus `json:"status"`
}

type LocationInfo struct {
	Pincode *string `json:"pincode"`
	City    *string `json:"city"`
	State   *string `json:"state"`
}

type NearbyOfferResponse struct {
	OfferID           uuid.UUID `json:"offer_id"`
	Title             string    `json:"title"`
	EstimatedServings int       `json:"estimated_servings"`
	Purpose           *string   `json:"purpose"`
	ReadyFrom         time.Time `json:"ready_from"`
	ExpiresAt         time.Time `json:"expires_at"`
	Score             float64   `json:"score"`
	ProximityTier     string    `json:"proximity_tier"`
}

type ClaimResponse struct {
	ClaimID           uuid.UUID   `json:"claim_id"`
	OfferID           uuid.UUID   `json:"offer_id"`
	Status            ClaimStatus `json:"status"`
	PriorityScore     float64     `json:"priority_score"`
	RequestedServings int         `json:"requested_servings"`
	CreatedAt         time.Time   `json:"created_at"`
}

type CreditResponse struct {
	Month            string `json:"month"`
	CreditsIssued    int    `json:"credits_issued"`
	CreditsRemaining int    `json:"credits_remaining"`
	ExpiresAt        string `json:"expires_at"`
}

type TokenResponse struct {
	Month        string `json:"month"`
	TokensEarned int    `json:"tokens_earned"`
}

type RedemptionResponse struct {
	RedemptionID  uuid.UUID `json:"redemption_id"`
	CreditsSpent  int       `json:"credits_spent"`
	TokensAwarded int       `json:"tokens_awarded"`
	OfferStatus   string    `json:"offer_status"`
}

type RemoteAssignmentResponse struct {
	AssignmentID uuid.UUID `json:"assignment_id"`
	RemoteOrgID  uuid.UUID `json:"remote_org_id"`
	ProxyOrgID   uuid.UUID `json:"proxy_org_id"`
	Active       bool      `json:"active"`
	Notes        *string   `json:"notes"`
}

type PaginatedResponse struct {
	Items []any `json:"items"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int   `json:"total"`
}

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// List Response DTOs
type OrgListResponse struct {
	Items []OrgResponse `json:"items"`
	Page  int           `json:"page"`
	Limit int           `json:"limit"`
	Total int           `json:"total"`
}

type CollabListResponse struct {
	Items []CollabResponse `json:"items"`
	Page  int              `json:"page"`
	Limit int              `json:"limit"`
	Total int              `json:"total"`
}

type OfferListResponse struct {
	Items []OfferResponse `json:"items"`
	Page  int             `json:"page"`
	Limit int             `json:"limit"`
	Total int             `json:"total"`
}

type ClaimListResponse struct {
	Items []ClaimResponse `json:"items"`
	Page  int             `json:"page"`
	Limit int             `json:"limit"`
	Total int             `json:"total"`
}

type RedemptionListResponse struct {
	Items []RedemptionResponse `json:"items"`
	Page  int                  `json:"page"`
	Limit int                  `json:"limit"`
	Total int                  `json:"total"`
}

type UserListResponse struct {
	Items []UserResponse `json:"items"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
	Total int            `json:"total"`
}

type RemoteOrgListResponse struct {
	Items []RemoteOrgResponse `json:"items"`
	Page  int                 `json:"page"`
	Limit int                 `json:"limit"`
	Total int                 `json:"total"`
}

// Additional Response DTOs
type CreditsResponse struct {
	TotalCredits   int              `json:"total_credits"`
	MonthlyCredits []CreditResponse `json:"monthly_credits"`
}

type CreditHistoryResponse struct {
	Items []CreditTransactionResponse `json:"items"`
	Page  int                         `json:"page"`
	Limit int                         `json:"limit"`
	Total int                         `json:"total"`
}

type CreditTransactionResponse struct {
	TransactionID uuid.UUID `json:"transaction_id"`
	Amount        int       `json:"amount"`
	Type          string    `json:"type"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"created_at"`
}

type TokensResponse struct {
	TotalTokens   int             `json:"total_tokens"`
	MonthlyTokens []TokenResponse `json:"monthly_tokens"`
}

type TokenHistoryResponse struct {
	Items []TokenTransactionResponse `json:"items"`
	Page  int                        `json:"page"`
	Limit int                        `json:"limit"`
	Total int                        `json:"total"`
}

type TokenTransactionResponse struct {
	TransactionID uuid.UUID `json:"transaction_id"`
	Amount        int       `json:"amount"`
	Type          string    `json:"type"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"created_at"`
}

type TokenRedemptionResponse struct {
	RedemptionID  uuid.UUID `json:"redemption_id"`
	TokensSpent   int       `json:"tokens_spent"`
	CreditsEarned int       `json:"credits_earned"`
	Status        string    `json:"status"`
}

type RemoteOrgResponse struct {
	RemoteOrgID     uuid.UUID `json:"remote_org_id"`
	Name            string    `json:"name"`
	ContactEmail    string    `json:"contact_email"`
	ContactPhone    *string   `json:"contact_phone"`
	Address         string    `json:"address"`
	Location        LocationInfo `json:"location"`
	ServingCapacity int       `json:"serving_capacity"`
	PurposeFocus    []string  `json:"purpose_focus"`
}

type RemoteOrgAssignmentResponse struct {
	AssignmentID uuid.UUID `json:"assignment_id"`
	RemoteOrgID  uuid.UUID `json:"remote_org_id"`
	ProxyOrgID   uuid.UUID `json:"proxy_org_id"`
	Active       bool      `json:"active"`
	Notes        *string   `json:"notes"`
	CreatedAt    time.Time `json:"created_at"`
}

// Admin Response DTOs
type SystemStatsResponse struct {
	TotalUsers        int `json:"total_users"`
	ActiveUsers       int `json:"active_users"`
	TotalOrganizations int `json:"total_organizations"`
	TotalCollaborators int `json:"total_collaborators"`
	TotalOffers       int `json:"total_offers"`
	ActiveOffers      int `json:"active_offers"`
	TotalClaims       int `json:"total_claims"`
	TotalRedemptions  int `json:"total_redemptions"`
}

type UserStatsResponse struct {
	NewUsersToday     int `json:"new_users_today"`
	NewUsersThisWeek  int `json:"new_users_this_week"`
	NewUsersThisMonth int `json:"new_users_this_month"`
	ActiveUsersToday  int `json:"active_users_today"`
	SuspendedUsers    int `json:"suspended_users"`
}

type DonationStatsResponse struct {
	OffersToday         int     `json:"offers_today"`
	OffersThisWeek      int     `json:"offers_this_week"`
	OffersThisMonth     int     `json:"offers_this_month"`
	ServingsOffered     int     `json:"servings_offered"`
	ServingsRedeemed    int     `json:"servings_redeemed"`
	RedemptionRate      float64 `json:"redemption_rate"`
	AvgServingsPerOffer float64 `json:"avg_servings_per_offer"`
}

type AuditLogResponse struct {
	Items []AuditLogEntry `json:"items"`
	Page  int             `json:"page"`
	Limit int             `json:"limit"`
	Total int             `json:"total"`
}

type AuditLogEntry struct {
	ID          uuid.UUID      `json:"id"`
	UserID      *uuid.UUID     `json:"user_id"`
	Action      string         `json:"action"`
	EntityType  string         `json:"entity_type"`
	EntityID    *uuid.UUID     `json:"entity_id"`
	Changes     map[string]any `json:"changes"`
	IPAddress   string         `json:"ip_address"`
	UserAgent   string         `json:"user_agent"`
	CreatedAt   time.Time      `json:"created_at"`
}

// Common error codes
const (
	ErrorCodeValidation     = "VALIDATION_ERROR"
	ErrorCodeUnauthorized   = "UNAUTHORIZED"
	ErrorCodeForbidden      = "FORBIDDEN"
	ErrorCodeNotFound       = "NOT_FOUND"
	ErrorCodeConflict       = "CONFLICT"
	ErrorCodeInternalError  = "INTERNAL_ERROR"
	ErrorCodeRateLimited    = "RATE_LIMITED"
	ErrorCodeIdempotencyKey = "IDEMPOTENCY_KEY_ERROR"
)
