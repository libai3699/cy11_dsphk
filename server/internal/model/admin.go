package model

import "time"

const (
	AdminUserStatusDisabled = 0
	AdminUserStatusEnabled  = 1
)

type AdminUser struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"size:64;not null;uniqueIndex" json:"username"`
	Password  string    `gorm:"size:255;not null" json:"-"`
	RealName  string    `gorm:"size:64;not null" json:"realName"`
	Avatar    string    `gorm:"size:255" json:"avatar"`
	Status    int       `gorm:"not null;default:1" json:"status"`
	HomePath  string    `gorm:"size:128;not null;default:/workspace" json:"homePath"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Role struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Code      string    `gorm:"size:64;not null;uniqueIndex" json:"code"`
	Name      string    `gorm:"size:64;not null" json:"name"`
	SortOrder int       `gorm:"not null;default:0" json:"sortOrder"`
	IsSystem  bool      `gorm:"not null;default:false" json:"isSystem"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Menu struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ParentID   *uint     `gorm:"index" json:"parentId"`
	Name       string    `gorm:"size:128;not null;uniqueIndex" json:"name"`
	Path       string    `gorm:"size:255;not null" json:"path"`
	Component  string    `gorm:"size:255" json:"component"`
	Redirect   string    `gorm:"size:255" json:"redirect"`
	Title      string    `gorm:"size:128;not null" json:"title"`
	Icon       string    `gorm:"size:128" json:"icon"`
	SortOrder  int       `gorm:"not null;default:0" json:"sortOrder"`
	AffixTab   bool      `gorm:"not null;default:false" json:"affixTab"`
	HideInMenu bool      `gorm:"not null;default:false" json:"hideInMenu"`
	Visible    bool      `gorm:"not null;default:true" json:"visible"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type AdminUserRole struct {
	UserID    uint64    `gorm:"primaryKey;autoIncrement:false" json:"userId"`
	RoleID    uint      `gorm:"primaryKey;autoIncrement:false" json:"roleId"`
	CreatedAt time.Time `json:"createdAt"`
}

type RoleMenu struct {
	RoleID    uint      `gorm:"primaryKey;autoIncrement:false" json:"roleId"`
	MenuID    uint      `gorm:"primaryKey;autoIncrement:false" json:"menuId"`
	CreatedAt time.Time `json:"createdAt"`
}

type AdminLoginLog struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	UserID    uint64    `gorm:"index" json:"userId"`
	Username  string    `gorm:"size:64;index" json:"username"`
	IP        string    `gorm:"size:64" json:"ip"`
	UserAgent string    `gorm:"size:512" json:"userAgent"`
	Result    string    `gorm:"size:32;not null" json:"result"`
	Message   string    `gorm:"size:255" json:"message"`
	CreatedAt time.Time `json:"createdAt"`
}
