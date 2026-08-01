package model

import "time"

const (
	AccountDiagnosisStatusPending   = "待执行"
	AccountDiagnosisStatusRunning   = "执行中"
	AccountDiagnosisStatusCompleted = "已完成"
	AccountDiagnosisStatusFailed    = "失败"
)

type AccountDiagnosisTask struct {
	ID            uint64    `gorm:"primaryKey" json:"id"`
	MerchantID    uint64    `gorm:"index;not null" json:"merchantId"`
	MerchantName  string    `gorm:"size:128;not null;default:''" json:"merchantName"`
	AccountAuthID uint64    `gorm:"index;not null;default:0" json:"accountAuthId"`
	AccountName   string    `gorm:"size:128;not null;default:''" json:"accountName"`
	Status        string    `gorm:"size:32;not null;default:'待执行'" json:"status"`
	InputSnapshot string    `gorm:"type:text" json:"inputSnapshot"`
	ResultJSON    string    `gorm:"type:text" json:"resultJson"`
	ErrorMessage  string    `gorm:"size:1000;not null;default:''" json:"errorMessage"`
	CreatedBy     uint64    `gorm:"index;not null;default:0" json:"createdBy"`
	UpdatedBy     uint64    `gorm:"index;not null;default:0" json:"updatedBy"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
