package usecase

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Senoue/aws-rds-apprunner-with-terraform/domain/models"
	"github.com/Senoue/aws-rds-apprunner-with-terraform/domain/repositories"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthUsecase handles business logic related to authentication.
type AuthUsecase struct {
	authRepo repository.AuthRepository
}

// NewAuthUsecase creates a new AuthUsecase instance.
func NewAuthUsecase(authRepo repository.AuthRepository) *AuthUsecase {
	return &AuthUsecase{authRepo: authRepo}
}

var SECRET_KEY = os.Getenv("TOKEN_SECRET")

// Login authenticates a user and returns JSON response.
func (a *AuthUsecase) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := a.authRepo.Login(c, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	tokenString, refreshTokenString, err := a.generateTokens(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Authentication successful",
		"data": map[string]string{
			"token":         tokenString,
			"refresh_token": refreshTokenString,
		},
	})
}

// Register creates a new user and returns JSON response.
func (a *AuthUsecase) Register(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := a.authRepo.Register(c, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

// UserInfo returns user information.
func (a *AuthUsecase) UserInfo(c *gin.Context) {
	u, err := a.getUserFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	user, err := a.authRepo.UserInfo(c, u.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user information"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (a *AuthUsecase) generateTokens(user *model.User) (string, string, error) {
	token, err := a.createToken(user, time.Hour*24, SECRET_KEY)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := a.createToken(user, time.Hour*24*30, SECRET_KEY)
	if err != nil {
		return "", "", err
	}

	return token, refreshToken, nil
}

func (a *AuthUsecase) createToken(user *model.User, duration time.Duration, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.Name,
		"exp":      time.Now().Add(duration).Unix(),
	})
	return token.SignedString([]byte(secretKey))
}

// JWTからユーザー情報を取得
func (a *AuthUsecase) getUserFromToken(c *gin.Context) (*model.User, error) {
	// Authorizationヘッダーからトークンを取得
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, jwt.ErrSignatureInvalid
	}

	// "Bearer "を取り除いてトークンを取得
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	user := &model.User{
		ID:   int(claims["id"].(float64)),
		Name: claims["username"].(string),
	}

	return user, nil
}

func (a *AuthUsecase) AuthMiddleware(c *gin.Context) {
	// Authorizationヘッダーからトークンを取得
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
		c.Abort()
		return
	}

	// "Bearer "を取り除いてトークンを取得
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// トークンの検証
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SECRET_KEY), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	c.Next()
}
