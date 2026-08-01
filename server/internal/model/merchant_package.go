package model

import "time"

const (
	MerchantPackageStatusDisabled = 0
	MerchantPackageStatusEnabled  = 1
)

type MerchantPackage struct {
	ID             uint64    `gorm:"primaryKey" json:"id"`
	MerchantID     uint64    `gorm:"index;not null" json:"merchantId"`
	MerchantName   string    `gorm:"size:128;not null;default:''" json:"merchantName"`
	Name           string    `gorm:"size:128;not null;index" json:"name"`
	OriginalPrice  float64   `gorm:"not null;default:0" json:"originalPrice"`
	SellingPrice   float64   `gorm:"not null;default:0" json:"sellingPrice"`
	CostPrice      float64   `gorm:"not null;default:0" json:"costPrice"`
	CommissionRate float64   `gorm:"not null;default:10" json:"commissionRate"`
	TrafficLabel   string    `gorm:"size:128;not null;default:''" json:"trafficLabel"`
	ProfitGuard    string    `gorm:"size:255;not null;default:''" json:"profitGuard"`
	Status         int       `gorm:"not null;default:1" json:"status"`
	Remark         string    `gorm:"size:1000;not null;default:''" json:"remark"`
	CreatedBy      uint64    `gorm:"index;not null;default:0" json:"createdBy"`
	UpdatedBy      uint64    `gorm:"index;not null;default:0" json:"updatedBy"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
