package model

import "time"

const (
	PainPointStatusPending   = "待拆"
	PainPointStatusAdopted   = "已采纳"
	PainPointStatusConverted = "已转化"

	CaseStudyStatusActive   = "有效"
	CaseStudyStatusInactive = "失效"

	AccountProfileStatusActive   = "启用"
	AccountProfileStatusInactive = "停用"

	PlatformRuleStatusActive   = "生效"
	PlatformRuleStatusInactive = "失效"

	ContentTemplateStatusActive   = "启用"
	ContentTemplateStatusInactive = "停用"
)

// PainPoint 用户痛点（倒推选题法中枢资产）：记录用户愿意为什么掏钱的真实痛点。
type PainPoint struct {
	ID            uint64    `gorm:"primaryKey" json:"id"`
	MerchantID    uint64    `gorm:"index;not null" json:"merchantId"`
	MerchantName  string    `gorm:"size:128;not null;default:''" json:"merchantName"`
	Source        string    `gorm:"size:32;not null;default:''" json:"source"`      // 私信/评论区/群聊/热点/对标
	Category      string    `gorm:"size:64;not null;default:''" json:"category"`    // 内容选题/拍摄剪辑/投放转化/账号定位
	Content       string    `gorm:"type:text" json:"content"`                       // 痛点描述
	UserQuote     string    `gorm:"type:text" json:"userQuote"`                     // 用户原话
	Emotion       string    `gorm:"size:32;not null;default:''" json:"emotion"`     // 焦虑/困惑/期待
	Product       string    `gorm:"size:128;not null;default:''" json:"product"`    // 对应产品
	DemandLevel   string    `gorm:"size:16;not null;default:'中'" json:"demandLevel"` // 高/中/低
	SuggestedTopic string   `gorm:"type:text" json:"suggestedTopic"`                // 建议选题方向
	Status        string    `gorm:"size:16;not null;default:'待拆'" json:"status"`
	CreatedBy     uint64    `gorm:"index;not null;default:0" json:"createdBy"`
	UpdatedBy     uint64    `gorm:"index;not null;default:0" json:"updatedBy"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func (PainPoint) TableName() string { return "pain_points" }

// CaseStudy 爆款案例库：拆解对标/自身爆款，提取可复用结构。
type CaseStudy struct {
	ID              uint64    `gorm:"primaryKey" json:"id"`
	MerchantID      uint64    `gorm:"index;not null" json:"merchantId"`
	MerchantName    string    `gorm:"size:128;not null;default:''" json:"merchantName"`
	Title           string    `gorm:"size:255;not null;index" json:"title"`
	Platform        string    `gorm:"size:32;not null;default:''" json:"platform"`     // 抖音/小红书/视频号
	AccountName     string    `gorm:"size:128;not null;default:''" json:"accountName"` // 来源账号
	Industry        string    `gorm:"size:64;not null;default:''" json:"industry"`     // 行业/赛道
	Form            string    `gorm:"size:32;not null;default:''" json:"form"`         // 口播/剧情/测评/干货
	HookType        string    `gorm:"size:32;not null;default:''" json:"hookType"`     // 钩子类型
	Structure       string    `gorm:"type:text" json:"structure"`                      // 结构(痛点-价值-行动)
	ConversionAction string   `gorm:"type:text" json:"conversionAction"`               // 转化动作
	ReusablePoint   string    `gorm:"type:text" json:"reusablePoint"`                  // 可复用点
	DataJSON        string    `gorm:"type:text;not null;default:'{}'" json:"dataJson"` // 数据快照
	Lesson          string    `gorm:"type:text" json:"lesson"`                         // 经验沉淀
	Status          string    `gorm:"size:16;not null;default:'有效'" json:"status"`
	CreatedBy       uint64    `gorm:"index;not null;default:0" json:"createdBy"`
	UpdatedBy       uint64    `gorm:"index;not null;default:0" json:"updatedBy"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

func (CaseStudy) TableName() string { return "case_studies" }

// AccountProfile 账号定位画像：选题依赖项，描述赚谁的钱/卖什么/人设与语气。
type AccountProfile struct {
	ID              uint64    `gorm:"primaryKey" json:"id"`
	MerchantID      uint64    `gorm:"index;not null" json:"merchantId"`
	MerchantName    string    `gorm:"size:128;not null;default:''" json:"merchantName"`
	AccountName     string    `gorm:"size:128;not null;default:''" json:"accountName"`
	Platform        string    `gorm:"size:32;not null;default:''" json:"platform"`
	Industry        string    `gorm:"size:64;not null;default:''" json:"industry"`
	TargetAudience  string    `gorm:"type:text" json:"targetAudience"` // 赚谁的钱
	Product         string    `gorm:"type:text" json:"product"`        // 卖什么
	Positioning     string    `gorm:"size:255;not null;default:''" json:"positioning"` // 一句话定位
	Persona         string    `gorm:"size:128;not null;default:''" json:"persona"`      // 人设
	ToneStyle       string    `gorm:"size:128;not null;default:''" json:"toneStyle"`    // 语气风格
	Monetization    string    `gorm:"size:128;not null;default:''" json:"monetization"` // 变现方式
	UpdateFrequency string    `gorm:"size:64;not null;default:''" json:"updateFrequency"`
	Strengths       string    `gorm:"type:text" json:"strengths"`
	Weaknesses      string    `gorm:"type:text" json:"weaknesses"`
	Status          string    `gorm:"size:16;not null;default:'启用'" json:"status"`
	CreatedBy       uint64    `gorm:"index;not null;default:0" json:"createdBy"`
	UpdatedBy       uint64    `gorm:"index;not null;default:0" json:"updatedBy"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

func (AccountProfile) TableName() string { return "account_profiles" }

// PlatformRule 抖音平台规则与算法常识：用于选题/脚本/互动的合规与效果判断。
type PlatformRule struct {
	ID           uint64    `gorm:"primaryKey" json:"id"`
	MerchantID   uint64    `gorm:"index;not null" json:"merchantId"`
	MerchantName string    `gorm:"size:128;not null;default:''" json:"merchantName"`
	Platform     string    `gorm:"size:32;not null;default:'抖音'" json:"platform"`
	Category     string    `gorm:"size:64;not null;default:''" json:"category"`  // 算法分发/合规红线/转化规则
	Title        string    `gorm:"size:255;not null;index" json:"title"`
	Content      string    `gorm:"type:text" json:"content"`
	RiskLevel    string    `gorm:"size:16;not null;default:'中'" json:"riskLevel"` // 高/中/低
	Source       string    `gorm:"size:128;not null;default:''" json:"source"`
	EffectiveDate string   `gorm:"size:32;not null;default:''" json:"effectiveDate"`
	Status       string    `gorm:"size:16;not null;default:'生效'" json:"status"`
	CreatedBy    uint64    `gorm:"index;not null;default:0" json:"createdBy"`
	UpdatedBy    uint64    `gorm:"index;not null;default:0" json:"updatedBy"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func (PlatformRule) TableName() string { return "platform_rules" }

// ContentTemplate 内容模板：选题卡/口播脚本/分镜脚本/复盘报告/对标报告等结构化模板。
type ContentTemplate struct {
	ID          uint64    `gorm:"primaryKey" json:"id"`
	MerchantID  uint64    `gorm:"index;not null" json:"merchantId"`
	MerchantName string   `gorm:"size:128;not null;default:''" json:"merchantName"`
	Name        string    `gorm:"size:255;not null;index" json:"name"`
	Type        string    `gorm:"size:32;not null;default:''" json:"type"`     // 选题卡/口播脚本/分镜脚本/复盘报告/对标报告
	Category    string    `gorm:"size:64;not null;default:''" json:"category"` // 栏目/场景
	StructureJSON string  `gorm:"type:text;not null;default:'{}'" json:"structureJson"`
	Content     string    `gorm:"type:text" json:"content"`
	UsageNote   string    `gorm:"type:text" json:"usageNote"`
	Status      string    `gorm:"size:16;not null;default:'启用'" json:"status"`
	CreatedBy   uint64    `gorm:"index;not null;default:0" json:"createdBy"`
	UpdatedBy   uint64    `gorm:"index;not null;default:0" json:"updatedBy"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (ContentTemplate) TableName() string { return "content_templates" }
