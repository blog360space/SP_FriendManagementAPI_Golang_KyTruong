package endpoints_test

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func TestCreateRelationshipWithInvalidAccounts(t *testing.T) {
	var invalidRequests = []string{
		`{"friends":"target@email.com"}`,
		`{"friends":["request","target@email.com"]}`,
		`{"friends":["target@email.com","request"]}`,
		`{"friends":["target@email.com"]}`,
		`{"friends":["request@email.com","request@email.com"]}`,
		`{"friends":["request@email.com","target@email.com","unknown@email.com"]}`}
	for _, request := range invalidRequests {

		var jsonStr = []byte(request)

		relationshipServiceMock := services.RelationshipServiceMock{}
		userServiceMock := services.UserServiceMock{}

		relationshipEndpoint := endpoints.RelationshipEndpoint{relationshipServiceMock, userServiceMock}
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/friends/add", bytes.NewBuffer(jsonStr))
		c.Request.Header.Set("Content-Type", "application/json")

		relationshipEndpoint.CreateRelationship(c)

		assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)

		var actualResult models.Failure
		body, _ := ioutil.ReadAll(w.Result().Body)
		json.Unmarshal(body, &actualResult)

		assert.Equal(t, false, actualResult.Success)
		assert.Equal(t, "Invalid request: incorrect info", actualResult.Message)
	}
}

func TestCreateRelationshipWithNotFoundAccounts(t *testing.T) {
	var invalidRequests = []string{
		`{"friends":["undefined@request.com","email@target.com"]}`,
		`{"friends":["email@request.com","undefined@target.com"]}`}

	for i, request := range invalidRequests {

		var jsonStr = []byte(request)

		friendCheckObj := models.FriendCheck{}
		json.Unmarshal(jsonStr, &friendCheckObj)

		var undefinedEmail = friendCheckObj.Friends[0]
		if i == 1 {
			undefinedEmail = friendCheckObj.Friends[1]
		}

		relationshipServiceMock := services.RelationshipServiceMock{}
		userServiceMock := services.UserServiceMock{}

		if i == 0 {
			userRepositoryMock := data.UserRepositoryMock{}
			userRepositoryMock.On("CheckUserExist", undefinedEmail).Return(int64(-1))

			userServiceMock.On("CheckUserExist", undefinedEmail).Return(int64(-1))
		} else {
			userRepositoryMock := data.UserRepositoryMock{}
			userRepositoryMock.On("CheckUserExist", friendCheckObj.Friends[0]).Return(int64(1))
			userRepositoryMock.On("CheckUserExist", undefinedEmail).Return(int64(-1))

			userServiceMock.On("CheckUserExist", friendCheckObj.Friends[0]).Return(int64(1))
			userServiceMock.On("CheckUserExist", undefinedEmail).Return(int64(-1))
		}

		relationshipEndpoint := endpoints.RelationshipEndpoint{relationshipServiceMock, userServiceMock}
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/friends/add", bytes.NewBuffer(jsonStr))
		c.Request.Header.Set("Content-Type", "application/json")

		relationshipEndpoint.CreateRelationship(c)

		assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)

		var actualResult models.Failure
		body, _ := ioutil.ReadAll(w.Result().Body)
		json.Unmarshal(body, &actualResult)

		assert.Equal(t, false, actualResult.Success)
		assert.Equal(t, fmt.Sprintf("Invalid request: User name %s is not found", undefinedEmail), actualResult.Message)
	}
}

func TestCreateRelationshipForAlreadyConnectedAccounts(t *testing.T) {
	var jsonStr = []byte(`{"friends":["email@request.com","email@target.com"]}`)

	friendCheckObj := models.FriendCheck{}
	json.Unmarshal(jsonStr, &friendCheckObj)

	relationshipRepositoryMock := data.RelationshipRepositoryMock{}
	userRepositoryMock := data.UserRepositoryMock{}
	relationshipServiceMock := services.RelationshipServiceMock{}
	userServiceMock := services.UserServiceMock{}

	requestUser := friendCheckObj.Friends[0]
	requestUserId := int64(1)
	targetUser := friendCheckObj.Friends[1]
	targetUserId := int64(2)
	status := int64(1)
	relationshipId := int64(1)

	userRepositoryMock.On("CheckUserExist", requestUser).Return(requestUserId)
	userServiceMock.On("CheckUserExist", requestUser).Return(requestUserId)

	userRepositoryMock.On("CheckUserExist", targetUser).Return(targetUserId)
	userServiceMock.On("CheckUserExist", targetUser).Return(targetUserId)

	relationshipRepositoryMock.On("CheckRelationshipTwoWay", requestUserId, targetUserId, status).Return([]int64{relationshipId})
	relationshipServiceMock.On("CheckConnected", requestUserId, targetUserId).Return([]int64{relationshipId})

	relationshipEndpoint := endpoints.RelationshipEndpoint{relationshipServiceMock, userServiceMock}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/friends/add", bytes.NewBuffer(jsonStr))
	c.Request.Header.Set("Content-Type", "application/json")

	relationshipEndpoint.CreateRelationship(c)

	assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)

	var actualResult models.Failure
	body, _ := ioutil.ReadAll(w.Result().Body)
	json.Unmarshal(body, &actualResult)

	assert.Equal(t, false, actualResult.Success)
	assert.Equal(t, "Invalid request: connected status is existed", actualResult.Message)
}

