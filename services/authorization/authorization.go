package authorization

import (
	"fmt"

	"github.com/WilliamTrojniak/TabAppBackend/models"
)

type Action string

type authorizeFn[T any] func(s *models.User, t *T) bool
type authorizeActionMap[T any] map[Action]authorizeFn[T]

func authorizeAction[T any](s *models.User, t *T, action Action, fnMap authorizeActionMap[T]) (bool, error) {
	fn, ok := fnMap[action]
	if !ok {
		return false, fmt.Errorf("Unknown action '%v' specified.", action)
	}
	return fn(s, t), nil

}
