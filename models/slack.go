package models

type SlackChannel struct {
	IsPrivate bool   `json:"is_private"`
	Name      string `json:"name"`
}