func TestCreateRelationshipForAlreadyBlockedAccounts(t *testing.T) {
	var jsonStr = []byte(`{"friends":["email@request.com","email@target.com"]}`)

	friendCheckObj := models.FriendCheck{}
	json.Unmarshal(jsonStr, &friendCheckObj)

	relationshipRepositoryMock := data.RelationshipRepositoryMock{}
	userRepositoryMock := data.UserRepositoryMock{}
	relationshipServiceMock := services.RelationshipServiceMock{}
	userServiceMock := services.UserServiceMock{}

	requestUser := friendCheckObj.Friends[0]
	requestUserId := int64(1)
	targetUser := friendCheckObj.Friends[1]
	targetUserId := int64(2)
	connectedStatus := int64(1)
	blockedStatus := int64(2)
	connectedIds := []int64{}
	blockedIds := []int64{int64(1)}

	userRepositoryMock.On("CheckUserExist", requestUser).Return(requestUserId)
	userServiceMock.On("CheckUserExist", requestUser).Return(requestUserId)

	userRepositoryMock.On("CheckUserExist", targetUser).Return(targetUserId)
	userServiceMock.On("CheckUserExist", targetUser).Return(targetUserId)

	relationshipRepositoryMock.On("CheckRelationshipTwoWay", requestUserId, targetUserId, connectedStatus).Return(connectedIds)
	relationshipServiceMock.On("CheckConnected", requestUserId, targetUserId).Return(connectedIds)

	relationshipRepositoryMock.On("CheckRelationshipTwoWay", requestUserId, targetUserId, blockedStatus).Return(blockedIds)
	relationshipServiceMock.On("CheckFullyBlocked", requestUserId, targetUserId).Return(blockedIds)

	relationshipEndpoint := endpoints.RelationshipEndpoint{relationshipServiceMock, userServiceMock}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/friends/add", bytes.NewBuffer(jsonStr))
	c.Request.Header.Set("Content-Type", "application/json")

	relationshipEndpoint.CreateRelationship(c)

	assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)

	var actualResult models.Failure
	body, _ := ioutil.ReadAll(w.Result().Body)
	json.Unmarshal(body, &actualResult)

	assert.Equal(t, false, actualResult.Success)
	assert.Equal(t, "Invalid request: blocked status is existed", actualResult.Message)
}

func TestCreateRelationshipForAlreadySubcribedAccounts(t *testing.T) {
	var jsonStr = []byte(`{"friends":["email@request.com","email@target.com"]}`)

	friendCheckObj := models.FriendCheck{}
	json.Unmarshal(jsonStr, &friendCheckObj)

	relationshipRepositoryMock := data.RelationshipRepositoryMock{}
	userRepositoryMock := data.UserRepositoryMock{}
	relationshipServiceMock := services.RelationshipServiceMock{}
	userServiceMock := services.UserServiceMock{}

	requestUser := friendCheckObj.Friends[0]
	requestUserId := int64(1)
	targetUser := friendCheckObj.Friends[1]
	targetUserId := int64(2)
	connectedStatus := int64(1)
	blockedStatus := int64(3)
	subcribedStatus := int64(2)
	connectedIds := []int64{}
	blockedIds := []int64{}
	subcribedIds := []int64{int64(1)}

	userRepositoryMock.On("CheckUserExist", requestUser).Return(requestUserId)
	userServiceMock.On("CheckUserExist", requestUser).Return(requestUserId)

	userRepositoryMock.On("CheckUserExist", targetUser).Return(targetUserId)
	userServiceMock.On("CheckUserExist", targetUser).Return(targetUserId)

	relationshipRepositoryMock.On("CheckRelationshipTwoWay", requestUserId, targetUserId, connectedStatus).Return(connectedIds)
	relationshipServiceMock.On("CheckConnected", requestUserId, targetUserId).Return(connectedIds)

	relationshipRepositoryMock.On("CheckRelationshipTwoWay", requestUserId, targetUserId, blockedStatus).Return(blockedIds)
	relationshipServiceMock.On("CheckFullyBlocked", requestUserId, targetUserId).Return(blockedIds)

	relationshipRepositoryMock.On("CheckRelationshipTwoWay", requestUserId, targetUserId, subcribedStatus).Return(subcribedIds)
	relationshipServiceMock.On("CheckFullySubcribed", requestUserId, targetUserId).Return(subcribedIds)

	relationshipRepositoryMock.On("DeleteRelationships", subcribedIds).Return(true)
	relationshipServiceMock.On("DeleteRelationships", subcribedIds).Return(true)

	relationshipModel := models.Relationship{Status: connectedStatus, RequestUserId: requestUserId, TargetUserId: targetUserId}
	relationshipRepositoryMock.On("CreateRelationship", &relationshipModel).Return(int64(10))
	relationshipServiceMock.On("CreateRelationship", &relationshipModel).Return(int64(10))

	relationshipEndpoint := endpoints.RelationshipEndpoint{relationshipServiceMock, userServiceMock}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/friends/add", bytes.NewBuffer(jsonStr))
	c.Request.Header.Set("Content-Type", "application/json")

	relationshipEndpoint.CreateRelationship(c)

	assert.Equal(t, w.Result().StatusCode, http.StatusOK)

	var actualResult models.Success
	body, _ := ioutil.ReadAll(w.Result().Body)
	json.Unmarshal(body, &actualResult)

	assert.Equal(t, true, actualResult.Success)
}

func TestCreateRelationshipReturnInternalError(t *testing.T) {
	var jsonStr = []byte(`{"friends":["email@request.com","email@target.com"]}`)

	friendCheckObj := models.FriendCheck{}
	json.Unmarshal(jsonStr, &friendCheckObj)

	relationshipRepositoryMock := data.RelationshipRepositoryMock{}
	userRepositoryMock := data.UserRepositoryMock{}
	relationshipServiceMock := services.RelationshipServiceMock{}
	userServiceMock := services.UserServiceMock{}

	requestUser := friendCheckObj.Friends[0]
	requestUserId := int64(1)
	targetUser := friendCheckObj.Friends[1]
	targetUserId := int64(2)
	connectedStatus := int64(1)
	blockedStatus := int64(3)
	subcribedStatus := int64(2)
	connectedIds := []int64{}
	blockedIds := []int64{}
	subcribedIds := []int64{}

	userRepositoryMock.On("CheckUserExist", requestUser).Return(requestUserId)
	userServiceMock.On("CheckUserExist", requestUser).Return(requestUserId)

	userRepositoryMock.On("CheckUserExist", targetUser).Return(targetUserId)
	userServiceMock.On("CheckUserExist", targetUser).Return(targetUserId)

	relationshipRepositoryMock.On("CheckRelationshipTwoWay", requestUserId, targetUserId, connectedStatus).Return(connectedIds)
	relationshipServiceMock.On("CheckConnected", requestUserId, targetUserId).Return(connectedIds)

	relationshipRepositoryMock.On("CheckRelationshipTwoWay", requestUserId, targetUserId, blockedStatus).Return(blockedIds)
	relationshipServiceMock.On("CheckFullyBlocked", requestUserId, targetUserId).Return(blockedIds)

	relationshipRepositoryMock.On("CheckRelationshipTwoWay", requestUserId, targetUserId, subcribedStatus).Return(subcribedIds)
	relationshipServiceMock.On("CheckFullySubcribed", requestUserId, targetUserId).Return(subcribedIds)

	relationshipModel := models.Relationship{Status: connectedStatus, RequestUserId: requestUserId, TargetUserId: targetUserId}
	relationshipRepositoryMock.On("CreateRelationship", &relationshipModel).Return(int64(-1))
	relationshipServiceMock.On("CreateRelationship", &relationshipModel).Return(int64(-1))

	relationshipEndpoint := endpoints.RelationshipEndpoint{relationshipServiceMock, userServiceMock}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/friends/add", bytes.NewBuffer(jsonStr))
	c.Request.Header.Set("Content-Type", "application/json")

	relationshipEndpoint.CreateRelationship(c)

	assert.Equal(t, w.Result().StatusCode, http.StatusInternalServerError)

	var actualResult models.Failure
	body, _ := ioutil.ReadAll(w.Result().Body)
	json.Unmarshal(body, &actualResult)

	assert.Equal(t, false, actualResult.Success)
	assert.Equal(t, "Oops! There is an error, please try again.", actualResult.Message)
}

