package endpoints

import (
	"fmt"
	"friendMgmt/common"
	"friendMgmt/models"
	"friendMgmt/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mcnijman/go-emailaddress"
)

type RelationshipEndpoint struct {
	IRelationshipService services.IRelationshipService
	IUserService         services.IUserService
}

// CreateRelationship godoc
// @Tags Friend
// @Summary API to create a friend connection between two users
// @Accept  json
// @Produce  json
// @Param model body models.FriendCheck true "Body"
// @Success 200 {object} models.Success "OK"
// @Failure 400 {object} models.Failure "Bad Request"
// @Failure 500 {object} models.Failure "Internal Error"
// @Router /friends/add [post]
func (r RelationshipEndpoint) CreateRelationship(c *gin.Context) {
	var friendCheck models.FriendCheck
	if err := c.BindJSON(&friendCheck); err != nil {
		responseError(c, http.StatusBadRequest, "Invalid request: incorrect info")
		return
	}

	if len(friendCheck.Friends) != 2 {
		responseError(c, http.StatusBadRequest, "Invalid request: incorrect info")
		return
	}

	var requestUser = friendCheck.Friends[0]
	var targetUser = friendCheck.Friends[1]

	if !common.IsValidEmail(requestUser) || !common.IsValidEmail(targetUser) || requestUser == targetUser {
		responseError(c, http.StatusBadRequest, "Invalid request: incorrect info")
		return
	}

	var requestUserId = r.IUserService.CheckUserExist(requestUser)
	if requestUserId <= 0 {
		responseError(c, http.StatusBadRequest, fmt.Sprintf("Invalid request: User name %s is not found", requestUser))
		return

	}

	var targetUserId = r.IUserService.CheckUserExist(targetUser)
	if targetUserId <= 0 {
		responseError(c, http.StatusBadRequest, fmt.Sprintf("Invalid request: User name %s is not found", targetUser))
		return
	}

	connectedRelationshipIds := r.IRelationshipService.CheckConnected(requestUserId, targetUserId)
	if len(connectedRelationshipIds) > 0 {
		responseError(c, http.StatusBadRequest, "Invalid request: connected status is existed")
		return
	}

	blockedRelationshipIds := r.IRelationshipService.CheckFullyBlocked(requestUserId, targetUserId)
	if len(blockedRelationshipIds) > 0 {
		responseError(c, http.StatusBadRequest, "Invalid request: blocked status is existed")
		return
	}

	subcribedRelationshipIds := r.IRelationshipService.CheckFullySubcribed(requestUserId, targetUserId)
	if len(subcribedRelationshipIds) > 0 {
		r.IRelationshipService.DeleteRelationships(subcribedRelationshipIds)
	}

	relationshipModel := models.Relationship{Status: 1, RequestUserId: requestUserId, TargetUserId: targetUserId}

	if insertedId := r.IRelationshipService.CreateRelationship(&relationshipModel); insertedId > 0 {
		success := models.Success{Success: true}
		responseOk(c, success)
		return
	}

	responseError(c, http.StatusInternalServerError, "Oops! There is an error, please try again.")
	return
}

// FriendList godoc
// @Tags Friend
// @Summary API to check list friends of an user
// @Accept  json
// @Produce  json
// @Param model body models.Email true "Body"
// @Success 200 {object} models.Success "OK"
// @Failure 400 {object} models.Failure "Bad Request"
// @Router /friends [post]
func (r RelationshipEndpoint) FriendList(c *gin.Context) {
	var email models.Email
	if err := c.BindJSON(&email); err != nil {
		responseError(c, http.StatusBadRequest, "Invalid request: incorrect info")
		return
	}

	if isValid := common.IsValidEmail(email.Email); !isValid {
		responseError(c, http.StatusBadRequest, "Invalid request: incorrect info")
		return
	}

	userId := r.IUserService.CheckUserExist(email.Email)
	if userId < 0 {
		responseError(c, http.StatusBadRequest, fmt.Sprintf("Invalid request: User name %s is not found", email.Email))
		return
	}

	friendList := r.IRelationshipService.GetFriendList(userId)

	friendModel := models.Friend{Friends: friendList, Count: len(friendList), Success: true}

	responseOk(c, friendModel)
}

