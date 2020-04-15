package services_test

import (
	"friendMgmt/data"
	"friendMgmt/models"
	"friendMgmt/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateRelationship(t *testing.T) {
	relationshipModel := models.Relationship{Status: int64(1), RequestUserId: int64(1), TargetUserId: int64(2)}

	relationshipRepositoryMock := data.RelationshipRepositoryMock{}
	relationshipRepositoryMock.On("CreateRelationship", &relationshipModel).Return(int64(1))

	relationshipService := services.RelationshipService{relationshipRepositoryMock}
	id := relationshipService.CreateRelationship(&relationshipModel)

	assert.Equal(t, int64(1), id)

	relationshipRepositoryMock.AssertExpectations(t)
}

func TestDeleteRelationships(t *testing.T) {
	ids := []int64{int64(1), int64(2)}

	relationshipRepositoryMock := data.RelationshipRepositoryMock{}
	relationshipRepositoryMock.On("DeleteRelationships", ids).Return(true)

	relationshipService := services.RelationshipService{relationshipRepositoryMock}
	isDeleted := relationshipService.DeleteRelationships(ids)

	assert.Equal(t, true, isDeleted)

	relationshipRepositoryMock.AssertExpectations(t)
}

func TestGetFriendList(t *testing.T) {
	expectedResult := []string{"user1@gmail.com", "user2@gmail.com"}

	relationshipRepositoryMock := data.RelationshipRepositoryMock{}
	relationshipRepositoryMock.On("GetFriendList", int64(1)).Return(expectedResult)

	relationshipService := services.RelationshipService{relationshipRepositoryMock}

	assert.Equal(t, expectedResult, relationshipService.GetFriendList(int64(1)))

	relationshipRepositoryMock.AssertExpectations(t)
}

func TestGetCommonFriendList(t *testing.T) {
	expectedResult := []string{"user1@gmail.com", "user2@gmail.com"}

	relationshipRepositoryMock := data.RelationshipRepositoryMock{}
	relationshipRepositoryMock.On("GetCommonFriendList", int64(1), int64(2)).Return(expectedResult)

	relationshipService := services.RelationshipService{relationshipRepositoryMock}

	assert.Equal(t, expectedResult, relationshipService.GetCommonFriendList(int64(1), int64(2)))

	relationshipRepositoryMock.AssertExpectations(t)
}

func TestGetValidUsersCanReceiveUpdates(t *testing.T) {
	expectedResult := []string{"user1@gmail.com", "user2@gmail.com"}
	senderId := int64(1)
	mentionedIds := []int64{int64(2), int64(3)}

	relationshipRepositoryMock := data.RelationshipRepositoryMock{}
	relationshipRepositoryMock.On("GetValidUsersCanReceiveUpdates", senderId, mentionedIds).Return(expectedResult)

	relationshipService := services.RelationshipService{relationshipRepositoryMock}

	assert.Equal(t, expectedResult, relationshipService.GetValidUsersCanReceiveUpdates(senderId, mentionedIds))

	relationshipRepositoryMock.AssertExpectations(t)
}

func TestCheckConnected(t *testing.T) {
	expectedResult := []int64{int64(3), int64(4)}
	requestUserId := int64(1)
	targetUserId := int64(2)

	relationshipRepositoryMock := data.RelationshipRepositoryMock{}
	relationshipRepositoryMock.On("CheckRelationshipTwoWay", requestUserId, targetUserId, int64(1)).Return(expectedResult)

	relationshipService := services.RelationshipService{relationshipRepositoryMock}

	assert.Equal(t, expectedResult, relationshipService.CheckConnected(requestUserId, targetUserId))

	relationshipRepositoryMock.AssertExpectations(t)
}

func TestCheckFullySubcribed(t *testing.T) {
	expectedResult := []int64{int64(3), int64(4)}
	requestUserId := int64(1)
	targetUserId := int64(2)

	relationshipRepositoryMock := data.RelationshipRepositoryMock{}
	relationshipRepositoryMock.On("CheckRelationshipTwoWay", requestUserId, targetUserId, int64(2)).Return(expectedResult)

	relationshipService := services.RelationshipService{relationshipRepositoryMock}

	assert.Equal(t, expectedResult, relationshipService.CheckFullySubcribed(requestUserId, targetUserId))

	relationshipRepositoryMock.AssertExpectations(t)
}

func TestCheckFullyBlocked(t *testing.T) {
	expectedResult := []int64{int64(3), int64(4)}
	requestUserId := int64(1)
	targetUserId := int64(2)

	relationshipRepositoryMock := data.RelationshipRepositoryMock{}
	relationshipRepositoryMock.On("CheckRelationshipTwoWay", requestUserId, targetUserId, int64(3)).Return(expectedResult)

	relationshipService := services.RelationshipService{relationshipRepositoryMock}

	assert.Equal(t, expectedResult, relationshipService.CheckFullyBlocked(requestUserId, targetUserId))

	relationshipRepositoryMock.AssertExpectations(t)
}

func TestCheckPartialSubcribed(t *testing.T) {
	expectedResult := []int64{int64(3), int64(4)}
	requestUserId := int64(1)
	targetUserId := int64(2)

	relationshipRepositoryMock := data.RelationshipRepositoryMock{}
	relationshipRepositoryMock.On("CheckRelationshipOneWay", requestUserId, targetUserId, int64(2)).Return(expectedResult)

	relationshipService := services.RelationshipService{relationshipRepositoryMock}

	assert.Equal(t, expectedResult, relationshipService.CheckPartialSubcribed(requestUserId, targetUserId))

	relationshipRepositoryMock.AssertExpectations(t)
}

func TestCheckPartialBlocked(t *testing.T) {
	expectedResult := []int64{int64(3), int64(4)}
	requestUserId := int64(1)
	targetUserId := int64(2)

	relationshipRepositoryMock := data.RelationshipRepositoryMock{}
	relationshipRepositoryMock.On("CheckRelationshipOneWay", requestUserId, targetUserId, int64(3)).Return(expectedResult)

	relationshipService := services.RelationshipService{relationshipRepositoryMock}

	assert.Equal(t, expectedResult, relationshipService.CheckPartialBlocked(requestUserId, targetUserId))

	relationshipRepositoryMock.AssertExpectations(t)
}
