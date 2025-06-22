package authorization

import (
	"slices"

	"github.com/willtrojniak/TabAppBackend/models"
)

type shopAuthorizeFn = authorizeFn[models.Shop]

func AuthorizeShopAction(subject *models.User, target *models.Shop, action Action) (bool, error) {
	return authorizeAction(subject, target, action, shopAuthorizeActionFns)
}

const (
	SHOP_ACTION_READ                Action = "SHOP_ACTION_READ"
	SHOP_ACTION_INVITE_USER         Action = "SHOP_ACTION_INVITE_USER"
	SHOP_ACTION_REMOVE_USER         Action = "SHOP_ACTION_REMOVE_USER"
	SHOP_ACTION_UPDATE              Action = "SHOP_ACTION_UPDATE"
	SHOP_ACTION_DELETE              Action = "SHOP_ACTION_DELETE"
	SHOP_ACTION_CREATE_LOCATION     Action = "SHOP_ACTION_CREATE_LOCATION"
	SHOP_ACTION_UPDATE_LOCATION     Action = "SHOP_ACTION_UPDATE_LOCATION"
	SHOP_ACTION_DELETE_LOCATION     Action = "SHOP_ACTION_DELETE_LOCATION"
	SHOP_ACTION_READ_CATEGORIES     Action = "SHOP_ACTION_READ_CATEGORIES"
	SHOP_ACTION_CREATE_CATEGORY     Action = "SHOP_ACTION_CREATE_CATEGORY"
	SHOP_ACTION_UPDATE_CATEGORY     Action = "SHOP_ACTION_UPDATE_CATEGORY"
	SHOP_ACTION_DELETE_CATEGORY     Action = "SHOP_ACTION_DELETE_CATEGORY"
	SHOP_ACTION_READ_ITEMS          Action = "SHOP_ACTION_READ_ITEMS"
	SHOP_ACTION_READ_ITEM           Action = "SHOP_ACTION_READ_ITEM"
	SHOP_ACTION_CREATE_ITEM         Action = "SHOP_ACTION_CREATE_ITEM"
	SHOP_ACTION_UPDATE_ITEM         Action = "SHOP_ACTION_UPDATE_ITEM"
	SHOP_ACTION_DELETE_ITEM         Action = "SHOP_ACTION_DELETE_ITEM"
	SHOP_ACTION_CREATE_VARIANT      Action = "SHOP_ACTION_CREATE_VARIANT"
	SHOP_ACTION_UPDATE_VARIANT      Action = "SHOP_ACTION_UPDATE_VARIANT"
	SHOP_ACTION_DELETE_VARIANT      Action = "SHOP_ACTION_DELETE_VARIANT"
	SHOP_ACTION_READ_SUBSTITUTIONS  Action = "SHOP_ACTION_READ_SUBSTITUTIONS"
	SHOP_ACTION_CREATE_SUBSTITUTION Action = "SHOP_ACTION_CREATE_SUBSTITUTION"
	SHOP_ACTION_UPDATE_SUBSTITUTION Action = "SHOP_ACTION_UPDATE_SUBSTITUTION"
	SHOP_ACTION_DELETE_SUBSTITUTION Action = "SHOP_ACTION_DELETE_SUBSTITUTION"
	SHOP_ACTION_READ_TABS           Action = "SHOP_ACTION_READ_TABS"
	SHOP_ACTION_REQUEST_TAB         Action = "SHOP_ACTION_REQUEST_TAB"
	SHOP_ACTION_CREATE_TAB          Action = "SHOP_ACTION_CREATE_TAB"
)

const (
	ROLE_SHOP_MANAGE_ITEMS     uint32 = 1 << 0 // 1
	ROLE_SHOP_MANAGE_TABS      uint32 = 1 << 1 // 2
	ROLE_SHOP_MANAGE_ORDERS    uint32 = 1 << 2 // 4
	ROLE_SHOP_READ_TABS        uint32 = 1 << 3 // 8
	ROLE_SHOP_MANAGE_LOCATIONS uint32 = 1 << 4 // 16
)

