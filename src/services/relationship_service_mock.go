package services

import (
	"friendMgmt/models"

	"github.com/stretchr/testify/mock"
)

type RelationshipServiceMock struct {
	mock.Mock
}

func (m RelationshipServiceMock) GetFriendList(id int64) []string {
	args := m.Called(id)

	return args.Get(0).([]string)
}

func (m RelationshipServiceMock) GetCommonFriendList(id int64, withId int64) []string {
	args := m.Called(id, withId)

	return args.Get(0).([]string)
}

func (m RelationshipServiceMock) CreateRelationship(relationship *models.Relationship) int64 {
	args := m.Called(relationship)

	return args.Get(0).(int64)
}

func (m RelationshipServiceMock) DeleteRelationships(ids []int64) bool {
	args := m.Called(ids)

	return args.Get(0).(bool)
}

func (m RelationshipServiceMock) CheckConnected(requestUserId int64, targetUserId int64) []int64 {
	args := m.Called(requestUserId, targetUserId)

	return args.Get(0).([]int64)
}

func (m RelationshipServiceMock) CheckFullySubcribed(requestUserId int64, targetUserId int64) []int64 {
	args := m.Called(requestUserId, targetUserId)

	return args.Get(0).([]int64)
}

func (m RelationshipServiceMock) CheckFullyBlocked(requestUserId int64, targetUserId int64) []int64 {
	args := m.Called(requestUserId, targetUserId)

	return args.Get(0).([]int64)
}

func (m RelationshipServiceMock) CheckPartialSubcribed(requestUserId int64, targetUserId int64) []int64 {
	args := m.Called(requestUserId, targetUserId)

	return args.Get(0).([]int64)
}

func (m RelationshipServiceMock) CheckPartialBlocked(requestUserId int64, targetUserId int64) []int64 {
	args := m.Called(requestUserId, targetUserId)

	return args.Get(0).([]int64)
}

func (m RelationshipServiceMock) GetValidUsersCanReceiveUpdates(senderId int64, mentionIds []int64) []string {
	args := m.Called(senderId, mentionIds)

	return args.Get(0).([]string)
}
