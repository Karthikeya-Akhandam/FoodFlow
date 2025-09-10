package core

// Matching and scoring constants
const (
	// Proximity weights
	ProximityWeightSamePincode = 1.0
	ProximityWeightSameCity    = 0.7
	ProximityWeightSameState   = 0.4

	// Scoring weights (must sum to 100)
	ScoreWeightProximity      = 40
	ScoreWeightPurposeMatch   = 25
	ScoreWeightCreditPressure = 20
	ScoreWeightReliability    = 10
	ScoreWeightRemoteProxy    = 5

	// Default values
	DefaultServingsPerCredit   = 10
	DefaultTokenMultiplier     = 5
	DefaultClaimCooloffSeconds = 600
	DefaultMonthlyCreditsBase  = 50
	DefaultMonthlyCreditsK     = 0.2
	DefaultRemoteProxyBonus    = 0.1

	// Pagination
	DefaultPageLimit = 20
	MaxPageLimit     = 100

	// JWT
	DefaultJWTExpiration = "24h"
	DefaultJWTIssuer     = "foodflow"

	// Rate limiting
	DefaultRateLimitRPS   = 10
	DefaultRateLimitBurst = 20

	// Idempotency
	DefaultIdempotencyTTL = "24h"
)

// Proximity tiers
type ProximityTier string

const (
	ProximityTierPincode ProximityTier = "pincode"
	ProximityTierCity    ProximityTier = "city"
	ProximityTierState   ProximityTier = "state"
)

// Organization types
const (
	OrgTypeNGO       = "NGO"
	OrgTypeOrphanage = "Orphanage"
	OrgTypeShelter   = "Shelter"
	OrgTypeOther     = "Other"
)

// Collaborator types
const (
	CollabTypeRestaurant = "Restaurant"
	CollabTypeHotel      = "Hotel"
	CollabTypeHostel     = "Hostel"
	CollabTypePG         = "PG"
	CollabTypeCaterer    = "Caterer"
	CollabTypeOther      = "Other"
)

// Entity types for audit events
const (
	EntityTypeUser       = "user"
	EntityTypeProfile    = "profile"
	EntityTypeOrg        = "organization"
	EntityTypeCollab     = "collaborator"
	EntityTypeOffer      = "offer"
	EntityTypeClaim      = "claim"
	EntityTypeCredit     = "credit"
	EntityTypeToken      = "token"
	EntityTypeRedemption = "redemption"
	EntityTypeRemote     = "remote_assignment"
)

// Event types for audit events
const (
	EventTypeCreated       = "created"
	EventTypeUpdated       = "updated"
	EventTypeDeleted       = "deleted"
	EventTypeStatusChanged = "status_changed"
	EventTypeRedeemed      = "redeemed"
	EventTypeClaimed       = "claimed"
	EventTypeAssigned      = "assigned"
	EventTypeCancelled     = "cancelled"
	EventTypeExpired       = "expired"
)

// Cache keys
const (
	CacheKeyUserSession  = "user_session:%s"
	CacheKeyRateLimit    = "rate_limit:%s"
	CacheKeyIdempotency  = "idempotency:%s"
	CacheKeyOrgCredits   = "org_credits:%s:%s"   // org_id:month
	CacheKeyCollabTokens = "collab_tokens:%s:%s" // collab_id:month
)

// Database constraints
const (
	MaxEmailLength       = 255
	MaxPasswordLength    = 255
	MaxNameLength        = 100
	MaxPhoneLength       = 15
	MaxPincodeLength     = 6
	MaxCityLength        = 50
	MaxDistrictLength    = 50
	MaxStateLength       = 50
	MaxAddressLength     = 200
	MaxTitleLength       = 100
	MaxDescriptionLength = 500
	MaxNotesLength       = 500
	MaxOrgTypeLength     = 50
	MaxCollabTypeLength  = 50
	MaxPurposeLength     = 20
	MaxMonthLength       = 6
	MaxServings          = 10000
	MaxCredits           = 10000
	MaxTokens            = 10000
)

// Validation rules
const (
	MinPasswordLength = 8
	MinNameLength     = 2
	MinTitleLength    = 5
	MinServings       = 1
	MinCredits        = 0
	MinTokens         = 0
	MinPeopleFed      = 0
)
