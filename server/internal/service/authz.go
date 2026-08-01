package service

import (
	"cy11dsphk/server/internal/database"
	"cy11dsphk/server/internal/model"
)

type MenuRoute struct {
	Name      string      `json:"name"`
	Path      string      `json:"path"`
	Component string      `json:"component,omitempty"`
	Redirect  string      `json:"redirect,omitempty"`
	Meta      MenuMeta    `json:"meta"`
	Children  []MenuRoute `json:"children,omitempty"`
}

type MenuMeta struct {
	Title      string `json:"title"`
	Icon       string `json:"icon,omitempty"`
	Order      int    `json:"order,omitempty"`
	AffixTab   bool   `json:"affixTab,omitempty"`
	HideInMenu bool   `json:"hideInMenu,omitempty"`
}

func GetUserRoles(userID uint64) ([]model.Role, error) {
	var roles []model.Role
	err := database.DB.
		Table("roles").
		Joins("JOIN admin_user_roles ON admin_user_roles.role_id = roles.id").
		Where("admin_user_roles.user_id = ?", userID).
		Order("roles.sort_order ASC, roles.id ASC").
		Find(&roles).Error
	return roles, err
}

func GetUserRoleCodes(userID uint64) ([]string, error) {
	roles, err := GetUserRoles(userID)
	if err != nil {
		return nil, err
	}
	codes := make([]string, 0, len(roles))
	for _, role := range roles {
		codes = append(codes, role.Code)
	}
	return codes, nil
}

func GetUserMenus(userID uint64) ([]model.Menu, error) {
	var menus []model.Menu
	err := database.DB.
		Table("menus").
		Distinct("menus.*").
		Joins("JOIN role_menus ON role_menus.menu_id = menus.id").
		Joins("JOIN admin_user_roles ON admin_user_roles.role_id = role_menus.role_id").
		Where("admin_user_roles.user_id = ? AND menus.visible = ?", userID, true).
		Order("menus.sort_order ASC, menus.id ASC").
		Find(&menus).Error
	return menus, err
}

func GetUserAccessCodes(userID uint64) ([]string, error) {
	roles, err := GetUserRoles(userID)
	if err != nil {
		return nil, err
	}
	menus, err := GetUserMenus(userID)
	if err != nil {
		return nil, err
	}

	codes := make([]string, 0, len(roles)+len(menus))
	for _, role := range roles {
		codes = append(codes, "ROLE_"+role.Code)
	}
	for _, menu := range menus {
		codes = append(codes, "MENU_"+menu.Name)
	}
	return codes, nil
}

func BuildMenuRoutes(menus []model.Menu) []MenuRoute {
	childrenByParent := map[uint][]model.Menu{}
	for _, menu := range menus {
		parentID := uint(0)
		if menu.ParentID != nil {
			parentID = *menu.ParentID
		}
		childrenByParent[parentID] = append(childrenByParent[parentID], menu)
	}
	return buildMenuRoutesFromParent(0, childrenByParent)
}

func buildMenuRoutesFromParent(parentID uint, childrenByParent map[uint][]model.Menu) []MenuRoute {
	children := childrenByParent[parentID]
	routes := make([]MenuRoute, 0, len(children))
	for _, menu := range children {
		route := MenuRoute{
			Name:      menu.Name,
			Path:      menu.Path,
			Component: menu.Component,
			Redirect:  menu.Redirect,
			Meta: MenuMeta{
				Title:      menu.Title,
				Icon:       menu.Icon,
				Order:      menu.SortOrder,
				AffixTab:   menu.AffixTab,
				HideInMenu: menu.HideInMenu,
			},
			Children: buildMenuRoutesFromParent(menu.ID, childrenByParent),
		}
		routes = append(routes, route)
	}
	return routes
}
