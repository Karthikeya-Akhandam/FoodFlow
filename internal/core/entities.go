package core

import (
	"time"

	"github.com/google/uuid"
)

// User represents a system user
type User struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	Email        string     `json:"email" db:"email"`
	PasswordHash string     `json:"-" db:"password_hash"`
	Role         UserRole   `json:"role" db:"role"`
	Status       UserStatus `json:"status" db:"status"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

// Profile represents user profile information
type Profile struct {
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	Name        string    `json:"name" db:"name"`
	Phone       *string   `json:"phone" db:"phone"`
	Pincode     *string   `json:"pincode" db:"pincode"`
	City        *string   `json:"city" db:"city"`
	District    *string   `json:"district" db:"district"`
	State       *string   `json:"state" db:"state"`
	Address     *string   `json:"address" db:"address"`
	IsRemoteOrg bool      `json:"is_remote_org" db:"is_remote_org"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Organization represents an organization entity
type Organization struct {
	ID                 uuid.UUID `json:"id" db:"id"`
	UserID             uuid.UUID `json:"user_id" db:"user_id"`
	OrgType            string    `json:"org_type" db:"org_type"`
	PurposeFocus       []string  `json:"purpose_focus" db:"purpose_focus"`
	LastMonthPeopleFed int       `json:"last_month_people_fed" db:"last_month_people_fed"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

// Collaborator represents a collaborator entity
type Collaborator struct {
	ID         uuid.UUID `json:"id" db:"id"`
	UserID     uuid.UUID `json:"user_id" db:"user_id"`
	CollabType string    `json:"collab_type" db:"collab_type"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// Credit represents monthly credits for organizations
type Credit struct {
	ID               uuid.UUID `json:"id" db:"id"`
	OrgID            uuid.UUID `json:"org_id" db:"org_id"`
	Month            string    `json:"month" db:"month"`
	CreditsIssued    int       `json:"credits_issued" db:"credits_issued"`
	CreditsRemaining int       `json:"credits_remaining" db:"credits_remaining"`
	IsRemoteBonus    bool      `json:"is_remote_bonus" db:"is_remote_bonus"`
	ExpiresAt        time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}

// Token represents monthly tokens for collaborators
type Token struct {
	ID           uuid.UUID `json:"id" db:"id"`
	CollabID     uuid.UUID `json:"collab_id" db:"collab_id"`
	Month        string    `json:"month" db:"month"`
	TokensEarned int       `json:"tokens_earned" db:"tokens_earned"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// DonationOffer represents a donation offer
type DonationOffer struct {
	ID                uuid.UUID      `json:"id" db:"id"`
	CollabID          uuid.UUID      `json:"collab_id" db:"collab_id"`
	Title             string         `json:"title" db:"title"`
	Description       *string        `json:"description" db:"description"`
	ReadyFrom         time.Time      `json:"ready_from" db:"ready_from"`
	ExpiresAt         time.Time      `json:"expires_at" db:"expires_at"`
	EstimatedServings int            `json:"estimated_servings" db:"estimated_servings"`
	Purpose           *string        `json:"purpose" db:"purpose"`
	Pincode           *string        `json:"pincode" db:"pincode"`
	City              *string        `json:"city" db:"city"`
	State             *string        `json:"state" db:"state"`
	Status            DonationStatus `json:"status" db:"status"`
	CreatedAt         time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at" db:"updated_at"`
}

// DonationClaim represents a claim on a donation offer
type DonationClaim struct {
	ID                uuid.UUID   `json:"id" db:"id"`
	OfferID           uuid.UUID   `json:"offer_id" db:"offer_id"`
	OrgID             uuid.UUID   `json:"org_id" db:"org_id"`
	RequestedServings int         `json:"requested_servings" db:"requested_servings"`
	PriorityScore     float64     `json:"priority_score" db:"priority_score"`
	Status            ClaimStatus `json:"status" db:"status"`
	CreatedAt         time.Time   `json:"created_at" db:"created_at"`
}

// Redemption represents a completed redemption
type Redemption struct {
	ID               uuid.UUID `json:"id" db:"id"`
	OfferID          uuid.UUID `json:"offer_id" db:"offer_id"`
	OrgID            uuid.UUID `json:"org_id" db:"org_id"`
	CollabID         uuid.UUID `json:"collab_id" db:"collab_id"`
	ServingsAccepted int       `json:"servings_accepted" db:"servings_accepted"`
	CreditsSpent     int       `json:"credits_spent" db:"credits_spent"`
	TokensAwarded    int       `json:"tokens_awarded" db:"tokens_awarded"`
	ConfirmedAt      time.Time `json:"confirmed_at" db:"confirmed_at"`
}

// RemoteOrgAssignment represents remote organization assignments
type RemoteOrgAssignment struct {
	ID          uuid.UUID `json:"id" db:"id"`
	RemoteOrgID uuid.UUID `json:"remote_org_id" db:"remote_org_id"`
	ProxyOrgID  uuid.UUID `json:"proxy_org_id" db:"proxy_org_id"`
	Active      bool      `json:"active" db:"active"`
	Notes       *string   `json:"notes" db:"notes"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// Event represents audit events
type Event struct {
	ID         uuid.UUID      `json:"id" db:"id"`
	EntityType string         `json:"entity_type" db:"entity_type"`
	EntityID   *uuid.UUID     `json:"entity_id" db:"entity_id"`
	EventType  string         `json:"event_type" db:"event_type"`
	Payload    map[string]any `json:"payload" db:"payload"`
	CreatedAt  time.Time      `json:"created_at" db:"created_at"`
}

// Enums
type UserRole string

const (
	UserRoleAdmin  UserRole = "ADMIN"
	UserRoleOrg    UserRole = "ORG"
	UserRoleCollab UserRole = "COLLAB"
)

type UserStatus string

const (
	UserStatusActive    UserStatus = "ACTIVE"
	UserStatusSuspended UserStatus = "SUSPENDED"
)

type DonationStatus string

const (
	DonationStatusOpen           DonationStatus = "OPEN"
	DonationStatusPendingConfirm DonationStatus = "PENDING_CONFIRM"
	DonationStatusClaimed        DonationStatus = "CLAIMED"
	DonationStatusExpired        DonationStatus = "EXPIRED"
	DonationStatusCancelled      DonationStatus = "CANCELLED"
)

type ClaimStatus string

const (
	ClaimStatusRequested ClaimStatus = "REQUESTED"
	ClaimStatusWon       ClaimStatus = "WON"
	ClaimStatusLost      ClaimStatus = "LOST"
	ClaimStatusCancelled ClaimStatus = "CANCELLED"
)

type Purpose string

const (
	PurposeChildren  Purpose = "CHILDREN"
	PurposeElderly   Purpose = "ELDERLY"
	PurposeWomen     Purpose = "WOMEN"
	PurposeGeneral   Purpose = "GENERAL"
	PurposeEmergency Purpose = "EMERGENCY"
)

type RedemptionStatus string

const (
	RedemptionStatusPending   RedemptionStatus = "PENDING"
	RedemptionStatusConfirmed RedemptionStatus = "CONFIRMED"
	RedemptionStatusCancelled RedemptionStatus = "CANCELLED"
)

type CreditTransactionType string

const (
	CreditTransactionTypeIssue CreditTransactionType = "ISSUE"
	CreditTransactionTypeSpend CreditTransactionType = "SPEND"
	CreditTransactionTypeRefund CreditTransactionType = "REFUND"
)

type TokenTransactionType string

const (
	TokenTransactionTypeEarn   TokenTransactionType = "EARN"
	TokenTransactionTypeRedeem TokenTransactionType = "REDEEM"
)

// Validation helpers
func (r UserRole) IsValid() bool {
	return r == UserRoleAdmin || r == UserRoleOrg || r == UserRoleCollab
}

func (s UserStatus) IsValid() bool {
	return s == UserStatusActive || s == UserStatusSuspended
}

func (s DonationStatus) IsValid() bool {
	return s == DonationStatusOpen || s == DonationStatusPendingConfirm ||
		s == DonationStatusClaimed || s == DonationStatusExpired ||
		s == DonationStatusCancelled
}

func (s ClaimStatus) IsValid() bool {
	return s == ClaimStatusRequested || s == ClaimStatusWon ||
		s == ClaimStatusLost || s == ClaimStatusCancelled
}

func (p Purpose) IsValid() bool {
	return p == PurposeChildren || p == PurposeElderly || p == PurposeWomen ||
		p == PurposeGeneral || p == PurposeEmergency
}
