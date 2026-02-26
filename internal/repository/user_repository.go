package repository

import (
	"errors"
	"time"

	"go-clean-architecture-boilerplate/internal/domain"

	"gorm.io/gorm"
)

// userModel represents the database schema for the users table.
// We keep tags here because this is the infrastructure implementation.
type userModel struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"type:varchar(100);not null"`
	Email     string `gorm:"type:varchar(150);uniqueIndex;not null"`
	Password  string `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName overrides the table name used by GORM.
func (userModel) TableName() string {
	return "users"
}

// ToDomain converts a DB model to a Domain entity.
func (m *userModel) ToDomain() *domain.User {
	return &domain.User{
		ID:        m.ID,
		Name:      m.Name,
		Email:     m.Email,
		Password:  m.Password,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// FromDomain creates a DB model from a Domain entity.
func fromUserDomain(u *domain.User) *userModel {
	return &userModel{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	// Ensure migration use the local model
	_ = db.AutoMigrate(&userModel{})
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	m := fromUserDomain(user)
	if err := r.db.Create(m).Error; err != nil {
		return err
	}
	user.ID = m.ID
	return nil
}

func (r *userRepository) FindByID(id uint) (*domain.User, error) {
	var m userModel
	if err := r.db.First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	var m userModel
	if err := r.db.Where("email = ?", email).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *userRepository) FindAll(page, limit int) ([]domain.User, int64, error) {
	var models []userModel
	var total int64
	offset := (page - 1) * limit

	if err := r.db.Model(&userModel{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := r.db.Offset(offset).Limit(limit).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	users := make([]domain.User, len(models))
	for i, m := range models {
		users[i] = *m.ToDomain()
	}
	return users, total, nil
}

func (r *userRepository) Update(user *domain.User) error {
	m := fromUserDomain(user)
	return r.db.Save(m).Error
}

func (r *userRepository) Delete(id uint) error {
	result := r.db.Delete(&userModel{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}
