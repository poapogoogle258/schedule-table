// Code generated by mockery v2.51.1. DO NOT EDIT.

package mocks

import (
	dao "schedule_table/internal/model/dao"

	mock "github.com/stretchr/testify/mock"

	service "schedule_table/internal/service"
)

// IManagerService is an autogenerated mock type for the IManagerService type
type IManagerService struct {
	mock.Mock
}

// NewManagerSchedule provides a mock function with given fields: schedule
func (_m *IManagerService) NewManagerSchedule(schedule *dao.Schedules) *service.Manager {
	ret := _m.Called(schedule)

	if len(ret) == 0 {
		panic("no return value specified for NewManagerSchedule")
	}

	var r0 *service.Manager
	if rf, ok := ret.Get(0).(func(*dao.Schedules) *service.Manager); ok {
		r0 = rf(schedule)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*service.Manager)
		}
	}

	return r0
}

// NewManagerScheduleWithQueue provides a mock function with given fields: schedule, queue
func (_m *IManagerService) NewManagerScheduleWithQueue(schedule *dao.Schedules, queue service.IQueue) *service.Manager {
	ret := _m.Called(schedule, queue)

	if len(ret) == 0 {
		panic("no return value specified for NewManagerScheduleWithQueue")
	}

	var r0 *service.Manager
	if rf, ok := ret.Get(0).(func(*dao.Schedules, service.IQueue) *service.Manager); ok {
		r0 = rf(schedule, queue)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*service.Manager)
		}
	}

	return r0
}

// NewIManagerService creates a new instance of IManagerService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIManagerService(t interface {
	mock.TestingT
	Cleanup(func())
}) *IManagerService {
	mock := &IManagerService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
