package utils

import (
	"fmt"
	"regexp"
)

func ValidateEmail(email string) error {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !re.MatchString(email) {
		return fmt.Errorf("invalid Email format")
	}
	return nil
}
func Validatepassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("Password must be at least 8 character")
	}
	uppercaseRegex := regexp.MustCompile(`[A-Z]`)
	if !uppercaseRegex.MatchString(password) {
		return fmt.Errorf("Password contain atleast on Upper character")
	}
	specialcaseRegex := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`)
	if !specialcaseRegex.MatchString(password) {
		return fmt.Errorf("Password must contain at lease on special character")
	}
	return nil

}
