package models

type Recipent struct {
	Recipents []string `json:"recipents" example:"johndoe@gmail.com,janedoe@gmail.com"`
	Success   bool     `json:"success" example:"true"`
}
