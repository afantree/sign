package sign

import (
	"embed"
	"time"
)

//go:embed sql/*.sql
var SignMigrateFS embed.FS

type SignDaily struct {
	ID        uint64    `json:"id" gorm:"column:id"`                 // 自增id
	Aid       int64     `json:"aid" gorm:"column:aid"`               // 角色id
	SignAt    int       `json:"sign_at" gorm:"column:sign_at"`       // 签到那天,20210101
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"` // 创建时间
}

func (m *SignDaily) TableName() string {
	return "sign_daily"
}
