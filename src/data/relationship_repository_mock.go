package data

import (
	"friendMgmt/models"

	"github.com/stretchr/testify/mock"
)

type RelationshipRepositoryMock struct {
	mock.Mock
}

func (m RelationshipRepositoryMock) CreateRelationship(relationship *models.Relationship) int64 {
	args := m.Called(relationship)

	return args.Get(0).(int64)
}

func (m RelationshipRepositoryMock) DeleteRelationships(ids []int64) bool {
	args := m.Called(ids)

	return args.Get(0).(bool)
}

func (m RelationshipRepositoryMock) GetFriendList(id int64) []string {
	args := m.Called(id)

	return args.Get(0).([]string)
}

func (m RelationshipRepositoryMock) GetCommonFriendList(id int64, withId int64) []string {
	args := m.Called(id, withId)

	return args.Get(0).([]string)
}

func (m RelationshipRepositoryMock) GetValidUsersCanReceiveUpdates(senderId int64, mentionIds []int64) []string {
	args := m.Called(senderId, mentionIds)

	return args.Get(0).([]string)
}

func (m RelationshipRepositoryMock) CheckRelationshipTwoWay(requestUserId int64, targetUserId int64, status int64) []int64 {
	args := m.Called(requestUserId, targetUserId, status)

	return args.Get(0).([]int64)
}

func (m RelationshipRepositoryMock) CheckRelationshipOneWay(requestUserId int64, targetUserId int64, status int64) []int64 {
	args := m.Called(requestUserId, targetUserId, status)

	return args.Get(0).([]int64)
}
