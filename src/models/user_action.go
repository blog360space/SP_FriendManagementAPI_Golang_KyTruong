package models

type UserAction struct {
	Requestor string `json:"requestor" example:"johndoe@gmail.com"`
	Target    string `json:"target" example:"janedoe@gmail.com"`
}