func TestFriendListWithInvalidAccount(t *testing.T) {
	var invalidRequests = []string{
		`{"email":}`,
		`{"email":"almost@valid.email'"}`,
		`{"email":"invalid_email"}`}
	for _, request := range invalidRequests {

		var jsonStr = []byte(request)

		relationshipServiceMock := services.RelationshipServiceMock{}
		userServiceMock := services.UserServiceMock{}

		relationshipEndpoint := endpoints.RelationshipEndpoint{relationshipServiceMock, userServiceMock}
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/friends", bytes.NewBuffer(jsonStr))
		c.Request.Header.Set("Content-Type", "application/json")

		relationshipEndpoint.FriendList(c)

		assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)

		var actualResult models.Failure
		body, _ := ioutil.ReadAll(w.Result().Body)
		json.Unmarshal(body, &actualResult)

		assert.Equal(t, false, actualResult.Success)
		assert.Equal(t, "Invalid request: incorrect info", actualResult.Message)
	}
}

func TestFriendListWithNotFoundAccount(t *testing.T) {
	var jsonStr = []byte(`{"email":"not@found.account"}`)

	email := models.Email{}
	json.Unmarshal(jsonStr, &email)

	relationshipServiceMock := services.RelationshipServiceMock{}
	userServiceMock := services.UserServiceMock{}
	userServiceMock.On("CheckUserExist", email.Email).Return(int64(-1))

	userRepositoryMock := data.UserRepositoryMock{}
	userRepositoryMock.On("CheckUserExist", email.Email).Return(int64(-1))

	relationshipEndpoint := endpoints.RelationshipEndpoint{relationshipServiceMock, userServiceMock}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/friends/add", bytes.NewBuffer(jsonStr))
	c.Request.Header.Set("Content-Type", "application/json")

	relationshipEndpoint.FriendList(c)

	assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)

	var actualResult models.Failure
	body, _ := ioutil.ReadAll(w.Result().Body)
	json.Unmarshal(body, &actualResult)

	assert.Equal(t, false, actualResult.Success)
	assert.Equal(t, fmt.Sprintf("Invalid request: User name %s is not found", email.Email), actualResult.Message)
}

func TestFriendListWithValidAccount(t *testing.T) {
	var jsonStr = []byte(`{"email":"not@found.account"}`)

	email := models.Email{}
	json.Unmarshal(jsonStr, &email)

	relationshipServiceMock := services.RelationshipServiceMock{}
	relationshipRepositoryMock := data.RelationshipRepositoryMock{}

	userServiceMock := services.UserServiceMock{}
	userRepositoryMock := data.UserRepositoryMock{}

	userServiceMock.On("CheckUserExist", email.Email).Return(int64(1))
	userRepositoryMock.On("CheckUserExist", email.Email).Return(int64(1))

	friendList := []string{"user1@email.com", "user2@email.com"}
	relationshipServiceMock.On("GetFriendList", int64(1)).Return(friendList)
	relationshipRepositoryMock.On("GetFriendList", int64(1)).Return(friendList)

	relationshipEndpoint := endpoints.RelationshipEndpoint{relationshipServiceMock, userServiceMock}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/friends/add", bytes.NewBuffer(jsonStr))
	c.Request.Header.Set("Content-Type", "application/json")

	relationshipEndpoint.FriendList(c)

	assert.Equal(t, w.Result().StatusCode, http.StatusOK)

	var actualResult models.Success
	body, _ := ioutil.ReadAll(w.Result().Body)
	json.Unmarshal(body, &actualResult)

	assert.Equal(t, true, actualResult.Success)
}

func TestCommonFriendListWithInvalidAccounts(t *testing.T) {
	var invalidRequests = []string{
		`{"friends":"target@email.com"}`,
		`{"friends":["request","target@email.com"]}`,
		`{"friends":["target@email.com","request"]}`,
		`{"friends":["target@email.com"]}`,
		`{"friends":["request@email.com","request@email.com"]}`,
		`{"friends":["request@email.com","target@email.com","unknown@email.com"]}`}
	for _, request := range invalidRequests {

		var jsonStr = []byte(request)

		relationshipServiceMock := services.RelationshipServiceMock{}
		userServiceMock := services.UserServiceMock{}

		relationshipEndpoint := endpoints.RelationshipEndpoint{relationshipServiceMock, userServiceMock}
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/friends/common-friends", bytes.NewBuffer(jsonStr))
		c.Request.Header.Set("Content-Type", "application/json")

		relationshipEndpoint.CommonFriendList(c)

		assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)

		var actualResult models.Failure
		body, _ := ioutil.ReadAll(w.Result().Body)
		json.Unmarshal(body, &actualResult)

		assert.Equal(t, false, actualResult.Success)
		assert.Equal(t, "Invalid request: incorrect info", actualResult.Message)
	}
}

