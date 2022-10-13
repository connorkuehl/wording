// Code generated by mockery v2.14.0. DO NOT EDIT.

package service

import (
	context "context"
	time "time"

	mock "github.com/stretchr/testify/mock"

	wording "github.com/connorkuehl/wording/internal/wording"
)

// MockStore is an autogenerated mock type for the Store type
type MockStore struct {
	mock.Mock
}

type MockStore_Expecter struct {
	mock *mock.Mock
}

func (_m *MockStore) EXPECT() *MockStore_Expecter {
	return &MockStore_Expecter{mock: &_m.Mock}
}

// CreateGame provides a mock function with given fields: ctx, adminToken, answer, guessLimit, expiresAt
func (_m *MockStore) CreateGame(ctx context.Context, adminToken string, answer string, guessLimit int, expiresAt time.Time) (*wording.Game, error) {
	ret := _m.Called(ctx, adminToken, answer, guessLimit, expiresAt)

	var r0 *wording.Game
	if rf, ok := ret.Get(0).(func(context.Context, string, string, int, time.Time) *wording.Game); ok {
		r0 = rf(ctx, adminToken, answer, guessLimit, expiresAt)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*wording.Game)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, int, time.Time) error); ok {
		r1 = rf(ctx, adminToken, answer, guessLimit, expiresAt)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockStore_CreateGame_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateGame'
type MockStore_CreateGame_Call struct {
	*mock.Call
}

// CreateGame is a helper method to define mock.On call
//  - ctx context.Context
//  - adminToken string
//  - answer string
//  - guessLimit int
//  - expiresAt time.Time
func (_e *MockStore_Expecter) CreateGame(ctx interface{}, adminToken interface{}, answer interface{}, guessLimit interface{}, expiresAt interface{}) *MockStore_CreateGame_Call {
	return &MockStore_CreateGame_Call{Call: _e.mock.On("CreateGame", ctx, adminToken, answer, guessLimit, expiresAt)}
}

func (_c *MockStore_CreateGame_Call) Run(run func(ctx context.Context, adminToken string, answer string, guessLimit int, expiresAt time.Time)) *MockStore_CreateGame_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(int), args[4].(time.Time))
	})
	return _c
}

func (_c *MockStore_CreateGame_Call) Return(_a0 *wording.Game, _a1 error) *MockStore_CreateGame_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

type mockConstructorTestingTNewMockStore interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockStore creates a new instance of MockStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockStore(t mockConstructorTestingTNewMockStore) *MockStore {
	mock := &MockStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
