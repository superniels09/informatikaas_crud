package endpoints

import (
	"database/sql"
	"errors"
	"fmt"
	"informatikaas_crud/config"
	"informatikaas_crud/sql_queries"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// JWTSecret is the secret key used for signing JWT tokens.
var JWTSecret []byte

// UserCredentials holds the username and password for login.
type UserCredentials struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserRegistration holds user registration data.
type UserRegistration struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Email       string `json:"email" binding:"required"`
	DisplayName string `json:"display_name"`
	Bio         string `json:"bio"`
}

// JWTToken holds the JWT token and its expiration.
type JWTToken struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}

// AuthMiddleware validates the JWT token in the Authorization header.
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			SendError(ctx, http.StatusUnauthorized, "Authorization header is required")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			SendError(ctx, http.StatusUnauthorized, "Invalid Authorization header format")
			return
		}

		tokenString := parts[1]
		claims := &jwt.RegisteredClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return JWTSecret, nil
		})

		if err != nil || !token.Valid {
			SendError(ctx, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		// Attach the user ID from the token to the context for further use.
		userID := claims.Subject
		ctx.Set("userID", userID)

		ctx.Next()
	}
}

// InitJWT initializes the JWT secret key.
func InitJWT(cfg *config.Config) {
	JWTSecret = []byte(cfg.JWT.Secret)
}

// Login handles user login and returns a JWT token.
func Login(ctx *gin.Context, ed *EndpointData) {
	var creds UserCredentials
	if _, ok := BindJSON[UserCredentials](ctx); !ok {
		return
	}

	user, err := ed.DB.Queries.GetUserByUsername(ed.DB.Ctx, creds.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			SendError(ctx, http.StatusUnauthorized, "Invalid username or password")
			return
		}
		SendError(ctx, http.StatusInternalServerError, "Failed to fetch user")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Passhash), []byte(creds.Password)); err != nil {
		SendError(ctx, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		Subject:   fmt.Sprintf("%d", user.ID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	})

	tokenString, err := token.SignedString(JWTSecret)
	if err != nil {
		SendError(ctx, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	ctx.JSON(http.StatusOK, JWTToken{
		Token:     tokenString,
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	})
}

// Register handles user registration.
func Register(ctx *gin.Context, ed *EndpointData) {
	var reg UserRegistration
	if _, ok := BindJSON[UserRegistration](ctx); !ok {
		return
	}

	// Check if username or email already exists.
	_, err := ed.DB.Queries.GetUserByUsername(ed.DB.Ctx, reg.Username)
	if err == nil {
		SendError(ctx, http.StatusConflict, "Username already exists")
		return
	}

	_, err = ed.DB.Queries.GetUserByEmail(ed.DB.Ctx, reg.Email)
	if err == nil {
		SendError(ctx, http.StatusConflict, "Email already exists")
		return
	}

	// Hash the password.
	passhash, err := bcrypt.GenerateFromPassword([]byte(reg.Password), bcrypt.DefaultCost)
	if err != nil {
		SendError(ctx, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	// Create the user.
	userID, err := ed.DB.Queries.CreateUser(ed.DB.Ctx, sql_queries.CreateUserParams{
		Username:    reg.Username,
		Passhash:    string(passhash),
		Email:       reg.Email,
		DisplayName: sql.NullString{String: reg.DisplayName, Valid: true},
		Bio:         sql.NullString{String: reg.Bio, Valid: true},
	})

	if err != nil {
		SendError(ctx, http.StatusInternalServerError, "Failed to create user")
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user_id": userID,
	})
}
