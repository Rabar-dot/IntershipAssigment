package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword hashes the password using bcrypt and returns the hashed password.
func HashPassword(pwd string) (string, error) {
    b, err := bcrypt.GenerateFromPassword([]byte(pwd), 14) // Hashing with a cost of 14
    return string(b), err
}

// CheckPassword compares a hashed password with a plain password.
func CheckPassword(hashed, pwd string) bool {
    return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(pwd)) == nil // Returns true if passwords match
}
