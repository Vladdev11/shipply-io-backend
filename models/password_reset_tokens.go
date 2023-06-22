package models

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"gorm.io/gorm"
)

type PasswordResetToken struct {
	ID        int
	UserID    int
	Token     string
	ExpiresAt time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	User User
}

func (prt *PasswordResetToken) Create() error {
	err := PGDB.Create(prt).Error
	return err
}

func (prt *PasswordResetToken) Update() error {
	err := PGDB.Save(prt).Error
	return err
}

func (prt *PasswordResetToken) Delete() error {
	err := PGDB.Delete(prt).Error
	return err
}

func GetPasswordResetTokenByToken(token string) (PasswordResetToken, error) {
	var passwordResetToken PasswordResetToken
	err := PGDB.Where("token = ?", token).First(&passwordResetToken).Error
	return passwordResetToken, err
}

func GenerateToken() string {
	tokenLength := 64
	tokenBytes := make([]byte, tokenLength)

	_, err := rand.Read(tokenBytes)
	if err != nil {
		return ""
	}

	token := base64.URLEncoding.EncodeToString(tokenBytes)
	return token
}
