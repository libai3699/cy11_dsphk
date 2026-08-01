package database

import (
	"cy11dsphk/server/internal/config"
	"cy11dsphk/server/internal/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func MigrateAndSeed(cfg config.Config) error {
	if err := DB.AutoMigrate(
		&model.AdminUser{},
		&model.Role{},
		&model.Menu{},
		&model.AdminUserRole{},
		&model.RoleMenu{},
		&model.AdminLoginLog{},
	); err != nil {
		return err
	}

	if err := seedRoles(); err != nil {
		return err
	}
	if err := seedMenus(); err != nil {
		return err
	}
	if err := seedRoleMenus(); err != nil {
		return err
	}
	return seedAdminUser(cfg)
}

func seedRoles() error {
	roles := []model.Role{
		{Code: "super_admin", Name: "超级管理员", SortOrder: 1, IsSystem: true},
		{Code: "operator", Name: "运营人员", SortOrder: 2, IsSystem: true},
		{Code: "shooter", Name: "拍摄剪辑", SortOrder: 3, IsSystem: true},
	}
	for _, role := range roles {
		if err := DB.Where("code = ?", role.Code).Assign(role).FirstOrCreate(&model.Role{}).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedMenus() error {
	type seedMenu struct {
		ParentName string
		Menu       model.Menu
	}

	items := []seedMenu{
		{Menu: model.Menu{Name: "WorkspacePage", Path: "/workspace", Component: "/dashboard/WorkspacePage", Title: "工作台", Icon: "carbon:dashboard", SortOrder: -1, AffixTab: true, Visible: true}},
		{Menu: model.Menu{Name: "UserMgmt", Path: "/users", Title: "商家管理", Icon: "carbon:user-multiple", SortOrder: 1, Visible: true}},
		{ParentName: "UserMgmt", Menu: model.Menu{Name: "UserListPage", Path: "/users/list", Component: "/users/UserListPage", Title: "商家列表", Icon: "carbon:user", SortOrder: 11, Visible: true}},
		{ParentName: "UserMgmt", Menu: model.Menu{Name: "UserDevicesPage", Path: "/users/devices", Component: "/users/UserDevicesPage", Title: "账号授权", Icon: "carbon:mobile", SortOrder: 12, Visible: true}},
		{ParentName: "UserMgmt", Menu: model.Menu{Name: "DurationLogPage", Path: "/users/duration-logs", Component: "/users/DurationLogPage", Title: "跟进记录", Icon: "carbon:time", SortOrder: 13, Visible: true}},
		{Menu: model.Menu{Name: "PlanMgmt", Path: "/plans", Title: "成交与分成", Icon: "carbon:purchase", SortOrder: 2, Visible: true}},
		{ParentName: "PlanMgmt", Menu: model.Menu{Name: "PlanListPage", Path: "/plans/list", Component: "/plans/PlanListPage", Title: "团购套餐", Icon: "carbon:list", SortOrder: 21, Visible: true}},
		{ParentName: "PlanMgmt", Menu: model.Menu{Name: "PlanOrdersPage", Path: "/plans/orders", Component: "/plans/PlanOrdersPage", Title: "分成订单", Icon: "carbon:document", SortOrder: 22, Visible: true}},
		{Menu: model.Menu{Name: "LineMgmt", Path: "/lines", Title: "对标分析", Icon: "carbon:network-4", SortOrder: 3, Visible: true}},
		{ParentName: "LineMgmt", Menu: model.Menu{Name: "LineListPage", Path: "/lines/list", Component: "/lines/LineListPage", Title: "对标账号库", Icon: "carbon:network-4", SortOrder: 31, Visible: true}},
		{Menu: model.Menu{Name: "ContentMgmt", Path: "/content", Title: "内容生产", Icon: "carbon:settings", SortOrder: 4, Visible: true}},
		{ParentName: "ContentMgmt", Menu: model.Menu{Name: "ContentNoticesPage", Path: "/content/notices", Component: "/content/ContentNoticesPage", Title: "选题中心", Icon: "carbon:notification", SortOrder: 41, Visible: true}},
		{ParentName: "ContentMgmt", Menu: model.Menu{Name: "ContentQuotesPage", Path: "/content/quotes", Component: "/content/ContentQuotesPage", Title: "文案脚本", Icon: "carbon:quotes", SortOrder: 42, Visible: true}},
		{ParentName: "ContentMgmt", Menu: model.Menu{Name: "ContentDiscoveriesPage", Path: "/content/discoveries", Component: "/content/ContentDiscoveriesPage", Title: "分镜脚本", Icon: "carbon:compass", SortOrder: 43, Visible: true}},
		{ParentName: "ContentMgmt", Menu: model.Menu{Name: "ContentPaymentsPage", Path: "/content/payments", Component: "/content/ContentPaymentsPage", Title: "发布排期", Icon: "carbon:wallet", SortOrder: 44, Visible: true}},
		{ParentName: "ContentMgmt", Menu: model.Menu{Name: "ContentConfigsPage", Path: "/content/configs", Component: "/content/ContentConfigsPage", Title: "系统规则", Icon: "carbon:settings-adjust", SortOrder: 45, Visible: true}},
		{Menu: model.Menu{Name: "LogMgmt", Path: "/logs", Redirect: "/logs/user", Title: "执行复盘", Icon: "carbon:document", SortOrder: 5, Visible: true}},
		{ParentName: "LogMgmt", Menu: model.Menu{Name: "LogUserPage", Path: "/logs/user", Component: "/logs/LogUserPage", Title: "拍摄任务", Icon: "carbon:user-activity", SortOrder: 51, Visible: true}},
		{ParentName: "LogMgmt", Menu: model.Menu{Name: "LogAdminPage", Path: "/logs/admin", Component: "/logs/LogAdminPage", Title: "数据复盘", Icon: "carbon:security", SortOrder: 52, Visible: true}},
	}

	nameToID := map[string]uint{}
	for _, item := range items {
		menu := item.Menu
		if item.ParentName != "" {
			parentID := nameToID[item.ParentName]
			if parentID == 0 {
				var parent model.Menu
				if err := DB.Where("name = ?", item.ParentName).First(&parent).Error; err != nil {
					return err
				}
				parentID = parent.ID
			}
			menu.ParentID = &parentID
		}

		var saved model.Menu
		if err := DB.Where("name = ?", menu.Name).Assign(menu).FirstOrCreate(&saved).Error; err != nil {
			return err
		}
		nameToID[saved.Name] = saved.ID
	}
	return nil
}

func seedRoleMenus() error {
	var allMenus []model.Menu
	if err := DB.Find(&allMenus).Error; err != nil {
		return err
	}

	allMenuNames := make([]string, 0, len(allMenus))
	for _, menu := range allMenus {
		allMenuNames = append(allMenuNames, menu.Name)
	}

	roleMenus := map[string][]string{
		"super_admin": allMenuNames,
		"operator": {
			"WorkspacePage",
			"UserMgmt", "UserListPage", "UserDevicesPage", "DurationLogPage",
			"LineMgmt", "LineListPage",
			"ContentMgmt", "ContentNoticesPage", "ContentQuotesPage", "ContentDiscoveriesPage", "ContentPaymentsPage",
			"LogMgmt", "LogUserPage", "LogAdminPage",
		},
		"shooter": {
			"WorkspacePage",
			"ContentMgmt", "ContentDiscoveriesPage", "ContentPaymentsPage",
			"LogMgmt", "LogUserPage",
		},
	}

	for roleCode, menuNames := range roleMenus {
		var role model.Role
		if err := DB.Where("code = ?", roleCode).First(&role).Error; err != nil {
			return err
		}
		for _, menuName := range menuNames {
			var menu model.Menu
			if err := DB.Where("name = ?", menuName).First(&menu).Error; err != nil {
				return err
			}
			rm := model.RoleMenu{RoleID: role.ID, MenuID: menu.ID}
			if err := DB.Where(rm).FirstOrCreate(&rm).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func seedAdminUser(cfg config.Config) error {
	var admin model.AdminUser
	err := DB.Where("username = ?", cfg.Seed.AdminUsername).First(&admin).Error
	if err == nil {
		return ensureUserRole(admin.ID, "super_admin")
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(cfg.Seed.AdminPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := model.AdminUser{
		Username: cfg.Seed.AdminUsername,
		Password: string(hash),
		RealName: cfg.Seed.AdminRealName,
		Avatar:   "/logo.png",
		Status:   model.AdminUserStatusEnabled,
		HomePath: "/workspace",
	}
	if err := DB.Create(&user).Error; err != nil {
		return err
	}
	return ensureUserRole(user.ID, "super_admin")
}

func ensureUserRole(userID uint64, roleCode string) error {
	var role model.Role
	if err := DB.Where("code = ?", roleCode).First(&role).Error; err != nil {
		return err
	}
	userRole := model.AdminUserRole{UserID: userID, RoleID: role.ID}
	return DB.Where(userRole).FirstOrCreate(&userRole).Error
}
