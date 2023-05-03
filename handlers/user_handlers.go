package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/reemployed/reemployed/models"
	"github.com/reemployed/reemployed/repositories"
)

// @BasePath /api/v1
// UserController handles HTTP requests related to the user endpoints.
type UserController struct {
	repo repositories.UserRepository
}

// NewUserController creates a new instance of UserController.
func NewUserController(repo repositories.UserRepository) *UserController {
	return &UserController{repo: repo}
}

// @Tags users
// @Summary List all users
// @Description List all users.
// @Produce json
// @Success 200 {array} models.User
// @Router /users [get]
func (ctrl *UserController) GetUserList(c *gin.Context) {
	users, err := ctrl.repo.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, users)
}

// @Tags users
// @Summary Get a user by ID
// @Description Get a user by ID.
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Failure 404 {string} string "User not found"
// @Router /users/{id} [get]
func (ctrl *UserController) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	user, err := ctrl.repo.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// @Tags users
// @Summary Create a new user
// @Description Create a new user.
// @Accept json
// @Produce json
// @Param user body models.User true "models.User object"
// @Success 201 {object} models.User
// @Failure 400 {string} string "Invalid JSON format"
// @Router /users [post]
func (ctrl *UserController) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}
	if err := ctrl.repo.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

// @Tags users
// @Summary Update a user
// @Description Update a user.
// @Accept json
// @Produce json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param id path string true "User ID"
// @Param user body models.User true "models.User object"
// @Success 200 {object} models.User
// @Failure 400 {string} string "Invalid JSON format"
// @Failure 404 {string} string "User not found"
// @Router /users/{id} [put]
func (ctrl *UserController) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	// Get the user from the JSON file
	user, err := ctrl.repo.GetUserByID(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Check if the user is authorized to update this user
	if email, ok := getUserEmailFromToken(c); !ok || email != user.Email {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var updatedUser models.User
	if err := c.BindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	updatedUser.ID = user.ID
	if err := ctrl.repo.UpdateUser(&updatedUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (ctrl *UserController) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.repo.DeleteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func getUserEmailFromToken(c *gin.Context) (string, bool) {
	// Get the authorization header from the request headers
	authHeader := c.GetHeader("Authorization")

	// Parse the token from the authorization header
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
		return "", false
	}
	tokenString := tokenParts[1]

	// Parse the token claims to get the email address
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid token signing method")
		}
		return []byte("secret"), nil // Use the same secret key as used to generate the token
	})
	if err != nil || !token.Valid {
		return "", false
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", false
	}
	email, ok := claims["email"].(string)
	if !ok {
		return "", false
	}

	return email, true
}
