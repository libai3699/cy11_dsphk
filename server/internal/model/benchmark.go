package model

import "time"

const (
	BenchmarkAccountStatusPending  = "待分析"
	BenchmarkAccountStatusAnalyzed = "已分析"
	BenchmarkAccountStatusDisabled = "停用"

	BenchmarkAnalysisStatusPending   = "待执行"
	BenchmarkAnalysisStatusRunning   = "执行中"
	BenchmarkAnalysisStatusCompleted = "已完成"
	BenchmarkAnalysisStatusFailed    = "失败"
)

type BenchmarkAccount struct {
	ID             uint64    `gorm:"primaryKey" json:"id"`
	MerchantID     uint64    `gorm:"index;not null" json:"merchantId"`
	MerchantName   string    `gorm:"size:128;not null;default:''" json:"merchantName"`
	AccountName    string    `gorm:"size:128;not null;index" json:"accountName"`
	Platform       string    `gorm:"size:64;not null;default:'抖音'" json:"platform"`
	City           string    `gorm:"size:64;not null;default:''" json:"city"`
	Industry       string    `gorm:"size:64;not null;default:''" json:"industry"`
	AccountURL     string    `gorm:"size:500;not null;default:''" json:"accountUrl"`
	FollowerCount  float64   `gorm:"not null;default:0" json:"followerCount"`
	BestPlayCount  float64   `gorm:"not null;default:0" json:"bestPlayCount"`
	LatestHitTitle string    `gorm:"size:255;not null;default:''" json:"latestHitTitle"`
	Takeaway       string    `gorm:"size:1000;not null;default:''" json:"takeaway"`
	Risk           string    `gorm:"size:1000;not null;default:''" json:"risk"`
	Status         string    `gorm:"size:32;not null;default:'待分析'" json:"status"`
	Remark         string    `gorm:"size:1000;not null;default:''" json:"remark"`
	CreatedBy      uint64    `gorm:"index;not null;default:0" json:"createdBy"`
	UpdatedBy      uint64    `gorm:"index;not null;default:0" json:"updatedBy"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type BenchmarkAnalysisTask struct {
	ID                 uint64    `gorm:"primaryKey" json:"id"`
	MerchantID         uint64    `gorm:"index;not null" json:"merchantId"`
	MerchantName       string    `gorm:"size:128;not null;default:''" json:"merchantName"`
	BenchmarkAccountID uint64    `gorm:"index;not null" json:"benchmarkAccountId"`
	BenchmarkName      string    `gorm:"size:128;not null;default:''" json:"benchmarkName"`
	Status             string    `gorm:"size:32;not null;default:'待执行'" json:"status"`
	InputSnapshot      string    `gorm:"type:text" json:"inputSnapshot"`
	ResultJSON         string    `gorm:"type:text" json:"resultJson"`
	ErrorMessage       string    `gorm:"size:1000;not null;default:''" json:"errorMessage"`
	CreatedBy          uint64    `gorm:"index;not null;default:0" json:"createdBy"`
	UpdatedBy          uint64    `gorm:"index;not null;default:0" json:"updatedBy"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
}
