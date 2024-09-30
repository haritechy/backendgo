package models

import (
	"time"

	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}
type Prompt struct {
	Prompt string `json:"prompt"`
}

type PromptResponse struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Prompt    string    `json:"prompt"`
	Response  string    `json:"response"`
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	gorm.Model

	FirstName string `json:"firstname"`
	LastName  string `json:"lastname" `
	FullName  string `json:"fullname"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Otp       string `json:"otp"`
}

type Claims struct {
	Email string `json:"email"  binding:"required" gorm:"not null"`

	jwt.StandardClaims
}
type Otp struct {
	Otp string
}
