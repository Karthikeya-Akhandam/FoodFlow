package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/robfig/cron/v3"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

/*
-------------------------------

	Models
	-------------------------------
*/
type User struct {
	ID           uint `gorm:"primaryKey"`
	Name         string
	Email        string `gorm:"uniqueIndex"`
	PasswordHash string
	Role         string // collaborator | organization | admin
	CreatedAt    time.Time
}

type Organization struct {
	ID                 uint `gorm:"primaryKey"`
	UserID             uint
	User               User
	Credits            int
	PeopleFedLastMonth int
	ProxyAssignedToID  *uint // optional: proxy NGO user id
	IsProxy            bool  // if this org is serving as a proxy for remote orgs
	CreatedAt          time.Time
}

type Collaborator struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	User      User
	Tokens    int
	CreatedAt time.Time
}

type Donation struct {
	ID             uint `gorm:"primaryKey"`
	CollaboratorID uint
	Collaborator   Collaborator
	Title          string
	Description    string
	Quantity       int
	Purpose        string // optional purpose string
	Status         string // available | claimed
	ExpiryAt       *time.Time
	CreatedAt      time.Time
}

type Claim struct {
	ID          uint `gorm:"primaryKey"`
	DonationID  uint
	Donation    Donation
	OrgID       uint
	Org         Organization
	CreditsUsed int
	Priority    int
	CreatedAt   time.Time
}

/* -------------------------------
   Globals & Config
   ------------------------------- */

var db *gorm.DB
var jwtSecret = []byte("your-secure-jwt-key")

const CREDIT_FACTOR = 1
const TOKEN_FACTOR = 1

/* -------------------------------
   Init DB
   ------------------------------- */

func initDB() {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		// default local postgres dsn
		dsn = "host=localhost user=postgres password=your-password dbname=foodshare port=5432 sslmode=disable TimeZone=Asia/Kolkata"
	}
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect db:", err)
	}
	err = db.AutoMigrate(&User{}, &Organization{}, &Collaborator{}, &Donation{}, &Claim{})
	if err != nil {
		log.Fatal("migrate failed:", err)
	}
}

/* -------------------------------
   Helpers: password & JWT
   ------------------------------- */

func hashPassword(p string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	return string(b), err
}
func checkPasswordHash(hash, p string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(p))
}

func generateToken(user User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func parseToken(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return jwtSecret, nil
	})
}

/* -------------------------------
   Middleware
   ------------------------------- */

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing auth header"})
			return
		}
		tokenStr := auth
		token, err := parseToken(tokenStr)
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		c.Set("user_id", uint(claims["user_id"].(float64)))
		c.Set("role", claims["role"].(string))
		c.Next()
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		if role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin only"})
			return
		}
		c.Next()
	}
}

/* -------------------------------
   Handlers
   ------------------------------- */

func SignupHandler(c *gin.Context) {
	var req struct {
		Name               string `json:"name" binding:"required"`
		Email              string `json:"email" binding:"required"`
		Password           string `json:"password" binding:"required"`
		Role               string `json:"role" binding:"required"`
		PeopleFedLastMonth int    `json:"people_fed_last_month"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pwHash, _ := hashPassword(req.Password)
	user := User{Name: req.Name, Email: req.Email, PasswordHash: pwHash, Role: req.Role}
	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already exists or invalid"})
		return
	}

	if req.Role == "organization" {
		org := Organization{UserID: user.ID, Credits: 0, PeopleFedLastMonth: req.PeopleFedLastMonth}
		db.Create(&org)
	} else if req.Role == "collaborator" {
		col := Collaborator{UserID: user.ID, Tokens: 0}
		db.Create(&col)
	}

	token, _ := generateToken(user)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func LoginHandler(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user User
	if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	if err := checkPasswordHash(user.PasswordHash, req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	token, _ := generateToken(user)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func CreateDonationHandler(c *gin.Context) {
	// collaborator-only
	role, _ := c.Get("role")
	if role != "collaborator" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only collaborators"})
		return
	}
	userID := c.GetUint("user_id")
	var col Collaborator
	if err := db.Where("user_id = ?", userID).First(&col).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "collaborator profile not found"})
		return
	}
	var req struct {
		Title       string  `json:"title" binding:"required"`
		Description string  `json:"description"`
		Quantity    int     `json:"quantity" binding:"required"`
		Purpose     string  `json:"purpose"`
		ExpiryAt    *string `json:"expiry_at"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var expiry *time.Time
	if req.ExpiryAt != nil {
		t, err := time.Parse(time.RFC3339, *req.ExpiryAt)
		if err == nil {
			expiry = &t
		}
	}
	d := Donation{
		CollaboratorID: col.ID,
		Title:          req.Title,
		Description:    req.Description,
		Quantity:       req.Quantity,
		Purpose:        req.Purpose,
		Status:         "available",
		ExpiryAt:       expiry,
	}
	db.Create(&d)
	c.JSON(http.StatusOK, d)
}

func ListAvailableDonations(c *gin.Context) {
	var donations []Donation
	db.Preload("Collaborator.User").Where("status = ?", "available").Find(&donations)
	c.JSON(http.StatusOK, donations)
}

func ClaimDonationHandler(c *gin.Context) {
	// organization-only
	role, _ := c.Get("role")
	if role != "organization" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only orgs can claim"})
		return
	}
	userID := c.GetUint("user_id")
	var org Organization
	if err := db.Where("user_id = ?", userID).First(&org).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "org profile missing"})
		return
	}
	id := c.Param("id")
	var donation Donation
	if err := db.Preload("Collaborator").First(&donation, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "donation not found"})
		return
	}
	if donation.Status != "available" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "already claimed"})
		return
	}
	// minimal credits used logic:
	creditsNeeded := 1
	if org.Credits < creditsNeeded {
		c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient credits"})
		return
	}
	// deduct credits
	org.Credits -= creditsNeeded
	db.Save(&org)

	// create claim record
	claim := Claim{
		DonationID:  donation.ID,
		OrgID:       org.ID,
		CreditsUsed: creditsNeeded,
		Priority:    1,
	}
	db.Create(&claim)

	// mark donation claimed
	donation.Status = "claimed"
	db.Save(&donation)

	// reward collaborator tokens
	var coll Collaborator
	db.First(&coll, donation.CollaboratorID)
	coll.Tokens += creditsNeeded * TOKEN_FACTOR
	db.Save(&coll)

	c.JSON(http.StatusOK, gin.H{"ok": true, "claim": claim})
}

