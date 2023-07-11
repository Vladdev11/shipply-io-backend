package models

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/shipply-io/shipply-io-backend/util"
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

func (prt *PasswordResetToken) Create(ctx context.Context) error {
	err := util.DBFromContext(ctx).Create(prt).Error
	return err
}

func (prt *PasswordResetToken) Update(ctx context.Context) error {
	err := util.DBFromContext(ctx).Save(prt).Error
	return err
}

func (prt *PasswordResetToken) Delete(ctx context.Context) error {
	err := util.DBFromContext(ctx).Delete(prt).Error
	return err
}

func GetPasswordResetTokenByToken(ctx context.Context, token string) (PasswordResetToken, error) {
	var passwordResetToken PasswordResetToken
	err := util.DBFromContext(ctx).Where("token = ?", token).First(&passwordResetToken).Error
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
