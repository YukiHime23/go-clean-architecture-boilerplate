package usecase

import (
	"errors"
	"fmt"
	"time"

	"go-clean-architecture-boilerplate/config"
	"go-clean-architecture-boilerplate/internal/domain"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type authUsecase struct {
	userRepo  domain.UserRepository
	jwtConfig config.JWTConfig
	log       *zap.Logger
}

// NewAuthUsecase creates a new AuthUsecase with injected JWT config.
func NewAuthUsecase(userRepo domain.UserRepository, jwtCfg config.JWTConfig, log *zap.Logger) domain.AuthUsecase {
	return &authUsecase{userRepo: userRepo, jwtConfig: jwtCfg, log: log}
}

// Register creates a new user with a hashed password.
func (u *authUsecase) Register(input *domain.RegisterInput) (*domain.User, error) {
	_, err := u.userRepo.FindByEmail(input.Email)
	if err == nil {
		return nil, domain.NewConflictError("email already registered")
	}
	if !errors.Is(err, domain.ErrNotFound) {
		u.log.Error("Register: unexpected DB error", zap.Error(err))
		return nil, domain.NewInternalError(err)
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		u.log.Error("Register: bcrypt error", zap.Error(err))
		return nil, domain.NewInternalError(err)
	}

	user := &domain.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashed),
	}
	if err := u.userRepo.Create(user); err != nil {
		u.log.Error("Register: create user error", zap.Error(err))
		return nil, domain.NewInternalError(err)
	}
	return user, nil
}

// Login validates credentials and returns a signed JWT.
func (u *authUsecase) Login(input *domain.LoginInput) (string, error) {
	user, err := u.userRepo.FindByEmail(input.Email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return "", domain.NewUnauthorizedError("invalid credentials")
		}
		u.log.Error("Login: DB error", zap.Error(err))
		return "", domain.NewInternalError(err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return "", domain.NewUnauthorizedError("invalid credentials")
	}

	token, err := u.generateJWT(user)
	if err != nil {
		u.log.Error("Login: JWT signing error", zap.Error(err))
		return "", domain.NewInternalError(err)
	}
	return token, nil
}

// JWTClaims represents the claims embedded in the JWT payload.
type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func (u *authUsecase) generateJWT(user *domain.User) (string, error) {
	if u.jwtConfig.Secret == "" {
		return "", fmt.Errorf("JWT secret is not configured")
	}

	claims := JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(u.jwtConfig.ExpirationHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(u.jwtConfig.Secret))
}
