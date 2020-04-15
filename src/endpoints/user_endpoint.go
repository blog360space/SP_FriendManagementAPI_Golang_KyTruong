package endpoints

import (
	"fmt"
	"friendMgmt/common"
	"friendMgmt/models"
	"friendMgmt/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserEndpoint struct {
	IUserService services.IUserService
}

// Users godoc
// @Tags User
// @Summary API to get all users in app
// @Accept  json
// @Produce  json
// @Success 200 {array} string
// @Router /users [get]
func (u UserEndpoint) Users(c *gin.Context) {
	emails := u.IUserService.FindAll()

	responseOk(c, emails)
}

// CreateUser godoc
// @Tags User
// @Summary API to create new user
// @Description create user
// @Accept  json
// @Produce  json
// @Param email body models.Email true "Body"
// @Success 200 {object} models.Success "OK"
// @Failure 400 {object} models.Failure "Bad Request"
// @Router /users [post]
func (u UserEndpoint) CreateUser(c *gin.Context) {
	var emailModel models.Email
	err := c.BindJSON(&emailModel)

	fmt.Println(emailModel.Email)

	if err != nil || !common.IsValidEmail(emailModel.Email) {
		responseError(c, http.StatusBadRequest, "Invalid request: incorrect info")
		return
	}

	if userId := u.IUserService.CheckUserExist(emailModel.Email); userId > 0 {
		responseError(c, http.StatusBadRequest, "Invalid request: the email is already in use")
		return
	}

	u.IUserService.Create(emailModel.Email)

	success := models.Success{Success: true}

	responseOk(c, success)
}
