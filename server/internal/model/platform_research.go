package model

import "time"

const (
	PlatformResearchTaskStatusPending   = "pending"
	PlatformResearchTaskStatusCompleted = "completed"
	PlatformResearchTaskStatusFailed    = "failed"
)

type PlatformResearchTask struct {
	ID                uint64    `gorm:"primaryKey" json:"id"`
	MerchantID        uint64    `gorm:"index;not null" json:"merchantId"`
	MerchantName      string    `gorm:"size:128;not null;default:''" json:"merchantName"`
	Industry          string    `gorm:"size:64;not null;default:''" json:"industry"`
	City              string    `gorm:"size:64;not null;default:''" json:"city"`
	SourcesJSON       string    `gorm:"type:text" json:"sourcesJson"`
	KeywordsJSON      string    `gorm:"type:text" json:"keywordsJson"`
	SearchResultsJSON string    `gorm:"type:text" json:"searchResultsJson"`
	GoodCasesJSON     string    `gorm:"type:text" json:"goodCasesJson"`
	BadCasesJSON      string    `gorm:"type:text" json:"badCasesJson"`
	InsightsJSON      string    `gorm:"type:text" json:"insightsJson"`
	Summary           string    `gorm:"type:text" json:"summary"`
	SuggestionsJSON   string    `gorm:"type:text" json:"suggestionsJson"`
	Status            string    `gorm:"size:32;not null;default:'pending';index" json:"status"`
	ErrorMessage      string    `gorm:"type:text" json:"errorMessage"`
	CreatedBy         uint64    `gorm:"index;not null;default:0" json:"createdBy"`
	UpdatedBy         uint64    `gorm:"index;not null;default:0" json:"updatedBy"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}
