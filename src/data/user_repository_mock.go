package data

import "github.com/stretchr/testify/mock"

type UserRepositoryMock struct {
	mock.Mock
}

func (m UserRepositoryMock) FindAll() []string {
	args := m.Called()

	return args.Get(0).([]string)
}

func (m UserRepositoryMock) Create(email string) bool {
	args := m.Called(email)

	return args.Get(0).(bool)
}

func (m UserRepositoryMock) CheckUserExist(email string) int64 {
	args := m.Called(email)

	return args.Get(0).(int64)
}

func (m UserRepositoryMock) CheckUsersExist(emails []string) []int64 {
	args := m.Called(emails)

	return args.Get(0).([]int64)
}
