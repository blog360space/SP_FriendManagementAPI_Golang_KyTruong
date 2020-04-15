package services

import (
	"friendMgmt/data"
	"friendMgmt/models"
)

type IRelationshipService interface {
	CreateRelationship(relationship *models.Relationship) int64
	DeleteRelationships(ids []int64) bool
	CheckConnected(requestUserId int64, targetUserId int64) []int64
	CheckFullySubcribed(requestUserId int64, targetUserId int64) []int64
	CheckFullyBlocked(requestUserId int64, targetUserId int64) []int64
	CheckPartialSubcribed(requestUserId int64, targetUserId int64) []int64
	CheckPartialBlocked(requestUserId int64, targetUserId int64) []int64
	GetFriendList(id int64) []string
	GetCommonFriendList(id int64, withId int64) []string
	GetValidUsersCanReceiveUpdates(senderId int64, mentionIds []int64) []string
}

type RelationshipService struct {
	IRelationshipRepository data.IRelationshipRepository
}

func (svc RelationshipService) GetFriendList(id int64) []string {
	return svc.IRelationshipRepository.GetFriendList(id)
}

func (svc RelationshipService) GetCommonFriendList(id int64, withId int64) []string {
	return svc.IRelationshipRepository.GetCommonFriendList(id, withId)
}

func (svc RelationshipService) CreateRelationship(relationship *models.Relationship) int64 {
	return svc.IRelationshipRepository.CreateRelationship(relationship)
}

func (svc RelationshipService) DeleteRelationships(ids []int64) bool {
	return svc.IRelationshipRepository.DeleteRelationships(ids)
}

func (svc RelationshipService) CheckConnected(requestUserId int64, targetUserId int64) []int64 {
	return svc.IRelationshipRepository.CheckRelationshipTwoWay(requestUserId, targetUserId, 1)
}

func (svc RelationshipService) CheckFullySubcribed(requestUserId int64, targetUserId int64) []int64 {
	return svc.IRelationshipRepository.CheckRelationshipTwoWay(requestUserId, targetUserId, 2)
}

func (svc RelationshipService) CheckFullyBlocked(requestUserId int64, targetUserId int64) []int64 {
	return svc.IRelationshipRepository.CheckRelationshipTwoWay(requestUserId, targetUserId, 3)
}

func (svc RelationshipService) CheckPartialSubcribed(requestUserId int64, targetUserId int64) []int64 {
	return svc.IRelationshipRepository.CheckRelationshipOneWay(requestUserId, targetUserId, 2)
}

func (svc RelationshipService) CheckPartialBlocked(requestUserId int64, targetUserId int64) []int64 {
	return svc.IRelationshipRepository.CheckRelationshipOneWay(requestUserId, targetUserId, 3)
}

func (svc RelationshipService) GetValidUsersCanReceiveUpdates(senderId int64, mentionIds []int64) []string {
	return svc.IRelationshipRepository.GetValidUsersCanReceiveUpdates(senderId, mentionIds)
}
