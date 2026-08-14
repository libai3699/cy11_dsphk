package router

import (
	"errors"
	"net/http"
	"strconv"

	"cy11dsphk/server/internal/database"
	"cy11dsphk/server/internal/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 以下泛型 handler 工厂复用同一套「列表/详情/新建/更新/删除」逻辑，
// 行为对齐现有 admin 路由：统一 JSON 信封 {code,message,data}，错误透传。
// 服务于运营知识库模块，挂在 AdminAuthRequired 中间件之下。

// knowledgeList 列表工厂：返回与既有列表接口一致的契约 {list,total,page,size}。
// svc 参数保留以兼容现有路由注册签名（router.go 调用点不变）；实际查询改为在本工厂内
// 基于 database.DB 完成，where 条件严格复用 service 层各模型的既有逻辑（merchantId + keyword）。
func knowledgeList[T any](svc func(uint64, string) ([]T, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		merchantID := parseUintParam(c.Query("merchantId"))
		keyword := c.Query("keyword")
		page := parseIntQuery(c.Query("page"), 1)
		size := parseIntQuery(c.Query("size"), 10)

		q := knowledgeListQuery[T](merchantID, keyword)

		var total int64
		if err := q.Count(&total).Error; err != nil {
			knowledgeError(c, http.StatusInternalServerError, err.Error())
			return
		}

		var list []T
		if err := q.Order("id DESC").Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
			knowledgeError(c, http.StatusInternalServerError, err.Error())
			return
		}

		knowledgeOK(c, gin.H{"list": list, "total": total, "page": page, "size": size})
	}
}

// knowledgeListQuery 按既有逻辑构建 where 条件。
// merchantId 对所有知识库模型通用；keyword 的列集合因模型而异，故按现有逻辑要求逐模型拼装，
// 与 service/knowledge.go 中 List* 函数的 where 保持一致。
func knowledgeListQuery[T any](merchantID uint64, keyword string) *gorm.DB {
	q := database.DB.Model(new(T))
	if merchantID > 0 {
		q = q.Where("merchant_id = ?", merchantID)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		var cond string
		var args []interface{}
		var t T
		switch any(t).(type) {
		case model.PainPoint:
			cond = "content LIKE ? OR user_quote LIKE ? OR category LIKE ? OR suggested_topic LIKE ?"
			args = []interface{}{like, like, like, like}
		case model.CaseStudy:
			cond = "title LIKE ? OR account_name LIKE ? OR reusable_point LIKE ? OR industry LIKE ?"
			args = []interface{}{like, like, like, like}
		case model.AccountProfile:
			cond = "account_name LIKE ? OR positioning LIKE ? OR industry LIKE ?"
			args = []interface{}{like, like, like}
		case model.PlatformRule:
			cond = "title LIKE ? OR content LIKE ? OR category LIKE ?"
			args = []interface{}{like, like, like}
		case model.ContentTemplate:
			cond = "name LIKE ? OR type LIKE ? OR category LIKE ? OR content LIKE ?"
			args = []interface{}{like, like, like, like}
		}
		if cond != "" {
			q = q.Where(cond, args...)
		}
	}
	return q
}

func parseIntQuery(raw string, fallback int) int {
	v, err := strconv.Atoi(raw)
	if err != nil || v <= 0 {
		return fallback
	}
	return v
}

func knowledgeGet[T any](svc func(uint64) (*T, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := parseIDParam(c)
		if err != nil {
			knowledgeError(c, http.StatusBadRequest, "invalid id")
			return
		}
		data, err := svc(id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				knowledgeError(c, http.StatusNotFound, "not found")
				return
			}
			knowledgeError(c, http.StatusInternalServerError, err.Error())
			return
		}
		knowledgeOK(c, data)
	}
}

func knowledgeCreate[T any](svc func(*T) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		var m T
		if err := c.ShouldBindJSON(&m); err != nil {
			knowledgeError(c, http.StatusBadRequest, err.Error())
			return
		}
		if err := svc(&m); err != nil {
			knowledgeError(c, http.StatusInternalServerError, err.Error())
			return
		}
		knowledgeOK(c, m)
	}
}

func knowledgeUpdate[T any](svc func(uint64, *T) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := parseIDParam(c)
		if err != nil {
			knowledgeError(c, http.StatusBadRequest, "invalid id")
			return
		}
		var m T
		if err := c.ShouldBindJSON(&m); err != nil {
			knowledgeError(c, http.StatusBadRequest, err.Error())
			return
		}
		if err := svc(id, &m); err != nil {
			knowledgeError(c, http.StatusInternalServerError, err.Error())
			return
		}
		knowledgeOK(c, m)
	}
}

func knowledgeDelete(svc func(uint64) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := parseIDParam(c)
		if err != nil {
			knowledgeError(c, http.StatusBadRequest, "invalid id")
			return
		}
		if err := svc(id); err != nil {
			knowledgeError(c, http.StatusInternalServerError, err.Error())
			return
		}
		knowledgeOK(c, gin.H{"id": id})
	}
}

func knowledgeOK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": data})
}

func knowledgeError(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{"code": 1, "message": msg})
}

func parseIDParam(c *gin.Context) (uint64, error) {
	return strconv.ParseUint(c.Param("id"), 10, 64)
}

func parseUintParam(s string) uint64 {
	v, _ := strconv.ParseUint(s, 10, 64)
	return v
}
