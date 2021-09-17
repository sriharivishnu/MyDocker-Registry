// Code generated by mockery 2.9.2. DO NOT EDIT.

package mocks

import (
	models "github.com/sriharivishnu/shopify-challenge/models"
	mock "github.com/stretchr/testify/mock"
)

// UserLayer is an autogenerated mock type for the UserLayer type
type UserLayer struct {
	mock.Mock
}

// Create provides a mock function with given fields: username, password
func (_m *UserLayer) Create(username string, password string) (models.User, error) {
	ret := _m.Called(username, password)

	var r0 models.User
	if rf, ok := ret.Get(0).(func(string, string) models.User); ok {
		r0 = rf(username, password)
	} else {
		r0 = ret.Get(0).(models.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(username, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateToken provides a mock function with given fields: user
func (_m *UserLayer) CreateToken(user models.User) (string, error) {
	ret := _m.Called(user)

	var r0 string
	if rf, ok := ret.Get(0).(func(models.User) string); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(models.User) error); ok {
		r1 = rf(user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: userId
func (_m *UserLayer) GetByID(userId string) (models.User, error) {
	ret := _m.Called(userId)

	var r0 models.User
	if rf, ok := ret.Get(0).(func(string) models.User); ok {
		r0 = rf(userId)
	} else {
		r0 = ret.Get(0).(models.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByUsername provides a mock function with given fields: username
func (_m *UserLayer) GetByUsername(username string) (models.User, error) {
	ret := _m.Called(username)

	var r0 models.User
	if rf, ok := ret.Get(0).(func(string) models.User); ok {
		r0 = rf(username)
	} else {
		r0 = ret.Get(0).(models.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}