func TestCommonFriendListWithNotFoundAccounts(t *testing.T) {
	var invalidRequests = []string{
		`{"friends":["undefined@request.com","email@target.com"]}`,
		`{"friends":["email@request.com","undefined@target.com"]}`}

	for i, request := range invalidRequests {

		var jsonStr = []byte(request)

		friendCheckObj := models.FriendCheck{}
		json.Unmarshal(jsonStr, &friendCheckObj)

		var undefinedEmail = friendCheckObj.Friends[0]
		if i == 1 {
			undefinedEmail = friendCheckObj.Friends[1]
		}

		relationshipServiceMock := services.RelationshipServiceMock{}
		userServiceMock := services.UserServiceMock{}

		if i == 0 {
			userRepositoryMock := data.UserRepositoryMock{}
			userRepositoryMock.On("CheckUserExist", undefinedEmail).Return(int64(-1))

			userServiceMock.On("CheckUserExist", undefinedEmail).Return(int64(-1))
		} else {
			userRepositoryMock := data.UserRepositoryMock{}
			userRepositoryMock.On("CheckUserExist", friendCheckObj.Friends[0]).Return(int64(1))
			userRepositoryMock.On("CheckUserExist", undefinedEmail).Return(int64(-1))

			userServiceMock.On("CheckUserExist", friendCheckObj.Friends[0]).Return(int64(1))
			userServiceMock.On("CheckUserExist", undefinedEmail).Return(int64(-1))
		}

		relationshipEndpoint := endpoints.RelationshipEndpoint{relationshipServiceMock, userServiceMock}
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/friends/common-friends", bytes.NewBuffer(jsonStr))
		c.Request.Header.Set("Content-Type", "application/json")

		relationshipEndpoint.CommonFriendList(c)

		assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)

		var actualResult models.Failure
		body, _ := ioutil.ReadAll(w.Result().Body)
		json.Unmarshal(body, &actualResult)

		assert.Equal(t, false, actualResult.Success)
		assert.Equal(t, fmt.Sprintf("Invalid request: User name %s is not found", undefinedEmail), actualResult.Message)
	}
}

func TestCommonFriendListWhichFriendsReturn(t *testing.T) {
	var jsonStr = []byte(`{"friends":["email@request.com","email@target.com"]}`)

	friendCheckObj := models.FriendCheck{}
	json.Unmarshal(jsonStr, &friendCheckObj)

	relationshipRepositoryMock := data.RelationshipRepositoryMock{}
	userRepositoryMock := data.UserRepositoryMock{}
	relationshipServiceMock := services.RelationshipServiceMock{}
	userServiceMock := services.UserServiceMock{}

	requestUser := friendCheckObj.Friends[0]
	requestUserId := int64(1)
	targetUser := friendCheckObj.Friends[1]
	targetUserId := int64(2)

	userRepositoryMock.On("CheckUserExist", requestUser).Return(requestUserId)
	userServiceMock.On("CheckUserExist", requestUser).Return(requestUserId)

	userRepositoryMock.On("CheckUserExist", targetUser).Return(targetUserId)
	userServiceMock.On("CheckUserExist", targetUser).Return(targetUserId)

	friendList := []string{"user1@email.com", "user2@email.com"}
	relationshipServiceMock.On("GetCommonFriendList", requestUserId, targetUserId).Return(friendList)
	relationshipRepositoryMock.On("GetCommonFriendList", requestUserId, targetUserId).Return(friendList)

	relationshipEndpoint := endpoints.RelationshipEndpoint{relationshipServiceMock, userServiceMock}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/friends/common-friends", bytes.NewBuffer(jsonStr))
	c.Request.Header.Set("Content-Type", "application/json")

	relationshipEndpoint.CommonFriendList(c)

	assert.Equal(t, w.Result().StatusCode, http.StatusOK)

	var actualResult models.Friend
	body, _ := ioutil.ReadAll(w.Result().Body)
	json.Unmarshal(body, &actualResult)

	assert.Equal(t, true, actualResult.Success)
	assert.Equal(t, len(friendList), actualResult.Count)
	assert.Equal(t, friendList, actualResult.Friends)
}

func TestSubcribeWithInvalidAccounts(t *testing.T) {
	var invalidRequests = []string{
		`{"requestor":"invalid_model}`,
		`{"requestor":"target@email.com"}`,
		`{"requestor":"request","target":"target@email.com"}`,
		`{"requestor":"target@email.com","target":"request"}`,
		`{"requestor":"","target":""}`}
	for _, request := range invalidRequests {

		var jsonStr = []byte(request)

		relationshipServiceMock := services.RelationshipServiceMock{}
		userServiceMock := services.UserServiceMock{}

		relationshipEndpoint := endpoints.RelationshipEndpoint{relationshipServiceMock, userServiceMock}
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/friends/subcribe", bytes.NewBuffer(jsonStr))
		c.Request.Header.Set("Content-Type", "application/json")

		relationshipEndpoint.Subscribe(c)

		assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)

		var actualResult models.Failure
		body, _ := ioutil.ReadAll(w.Result().Body)
		json.Unmarshal(body, &actualResult)

		assert.Equal(t, false, actualResult.Success)
		assert.Equal(t, "Invalid request: incorrect info", actualResult.Message)
	}
}

func TestSubcribeUpdateWithNotFoundAccounts(t *testing.T) {
	var invalidRequests = []string{
		`{"requestor":"request@notfound.com","target":"target@email.com"}`,
		`{"requestor":"request@email.com","target":"target@notfound.com"}`}

	for i, request := range invalidRequests {

		var jsonStr = []byte(request)

		userActionObj := models.UserAction{}
		json.Unmarshal(jsonStr, &userActionObj)

		var undefinedEmail = userActionObj.Requestor
		if i == 1 {
			undefinedEmail = userActionObj.Target
		}

		relationshipServiceMock := services.RelationshipServiceMock{}
		userServiceMock := services.UserServiceMock{}

		if i == 0 {
			userRepositoryMock := data.UserRepositoryMock{}
			userRepositoryMock.On("CheckUserExist", undefinedEmail).Return(int64(-1))

			userServiceMock.On("CheckUserExist", undefinedEmail).Return(int64(-1))
		} else {
			userRepositoryMock := data.UserRepositoryMock{}
			userRepositoryMock.On("CheckUserExist", userActionObj.Requestor).Return(int64(1))
			userRepositoryMock.On("CheckUserExist", undefinedEmail).Return(int64(-1))

			userServiceMock.On("CheckUserExist", userActionObj.Requestor).Return(int64(1))
			userServiceMock.On("CheckUserExist", undefinedEmail).Return(int64(-1))
		}

		relationshipEndpoint := endpoints.RelationshipEndpoint{relationshipServiceMock, userServiceMock}
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/friends/subcribe", bytes.NewBuffer(jsonStr))
		c.Request.Header.Set("Content-Type", "application/json")

		relationshipEndpoint.Subscribe(c)

		assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)

		var actualResult models.Failure
		body, _ := ioutil.ReadAll(w.Result().Body)
		json.Unmarshal(body, &actualResult)

		assert.Equal(t, false, actualResult.Success)
		assert.Equal(t, fmt.Sprintf("Invalid request: User name %s is not found", undefinedEmail), actualResult.Message)
	}
}

