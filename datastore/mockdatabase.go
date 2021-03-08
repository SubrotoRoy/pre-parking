package datastore

import (
	"github.com/SubrotoRoy/pre-parking/model"
	"github.com/stretchr/testify/mock"
)

//MockDBRepo is a mock implementation of a datastore client for testing purposes
type MockDBRepo struct {
	mock.Mock
}

//DBConnect is a mock implementation
func (m *MockDBRepo) DBConnect(config model.DbConfig) error {
	args := m.Mock.Called(config)
	return args.Error(0)
}

//IsSlotAvailable is a mock implementation
func (m *MockDBRepo) IsSlotAvailable() bool {
	args := m.Mock.Called()
	return args.Bool(0)
}

//IsCarParked is a mock implementation
func (m *MockDBRepo) IsCarParked(carNumber string) bool {
	args := m.Mock.Called(carNumber)
	return args.Bool(0)
}