func MeHandler(c *gin.Context) {
	userID := c.GetUint("user_id")
	var user User
	db.First(&user, userID)
	var resp = gin.H{"user": user}
	if user.Role == "organization" {
		var org Organization
		db.Where("user_id = ?", userID).First(&org)
		resp["organization"] = org
	} else if user.Role == "collaborator" {
		var coll Collaborator
		db.Where("user_id = ?", userID).First(&coll)
		resp["collaborator"] = coll
	}
	c.JSON(http.StatusOK, resp)
}

/* Admin endpoints */

func AdminAssignProxy(c *gin.Context) {
	var req struct {
		OrgID        uint `json:"org_id" binding:"required"`
		ProxyOrgID   uint `json:"proxy_org_id" binding:"required"`
		ExtraCredits int  `json:"extra_credits"` // optional
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var org Organization
	if err := db.First(&org, req.OrgID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "org not found"})
		return
	}
	org.ProxyAssignedToID = &req.ProxyOrgID
	db.Save(&org)

	// give extra credits to proxy as reward if elected
	if req.ExtraCredits > 0 {
		var proxy Organization
		if err := db.First(&proxy, req.ProxyOrgID).Error; err == nil {
			proxy.Credits += req.ExtraCredits
			db.Save(&proxy)
		}
	}
	c.JSON(http.StatusOK, org)
}

func AdminResetCredits(c *gin.Context) {
	// immediate manual reset + assignment for all orgs
	var orgs []Organization
	db.Find(&orgs)
	for _, o := range orgs {
		o.Credits = o.PeopleFedLastMonth * CREDIT_FACTOR
		db.Save(&o)
	}
	c.JSON(http.StatusOK, gin.H{"ok": true, "count": len(orgs)})
}

/* -------------------------------
   Monthly cron for credits reset
   ------------------------------- */

func startCron() {
	c := cron.New(cron.WithLocation(time.FixedZone("IST", 5*3600+1800))) // IST +05:30
	// Run at midnight on first day of month (cron spec: second minute hour day month weekday)
	// robfig/cron uses 6-field by default: sec min hour day month weekday
	_, err := c.AddFunc("0 0 0 1 * *", func() {
		log.Println("[cron] monthly reset running")
		var orgs []Organization
		db.Find(&orgs)
		for _, o := range orgs {
			o.Credits = o.PeopleFedLastMonth * CREDIT_FACTOR
			db.Save(&o)
		}
	})
	if err != nil {
		log.Println("cron err:", err)
	} else {
		c.Start()
	}
}

/* -------------------------------
   Main
   ------------------------------- */

func main() {
	initDB()

	// create default admin if none
	var admin User
	if err := db.Where("role = ?", "admin").First(&admin).Error; err != nil {
		pw, _ := hashPassword("admin123")
		admin = User{Name: "admin", Email: "admin@local.in", PasswordHash: pw, Role: "admin"}
		db.Create(&admin)
	}

	startCron()

	r := gin.Default()

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// public
	r.POST("/auth/signup", SignupHandler)
	r.POST("/auth/login", LoginHandler)
	r.GET("/donations", ListAvailableDonations)

	// protected
	auth := r.Group("/", AuthMiddleware())
	auth.POST("/donations", CreateDonationHandler)
	auth.POST("/donations/:id/claim", ClaimDonationHandler)
	auth.GET("/me", MeHandler)

	// admin
	adminGroup := r.Group("/admin", AuthMiddleware(), AdminOnly())
	adminGroup.POST("/assign-proxy", AdminAssignProxy)
	adminGroup.POST("/reset-credits", AdminResetCredits)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("listening on", port)
	r.Run(":" + port)
}
