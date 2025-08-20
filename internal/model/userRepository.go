package model

import (
	"github.com/go-dev-frame/sponge/pkg/sgorm"
)

type UserRepository struct {
	sgorm.Model `gorm:"embedded"` // embed id and time

	Identity           string `gorm:"column:identity;type:varchar(36)" json:"identity"`
	UserIdentity       string `gorm:"column:user_identity;type:varchar(36)" json:"userIdentity"`
	ParentID           int    `gorm:"column:parent_id;type:int(11)" json:"parentID"`
	RepositoryIdentity string `gorm:"column:repository_identity;type:varchar(36)" json:"repositoryIdentity"`
	Ext                string `gorm:"column:ext;type:varchar(255)" json:"ext"` // 文件或文件夹类型
	Name               string `gorm:"column:name;type:varchar(255)" json:"name"`
}

// TableName table name
func (m *UserRepository) TableName() string {
	return "user_repository"
}

// UserRepositoryColumnNames Whitelist for custom query fields to prevent sql injection attacks
var UserRepositoryColumnNames = map[string]bool{
	"id":                  true,
	"created_at":          true,
	"updated_at":          true,
	"deleted_at":          true,
	"identity":            true,
	"user_identity":       true,
	"parent_id":           true,
	"repository_identity": true,
	"ext":                 true,
	"name":                true,
}
