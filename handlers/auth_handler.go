package handlers

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/reemployed/reemployed/repositories"
	"golang.org/x/crypto/bcrypt"
)

type AuthController interface {
	Login(ctx *gin.Context)
}

type authController struct {
	userRepo repositories.UserRepository
}

func NewAuthController(userRepo repositories.UserRepository) AuthController {
	return &authController{userRepo: userRepo}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @BasePath /api/v1/
// Login godoc
// @Summary Login with email and password
// @Description Login with email and password
// @Tags auth
// @Accept  json
// @Produce  json
// @Param request body LoginRequest true "Login Request"
// @Success 200 {string} string "Success"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /auth/login [post]
func (c *authController) Login(ctx *gin.Context) {
	// Get the username and password from the request body
	var loginReq LoginRequest
	if err := ctx.ShouldBindJSON(&loginReq); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Find the user in the JSON file
	user, err := c.userRepo.GetUserByEmail(loginReq.Email)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Check the password against the hashed password in the JSON file
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate a new JWT token
	token, err := generateToken(user.Email)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Return the token as a response
	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func generateToken(email string) (string, error) {
	// Set the token claims
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}

	// Generate the JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("secret"))
}
