// Code generated by mockery v2.14.0. DO NOT EDIT.

package server

import (
	context "context"
	time "time"

	mock "github.com/stretchr/testify/mock"

	wording "github.com/connorkuehl/wording/internal/wording"
)

// MockService is an autogenerated mock type for the Service type
type MockService struct {
	mock.Mock
}

type MockService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockService) EXPECT() *MockService_Expecter {
	return &MockService_Expecter{mock: &_m.Mock}
}

// CreateGame provides a mock function with given fields: ctx, answer, guessLimit, expiresAt
func (_m *MockService) CreateGame(ctx context.Context, answer string, guessLimit int, expiresAt time.Time) (*wording.Game, error) {
	ret := _m.Called(ctx, answer, guessLimit, expiresAt)

	var r0 *wording.Game
	if rf, ok := ret.Get(0).(func(context.Context, string, int, time.Time) *wording.Game); ok {
		r0 = rf(ctx, answer, guessLimit, expiresAt)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*wording.Game)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, int, time.Time) error); ok {
		r1 = rf(ctx, answer, guessLimit, expiresAt)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockService_CreateGame_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateGame'
type MockService_CreateGame_Call struct {
	*mock.Call
}

// CreateGame is a helper method to define mock.On call
//  - ctx context.Context
//  - answer string
//  - guessLimit int
//  - expiresAt time.Time
func (_e *MockService_Expecter) CreateGame(ctx interface{}, answer interface{}, guessLimit interface{}, expiresAt interface{}) *MockService_CreateGame_Call {
	return &MockService_CreateGame_Call{Call: _e.mock.On("CreateGame", ctx, answer, guessLimit, expiresAt)}
}

func (_c *MockService_CreateGame_Call) Run(run func(ctx context.Context, answer string, guessLimit int, expiresAt time.Time)) *MockService_CreateGame_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(int), args[3].(time.Time))
	})
	return _c
}

func (_c *MockService_CreateGame_Call) Return(_a0 *wording.Game, _a1 error) *MockService_CreateGame_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

type mockConstructorTestingTNewMockService interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockService creates a new instance of MockService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockService(t mockConstructorTestingTNewMockService) *MockService {
	mock := &MockService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}