func TestSubcribeUpdateWithAlreadySubcribedAccounts(t *testing.T) {
	var jsonStr = []byte(`{"requestor":"email@request.com","target":"email@target.com"}`)

	userActionObj := models.UserAction{}
	json.Unmarshal(jsonStr, &userActionObj)

	relationshipRepositoryMock := data.RelationshipRepositoryMock{}
	userRepositoryMock := data.UserRepositoryMock{}
	relationshipServiceMock := services.RelationshipServiceMock{}
	userServiceMock := services.UserServiceMock{}

	requestUser := userActionObj.Requestor
	requestUserId := int64(1)
	targetUser := userActionObj.Target
	targetUserId := int64(2)
	subcribedStatus := int64(2)
	subcribedIds := []int64{int64(1)}

	userRepositoryMock.On("CheckUserExist", requestUser).Return(requestUserId)
	userServiceMock.On("CheckUserExist", requestUser).Return(requestUserId)

	userRepositoryMock.On("CheckUserExist", targetUser).Return(targetUserId)
	userServiceMock.On("CheckUserExist", targetUser).Return(targetUserId)

	relationshipRepositoryMock.On("CheckRelationshipOneWay", requestUserId, targetUserId, subcribedStatus).Return(subcribedIds)
	relationshipServiceMock.On("CheckPartialSubcribed", requestUserId, targetUserId).Return(subcribedIds)

	relationshipEndpoint := endpoints.RelationshipEndpoint{relationshipServiceMock, userServiceMock}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/friends/subcribe", bytes.NewBuffer(jsonStr))
	c.Request.Header.Set("Content-Type", "application/json")

	relationshipEndpoint.Subscribe(c)

	assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)

	var actualResult models.Failure
	body, _ := ioutil.ReadAll(w.Result().Body)
	json.Unmarshal(body, &actualResult)

	assert.Equal(t, false, actualResult.Success)
	assert.Equal(t, "Invalid request: subcribed status is existed", actualResult.Message)
}

func TestSubcribeUpdateWithAlreadyBlockedAccounts(t *testing.T) {
	var jsonStr = []byte(`{"requestor":"email@request.com","target":"email@target.com"}`)

	userActionObj := models.UserAction{}
	json.Unmarshal(jsonStr, &userActionObj)

	relationshipRepositoryMock := data.RelationshipRepositoryMock{}
	userRepositoryMock := data.UserRepositoryMock{}
	relationshipServiceMock := services.RelationshipServiceMock{}
	userServiceMock := services.UserServiceMock{}

	requestUser := userActionObj.Requestor
	requestUserId := int64(1)
	targetUser := userActionObj.Target
	targetUserId := int64(2)
	//connectedStatus := int64(1)
	blockedStatus := int64(3)
	subcribedStatus := int64(2)
	//connectedIds := []int64{}
	blockedIds := []int64{int64(1)}
	subcribedIds := []int64{}

	userRepositoryMock.On("CheckUserExist", requestUser).Return(requestUserId)
	userServiceMock.On("CheckUserExist", requestUser).Return(requestUserId)

	userRepositoryMock.On("CheckUserExist", targetUser).Return(targetUserId)
	userServiceMock.On("CheckUserExist", targetUser).Return(targetUserId)

	relationshipRepositoryMock.On("CheckRelationshipOneWay", requestUserId, targetUserId, subcribedStatus).Return(subcribedIds)
	relationshipServiceMock.On("CheckPartialSubcribed", requestUserId, targetUserId).Return(subcribedIds)

	relationshipRepositoryMock.On("CheckRelationshipOneWay", requestUserId, targetUserId, blockedStatus).Return(blockedIds)
	relationshipServiceMock.On("CheckPartialBlocked", requestUserId, targetUserId).Return(blockedIds)

	relationshipEndpoint := endpoints.RelationshipEndpoint{relationshipServiceMock, userServiceMock}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/friends/subcribe", bytes.NewBuffer(jsonStr))
	c.Request.Header.Set("Content-Type", "application/json")

	relationshipEndpoint.Subscribe(c)

	assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)

	var actualResult models.Failure
	body, _ := ioutil.ReadAll(w.Result().Body)
	json.Unmarshal(body, &actualResult)

	assert.Equal(t, false, actualResult.Success)
	assert.Equal(t, "Invalid request: blocked status is existed", actualResult.Message)
}

func TestSubcribeUpdateWithAlreadyConnectedAccounts(t *testing.T) {
	var jsonStr = []byte(`{"requestor":"email@request.com","target":"email@target.com"}`)

	userActionObj := models.UserAction{}
	json.Unmarshal(jsonStr, &userActionObj)

	relationshipRepositoryMock := data.RelationshipRepositoryMock{}
	userRepositoryMock := data.UserRepositoryMock{}
	relationshipServiceMock := services.RelationshipServiceMock{}
	userServiceMock := services.UserServiceMock{}

	requestUser := userActionObj.Requestor
	requestUserId := int64(1)
	targetUser := userActionObj.Target
	targetUserId := int64(2)
	connectedStatus := int64(1)
	blockedStatus := int64(3)
	subcribedStatus := int64(2)
	connectedIds := []int64{int64(1)}
	blockedIds := []int64{}
	subcribedIds := []int64{}

	userRepositoryMock.On("CheckUserExist", requestUser).Return(requestUserId)
	userServiceMock.On("CheckUserExist", requestUser).Return(requestUserId)

	userRepositoryMock.On("CheckUserExist", targetUser).Return(targetUserId)
	userServiceMock.On("CheckUserExist", targetUser).Return(targetUserId)

	relationshipRepositoryMock.On("CheckRelationshipOneWay", requestUserId, targetUserId, subcribedStatus).Return(subcribedIds)
	relationshipServiceMock.On("CheckPartialSubcribed", requestUserId, targetUserId).Return(subcribedIds)

	relationshipRepositoryMock.On("CheckRelationshipOnetWay", requestUserId, targetUserId, blockedStatus).Return(blockedIds)
	relationshipServiceMock.On("CheckPartialBlocked", requestUserId, targetUserId).Return(blockedIds)

	relationshipRepositoryMock.On("CheckRelationshipTwoWay", requestUserId, targetUserId, connectedStatus).Return(connectedIds)
	relationshipServiceMock.On("CheckConnected", requestUserId, targetUserId).Return(connectedIds)

	relationshipEndpoint := endpoints.RelationshipEndpoint{relationshipServiceMock, userServiceMock}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/friends/subcribe", bytes.NewBuffer(jsonStr))
	c.Request.Header.Set("Content-Type", "application/json")

	relationshipEndpoint.Subscribe(c)

	assert.Equal(t, w.Result().StatusCode, http.StatusOK)

	var actualResult models.Success
	body, _ := ioutil.ReadAll(w.Result().Body)
	json.Unmarshal(body, &actualResult)

	assert.Equal(t, true, actualResult.Success)
}

