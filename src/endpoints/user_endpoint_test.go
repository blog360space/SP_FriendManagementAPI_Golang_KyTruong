package endpoints_test

import (
	"bytes"
	"encoding/json"
	"friendMgmt/data"
	"friendMgmt/endpoints"
	"friendMgmt/models"
	"friendMgmt/services"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUsers(t *testing.T) {
	expectedResult := []string{"user1@gmail.com", "user2@gmail.com"}

	userRepositoryMock := data.UserRepositoryMock{}
	userRepositoryMock.On("FindAll").Return(expectedResult)

	userServiceMock := services.UserServiceMock{}
	userServiceMock.On("FindAll").Return(expectedResult)

	userEndpoint := endpoints.UserEndpoint{userServiceMock}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	userEndpoint.Users(c)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	var actualResult []string
	body, _ := ioutil.ReadAll(w.Result().Body)
	err := json.Unmarshal(body, &actualResult)

	assert.Equal(t, err, nil)
	assert.Equal(t, actualResult, expectedResult)
}

func TestCreateWithExistedEmail(t *testing.T) {
	var jsonStr = []byte(`{"email": "user@test.com"}`)

	userRepositoryMock := data.UserRepositoryMock{}
	userRepositoryMock.On("CheckUserExist", "user@test.com").Return(int64(1))

	userServiceMock := services.UserServiceMock{}
	userServiceMock.On("CheckUserExist", "user@test.com").Return(int64(1))

	userEndpoint := endpoints.UserEndpoint{userServiceMock}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/users", bytes.NewBuffer(jsonStr))
	c.Request.Header.Set("Content-Type", "application/json")

	userEndpoint.CreateUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)

	var actualResult models.Failure
	body, _ := ioutil.ReadAll(w.Result().Body)
	json.Unmarshal(body, &actualResult)

	assert.Equal(t, actualResult.Success, false)
	assert.Equal(t, actualResult.Message, "Invalid request: the email is already in use")
}

func TestCreateWithInvalidEmail(t *testing.T) {
	var invalidRequests = []string{`{"email": "fake_email"}`, `{"email": ""}`, `{"email": "almost@valid.email'"}`}
	for _, request := range invalidRequests {

		var jsonStr = []byte(request)

		userServiceMock := services.UserServiceMock{}

		userEndpoint := endpoints.UserEndpoint{userServiceMock}
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/users", bytes.NewBuffer(jsonStr))
		c.Request.Header.Set("Content-Type", "application/json")

		userEndpoint.CreateUser(c)

		assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)

		var actualResult models.Failure
		body, _ := ioutil.ReadAll(w.Result().Body)
		json.Unmarshal(body, &actualResult)

		assert.Equal(t, false, actualResult.Success)
		assert.Equal(t, "Invalid request: incorrect info", actualResult.Message)
	}
}

func TestCreateWithValidEmail(t *testing.T) {
	var jsonStr = []byte(`{"email": "user@test.com"}`)

	userRepositoryMock := data.UserRepositoryMock{}
	userRepositoryMock.On("CheckUserExist", "user@test.com").Return(int64(-1))
	userRepositoryMock.On("Create", "user@test.com").Return(true)

	userServiceMock := services.UserServiceMock{}
	userServiceMock.On("CheckUserExist", "user@test.com").Return(int64(-1))
	userServiceMock.On("Create", "user@test.com").Return(true)

	userEndpoint := endpoints.UserEndpoint{userServiceMock}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/users", bytes.NewBuffer(jsonStr))
	c.Request.Header.Set("Content-Type", "application/json")

	userEndpoint.CreateUser(c)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	var actualResult models.Success
	body, _ := ioutil.ReadAll(w.Result().Body)
	json.Unmarshal(body, &actualResult)

	assert.Equal(t, actualResult.Success, true)
}
