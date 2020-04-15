package endpoints

import (
	"friendMgmt/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func responseOk(c *gin.Context, body interface{}) {
	c.JSON(http.StatusOK, body)
	return
}

func responseError(c *gin.Context, code int, message string) {
	var failure models.Failure
	failure.Success = false
	failure.Message = message
	c.JSON(code, failure)
}
