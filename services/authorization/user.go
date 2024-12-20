package authorization

import "github.com/WilliamTrojniak/TabAppBackend/models"

type userAuthorizeFn = authorizeFn[models.User]

func AuthorizeUserAction(subject *models.User, target *models.User, action Action) (bool, error) {
	return authorizeAction(subject, target, action, userAuthorizeActionFns)
}

const (
	USER_ACTION_UPDATE      Action = "USER_ACTION_UPDATE"
	USER_ACTION_CREATE_SHOP Action = "USER_ACTION_CREATE_SHOP"
)

var userAuthorizeActionFns authorizeActionMap[models.User] = authorizeActionMap[models.User]{
	USER_ACTION_CREATE_SHOP: func(s, t *models.User) bool { return s.Id == t.Id },
	USER_ACTION_UPDATE:      func(s *models.User, t *models.User) bool { return s.Id == t.Id },
}
