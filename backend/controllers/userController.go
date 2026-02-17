package controllers

import (
	"backend/initializers"
	"backend/models"
	"net/http"
	"net/mail"
	"os"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	var body struct {
		Username string
		Password string
		Email    string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	if body.Username == "" || body.Password == "" || body.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "All fields should be completed",
		})
		return
	}

	var existent_user models.User
	initializers.DB.First(&existent_user, "username = ?", body.Username)
	if existent_user.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Username already exists",
		})
		return
	}

	_, err := mail.ParseAddress(body.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email is in an incorrect format",
		})
		return
	}

	if len(body.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Password must contain at least 8 characters",
		})
		return
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool

	for _, char := range body.Password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must contain at least one uppercase letter"})
		return
	}

	if !hasLower {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must contain at least one lowercase letter"})
		return
	}

	if !hasNumber {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must contain at least one number"})
		return
	}

	if !hasSpecial {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must contain at least one special character"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to generate password hash",
		})
		return
	}

	user := models.User{
		Username: body.Username,
		Password: string(hash),
		Email:    body.Email,
	}

	res := initializers.DB.Create(&user)

	if res.Error != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})

}

func Login(c *gin.Context) {
	var body struct {
		Username string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	var user models.User
	initializers.DB.First(&user, "username = ?", body.Username)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User not found",
		})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Incorrect password",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}

func ShowUser(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func GetUsers(c *gin.Context) {
	currUser, _ := c.Get("user")

	// For admins only
	if currUser.(models.User).Role != "admin" {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	type returnedAccount struct {
		ID       uint
		Username string
		Email    string
		Role     string
	}

	var users []models.User
	initializers.DB.Find(&users, "id != ?", currUser.(models.User).ID)

	var accounts []returnedAccount
	for i := range users {
		account := returnedAccount{
			ID:       users[i].ID,
			Username: users[i].Username,
			Email:    users[i].Email,
			Role:     users[i].Role,
		}
		accounts = append(accounts, account)
	}

	c.JSON(http.StatusOK, gin.H{
		"accounts": accounts,
	})
}

func ModifyUser(c *gin.Context) {
	currUser, _ := c.Get("user")

	var body struct {
		Username string
		Password string
		Email    string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	var existent_user models.User
	initializers.DB.First(&existent_user, "username = ?", body.Username)
	if existent_user.ID != 0 && existent_user.ID != currUser.(models.User).ID {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Username already exists",
		})
		return
	}

	_, err1 := mail.ParseAddress(body.Email)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email is in an incorrect format",
		})
		return
	}

	var user models.User

	initializers.DB.First(&user, currUser.(models.User).ID)

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Incorrect password",
		})
		return
	}

	user.Email = body.Email
	user.Username = body.Username

	initializers.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})

}

func ToggleRole(c *gin.Context) {
	currUser, _ := c.Get("user")

	// For admins only
	if currUser.(models.User).Role != "admin" {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	id := c.Param("id")

	var user models.User
	initializers.DB.First(&user, "id = ?", id)

	if user.Role == "admin" {
		user.Role = "user"
	} else {
		user.Role = "admin"
	}

	initializers.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func DeleteUser(c *gin.Context) {
	currUser, _ := c.Get("user")

	// For admins only
	if currUser.(models.User).Role != "admin" {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	id := c.Param("id")

	var devices []models.Device
	var dashboards []models.Dashboard

	initializers.DB.Find(&devices, "user_id", id)
	for i := range devices {
		initializers.DB.Unscoped().Delete(&models.Alert{}, "device_id = ?", devices[i].ID)
		initializers.DB.Unscoped().Delete(&models.Widget{}, "device_id = ?", devices[i].ID)
		initializers.DB.Unscoped().Delete(&models.AnomalyModel{}, "device_id = ?", devices[i].ID)
	}

	initializers.DB.Find(&dashboards, "user_id", id)
	for i := range dashboards {
		initializers.DB.Unscoped().Delete(&models.Widget{}, "dashboard_id = ?", dashboards[i].ID)
	}

	initializers.DB.Unscoped().Delete(&models.Dashboard{}, "user_id = ?", id)
	initializers.DB.Unscoped().Delete(&models.Device{}, "user_id = ?", id)

	initializers.DB.Unscoped().Delete(&models.User{}, id)

	c.JSON(http.StatusOK, gin.H{
		"message": "user deleted",
	})
}
