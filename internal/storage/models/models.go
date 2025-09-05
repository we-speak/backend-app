package models

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Username     string    `json:"username" validate:"required" gorm:"unique;not null"`
	Password     string    `json:"password" validate:"required" gorm:"not null"`
	Email        string    `json:"email" validate:"required,email" gorm:"unique;not null"`
	Role         string    `json:"role" validate:"required,oneof=user creator combined admin" gorm:"default:'user'"`
	Country      string    `json:"country" gorm:"not null"`
	RefreshToken string    `json:"-"`
	TokenExpiry  time.Time `json:"-"`
	CreatedAt    time.Time `json:"createdAt,omitempty" gorm:"autoCreateTime:true"`
	UpdatedAt    time.Time `json:"updatedAt,omitempty" gorm:"autoUpdateTime:true"`
}

func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
