package service

import (
	"cy11dsphk/server/internal/database"
	"cy11dsphk/server/internal/model"
)

// ---------- PainPoint ----------

func ListPainPoints(merchantID uint64, keyword string) ([]model.PainPoint, error) {
	var list []model.PainPoint
	q := database.DB.Model(&model.PainPoint{})
	if merchantID > 0 {
		q = q.Where("merchant_id = ?", merchantID)
	}
	if keyword != "" {
		q = q.Where("content LIKE ? OR user_quote LIKE ? OR category LIKE ? OR suggested_topic LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	err := q.Order("id DESC").Find(&list).Error
	return list, err
}

func GetPainPoint(id uint64) (*model.PainPoint, error) {
	var m model.PainPoint
	if err := database.DB.First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func CreatePainPoint(m *model.PainPoint) error {
	return database.DB.Create(m).Error
}

func UpdatePainPoint(id uint64, m *model.PainPoint) error {
	m.ID = id
	return database.DB.Save(m).Error
}

func DeletePainPoint(id uint64) error {
	return database.DB.Delete(&model.PainPoint{}, id).Error
}

// ---------- CaseStudy ----------

func ListCaseStudies(merchantID uint64, keyword string) ([]model.CaseStudy, error) {
	var list []model.CaseStudy
	q := database.DB.Model(&model.CaseStudy{})
	if merchantID > 0 {
		q = q.Where("merchant_id = ?", merchantID)
	}
	if keyword != "" {
		q = q.Where("title LIKE ? OR account_name LIKE ? OR reusable_point LIKE ? OR industry LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	err := q.Order("id DESC").Find(&list).Error
	return list, err
}

func GetCaseStudy(id uint64) (*model.CaseStudy, error) {
	var m model.CaseStudy
	if err := database.DB.First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func CreateCaseStudy(m *model.CaseStudy) error {
	return database.DB.Create(m).Error
}

func UpdateCaseStudy(id uint64, m *model.CaseStudy) error {
	m.ID = id
	return database.DB.Save(m).Error
}

func DeleteCaseStudy(id uint64) error {
	return database.DB.Delete(&model.CaseStudy{}, id).Error
}

// ---------- AccountProfile ----------

func ListAccountProfiles(merchantID uint64, keyword string) ([]model.AccountProfile, error) {
	var list []model.AccountProfile
	q := database.DB.Model(&model.AccountProfile{})
	if merchantID > 0 {
		q = q.Where("merchant_id = ?", merchantID)
	}
	if keyword != "" {
		q = q.Where("account_name LIKE ? OR positioning LIKE ? OR industry LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	err := q.Order("id DESC").Find(&list).Error
	return list, err
}

func GetAccountProfile(id uint64) (*model.AccountProfile, error) {
	var m model.AccountProfile
	if err := database.DB.First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func CreateAccountProfile(m *model.AccountProfile) error {
	return database.DB.Create(m).Error
}

func UpdateAccountProfile(id uint64, m *model.AccountProfile) error {
	m.ID = id
	return database.DB.Save(m).Error
}

func DeleteAccountProfile(id uint64) error {
	return database.DB.Delete(&model.AccountProfile{}, id).Error
}

// ---------- PlatformRule ----------

func ListPlatformRules(merchantID uint64, keyword string) ([]model.PlatformRule, error) {
	var list []model.PlatformRule
	q := database.DB.Model(&model.PlatformRule{})
	if merchantID > 0 {
		q = q.Where("merchant_id = ?", merchantID)
	}
	if keyword != "" {
		q = q.Where("title LIKE ? OR content LIKE ? OR category LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	err := q.Order("id DESC").Find(&list).Error
	return list, err
}

func GetPlatformRule(id uint64) (*model.PlatformRule, error) {
	var m model.PlatformRule
	if err := database.DB.First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func CreatePlatformRule(m *model.PlatformRule) error {
	return database.DB.Create(m).Error
}

func UpdatePlatformRule(id uint64, m *model.PlatformRule) error {
	m.ID = id
	return database.DB.Save(m).Error
}

func DeletePlatformRule(id uint64) error {
	return database.DB.Delete(&model.PlatformRule{}, id).Error
}

// ---------- ContentTemplate ----------

func ListContentTemplates(merchantID uint64, keyword string) ([]model.ContentTemplate, error) {
	var list []model.ContentTemplate
	q := database.DB.Model(&model.ContentTemplate{})
	if merchantID > 0 {
		q = q.Where("merchant_id = ?", merchantID)
	}
	if keyword != "" {
		q = q.Where("name LIKE ? OR type LIKE ? OR category LIKE ? OR content LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	err := q.Order("id DESC").Find(&list).Error
	return list, err
}

func GetContentTemplate(id uint64) (*model.ContentTemplate, error) {
	var m model.ContentTemplate
	if err := database.DB.First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func CreateContentTemplate(m *model.ContentTemplate) error {
	return database.DB.Create(m).Error
}

func UpdateContentTemplate(id uint64, m *model.ContentTemplate) error {
	m.ID = id
	return database.DB.Save(m).Error
}

func DeleteContentTemplate(id uint64) error {
	return database.DB.Delete(&model.ContentTemplate{}, id).Error
}
