package services

import (
	"friendMgmt/data"
)

type IUserService interface {
	FindAll() []string
	Create(email string) bool
	CheckUserExist(email string) int64
	CheckUsersExist(emails []string) []int64
}

type UserService struct {
	IUserRepository data.IUserRepository
}

func (svc UserService) FindAll() []string {
	return svc.IUserRepository.FindAll()
}

func (svc UserService) Create(email string) bool {
	return svc.IUserRepository.Create(email)
}

func (svc UserService) CheckUserExist(email string) int64 {
	return svc.IUserRepository.CheckUserExist(email)
}

func (svc UserService) CheckUsersExist(emails []string) []int64 {
	return svc.IUserRepository.CheckUsersExist(emails)
}
