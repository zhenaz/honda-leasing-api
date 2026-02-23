package service

import (
	"context"
	"errors"
	"honda-leasing-api/internal/domain/query"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type AuthService interface {
	Login(ctx context.Context, phone string) (string, error)
}

type authService struct {
	db        *gorm.DB
	jwtSecret string
}

func NewAuthService(db *gorm.DB, secret string) AuthService {
	return &authService{
		db:        db,
		jwtSecret: secret,
	}
}

func (s *authService) Login(ctx context.Context, phone string) (string, error) {

	q := query.Use(s.db)

	// 1️⃣ Cari user
	user, err := q.User.WithContext(ctx).
		Where(q.User.PhoneNumber.Eq(phone)).
		First()

	if err != nil {
		return "", errors.New("user not found")
	}

	// 2️⃣ Ambil roles user
	userRoles, err := q.UserRole.WithContext(ctx).
		Where(q.UserRole.UserID.Eq(user.UserID)).
		Find()

	if err != nil {
		return "", err
	}

	// 3️⃣ Ambil nama role
	var roles []string

	for _, ur := range userRoles {

		role, err := q.Role.WithContext(ctx).
			Where(q.Role.RoleID.Eq(ur.RoleID)).
			First()

		if err == nil {
			roles = append(roles, role.RoleName)
		}
	}

	// 4️⃣ Generate JWT
	claims := jwt.MapClaims{
		"user_id": user.UserID,
		"roles":   roles,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
