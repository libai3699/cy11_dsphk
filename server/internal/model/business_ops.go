package model

import "time"

const (
	MerchantFollowUpStageContacting  = "沟通中"
	MerchantFollowUpStageNegotiating = "方案确认"
	MerchantFollowUpStageContracted  = "已签约"
	MerchantFollowUpStagePaused      = "暂缓"
	MerchantFollowUpStageLost        = "已流失"

	SettlementOrderStatusPending   = "待核对"
	SettlementOrderStatusConfirmed = "已确认"
	SettlementOrderStatusPaid      = "已结算"
)

type MerchantFollowUpLog struct {
	ID             uint64     `gorm:"primaryKey" json:"id"`
	MerchantID     uint64     `gorm:"index;not null" json:"merchantId"`
	MerchantName   string     `gorm:"size:128;not null;default:''" json:"merchantName"`
	Stage          string     `gorm:"size:32;not null;default:'沟通中'" json:"stage"`
	LatestTalk     string     `gorm:"type:text" json:"latestTalk"`
	Objection      string     `gorm:"size:1000;not null;default:''" json:"objection"`
	NextStep       string     `gorm:"size:1000;not null;default:''" json:"nextStep"`
	Owner          string     `gorm:"size:64;not null;default:''" json:"owner"`
	FollowTime     *time.Time `json:"followTime"`
	NextFollowTime *time.Time `json:"nextFollowTime"`
	CreatedBy      uint64     `gorm:"index;not null;default:0" json:"createdBy"`
	UpdatedBy      uint64     `gorm:"index;not null;default:0" json:"updatedBy"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
}

type SettlementOrder struct {
	ID             uint64     `gorm:"primaryKey" json:"id"`
	MerchantID     uint64     `gorm:"index;not null" json:"merchantId"`
	MerchantName   string     `gorm:"size:128;not null;default:''" json:"merchantName"`
	ScheduleID     uint64     `gorm:"index;not null;default:0" json:"scheduleId"`
	VideoTitle     string     `gorm:"size:255;not null;default:''" json:"videoTitle"`
	SourceVideo    string     `gorm:"size:500;not null;default:''" json:"sourceVideo"`
	OrderWindow    string     `gorm:"size:128;not null;default:''" json:"orderWindow"`
	PeriodStart    *time.Time `json:"periodStart"`
	PeriodEnd      *time.Time `json:"periodEnd"`
	RedeemedAmount float64    `gorm:"not null;default:0" json:"redeemedAmount"`
	CommissionRate float64    `gorm:"not null;default:10" json:"commissionRate"`
	Commission     float64    `gorm:"not null;default:0" json:"commission"`
	Status         string     `gorm:"size:32;not null;default:'待核对'" json:"status"`
	Remark         string     `gorm:"size:1000;not null;default:''" json:"remark"`
	CreatedBy      uint64     `gorm:"index;not null;default:0" json:"createdBy"`
	UpdatedBy      uint64     `gorm:"index;not null;default:0" json:"updatedBy"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
}
