package sign

import "time"

type SignCumulative struct {
	ID        uint64    `json:"id" gorm:"column:id"`                 // 自增id
	Aid       int64     `json:"aid" gorm:"column:aid"`               // 角色id
	Num       int8      `json:"num" gorm:"column:num"`               // 领取次数
	SignMonth int       `json:"sign_month" gorm:"column:sign_month"` // 签到月份,202101
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"` // 创建时间
}

func (m *SignCumulative) TableName() string {
	return "sign_cumulative"
}
