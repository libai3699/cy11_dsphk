package model

import "time"

const (
	MerchantAccountAuthStatusPending = "待授权"
	MerchantAccountAuthStatusActive  = "已授权"
	MerchantAccountAuthStatusRenewal = "待续期"
	MerchantAccountAuthStatusExpired = "已失效"
)

type MerchantAccountAuth struct {
	ID           uint64    `gorm:"primaryKey" json:"id"`
	MerchantID   uint64    `gorm:"index;not null" json:"merchantId"`
	MerchantName string    `gorm:"size:128;not null;default:''" json:"merchantName"`
	Platform     string    `gorm:"size:64;not null;default:'抖音号'" json:"platform"`
	AuthMethod   string    `gorm:"size:64;not null;default:'验证码代登'" json:"authMethod"`
	AccountName  string    `gorm:"size:128;not null;default:''" json:"accountName"`
	AccountUID   string    `gorm:"size:128;not null;default:''" json:"accountUid"`
	AuthStatus   string    `gorm:"size:32;not null;default:'待授权'" json:"authStatus"`
	LastLoginAt  string    `gorm:"size:32;not null;default:''" json:"lastLoginAt"`
	Remark       string    `gorm:"size:1000;not null;default:''" json:"remark"`
	CreatedBy    uint64    `gorm:"index;not null;default:0" json:"createdBy"`
	UpdatedBy    uint64    `gorm:"index;not null;default:0" json:"updatedBy"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
