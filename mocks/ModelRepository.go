// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	context "context"
	domain "display_parser/internal/domain"

	mock "github.com/stretchr/testify/mock"
)

// ModelRepository is an autogenerated mock type for the ModelRepository type
type ModelRepository struct {
	mock.Mock
}

type ModelRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *ModelRepository) EXPECT() *ModelRepository_Expecter {
	return &ModelRepository_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: ctx, item
func (_m *ModelRepository) Create(ctx context.Context, item domain.ModelEntity) error {
	ret := _m.Called(ctx, item)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.ModelEntity) error); ok {
		r0 = rf(ctx, item)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ModelRepository_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type ModelRepository_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - item domain.ModelEntity
func (_e *ModelRepository_Expecter) Create(ctx interface{}, item interface{}) *ModelRepository_Create_Call {
	return &ModelRepository_Create_Call{Call: _e.mock.On("Create", ctx, item)}
}

func (_c *ModelRepository_Create_Call) Run(run func(ctx context.Context, item domain.ModelEntity)) *ModelRepository_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(domain.ModelEntity))
	})
	return _c
}

func (_c *ModelRepository_Create_Call) Return(_a0 error) *ModelRepository_Create_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ModelRepository_Create_Call) RunAndReturn(run func(context.Context, domain.ModelEntity) error) *ModelRepository_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Find provides a mock function with given fields: ctx, url
func (_m *ModelRepository) Find(ctx context.Context, url string) (domain.ModelEntity, bool, error) {
	ret := _m.Called(ctx, url)

	var r0 domain.ModelEntity
	var r1 bool
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (domain.ModelEntity, bool, error)); ok {
		return rf(ctx, url)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) domain.ModelEntity); ok {
		r0 = rf(ctx, url)
	} else {
		r0 = ret.Get(0).(domain.ModelEntity)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) bool); ok {
		r1 = rf(ctx, url)
	} else {
		r1 = ret.Get(1).(bool)
	}

	if rf, ok := ret.Get(2).(func(context.Context, string) error); ok {
		r2 = rf(ctx, url)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// ModelRepository_Find_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Find'
type ModelRepository_Find_Call struct {
	*mock.Call
}

// Find is a helper method to define mock.On call
//   - ctx context.Context
//   - url string
func (_e *ModelRepository_Expecter) Find(ctx interface{}, url interface{}) *ModelRepository_Find_Call {
	return &ModelRepository_Find_Call{Call: _e.mock.On("Find", ctx, url)}
}

func (_c *ModelRepository_Find_Call) Run(run func(ctx context.Context, url string)) *ModelRepository_Find_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *ModelRepository_Find_Call) Return(_a0 domain.ModelEntity, _a1 bool, _a2 error) *ModelRepository_Find_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *ModelRepository_Find_Call) RunAndReturn(run func(context.Context, string) (domain.ModelEntity, bool, error)) *ModelRepository_Find_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewModelRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewModelRepository creates a new instance of ModelRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewModelRepository(t mockConstructorTestingTNewModelRepository) *ModelRepository {
	mock := &ModelRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}