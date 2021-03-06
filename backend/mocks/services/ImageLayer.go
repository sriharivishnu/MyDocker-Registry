// Code generated by mockery 2.9.2. DO NOT EDIT.

package mocks

import (
	models "github.com/sriharivishnu/shopify-challenge/models"
	mock "github.com/stretchr/testify/mock"
)

// ImageLayer is an autogenerated mock type for the ImageLayer type
type ImageLayer struct {
	mock.Mock
}

// Create provides a mock function with given fields: repoId, tag, description, fileKey
func (_m *ImageLayer) Create(repoId string, tag string, description string, fileKey string) (models.ImageTag, error) {
	ret := _m.Called(repoId, tag, description, fileKey)

	var r0 models.ImageTag
	if rf, ok := ret.Get(0).(func(string, string, string, string) models.ImageTag); ok {
		r0 = rf(repoId, tag, description, fileKey)
	} else {
		r0 = ret.Get(0).(models.ImageTag)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string, string) error); ok {
		r1 = rf(repoId, tag, description, fileKey)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetImageTagByRepoAndTag provides a mock function with given fields: repository_id, tagName
func (_m *ImageLayer) GetImageTagByRepoAndTag(repository_id string, tagName string) (models.ImageTag, error) {
	ret := _m.Called(repository_id, tagName)

	var r0 models.ImageTag
	if rf, ok := ret.Get(0).(func(string, string) models.ImageTag); ok {
		r0 = rf(repository_id, tagName)
	} else {
		r0 = ret.Get(0).(models.ImageTag)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(repository_id, tagName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetImageTagsForRepo provides a mock function with given fields: repository_id
func (_m *ImageLayer) GetImageTagsForRepo(repository_id string) ([]models.ImageTag, error) {
	ret := _m.Called(repository_id)

	var r0 []models.ImageTag
	if rf, ok := ret.Get(0).(func(string) []models.ImageTag); ok {
		r0 = rf(repository_id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.ImageTag)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(repository_id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
