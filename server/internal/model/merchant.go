package model

import "time"

const (
	MerchantStatusDisabled = 0
	MerchantStatusEnabled  = 1
)

type Merchant struct {
	ID                 uint64    `gorm:"primaryKey" json:"id"`
	Name               string    `gorm:"size:128;not null;index" json:"name"`
	Industry           string    `gorm:"size:64;not null;default:''" json:"industry"`
	City               string    `gorm:"size:64;not null;default:''" json:"city"`
	Address            string    `gorm:"size:255;not null;default:''" json:"address"`
	ContactName        string    `gorm:"size:64;not null;default:''" json:"contactName"`
	ContactPhone       string    `gorm:"size:32;not null;default:''" json:"contactPhone"`
	DouyinAccount      string    `gorm:"size:128;not null;default:''" json:"douyinAccount"`
	DouyinLaikeAccount string    `gorm:"size:128;not null;default:''" json:"douyinLaikeAccount"`
	CooperationType    string    `gorm:"size:64;not null;default:'成交提点'" json:"cooperationType"`
	CommissionRate     float64   `gorm:"not null;default:10" json:"commissionRate"`
	Stage              string    `gorm:"size:64;not null;default:'已建档'" json:"stage"`
	Status             int       `gorm:"not null;default:1" json:"status"`
	Remark             string    `gorm:"size:1000;not null;default:''" json:"remark"`
	CreatedBy          uint64    `gorm:"index;not null;default:0" json:"createdBy"`
	UpdatedBy          uint64    `gorm:"index;not null;default:0" json:"updatedBy"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
}
