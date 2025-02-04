// Code generated by mockery v2.51.1. DO NOT EDIT.

package mocks

import (
	jwt "github.com/golang-jwt/jwt"
	mock "github.com/stretchr/testify/mock"
)

// JWTService is an autogenerated mock type for the JWTService type
type JWTService struct {
	mock.Mock
}

// GenerateToken provides a mock function with given fields: userId, name, email
func (_m *JWTService) GenerateToken(userId string, name string, email string) string {
	ret := _m.Called(userId, name, email)

	if len(ret) == 0 {
		panic("no return value specified for GenerateToken")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string, string) string); ok {
		r0 = rf(userId, name, email)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// ValidateToken provides a mock function with given fields: token
func (_m *JWTService) ValidateToken(token string) (*jwt.Token, error) {
	ret := _m.Called(token)

	if len(ret) == 0 {
		panic("no return value specified for ValidateToken")
	}

	var r0 *jwt.Token
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*jwt.Token, error)); ok {
		return rf(token)
	}
	if rf, ok := ret.Get(0).(func(string) *jwt.Token); ok {
		r0 = rf(token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*jwt.Token)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewJWTService creates a new instance of JWTService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewJWTService(t interface {
	mock.TestingT
	Cleanup(func())
}) *JWTService {
	mock := &JWTService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
