package helper

import (
	"crypto/md5"
	"encoding/hex"
	"regexp"
)

func ValidateUserName(userName string) bool {
	matched, err := regexp.MatchString(`^[a-zа-яёA-ZA-ЯЁ0-9]{2,}$`, userName)
	if err != nil || !matched {
		return false
	}
	return true
}

func ValidateEmail(email string) bool {
	matched, err := regexp.MatchString(`^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`, email)
	if err != nil || !matched {
		return false
	}
	return true
}

func ValidatePassword(password string) bool {
	matched, err := regexp.MatchString(`^[a-zA-Z0-9$~@#%*!&?=()]{8,}$`, password)
	if err != nil || !matched {
		return false
	}
	return true
}

func HashPassword(password string) string {
	hash := md5.Sum([]byte(password))
	return hex.EncodeToString(hash[:])
}

func Hash(toHash []byte) string {
	hash := md5.Sum(toHash)
	return hex.EncodeToString(hash[:])
}
