package mysql

import (
	"go-clean-api/internal/domain/entity"
	"go-clean-api/pkg/apperror"

	"gorm.io/gorm"
)

type userModel struct {
	gorm.Model
	Name     string `gorm:"size:255;not null"`
	Email    string `gorm:"size:255;uniqueIndex;not null"`
	Password string `gorm:"size:255;not null"`
}

func (userModel) TableName() string { return "users" }

func toUserEntity(m *userModel) *entity.User {
	return &entity.User{
		ID:        m.ID,
		Name:      m.Name,
		Email:     m.Email,
		Password:  m.Password,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByID(id uint) (*entity.User, error) {
	var m userModel
	if err := r.db.First(&m, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}
	return toUserEntity(&m), nil
}

func (r *userRepository) FindByEmail(email string) (*entity.User, error) {
	var m userModel
	if err := r.db.Where("email = ?", email).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}
	return toUserEntity(&m), nil
}

func (r *userRepository) FindAll() ([]*entity.User, error) {
	var models []userModel
	if err := r.db.Find(&models).Error; err != nil {
		return nil, err
	}
	users := make([]*entity.User, len(models))
	for i, m := range models {
		m := m
		users[i] = toUserEntity(&m)
	}
	return users, nil
}

func (r *userRepository) Create(user *entity.User) error {
	m := &userModel{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}
	if err := r.db.Create(m).Error; err != nil {
		return err
	}
	user.ID = m.ID
	user.CreatedAt = m.CreatedAt
	user.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *userRepository) Update(user *entity.User) error {
	return r.db.Model(&userModel{}).Where("id = ?", user.ID).
		Updates(map[string]interface{}{
			"name":  user.Name,
			"email": user.Email,
		}).Error
}