func TestSubcribeUpdateWhichOkResult(t *testing.T) {
	var jsonStr = []byte(`{"requestor":"email@request.com","target":"email@target.com"}`)

	userActionObj := models.UserAction{}
	json.Unmarshal(jsonStr, &userActionObj)

	relationshipRepositoryMock := data.RelationshipRepositoryMock{}
	userRepositoryMock := data.UserRepositoryMock{}
	relationshipServiceMock := services.RelationshipServiceMock{}
	userServiceMock := services.UserServiceMock{}

	requestUser := userActionObj.Requestor
	requestUserId := int64(1)
	targetUser := userActionObj.Target
	targetUserId := int64(2)
	connectedStatus := int64(1)
	blockedStatus := int64(3)
	subcribedStatus := int64(2)
	connectedIds := []int64{}
	blockedIds := []int64{}
	subcribedIds := []int64{}

	userRepositoryMock.On("CheckUserExist", requestUser).Return(requestUserId)
	userServiceMock.On("CheckUserExist", requestUser).Return(requestUserId)

	userRepositoryMock.On("CheckUserExist", targetUser).Return(targetUserId)
	userServiceMock.On("CheckUserExist", targetUser).Return(targetUserId)

	relationshipRepositoryMock.On("CheckRelationshipOneWay", requestUserId, targetUserId, subcribedStatus).Return(subcribedIds)
	relationshipServiceMock.On("CheckPartialSubcribed", requestUserId, targetUserId).Return(subcribedIds)

	relationshipRepositoryMock.On("CheckRelationshipOnetWay", requestUserId, targetUserId, blockedStatus).Return(blockedIds)
	relationshipServiceMock.On("CheckPartialBlocked", requestUserId, targetUserId).Return(blockedIds)

	relationshipRepositoryMock.On("CheckRelationshipTwoWay", requestUserId, targetUserId, connectedStatus).Return(connectedIds)
	relationshipServiceMock.On("CheckConnected", requestUserId, targetUserId).Return(connectedIds)

	relationshipModel := models.Relationship{Status: subcribedStatus, RequestUserId: requestUserId, TargetUserId: targetUserId}
	relationshipRepositoryMock.On("CreateRelationship", &relationshipModel).Return(int64(10))
	relationshipServiceMock.On("CreateRelationship", &relationshipModel).Return(int64(10))

	relationshipEndpoint := endpoints.RelationshipEndpoint{relationshipServiceMock, userServiceMock}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/friends/subcribe", bytes.NewBuffer(jsonStr))
	c.Request.Header.Set("Content-Type", "application/json")

	relationshipEndpoint.Subscribe(c)

	assert.Equal(t, w.Result().StatusCode, http.StatusOK)

	var actualResult models.Success
	body, _ := ioutil.ReadAll(w.Result().Body)
	json.Unmarshal(body, &actualResult)

	assert.Equal(t, true, actualResult.Success)
}

func TestBlockUpdateWithInvalidAccounts(t *testing.T) {
	var invalidRequests = []string{
		`{"requestor":"invalid_model`,
		`{"requestor":"target@email.com"}`,
		`{"requestor":"request","target":"target@email.com"}`,
		`{"requestor":"target@email.com","target":"request"}`,
		`{"requestor":"","target":""}`}
	for _, request := range invalidRequests {

		var jsonStr = []byte(request)

		relationshipServiceMock := services.RelationshipServiceMock{}
		userServiceMock := services.UserServiceMock{}

		relationshipEndpoint := endpoints.RelationshipEndpoint{relationshipServiceMock, userServiceMock}
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/friends/block", bytes.NewBuffer(jsonStr))
		c.Request.Header.Set("Content-Type", "application/json")

		relationshipEndpoint.Block(c)

		assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)

		var actualResult models.Failure
		body, _ := ioutil.ReadAll(w.Result().Body)
		json.Unmarshal(body, &actualResult)

		assert.Equal(t, false, actualResult.Success)
		assert.Equal(t, "Invalid request: incorrect info", actualResult.Message)
	}
}

func TestBlockUpdateWithNotFoundAccounts(t *testing.T) {
	var invalidRequests = []string{
		`{"requestor":"request@notfound.com","target":"target@email.com"}`,
		`{"requestor":"request@email.com","target":"target@notfound.com"}`}

	for i, request := range invalidRequests {

		var jsonStr = []byte(request)

		userActionObj := models.UserAction{}
		json.Unmarshal(jsonStr, &userActionObj)

		var undefinedEmail = userActionObj.Requestor
		if i == 1 {
			undefinedEmail = userActionObj.Target
		}

		relationshipServiceMock := services.RelationshipServiceMock{}
		userServiceMock := services.UserServiceMock{}

		if i == 0 {
			userRepositoryMock := data.UserRepositoryMock{}
			userRepositoryMock.On("CheckUserExist", undefinedEmail).Return(int64(-1))

			userServiceMock.On("CheckUserExist", undefinedEmail).Return(int64(-1))
		} else {
			userRepositoryMock := data.UserRepositoryMock{}
			userRepositoryMock.On("CheckUserExist", userActionObj.Requestor).Return(int64(1))
			userRepositoryMock.On("CheckUserExist", undefinedEmail).Return(int64(-1))

			userServiceMock.On("CheckUserExist", userActionObj.Requestor).Return(int64(1))
			userServiceMock.On("CheckUserExist", undefinedEmail).Return(int64(-1))
		}

		relationshipEndpoint := endpoints.RelationshipEndpoint{relationshipServiceMock, userServiceMock}
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/friends/block", bytes.NewBuffer(jsonStr))
		c.Request.Header.Set("Content-Type", "application/json")

		relationshipEndpoint.Block(c)

		assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)

		var actualResult models.Failure
		body, _ := ioutil.ReadAll(w.Result().Body)
		json.Unmarshal(body, &actualResult)

		assert.Equal(t, false, actualResult.Success)
		assert.Equal(t, fmt.Sprintf("Invalid request: User name %s is not found", undefinedEmail), actualResult.Message)
	}
}

