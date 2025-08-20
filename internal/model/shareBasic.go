package model

import (
	"github.com/go-dev-frame/sponge/pkg/sgorm"
)

type ShareBasic struct {
	sgorm.Model `gorm:"embedded"` // embed id and time

	Identity               string `gorm:"column:identity;type:varchar(36)" json:"identity"`
	UserIdentity           string `gorm:"column:user_identity;type:varchar(36)" json:"userIdentity"`
	RepositoryIdentity     string `gorm:"column:repository_identity;type:varchar(36)" json:"repositoryIdentity"`          // 公共池中的唯一标识
	UserRepositoryIdentity string `gorm:"column:user_repository_identity;type:varchar(36)" json:"userRepositoryIdentity"` // 用户池子中的唯一标识
	ExpiredTime            int    `gorm:"column:expired_time;type:int(11)" json:"expiredTime"`                            // 失效时间，单位秒, 【0-永不失效】
	ClickNum               int    `gorm:"column:click_num;type:int(11);default:0" json:"clickNum"`                        // 点击次数
}

// TableName table name
func (m *ShareBasic) TableName() string {
	return "share_basic"
}

// ShareBasicColumnNames Whitelist for custom query fields to prevent sql injection attacks
var ShareBasicColumnNames = map[string]bool{
	"id":                       true,
	"created_at":               true,
	"updated_at":               true,
	"deleted_at":               true,
	"identity":                 true,
	"user_identity":            true,
	"repository_identity":      true,
	"user_repository_identity": true,
	"expired_time":             true,
	"click_num":                true,
}
