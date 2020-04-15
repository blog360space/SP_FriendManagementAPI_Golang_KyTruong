package services_test

import (
	"friendMgmt/data"
	"friendMgmt/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindAll(t *testing.T) {
	userRepositoryMock := data.UserRepositoryMock{}

	expectedResult := []string{"user1@gmail.com", "user2@gmail.com"}

	userRepositoryMock.On("FindAll").Return(expectedResult)

	userService := services.UserService{userRepositoryMock}

	assert.Equal(t, expectedResult, userService.FindAll())

	userRepositoryMock.AssertExpectations(t)
}

func TestCreate(t *testing.T) {
	userRepositoryMock := data.UserRepositoryMock{}

	createSuccess := true

	userRepositoryMock.On("Create", "user@test.com").Return(createSuccess)

	userService := services.UserService{userRepositoryMock}

	assert.Equal(t, createSuccess, userService.Create("user@test.com"))

	userRepositoryMock.AssertExpectations(t)
}

func TestCheckUserExist(t *testing.T) {
	userRepositoryMock := data.UserRepositoryMock{}

	userRepositoryMock.On("CheckUserExist", "user@test.com").Return(int64(1))

	userService := services.UserService{userRepositoryMock}

	assert.Equal(t, int64(1), userService.CheckUserExist("user@test.com"))

	userRepositoryMock.AssertExpectations(t)
}

func TestCheckUsersExist(t *testing.T) {
	userRepositoryMock := data.UserRepositoryMock{}

	emails := []string{"user1@gmail.com", "user2@gmail.com"}
	idsResult := []int64{int64(1), int64(2)}

	userRepositoryMock.On("CheckUsersExist", emails).Return(idsResult)

	userService := services.UserService{userRepositoryMock}

	assert.Equal(t, idsResult, userService.CheckUsersExist(emails))

	userRepositoryMock.AssertExpectations(t)
}