func TestBlockUpdateWithAlreadyBlockedAccounts(t *testing.T) {
	var jsonStr = []byte(`{"requestor":"email@request.com","target":"email@target.com"}`)

	userActionObj := models.UserAction{}
	json.Unmarshal(jsonStr, &userActionObj)

	relationshipRepositoryMock := data.RelationshipRepositoryMock{}
	userRepositoryMock := data.UserRepositoryMock{}
	relationshipServiceMock := services.RelationshipServiceMock{}
	userServiceMock := services.UserServiceMock{}

	requestUser := userActionObj.Requestor
	requestUserId := int64(1)
	targetUser := userActionObj.Target
	targetUserId := int64(2)
	blockedStatus := int64(3)
	blockedIds := []int64{int64(1)}

	userRepositoryMock.On("CheckUserExist", requestUser).Return(requestUserId)
	userServiceMock.On("CheckUserExist", requestUser).Return(requestUserId)

	userRepositoryMock.On("CheckUserExist", targetUser).Return(targetUserId)
	userServiceMock.On("CheckUserExist", targetUser).Return(targetUserId)

	relationshipRepositoryMock.On("CheckRelationshipOneWay", requestUserId, targetUserId, blockedStatus).Return(blockedIds)
	relationshipServiceMock.On("CheckPartialBlocked", requestUserId, targetUserId).Return(blockedIds)

	relationshipEndpoint := endpoints.RelationshipEndpoint{relationshipServiceMock, userServiceMock}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/friends/block", bytes.NewBuffer(jsonStr))
	c.Request.Header.Set("Content-Type", "application/json")

	relationshipEndpoint.Block(c)

	assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)

	var actualResult models.Failure
	body, _ := ioutil.ReadAll(w.Result().Body)
	json.Unmarshal(body, &actualResult)

	assert.Equal(t, false, actualResult.Success)
	assert.Equal(t, "Invalid request: blocked status is existed", actualResult.Message)
}

func TestBlockUpdateWithAlreadySubcribedAccounts(t *testing.T) {
	var jsonStr = []byte(`{"requestor":"email@request.com","target":"email@target.com"}`)

	userActionObj := models.UserAction{}
	json.Unmarshal(jsonStr, &userActionObj)

	relationshipRepositoryMock := data.RelationshipRepositoryMock{}
	userRepositoryMock := data.UserRepositoryMock{}
	relationshipServiceMock := services.RelationshipServiceMock{}
	userServiceMock := services.UserServiceMock{}

	requestUser := userActionObj.Requestor
	requestUserId := int64(1)
	targetUser := userActionObj.Target
	targetUserId := int64(2)
	connectedStatus := int64(1)
	blockedStatus := int64(3)
	subcribedStatus := int64(2)
	connectedIds := []int64{}
	blockedIds := []int64{}
	subcribedIds := []int64{int64(1)}

	userRepositoryMock.On("CheckUserExist", requestUser).Return(requestUserId)
	userServiceMock.On("CheckUserExist", requestUser).Return(requestUserId)

	userRepositoryMock.On("CheckUserExist", targetUser).Return(targetUserId)
	userServiceMock.On("CheckUserExist", targetUser).Return(targetUserId)

	relationshipRepositoryMock.On("CheckRelationshipOneWay", requestUserId, targetUserId, blockedStatus).Return(blockedIds)
	relationshipServiceMock.On("CheckPartialBlocked", requestUserId, targetUserId).Return(blockedIds)

	relationshipRepositoryMock.On("CheckRelationshipOneWay", requestUserId, targetUserId, subcribedStatus).Return(subcribedIds)
	relationshipServiceMock.On("CheckPartialSubcribed", requestUserId, targetUserId).Return(subcribedIds)

	relationshipRepositoryMock.On("DeleteRelationships", subcribedIds).Return(true)
	relationshipServiceMock.On("DeleteRelationships", subcribedIds).Return(true)

	relationshipRepositoryMock.On("CheckRelationshipTwoWay", requestUserId, targetUserId, connectedStatus).Return(connectedIds)
	relationshipServiceMock.On("CheckConnected", requestUserId, targetUserId).Return(connectedIds)

	relationshipModel := models.Relationship{Status: blockedStatus, RequestUserId: requestUserId, TargetUserId: targetUserId}
	relationshipRepositoryMock.On("CreateRelationship", &relationshipModel).Return(int64(10))
	relationshipServiceMock.On("CreateRelationship", &relationshipModel).Return(int64(10))

	relationshipEndpoint := endpoints.RelationshipEndpoint{relationshipServiceMock, userServiceMock}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/friends/subcribe", bytes.NewBuffer(jsonStr))
	c.Request.Header.Set("Content-Type", "application/json")

	relationshipEndpoint.Block(c)

	assert.Equal(t, w.Result().StatusCode, http.StatusOK)

	var actualResult models.Success
	body, _ := ioutil.ReadAll(w.Result().Body)
	json.Unmarshal(body, &actualResult)

	assert.Equal(t, true, actualResult.Success)
}

