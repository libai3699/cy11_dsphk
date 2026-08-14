package model

import "time"

const (
	ContentTopicStatusPending  = "待确认"
	ContentTopicStatusAccepted = "已采用"
	ContentTopicStatusDisabled = "停用"

	HotspotTopicTaskStatusRunning   = "执行中"
	HotspotTopicTaskStatusCompleted = "已完成"
	HotspotTopicTaskStatusFailed    = "失败"
)

type ContentTopic struct {
	ID              uint64    `gorm:"primaryKey" json:"id"`
	MerchantID      uint64    `gorm:"index;not null" json:"merchantId"`
	MerchantName    string    `gorm:"size:128;not null;default:''" json:"merchantName"`
	BenchmarkID     uint64    `gorm:"index;not null;default:0" json:"benchmarkId"`
	BenchmarkName   string    `gorm:"size:500;not null;default:''" json:"benchmarkName"`
	Title           string    `gorm:"size:255;not null;index" json:"title"`
	Hook            string    `gorm:"size:500;not null;default:''" json:"hook"`
	Angle           string    `gorm:"size:255;not null;default:''" json:"angle"`
	Target          string    `gorm:"size:128;not null;default:''" json:"target"`
	RiskLevel       string    `gorm:"size:32;not null;default:'low'" json:"riskLevel"`
	RecommendReason string    `gorm:"size:1000;not null;default:''" json:"recommendReason"`
	TagsJSON        string    `gorm:"size:1000;not null;default:'[]'" json:"tagsJson"`
	PublishWindow   string    `gorm:"size:64;not null;default:''" json:"publishWindow"`
	Status          string    `gorm:"size:32;not null;default:'待确认'" json:"status"`
	SourceTaskID    uint64    `gorm:"index;not null;default:0" json:"sourceTaskId"`
	CreatedBy       uint64    `gorm:"index;not null;default:0" json:"createdBy"`
	UpdatedBy       uint64    `gorm:"index;not null;default:0" json:"updatedBy"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type HotspotTopicTask struct {
	ID            uint64    `gorm:"primaryKey" json:"id"`
	MerchantID    uint64    `gorm:"index;not null" json:"merchantId"`
	MerchantName  string    `gorm:"size:128;not null;default:''" json:"merchantName"`
	BenchmarkID   uint64    `gorm:"index;not null;default:0" json:"benchmarkId"`
	BenchmarkName string    `gorm:"size:500;not null;default:''" json:"benchmarkName"`
	Status        string    `gorm:"size:32;not null;default:'执行中'" json:"status"`
	InputSnapshot string    `gorm:"type:text" json:"inputSnapshot"`
	ResultJSON    string    `gorm:"type:text" json:"resultJson"`
	ErrorMessage  string    `gorm:"size:1000;not null;default:''" json:"errorMessage"`
	CreatedBy     uint64    `gorm:"index;not null;default:0" json:"createdBy"`
	UpdatedBy     uint64    `gorm:"index;not null;default:0" json:"updatedBy"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
