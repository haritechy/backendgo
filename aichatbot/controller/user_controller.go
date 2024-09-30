package controller

import (
	"fmt"
	"slack-chatbot/database"
	"slack-chatbot/models"
	"slack-chatbot/utils"

	"net/http"
	"os"
	"strings"

	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

var jwtkey = []byte(os.Getenv("JWT_KEY"))
var logger = logrus.New()

func init() {
	// Set up logging configuration
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logger.SetOutput(os.Stdout)
}

func GenerateJwt(email string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &models.Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtkey)
}

func UserRegister(c *gin.Context) {
	var UserRegister models.User

	if err := c.BindJSON(&UserRegister); err != nil {
		logger.Errorf("Error binding JSON: %v", err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := utils.ValidateEmail(UserRegister.Email); err != nil {

		logger.Errorf("Invalid email fomartd %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email"})
		return

	}

	var existingUser models.User
	if err := database.DB.Where("email = ?", UserRegister.Email).First(&existingUser).Error; err == nil {
		logger.Errorf("Email already registered: %v", UserRegister.Email)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is already registered"})
		return
	}

	if err := utils.Validatepassword(strings.TrimSpace(UserRegister.Password)); err != nil {

		logger.Errorf("Error hashing password:%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(UserRegister.Password), bcrypt.DefaultCost)

	if err != nil {
		logger.Errorf("Error hashing password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	UserRegister.Password = string(hashedPassword)
	var ExistingPassowrd models.User
	if err := database.DB.Where("password = ?", string(hashedPassword)).First(&ExistingPassowrd).Error; err == nil {

		logger.Errorf("Password is already taken by  user: %v", UserRegister.Password)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password is already taken by  user"})
		return
	}
	result := database.DB.Create(&UserRegister)
	if result.Error != nil {
		logger.Errorf("Error creating user: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	logger.Infof("User registered successfully: %v", UserRegister)
	c.JSON(http.StatusOK, &UserRegister)
}

func UserGet(c *gin.Context) {
	var UserGet []models.User

	result := database.DB.Find(&UserGet)
	if result.Error != nil {
		logger.Errorf("Error retrieving users: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if len(UserGet) == 0 {
		logger.Warn("No users found")
		c.JSON(http.StatusBadRequest, gin.H{"error": "No data found"})
		return
	}

	logger.Infof("Retrieved users: %v", UserGet)
	c.JSON(http.StatusOK, &UserGet)
}

func UseDelete(c *gin.Context) {
	param := c.Param("id")
	var userData models.User

	result := database.DB.First(&userData, param)
	if result.Error != nil {
		logger.Errorf("Error finding user: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	database.DB.Delete(&userData, param)
	logger.Infof("User deleted successfully: %s", param)
	c.JSON(http.StatusOK, "User delete successful")
}

func UserUpdate(c *gin.Context) {
	id := c.Param("id")
	body := models.User{}
	var updateUser models.User

	result := database.DB.First(&updateUser, id)
	if result.Error != nil {
		logger.Errorf("Error finding user for update: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if err := c.BindJSON(&body); err != nil {
		logger.Errorf("Error binding JSON for update: %v", err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	updateUser.FirstName = body.FirstName
	updateUser.LastName = body.LastName
	updateUser.Email = body.Email
	updateUser.Password = body.Password

	database.DB.Save(&updateUser)
	logger.Infof("User updated successfully: %v", updateUser)
	c.JSON(http.StatusOK, &updateUser)
}

func UserGetbyEmail(c *gin.Context) {
	var user []models.User
	id := c.Param("id")

	getbyId := database.DB.First(&user, id)
	if getbyId.Error != nil {
		logger.Errorf("Error finding user by email: %v", getbyId.Error)
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	logger.Infof("User retrieved by email: %v", user)
	c.JSON(http.StatusOK, &user)
}

func UserLogin(c *gin.Context) {
	var user models.User
	var logindata struct {
		Email    string
		Password string
	}

	if err := c.BindJSON(&logindata); err != nil {
		logger.Errorf("Error binding JSON for login: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	result := database.DB.Where("email=?", logindata.Email).First(&user)
	if result.Error != nil {
		logger.Errorf("Error finding user for login: %v", result.Error)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(logindata.Password)); err != nil {
		logger.Errorf("Password mismatch for user: %v", logindata.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := GenerateJwt(user.Email)
	if err != nil {
		logger.Errorf("Error generating JWT: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	logger.Infof("User logged in successfully: %v", user.Email)
	c.JSON(http.StatusOK, gin.H{"message": "userlogin succesful",
		"login_token": token,
	})
}
func UserForgot(c *gin.Context) {
	var Users models.User

	var Forgotdata struct {
		Email string
	}
	var randowCodes = [...]byte{
		'1', '2', '3', '4', '5', '6', '7', '8', '9', '0',
	}
	var r *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	var pwd []byte = make([]byte, 6)
	for i := 0; i < 3; i++ {

		for j := 0; j < 6; j++ {
			index := r.Int() % len(randowCodes)

			pwd[j] = randowCodes[index]
		}

		fmt.Printf("%s\n", string(pwd))
	}

	if err := c.BindJSON(&Forgotdata); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json type"})
		return

	}
	result := database.DB.Where("email=?", Forgotdata.Email).First(&Users)
	if result.Error != nil {
		logger.Errorf("Error sending email: %v", result.Error)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email does not find your db"})
		return
	}

	otp := string(pwd)
	utils.SendEmail(Users.Email, "Forgot password", otp)
	Users.Otp = otp

	if err := database.DB.Model(&Users).Update("Otp", Users.Otp).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update OTP"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent to email"})

}

func ChnagePassword(c *gin.Context) {
	var Users models.User
	var PasswordChange struct {
		Otp        string `json:"otp"`
		NewPasword string `json:"newpassword"`
	}

	if err := c.BindJSON(&PasswordChange); err != nil {
		logger.Errorf("error invalid body  or filed missing %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json type"})
		return

	}

	result := database.DB.Where("Otp=?", PasswordChange.Otp).First(&Users)
	if result.Error != nil {
		logger.Errorf("Error : %v", result.Error)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Otp not found or invalid Otp"})
		return
	}
	utils.Validatepassword(PasswordChange.NewPasword)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(PasswordChange.NewPasword), bcrypt.DefaultCost)
	if err != nil {
		logger.Errorf("Error hashing password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	if err := database.DB.Model(&Users).Update("password", hashedPassword).Error; err != nil {
		c.JSON(http.StatusBadRequest, "error updating password")
		return
	}

	if err := database.DB.Model(&Users).Update("Otp", "").Error; err != nil {
		logger.Errorf("Error clearing OTP: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear OTP"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "password changed succeful"})

}
