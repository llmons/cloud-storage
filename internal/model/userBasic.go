package model

import (
	"github.com/go-dev-frame/sponge/pkg/sgorm"
)

type UserBasic struct {
	sgorm.Model `gorm:"embedded"` // embed id and time

	Identity string `gorm:"column:identity;type:varchar(36)" json:"identity"`
	Name     string `gorm:"column:name;type:varchar(60)" json:"name"`
	Password string `gorm:"column:password;type:varchar(32)" json:"password"`
	Email    string `gorm:"column:email;type:varchar(100)" json:"email"`
}

// TableName table name
func (m *UserBasic) TableName() string {
	return "user_basic"
}

// UserBasicColumnNames Whitelist for custom query fields to prevent sql injection attacks
var UserBasicColumnNames = map[string]bool{
	"id":         true,
	"created_at": true,
	"updated_at": true,
	"deleted_at": true,
	"identity":   true,
	"name":       true,
	"password":   true,
	"email":      true,
}
