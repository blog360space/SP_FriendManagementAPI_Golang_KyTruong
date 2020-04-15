package models

type Friend struct {
	Friends []string `json:"friends" example:"johndoe@gmail.com,janedoe@gmail.com"`
	Count   int      `json:"count" example:"2"`
	Success bool     `json:"success" example:"true"`
}
