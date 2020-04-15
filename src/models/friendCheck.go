package models

type FriendCheck struct {
	Friends []string `json:"friends" example:"johndoe@gmail.com,janedoe@gmail.com"`
}
