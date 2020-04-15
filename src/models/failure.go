package models

type Failure struct {
	Message string `json:"message" example:"error message"`
	Success bool   `json:"success" example:"false"`
}
