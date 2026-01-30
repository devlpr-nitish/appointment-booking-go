package handlers

import (
	"fmt"
	"net/http"

	"github.com/devlpr-nitish/appointment-booking-go/internal/models"
	"github.com/devlpr-nitish/appointment-booking-go/internal/services"
	"github.com/devlpr-nitish/appointment-booking-go/internal/utils"
	"github.com/labstack/echo/v4"
)

type RegisterRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
	Name     string `json:"name"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
}

type ProfileUpdateRequest struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

func Register(c echo.Context) error {
	var req RegisterRequest

	if err := c.Bind(&req); err != nil {
		return utils.RespondError(c, http.StatusBadRequest, err, "Invalid request format")
	}

	if req.Email == "" || req.Password == "" {
		return utils.RespondError(c, http.StatusBadRequest, echo.NewHTTPError(http.StatusBadRequest, "Missing required field"), "email and password are required")
	}

	// Validate role
	if req.Role != "" && req.Role != "user" && req.Role != "expert" {
		return utils.RespondError(c, http.StatusBadRequest, echo.NewHTTPError(http.StatusBadRequest, "Invalid role"), "role must be either 'user' or 'expert'")
	}

	// Default to "user" if no role specified
	if req.Role == "" {
		req.Role = "user"
	}

	user, err := services.RegisterUser(req.Email, req.Password, req.Name, req.Role)

	if err != nil {
		return utils.RespondError(c, http.StatusInternalServerError, err, "Registration failed")
	}

	return utils.RespondSuccess(c, http.StatusCreated, "user registered successfully", user)
}

func Login(c echo.Context) error {
	var req LoginRequest

	if err := c.Bind(&req); err != nil {
		return utils.RespondError(c, http.StatusBadRequest, err, "Invalid request format")
	}

	if req.Email == "" || req.Password == "" {
		return utils.RespondError(c, http.StatusBadRequest, echo.NewHTTPError(http.StatusBadRequest, "Missing required field"), "email and password are required")
	}

	token, err := services.LoginUser(req.Email, req.Password)

	if err != nil {
		return utils.RespondError(c, http.StatusUnauthorized, err, "invalid email or password")
	}

	return utils.RespondSuccess(c, http.StatusOK, "user loggedin successfully", map[string]string{"token": token})
}

func UpdateProfile(c echo.Context) error {
	fmt.Println("-------------------------------------------")
	fmt.Println("UpdateProfile Request Received")
	fmt.Println("Content-Type:", c.Request().Header.Get("Content-Type"))

	// Parse multipart form
	// Max memory 10MB
	if err := c.Request().ParseMultipartForm(10 << 20); err != nil {
		fmt.Println("ParseMultipartForm error:", err)
	}

	name := c.FormValue("name")
	fmt.Println("Form name:", name)

	// Handle file upload
	file, header, err := c.Request().FormFile("image")
	var imageURL string

	if err == nil {
		fmt.Println("File received:", header.Filename)
		fmt.Println("File size:", header.Size)
		defer file.Close()
		// Upload to Cloudinary
		url, err := services.UploadToCloudinary(file, header.Filename)
		if err != nil {
			fmt.Println("Cloudinary upload error:", err)
			return utils.RespondError(c, http.StatusInternalServerError, err, "Failed to upload image")
		}
		imageURL = url
		fmt.Println("Cloudinary URL:", imageURL)
	} else {
		fmt.Println("FormFile error:", err)
		// If no file uploaded, check if image_url is provided in form (e.g. keeping existing)
		imageURL = c.FormValue("image_url")
		fmt.Println("Form image_url:", imageURL)
	}

	fmt.Println("imageURL: ", imageURL)
	if name == "" && imageURL == "" {
		return utils.RespondError(c, http.StatusBadRequest, echo.NewHTTPError(http.StatusBadRequest, "Missing required field"), "name or image/image_url is required")
	}

	authUser := c.Get("user").(*models.User)

	// If no new image, keep existing
	if imageURL == "" {
		imageURL = authUser.ImageURL
	}
	// If no new name, keep existing
	if name == "" {
		name = authUser.Name
	}

	user, err := services.UpdateProfile(authUser.ID, name, imageURL)

	if err != nil {
		return utils.RespondError(c, http.StatusInternalServerError, err, "Profile update failed")
	}

	// Generate new token with updated info
	token, err := utils.GenerateJWT(user.ID, user.Email, user.Name, string(user.Role), user.ImageURL)
	if err != nil {
		return utils.RespondError(c, http.StatusInternalServerError, err, "Failed to generate new token")
	}

	return utils.RespondSuccess(c, http.StatusOK, "profile updated successfully", map[string]interface{}{
		"user":  user,
		"token": token,
	})
}