// CommonFriendList godoc
// @Tags Friend
// @Summary API to check common friends of two users
// @Accept  json
// @Produce  json
// @Param model body models.FriendCheck true "Body"
// @Success 200 {object} models.Success "OK"
// @Failure 400 {object} models.Failure "Bad Request"
// @Router /friends/common-friends [post]
func (r RelationshipEndpoint) CommonFriendList(c *gin.Context) {

	var friendCheck models.FriendCheck
	if err := c.BindJSON(&friendCheck); err != nil {
		responseError(c, http.StatusBadRequest, "Invalid request: incorrect info")
		return
	}

	if len(friendCheck.Friends) != 2 {
		responseError(c, http.StatusBadRequest, "Invalid request: incorrect info")
		return
	}

	var requestUser = friendCheck.Friends[0]
	var targetUser = friendCheck.Friends[1]

	if !common.IsValidEmail(requestUser) || !common.IsValidEmail(targetUser) || requestUser == targetUser {
		responseError(c, http.StatusBadRequest, "Invalid request: incorrect info")
		return
	}

	var requestUserId = r.IUserService.CheckUserExist(requestUser)
	if requestUserId <= 0 {
		responseError(c, http.StatusBadRequest, fmt.Sprintf("Invalid request: User name %s is not found", requestUser))
		return

	}

	var targetUserId = r.IUserService.CheckUserExist(targetUser)
	if targetUserId <= 0 {
		responseError(c, http.StatusBadRequest, fmt.Sprintf("Invalid request: User name %s is not found", targetUser))
		return
	}

	commonFriends := r.IRelationshipService.GetCommonFriendList(requestUserId, targetUserId)

	friendModel := models.Friend{Friends: commonFriends, Count: len(commonFriends), Success: true}

	responseOk(c, friendModel)
}

// Subscribe godoc
// @Tags Friend
// @Summary API to allow an user can subscribe another user
// @Accept  json
// @Produce  json
// @Param model body models.UserAction true "Body"
// @Success 200 {object} models.Success "OK"
// @Failure 400 {object} models.Failure "Bad Request"
// @Router /friends/subcribe [post]
func (r RelationshipEndpoint) Subscribe(c *gin.Context) {
	var userAction models.UserAction

	if err := c.BindJSON(&userAction); err != nil {
		responseError(c, http.StatusBadRequest, "Invalid request: incorrect info")
		return
	}

	var requestUser = userAction.Requestor
	var targetUser = userAction.Target

	if !common.IsValidEmail(requestUser) || !common.IsValidEmail(targetUser) || requestUser == targetUser {
		responseError(c, http.StatusBadRequest, "Invalid request: incorrect info")
		return
	}

	var requestUserId = r.IUserService.CheckUserExist(requestUser)
	if requestUserId <= 0 {
		responseError(c, http.StatusBadRequest, fmt.Sprintf("Invalid request: User name %s is not found", requestUser))
		return
	}

	var targetUserId = r.IUserService.CheckUserExist(targetUser)
	if targetUserId <= 0 {
		responseError(c, http.StatusBadRequest, fmt.Sprintf("Invalid request: User name %s is not found", targetUser))
		return
	}

	if subcribedRelationshipId := r.IRelationshipService.CheckPartialSubcribed(requestUserId, targetUserId); len(subcribedRelationshipId) > 0 {
		responseError(c, http.StatusBadRequest, "Invalid request: subcribed status is existed")
		return
	}

	if blockedRelationshipId := r.IRelationshipService.CheckPartialBlocked(requestUserId, targetUserId); len(blockedRelationshipId) > 0 {
		responseError(c, http.StatusBadRequest, "Invalid request: blocked status is existed")
		return
	}

	if connectedRelationshipIds := r.IRelationshipService.CheckConnected(requestUserId, targetUserId); len(connectedRelationshipIds) > 0 {
		success := models.Success{Success: true}
		responseOk(c, success)
		return
	}

	relationshipModel := models.Relationship{Status: 2, RequestUserId: requestUserId, TargetUserId: targetUserId}

	r.IRelationshipService.CreateRelationship(&relationshipModel)

	success := models.Success{Success: true}
	responseOk(c, success)
}

