// Code generated by mockery v2.36.1. DO NOT EDIT.

package mocks

import (
	models "job-services/internal/models"

	mock "github.com/stretchr/testify/mock"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// GetTransactionByTitle provides a mock function with given fields: title
func (_m *Repository) GetTransactionByTitle(title string) ([]models.TransactionModel, error) {
	ret := _m.Called(title)

	var r0 []models.TransactionModel
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]models.TransactionModel, error)); ok {
		return rf(title)
	}
	if rf, ok := ret.Get(0).(func(string) []models.TransactionModel); ok {
		r0 = rf(title)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.TransactionModel)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(title)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Insert provides a mock function with given fields: entity
func (_m *Repository) Insert(entity []*models.TransactionModel) error {
	ret := _m.Called(entity)

	var r0 error
	if rf, ok := ret.Get(0).(func([]*models.TransactionModel) error); ok {
		r0 = rf(entity)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewRepository creates a new instance of Repository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *Repository {
	mock := &Repository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
