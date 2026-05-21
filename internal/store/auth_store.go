package store

import (
	"context"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/ariefzainuri96/ayo-test/cmd/api/dto/entity"
	"github.com/ariefzainuri96/ayo-test/cmd/api/dto/request"
	db "github.com/ariefzainuri96/ayo-test/internal/db"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthStore struct {
	gormDb *db.GormDB
}

func (store *AuthStore) Login(ctx context.Context, body request.LoginRequest) (entity.User, string, error) {
	var user entity.User

	err := store.gormDb.ExecWithTimeoutErr(ctx, func(tx *gorm.DB) error {
		return tx.
			Model(&entity.User{}).
			Where("email = ?", body.Email).
			First(&user).Error
	})

	if err != nil {
		return entity.User{}, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		return entity.User{}, "", errors.New("invalid email or password")
	}

	token, err := generateToken(user.Email, user.Role, int(user.ID))

	if err != nil {
		return entity.User{}, "", err
	}

	return user, token, nil
}

func generateToken(email string, role string, id int) (string, error) {
	jwtSecret := strings.TrimSpace(os.Getenv("SECRET_KEY"))

	claims := jwt.MapClaims{
		"user_id": id,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24 * 30).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func (store *AuthStore) Register(ctx context.Context, body request.RegisterRequest) (uint, error) {
	var emaiExists bool

	err := store.gormDb.ExecWithTimeoutErr(ctx, func(tx *gorm.DB) error {
		return tx.
			Model(&entity.User{}).
			Select("1").
			Where("email = ?", body.Email).
			Limit(1).
			Scan(&emaiExists).
			Error
	})

	if err != nil {
		return 0, err
	}

	if emaiExists {
		return 0, errors.New("email sudah terdaftar")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)

	if err != nil {
		return 0, err
	}

	user := entity.User{
		Email:    body.Email,
		Password: string(hashedPassword),
		Role:     body.Role,
	}

	err = store.gormDb.ExecWithTimeoutErr(ctx, func(tx *gorm.DB) error {
		return tx.Create(&user).Error
	})

	if err != nil {
		return 0, err
	}

	return user.ID, nil
}