// Block godoc
// @Tags Friend
// @Summary API to allow an user can block another user
// @Accept  json
// @Produce  json
// @Param model body models.UserAction true "Body"
// @Success 200 {object} models.Success "OK"
// @Failure 400 {object} models.Failure "Bad Request"
// @Router /friends/block [post]
func (r RelationshipEndpoint) Block(c *gin.Context) {
	var userAction models.UserAction

	if err := c.BindJSON(&userAction); err != nil {
		responseError(c, http.StatusBadRequest, "Invalid request: incorrect info")
		return
	}

	var requestUser = userAction.Requestor
	var targetUser = userAction.Target

	if !common.IsValidEmail(requestUser) || !common.IsValidEmail(targetUser) || requestUser == targetUser {
		responseError(c, http.StatusBadRequest, "Invalid request: incorrect info")
		return
	}

	var requestUserId = r.IUserService.CheckUserExist(requestUser)
	if requestUserId <= 0 {
		responseError(c, http.StatusBadRequest, fmt.Sprintf("Invalid request: User name %s is not found", requestUser))
		return
	}

	var targetUserId = r.IUserService.CheckUserExist(targetUser)
	if targetUserId <= 0 {
		responseError(c, http.StatusBadRequest, fmt.Sprintf("Invalid request: User name %s is not found", targetUser))
		return
	}

	if blockedRelationshipId := r.IRelationshipService.CheckPartialBlocked(requestUserId, targetUserId); len(blockedRelationshipId) > 0 {
		responseError(c, http.StatusBadRequest, "Invalid request: blocked status is existed")
		return
	}

	if subcribedRelationshipId := r.IRelationshipService.CheckPartialSubcribed(requestUserId, targetUserId); len(subcribedRelationshipId) > 0 {
		r.IRelationshipService.DeleteRelationships(subcribedRelationshipId)
	}

	if connectedRelationshipId := r.IRelationshipService.CheckConnected(requestUserId, targetUserId); len(connectedRelationshipId) > 0 {
		r.IRelationshipService.DeleteRelationships(connectedRelationshipId)
	}

	relationshipModel := models.Relationship{Status: 3, RequestUserId: requestUserId, TargetUserId: targetUserId}

	r.IRelationshipService.CreateRelationship(&relationshipModel)

	success := models.Success{Success: true}
	responseOk(c, success)
}

// ReceiveUpdates godoc
// @Tags Friend
// @Summary API to return list of users can receive update from an user
// @Accept  json
// @Produce  json
// @Param model body models.UserPost true "Body"
// @Success 200 {object} models.Success "OK"
// @Failure 400 {object} models.Failure "Bad Request"
// @Router /friends/receive-updates [post]
func (r RelationshipEndpoint) ReceiveUpdates(c *gin.Context) {
	var userPost models.UserPost

	if err := c.BindJSON(&userPost); err != nil {
		responseError(c, http.StatusBadRequest, "Invalid request: incorrect info")
		return
	}

	var sender = userPost.Sender
	var text = userPost.Text

	if !common.IsValidEmail(sender) || len(text) == 0 {
		responseError(c, http.StatusBadRequest, "Invalid request: incorrect info")
		return
	}

	var senderId = r.IUserService.CheckUserExist(sender)
	if senderId <= 0 {
		responseError(c, http.StatusBadRequest, fmt.Sprintf("Invalid request: User name %s is not found", sender))
		return
	}

	emails := emailaddress.Find([]byte(text), false)
	var mentionedEmails []string

	mentionedEmails = make([]string, len(emails))
	for i, strEmail := range emails {
		mentionedEmails[i] = strEmail.String()
	}

	var mentionedIds []int64
	if len(mentionedEmails) > 0 {
		senderIndex := common.GetIndex(sender, mentionedEmails)
		if senderIndex >= 0 {
			mentionedEmails = common.RemoveItemInStringSlice(mentionedEmails, senderIndex)
		}

		mentionedIds = r.IUserService.CheckUsersExist(mentionedEmails)
	}

	result := r.IRelationshipService.GetValidUsersCanReceiveUpdates(senderId, mentionedIds)

	recipent := models.Recipent{Success: true, Recipents: result}

	responseOk(c, recipent)
	return
}
