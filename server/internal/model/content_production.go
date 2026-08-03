package model

import "time"

const (
	ContentScriptStatusDraft     = "草稿"
	ContentScriptStatusConfirmed = "已确认"
	ContentScriptStatusDisabled  = "停用"

	ContentStoryboardStatusDraft     = "草稿"
	ContentStoryboardStatusConfirmed = "已确认"

	ShootingTaskStatusPending  = "待拍摄"
	ShootingTaskStatusShooting = "拍摄中"
	ShootingTaskStatusShot     = "已拍摄"
	ShootingTaskStatusEdited   = "已剪辑"
	ShootingTaskStatusDone     = "已完成"

	PublishScheduleStatusPending   = "待发布"
	PublishScheduleStatusPublished = "已发布"
	PublishScheduleStatusReviewed  = "已复盘"

	ContentReviewStatusCompleted = "已完成"
	ContentReviewStatusFailed    = "失败"
)

type ContentScript struct {
	ID                uint64    `gorm:"primaryKey" json:"id"`
	MerchantID        uint64    `gorm:"index;not null" json:"merchantId"`
	MerchantName      string    `gorm:"size:128;not null;default:''" json:"merchantName"`
	TopicID           uint64    `gorm:"index;not null;default:0" json:"topicId"`
	TopicTitle        string    `gorm:"size:255;not null;default:''" json:"topicTitle"`
	Title             string    `gorm:"size:255;not null;index" json:"title"`
	Opening           string    `gorm:"size:1000;not null;default:''" json:"opening"`
	Body              string    `gorm:"type:text" json:"body"`
	Ending            string    `gorm:"size:1000;not null;default:''" json:"ending"`
	CTA               string    `gorm:"size:1000;not null;default:''" json:"cta"`
	FullScript        string    `gorm:"type:text" json:"fullScript"`
	ShootingNotesJSON string    `gorm:"type:text" json:"shootingNotesJson"`
	Status            string    `gorm:"size:32;not null;default:'草稿'" json:"status"`
	InputSnapshot     string    `gorm:"type:text" json:"inputSnapshot"`
	ResultJSON        string    `gorm:"type:text" json:"resultJson"`
	ErrorMessage      string    `gorm:"size:1000;not null;default:''" json:"errorMessage"`
	CreatedBy         uint64    `gorm:"index;not null;default:0" json:"createdBy"`
	UpdatedBy         uint64    `gorm:"index;not null;default:0" json:"updatedBy"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

type ContentStoryboard struct {
	ID            uint64    `gorm:"primaryKey" json:"id"`
	MerchantID    uint64    `gorm:"index;not null" json:"merchantId"`
	MerchantName  string    `gorm:"size:128;not null;default:''" json:"merchantName"`
	TopicID       uint64    `gorm:"index;not null;default:0" json:"topicId"`
	TopicTitle    string    `gorm:"size:255;not null;default:''" json:"topicTitle"`
	ScriptID      uint64    `gorm:"index;not null;default:0" json:"scriptId"`
	ScriptTitle   string    `gorm:"size:255;not null;default:''" json:"scriptTitle"`
	ShotsJSON     string    `gorm:"type:text" json:"shotsJson"`
	Status        string    `gorm:"size:32;not null;default:'草稿'" json:"status"`
	InputSnapshot string    `gorm:"type:text" json:"inputSnapshot"`
	ResultJSON    string    `gorm:"type:text" json:"resultJson"`
	ErrorMessage  string    `gorm:"size:1000;not null;default:''" json:"errorMessage"`
	CreatedBy     uint64    `gorm:"index;not null;default:0" json:"createdBy"`
	UpdatedBy     uint64    `gorm:"index;not null;default:0" json:"updatedBy"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type PublishSchedule struct {
	ID             uint64     `gorm:"primaryKey" json:"id"`
	MerchantID     uint64     `gorm:"index;not null" json:"merchantId"`
	MerchantName   string     `gorm:"size:128;not null;default:''" json:"merchantName"`
	TopicID        uint64     `gorm:"index;not null;default:0" json:"topicId"`
	TopicTitle     string     `gorm:"size:255;not null;default:''" json:"topicTitle"`
	ScriptID       uint64     `gorm:"index;not null;default:0" json:"scriptId"`
	ScriptTitle    string     `gorm:"size:255;not null;default:''" json:"scriptTitle"`
	StoryboardID   uint64     `gorm:"index;not null;default:0" json:"storyboardId"`
	VideoTitle     string     `gorm:"size:255;not null;index" json:"videoTitle"`
	PublishTime    *time.Time `json:"publishTime"`
	Owner          string     `gorm:"size:64;not null;default:''" json:"owner"`
	DouyinAccount  string     `gorm:"size:128;not null;default:''" json:"douyinAccount"`
	MaterialStatus string     `gorm:"size:64;not null;default:'待拍摄'" json:"materialStatus"`
	Status         string     `gorm:"size:32;not null;default:'待发布'" json:"status"`
	Remark         string     `gorm:"size:1000;not null;default:''" json:"remark"`
	CreatedBy      uint64     `gorm:"index;not null;default:0" json:"createdBy"`
	UpdatedBy      uint64     `gorm:"index;not null;default:0" json:"updatedBy"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
}

type ShootingTask struct {
	ID           uint64     `gorm:"primaryKey" json:"id"`
	MerchantID   uint64     `gorm:"index;not null" json:"merchantId"`
	MerchantName string     `gorm:"size:128;not null;default:''" json:"merchantName"`
	TopicID      uint64     `gorm:"index;not null;default:0" json:"topicId"`
	TopicTitle   string     `gorm:"size:255;not null;default:''" json:"topicTitle"`
	ScriptID     uint64     `gorm:"index;not null;default:0" json:"scriptId"`
	ScriptTitle  string     `gorm:"size:255;not null;default:''" json:"scriptTitle"`
	StoryboardID uint64     `gorm:"index;not null;default:0" json:"storyboardId"`
	TaskTitle    string     `gorm:"size:255;not null;index" json:"taskTitle"`
	ShotCount    int        `gorm:"not null;default:0" json:"shotCount"`
	ShotsJSON    string     `gorm:"type:text" json:"shotsJson"`
	Assignee     string     `gorm:"size:64;not null;default:''" json:"assignee"`
	ShootTime    *time.Time `json:"shootTime"`
	Deadline     *time.Time `json:"deadline"`
	Status       string     `gorm:"size:32;not null;default:'待拍摄'" json:"status"`
	MaterialURL  string     `gorm:"size:500;not null;default:''" json:"materialUrl"`
	Remark       string     `gorm:"size:1000;not null;default:''" json:"remark"`
	CreatedBy    uint64     `gorm:"index;not null;default:0" json:"createdBy"`
	UpdatedBy    uint64     `gorm:"index;not null;default:0" json:"updatedBy"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}

type ContentReviewTask struct {
	ID             uint64     `gorm:"primaryKey" json:"id"`
	MerchantID     uint64     `gorm:"index;not null" json:"merchantId"`
	MerchantName   string     `gorm:"size:128;not null;default:''" json:"merchantName"`
	ScheduleID     uint64     `gorm:"index;not null;default:0" json:"scheduleId"`
	VideoTitle     string     `gorm:"size:255;not null;default:''" json:"videoTitle"`
	PeriodStart    *time.Time `json:"periodStart"`
	PeriodEnd      *time.Time `json:"periodEnd"`
	Status         string     `gorm:"size:32;not null;default:'已完成'" json:"status"`
	InputSnapshot  string     `gorm:"type:text" json:"inputSnapshot"`
	ResultJSON     string     `gorm:"type:text" json:"resultJson"`
	ErrorMessage   string     `gorm:"size:1000;not null;default:''" json:"errorMessage"`
	PlayCount      int64      `gorm:"not null;default:0" json:"playCount"`
	LikeCount      int64      `gorm:"not null;default:0" json:"likeCount"`
	CommentCount   int64      `gorm:"not null;default:0" json:"commentCount"`
	ShareCount     int64      `gorm:"not null;default:0" json:"shareCount"`
	DealCount      int64      `gorm:"not null;default:0" json:"dealCount"`
	WriteOffAmount float64    `gorm:"not null;default:0" json:"writeOffAmount"`
	CreatedBy      uint64     `gorm:"index;not null;default:0" json:"createdBy"`
	UpdatedBy      uint64     `gorm:"index;not null;default:0" json:"updatedBy"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
}
