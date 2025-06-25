package  utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(pwd string) (string, error) {
    b, err := bcrypt.GenerateFromPassword([]byte(pwd), 14)
    return string(b), err
}

func CheckPassword(hashed, pwd string) bool {
    return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(pwd)) == nil
}
