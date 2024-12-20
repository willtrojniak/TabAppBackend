package authorization

import "github.com/WilliamTrojniak/TabAppBackend/models"

type TabTarget struct {
	Shop *models.Shop
	Tab  *models.Tab
}

type tabAuthorizeFn = authorizeFn[TabTarget]

func AuthorizeTabAction(subject *models.User, target *TabTarget, action Action) (bool, error) {
	return authorizeAction(subject, target, action, tabAuthorizeActionFns)
}

const (
	TAB_ACTION_READ           Action = "TAB_ACTION_READ"
	TAB_ACTION_REQUEST_UPDATE Action = "TAB_ACTION_REQUEST_UPDATE"
	TAB_ACTION_UPDATE         Action = "TAB_ACTION_UPDATE"
	TAB_ACTION_APPROVE        Action = "TAB_ACTION_APPROVE"
	TAB_ACTION_CLOSE          Action = "TAB_ACTION_CLOSE"
	TAB_ACTION_CLOSE_BILL     Action = "TAB_ACTION_CLOSE_BILL"
	TAB_ACTION_ADD_ORDER      Action = "TAB_ACTION_ADD_ORDER"
	TAB_ACTION_REMOVE_ORDER   Action = "TAB_ACTION_REMOVE_ORDER"
)

var tabAuthorizeActionFns authorizeActionMap[TabTarget] = authorizeActionMap[TabTarget]{
	TAB_ACTION_READ:         func(s *models.User, t *TabTarget) bool { return true },
	TAB_ACTION_UPDATE:       func(s *models.User, t *TabTarget) bool { return true },
	TAB_ACTION_APPROVE:      func(s *models.User, t *TabTarget) bool { return true },
	TAB_ACTION_CLOSE:        func(s *models.User, t *TabTarget) bool { return true },
	TAB_ACTION_CLOSE_BILL:   func(s *models.User, t *TabTarget) bool { return true },
	TAB_ACTION_ADD_ORDER:    func(s *models.User, t *TabTarget) bool { return true },
	TAB_ACTION_REMOVE_ORDER: func(s *models.User, t *TabTarget) bool { return true },
}