func TestBlockUpdateWithAlreadyConnectedAccounts(t *testing.T) {
	var jsonStr = []byte(`{"requestor":"email@request.com","target":"email@target.com"}`)

	userActionObj := models.UserAction{}
	json.Unmarshal(jsonStr, &userActionObj)

	relationshipRepositoryMock := data.RelationshipRepositoryMock{}
	userRepositoryMock := data.UserRepositoryMock{}
	relationshipServiceMock := services.RelationshipServiceMock{}
	userServiceMock := services.UserServiceMock{}

	requestUser := userActionObj.Requestor
	requestUserId := int64(1)
	targetUser := userActionObj.Target
	targetUserId := int64(2)
	connectedStatus := int64(1)
	blockedStatus := int64(3)
	subcribedStatus := int64(2)
	connectedIds := []int64{int64(1)}
	blockedIds := []int64{}
	subcribedIds := []int64{}

	userRepositoryMock.On("CheckUserExist", requestUser).Return(requestUserId)
	userServiceMock.On("CheckUserExist", requestUser).Return(requestUserId)

	userRepositoryMock.On("CheckUserExist", targetUser).Return(targetUserId)
	userServiceMock.On("CheckUserExist", targetUser).Return(targetUserId)

	relationshipRepositoryMock.On("CheckRelationshipOneWay", requestUserId, targetUserId, blockedStatus).Return(blockedIds)
	relationshipServiceMock.On("CheckPartialBlocked", requestUserId, targetUserId).Return(blockedIds)

	relationshipRepositoryMock.On("CheckRelationshipOneWay", requestUserId, targetUserId, subcribedStatus).Return(subcribedIds)
	relationshipServiceMock.On("CheckPartialSubcribed", requestUserId, targetUserId).Return(subcribedIds)

	relationshipRepositoryMock.On("CheckRelationshipTwoWay", requestUserId, targetUserId, connectedStatus).Return(connectedIds)
	relationshipServiceMock.On("CheckConnected", requestUserId, targetUserId).Return(connectedIds)

	relationshipRepositoryMock.On("DeleteRelationships", connectedIds).Return(true)
	relationshipServiceMock.On("DeleteRelationships", connectedIds).Return(true)

	relationshipModel := models.Relationship{Status: blockedStatus, RequestUserId: requestUserId, TargetUserId: targetUserId}
	relationshipRepositoryMock.On("CreateRelationship", &relationshipModel).Return(int64(10))
	relationshipServiceMock.On("CreateRelationship", &relationshipModel).Return(int64(10))

	relationshipEndpoint := endpoints.RelationshipEndpoint{relationshipServiceMock, userServiceMock}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/friends/block", bytes.NewBuffer(jsonStr))
	c.Request.Header.Set("Content-Type", "application/json")

	relationshipEndpoint.Block(c)

	assert.Equal(t, w.Result().StatusCode, http.StatusOK)

	var actualResult models.Success
	body, _ := ioutil.ReadAll(w.Result().Body)
	json.Unmarshal(body, &actualResult)

	assert.Equal(t, true, actualResult.Success)
}

func TestReceiveUpdateWithInvalidAccount(t *testing.T) {
	var invalidRequests = []string{
		`{"sender":"invalid_model`,
		`{"sender":"sender@email.com"}`,
		`{"sender":"sender@email.com","text":""}`,
		`{"sender":"","text":"hello world!!!"}`,
		`{"sender":"invalid_email","text":"hello world!!!"}`}
	for _, request := range invalidRequests {

		var jsonStr = []byte(request)

		relationshipServiceMock := services.RelationshipServiceMock{}
		userServiceMock := services.UserServiceMock{}

		relationshipEndpoint := endpoints.RelationshipEndpoint{relationshipServiceMock, userServiceMock}
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/friends/receive-updates", bytes.NewBuffer(jsonStr))
		c.Request.Header.Set("Content-Type", "application/json")

		relationshipEndpoint.ReceiveUpdates(c)

		assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)

		var actualResult models.Failure
		body, _ := ioutil.ReadAll(w.Result().Body)
		json.Unmarshal(body, &actualResult)

		assert.Equal(t, false, actualResult.Success)
		assert.Equal(t, "Invalid request: incorrect info", actualResult.Message)
	}
}

func TestReceiveUpdateWithNotFoundAccount(t *testing.T) {
	var jsonStr = []byte(`{"sender":"request@notfound.com","text":"hello world!!!"}`)

	userPostObj := models.UserPost{}
	json.Unmarshal(jsonStr, &userPostObj)

	relationshipServiceMock := services.RelationshipServiceMock{}
	userServiceMock := services.UserServiceMock{}

	userRepositoryMock := data.UserRepositoryMock{}
	userRepositoryMock.On("CheckUserExist", userPostObj.Sender).Return(int64(-1))

	userServiceMock.On("CheckUserExist", userPostObj.Sender).Return(int64(-1))

	relationshipEndpoint := endpoints.RelationshipEndpoint{relationshipServiceMock, userServiceMock}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/friends/receive-updates", bytes.NewBuffer(jsonStr))
	c.Request.Header.Set("Content-Type", "application/json")

	relationshipEndpoint.ReceiveUpdates(c)

	assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)

	var actualResult models.Failure
	body, _ := ioutil.ReadAll(w.Result().Body)
	json.Unmarshal(body, &actualResult)

	assert.Equal(t, false, actualResult.Success)
	assert.Equal(t, fmt.Sprintf("Invalid request: User name %s is not found", userPostObj.Sender), actualResult.Message)
}

func TestReceiveUpdateReturnOk(t *testing.T) {
	var jsonStr = []byte(`{"sender":"sender@email.com","text":"hello world sender@email.com johndoe@gmail.com"}`)

	userPostObj := models.UserPost{}
	json.Unmarshal(jsonStr, &userPostObj)

	relationshipServiceMock := services.RelationshipServiceMock{}
	relationshipRepositoryMock := data.RelationshipRepositoryMock{}
	userServiceMock := services.UserServiceMock{}

	userRepositoryMock := data.UserRepositoryMock{}
	userRepositoryMock.On("CheckUserExist", userPostObj.Sender).Return(int64(1))
	userServiceMock.On("CheckUserExist", userPostObj.Sender).Return(int64(1))

	existedIds := []int64{int64(10)}
	userRepositoryMock.On("CheckUsersExist", []string{"johndoe@gmail.com"}).Return(existedIds)
	userServiceMock.On("CheckUsersExist", []string{"johndoe@gmail.com"}).Return(existedIds)

	senderId := int64(1)
	mentionedIds := []int64{int64(10)}
	receiveUpdateEmails := []string{"user1@email.com", "user2@email.com"}

	relationshipRepositoryMock.On("GetValidUsersCanReceiveUpdates", senderId, mentionedIds).Return(receiveUpdateEmails)
	relationshipServiceMock.On("GetValidUsersCanReceiveUpdates", senderId, mentionedIds).Return(receiveUpdateEmails)

	relationshipEndpoint := endpoints.RelationshipEndpoint{relationshipServiceMock, userServiceMock}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/friends/receive-updates", bytes.NewBuffer(jsonStr))
	c.Request.Header.Set("Content-Type", "application/json")

	relationshipEndpoint.ReceiveUpdates(c)

	assert.Equal(t, w.Result().StatusCode, http.StatusOK)

	var actualResult models.Success
	body, _ := ioutil.ReadAll(w.Result().Body)
	json.Unmarshal(body, &actualResult)

	assert.Equal(t, true, actualResult.Success)
}
