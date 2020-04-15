package models

type UserPost struct {
	Sender string `json:"sender" example:"janedoe@gmail.com"`
	Text   string `json:"text" example:"hello johndoe@gmail.com"`
}
