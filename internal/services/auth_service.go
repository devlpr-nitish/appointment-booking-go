package services

import (
	"errors"

	"github.com/devlpr-nitish/appointment-booking-go/internal/database"
	"github.com/devlpr-nitish/appointment-booking-go/internal/models"
	"github.com/devlpr-nitish/appointment-booking-go/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(email, password, name, role string) (*models.User, error) {
	db := database.GetDB()

	// Check if user already exists
	var existingUser models.User
	if err := db.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return nil, errors.New("user with this email or username already exists")
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	// Set role, default to "user" if empty
	var userRole models.UserRole
	switch role {
	case "expert":
		userRole = models.RoleExpert
	case "admin":
		userRole = models.RoleAdmin
	default:
		userRole = models.RoleUser
	}

	user := models.User{
		Email:    email,
		Password: hashedPassword,
		Name:     name,
		Role:     userRole,
	}

	if err := db.Create(&user).Error; err != nil {
		return nil, err
	}

	// If user is registering as expert, create expert profile
	if userRole == models.RoleExpert {
		expert := models.Expert{
			UserID:     user.ID,
			Bio:        "",
			Expertise:  "",
			HourlyRate: 0,
			IsVerified: false,
		}
		if err := db.Create(&expert).Error; err != nil {
			// Log error but don't fail registration
			// The expert can complete their profile later
		}
	}

	return &user, nil
}

func LoginUser(identifier, password string) (string, error) {
	db := database.GetDB()
	var user models.User

	// Find user by email
	if err := db.Where("email = ?", identifier).First(&user).Error; err != nil {
		return "", errors.New("invalid credentials")
	}

	if !checkPasswordHash(password, user.Password) {
		return "", errors.New("invalid credentials")
	}

	token, err := utils.GenerateJWT(user.ID, user.Email, user.Name, string(user.Role))
	if err != nil {
		return "", err
	}

	return token, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
