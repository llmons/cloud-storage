package model

import (
	"github.com/go-dev-frame/sponge/pkg/sgorm"
)

type RepositoryPool struct {
	sgorm.Model `gorm:"embedded"` // embed id and time

	Identity string `gorm:"column:identity;type:varchar(36)" json:"identity"`
	Hash     string `gorm:"column:hash;type:varchar(32)" json:"hash"` // 文件的唯一标识
	Name     string `gorm:"column:name;type:varchar(255)" json:"name"`
	Ext      string `gorm:"column:ext;type:varchar(30)" json:"ext"`    // 文件扩展名
	Size     int    `gorm:"column:size;type:int(11)" json:"size"`      // 文件大小
	Path     string `gorm:"column:path;type:varchar(255)" json:"path"` // 文件路径
}

// TableName table name
func (m *RepositoryPool) TableName() string {
	return "repository_pool"
}

// RepositoryPoolColumnNames Whitelist for custom query fields to prevent sql injection attacks
var RepositoryPoolColumnNames = map[string]bool{
	"id":         true,
	"created_at": true,
	"updated_at": true,
	"deleted_at": true,
	"identity":   true,
	"hash":       true,
	"name":       true,
	"ext":        true,
	"size":       true,
	"path":       true,
}
