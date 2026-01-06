package models

import "time"

type Users struct {
	ID                string    `gorm:"primaryKey;type:text;default:generate_object_id()" json:"id"`
	UserName          string    `gorm:"type:varchar(100);column:userName"`
	Email             string    `gorm:"type:varchar(150);unique;not null;column:email"`
	NoHandphone       string    `gorm:"type:varchar(20);column:noHandphone"`
	Password          string    `gorm:"type:text;not null;column:password"`
	ActiveUser        bool      `gorm:"default:false;column:activeUser"`
	Role              string    `gorm:"type:varchar(20);default:user;column:role"`
	VerificationToken string    `gorm:"type:text;column:verificationToken"`
	ApiKeyAI          string    `gorm:"type:text;column:apiKeyAI"`
	ProfilePicture    string    `gorm:"type:text;column:profilePicture"`
	UserBilling       int64     `gorm:"default:0;column:userBilling"`
	CreatedAt         time.Time `gorm:"column:createdAt"`
}

func (Users) TableName() string {
	return "users"
}