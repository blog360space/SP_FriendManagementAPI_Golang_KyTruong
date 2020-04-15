package services

import "github.com/stretchr/testify/mock"

type UserServiceMock struct {
	mock.Mock
}

func (m UserServiceMock) FindAll() []string {
	args := m.Called()

	return args.Get(0).([]string)
}

func (m UserServiceMock) Create(email string) bool {
	args := m.Called(email)

	return args.Get(0).(bool)
}

func (m UserServiceMock) CheckUserExist(email string) int64 {
	args := m.Called(email)

	return args.Get(0).(int64)
}

func (m UserServiceMock) CheckUsersExist(emails []string) []int64 {
	args := m.Called(emails)

	return args.Get(0).([]int64)
}
