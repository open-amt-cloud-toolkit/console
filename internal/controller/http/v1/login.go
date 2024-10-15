package v1

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/open-amt-cloud-toolkit/console/config"
	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto/v1"
	"github.com/open-amt-cloud-toolkit/console/pkg/consoleerrors"
)

var ErrLogin = consoleerrors.CreateConsoleError("LoginHandler")

const hoursInADay = 24

type LoginRoute struct {
	Config *config.Config
}

// NewVersionRoute creates a new version route
func NewLoginRoute(configData *config.Config) *LoginRoute {
	return &LoginRoute{
		Config: configData,
	}
}

// FetchLatestRelease fetches the latest release information from GitHub API
func (lr LoginRoute) Login(c *gin.Context) {
	var creds dto.Credentials

	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})

		return
	}

	if creds.Username != lr.Config.AdminUsername || creds.Password != lr.Config.AdminPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})

		return
	}

	// Create JWT token
	expirationTime := time.Now().Add(hoursInADay * time.Hour)
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(lr.Config.App.JWTKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create token"})

		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// JWT Middleware
func (lr LoginRoute) JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "request does not contain an access token"})
			c.Abort()

			return
		}

		claims := &jwt.MapClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(_ *jwt.Token) (interface{}, error) {
			return []byte(lr.Config.App.JWTKey), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid access token"})
			c.Abort()

			return
		}

		c.Next()
	}
}