var shopAuthorizeActionFns authorizeActionMap[models.Shop] = authorizeActionMap[models.Shop]{
	SHOP_ACTION_READ:                func(s *models.User, t *models.Shop) bool { return true },
	SHOP_ACTION_INVITE_USER:         func(s *models.User, t *models.Shop) bool { return s.Id == t.OwnerId },
	SHOP_ACTION_REMOVE_USER:         func(s *models.User, t *models.Shop) bool { return s.Id == t.OwnerId },
	SHOP_ACTION_UPDATE:              func(s *models.User, t *models.Shop) bool { return s.Id == t.OwnerId },
	SHOP_ACTION_DELETE:              func(s *models.User, t *models.Shop) bool { return s.Id == t.OwnerId },
	SHOP_ACTION_CREATE_LOCATION:     func(s *models.User, t *models.Shop) bool { return hasRole(s, t, ROLE_SHOP_MANAGE_LOCATIONS) },
	SHOP_ACTION_UPDATE_LOCATION:     func(s *models.User, t *models.Shop) bool { return hasRole(s, t, ROLE_SHOP_MANAGE_LOCATIONS) },
	SHOP_ACTION_DELETE_LOCATION:     func(s *models.User, t *models.Shop) bool { return hasRole(s, t, ROLE_SHOP_MANAGE_LOCATIONS) },
	SHOP_ACTION_READ_CATEGORIES:     func(s *models.User, t *models.Shop) bool { return true },
	SHOP_ACTION_CREATE_CATEGORY:     func(s *models.User, t *models.Shop) bool { return hasRole(s, t, ROLE_SHOP_MANAGE_ITEMS) },
	SHOP_ACTION_UPDATE_CATEGORY:     func(s *models.User, t *models.Shop) bool { return hasRole(s, t, ROLE_SHOP_MANAGE_ITEMS) },
	SHOP_ACTION_DELETE_CATEGORY:     func(s *models.User, t *models.Shop) bool { return hasRole(s, t, ROLE_SHOP_MANAGE_ITEMS) },
	SHOP_ACTION_READ_ITEMS:          func(s *models.User, t *models.Shop) bool { return true },
	SHOP_ACTION_READ_ITEM:           func(s *models.User, t *models.Shop) bool { return true },
	SHOP_ACTION_CREATE_ITEM:         func(s *models.User, t *models.Shop) bool { return hasRole(s, t, ROLE_SHOP_MANAGE_ITEMS) },
	SHOP_ACTION_UPDATE_ITEM:         func(s *models.User, t *models.Shop) bool { return hasRole(s, t, ROLE_SHOP_MANAGE_ITEMS) },
	SHOP_ACTION_DELETE_ITEM:         func(s *models.User, t *models.Shop) bool { return hasRole(s, t, ROLE_SHOP_MANAGE_ITEMS) },
	SHOP_ACTION_CREATE_VARIANT:      func(s *models.User, t *models.Shop) bool { return hasRole(s, t, ROLE_SHOP_MANAGE_ITEMS) },
	SHOP_ACTION_UPDATE_VARIANT:      func(s *models.User, t *models.Shop) bool { return hasRole(s, t, ROLE_SHOP_MANAGE_ITEMS) },
	SHOP_ACTION_DELETE_VARIANT:      func(s *models.User, t *models.Shop) bool { return hasRole(s, t, ROLE_SHOP_MANAGE_ITEMS) },
	SHOP_ACTION_READ_SUBSTITUTIONS:  func(s *models.User, t *models.Shop) bool { return true },
	SHOP_ACTION_CREATE_SUBSTITUTION: func(s *models.User, t *models.Shop) bool { return hasRole(s, t, ROLE_SHOP_MANAGE_ITEMS) },
	SHOP_ACTION_UPDATE_SUBSTITUTION: func(s *models.User, t *models.Shop) bool { return hasRole(s, t, ROLE_SHOP_MANAGE_ITEMS) },
	SHOP_ACTION_DELETE_SUBSTITUTION: func(s *models.User, t *models.Shop) bool { return hasRole(s, t, ROLE_SHOP_MANAGE_ITEMS) },
	SHOP_ACTION_READ_TABS:           func(s *models.User, t *models.Shop) bool { return hasRole(s, t, ROLE_SHOP_READ_TABS) },
	SHOP_ACTION_REQUEST_TAB:         func(s *models.User, t *models.Shop) bool { return true },
	SHOP_ACTION_CREATE_TAB:          func(s *models.User, t *models.Shop) bool { return hasRole(s, t, ROLE_SHOP_MANAGE_TABS) },
}

func hasRole(s *models.User, t *models.Shop, role uint32) bool {
	return s.Id == t.OwnerId || slices.ContainsFunc(t.Users, func(u models.ShopUser) bool {
		return u.Id == s.Id && u.Roles&role == role
	})
}
