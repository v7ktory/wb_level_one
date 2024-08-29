// Code generated by mockery v2.45.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	entity "github.com/v7ktory/wb_task_one/internal/entity"
)

// Order is an autogenerated mock type for the Order type
type Order struct {
	mock.Mock
}

// GetLRUOrders provides a mock function with given fields: ctx
func (_m *Order) GetLRUOrders(ctx context.Context) ([]*entity.Order, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetLRUOrders")
	}

	var r0 []*entity.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]*entity.Order, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []*entity.Order); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*entity.Order)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SaveOrder provides a mock function with given fields: ctx, order
func (_m *Order) SaveOrder(ctx context.Context, order *entity.Order) (string, error) {
	ret := _m.Called(ctx, order)

	if len(ret) == 0 {
		panic("no return value specified for SaveOrder")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *entity.Order) (string, error)); ok {
		return rf(ctx, order)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *entity.Order) string); ok {
		r0 = rf(ctx, order)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *entity.Order) error); ok {
		r1 = rf(ctx, order)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateOrderTime provides a mock function with given fields: ctx, uid
func (_m *Order) UpdateOrderTime(ctx context.Context, uid string) error {
	ret := _m.Called(ctx, uid)

	if len(ret) == 0 {
		panic("no return value specified for UpdateOrderTime")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, uid)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewOrder creates a new instance of Order. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewOrder(t interface {
	mock.TestingT
	Cleanup(func())
}) *Order {
	mock := &Order{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
