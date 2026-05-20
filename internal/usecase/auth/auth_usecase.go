package auth

import (
	"go-clean-api/internal/domain/entity"
	"go-clean-api/pkg/apperror"
	pkgjwt "go-clean-api/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
)

type userCreator interface {
	Create(user *entity.User) error
}

type userByEmailFinder interface {
	FindByEmail(email string) (*entity.User, error)
}

type userAuthRepo interface {
	userCreator
	userByEmailFinder
}

type RegisterInput struct {
	Name     string `json:"name" binding:"required,min=2"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type Output struct {
	Token string       `json:"token"`
	User  *entity.User `json:"user"`
}

type UseCase struct {
	userRepo    userAuthRepo
	jwtSecret   string
	jwtExpireHr int
}

func New(userRepo userAuthRepo, jwtSecret string, jwtExpireHr int) *UseCase {
	return &UseCase{
		userRepo:    userRepo,
		jwtSecret:   jwtSecret,
		jwtExpireHr: jwtExpireHr,
	}
}

func (u *UseCase) Register(input RegisterInput) (*Output, error) {
	if _, err := u.userRepo.FindByEmail(input.Email); err == nil {
		return nil, apperror.ErrConflict
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperror.ErrInternal
	}

	user := &entity.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hash),
	}
	if err := u.userRepo.Create(user); err != nil {
		return nil, apperror.ErrInternal
	}

	token, err := pkgjwt.Generate(user.ID, u.jwtSecret, u.jwtExpireHr)
	if err != nil {
		return nil, apperror.ErrInternal
	}

	return &Output{Token: token, User: user}, nil
}

func (u *UseCase) Login(input LoginInput) (*Output, error) {
	user, err := u.userRepo.FindByEmail(input.Email)
	if err != nil {
		return nil, apperror.New(401, "invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, apperror.New(401, "invalid email or password")
	}

	token, err := pkgjwt.Generate(user.ID, u.jwtSecret, u.jwtExpireHr)
	if err != nil {
		return nil, apperror.ErrInternal
	}

	return &Output{Token: token, User: user}, nil
}
