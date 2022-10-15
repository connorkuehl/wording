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

// CreateGame provides a mock function with given fields: ctx, answer, guessLimit, expiresAfter
func (_m *MockService) CreateGame(ctx context.Context, answer string, guessLimit int, expiresAfter time.Duration) (*wording.Game, error) {
	ret := _m.Called(ctx, answer, guessLimit, expiresAfter)

	var r0 *wording.Game
	if rf, ok := ret.Get(0).(func(context.Context, string, int, time.Duration) *wording.Game); ok {
		r0 = rf(ctx, answer, guessLimit, expiresAfter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*wording.Game)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, int, time.Duration) error); ok {
		r1 = rf(ctx, answer, guessLimit, expiresAfter)
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
//  - expiresAfter time.Duration
func (_e *MockService_Expecter) CreateGame(ctx interface{}, answer interface{}, guessLimit interface{}, expiresAfter interface{}) *MockService_CreateGame_Call {
	return &MockService_CreateGame_Call{Call: _e.mock.On("CreateGame", ctx, answer, guessLimit, expiresAfter)}
}

func (_c *MockService_CreateGame_Call) Run(run func(ctx context.Context, answer string, guessLimit int, expiresAfter time.Duration)) *MockService_CreateGame_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(int), args[3].(time.Duration))
	})
	return _c
}

func (_c *MockService_CreateGame_Call) Return(_a0 *wording.Game, _a1 error) *MockService_CreateGame_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// Game provides a mock function with given fields: ctx, adminToken
func (_m *MockService) Game(ctx context.Context, adminToken string) (*wording.Game, error) {
	ret := _m.Called(ctx, adminToken)

	var r0 *wording.Game
	if rf, ok := ret.Get(0).(func(context.Context, string) *wording.Game); ok {
		r0 = rf(ctx, adminToken)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*wording.Game)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, adminToken)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockService_Game_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Game'
type MockService_Game_Call struct {
	*mock.Call
}

// Game is a helper method to define mock.On call
//  - ctx context.Context
//  - adminToken string
func (_e *MockService_Expecter) Game(ctx interface{}, adminToken interface{}) *MockService_Game_Call {
	return &MockService_Game_Call{Call: _e.mock.On("Game", ctx, adminToken)}
}

func (_c *MockService_Game_Call) Run(run func(ctx context.Context, adminToken string)) *MockService_Game_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockService_Game_Call) Return(_a0 *wording.Game, _a1 error) *MockService_Game_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// GameByToken provides a mock function with given fields: ctx, token
func (_m *MockService) GameByToken(ctx context.Context, token string) (*wording.Game, error) {
	ret := _m.Called(ctx, token)

	var r0 *wording.Game
	if rf, ok := ret.Get(0).(func(context.Context, string) *wording.Game); ok {
		r0 = rf(ctx, token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*wording.Game)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockService_GameByToken_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GameByToken'
type MockService_GameByToken_Call struct {
	*mock.Call
}

// GameByToken is a helper method to define mock.On call
//  - ctx context.Context
//  - token string
func (_e *MockService_Expecter) GameByToken(ctx interface{}, token interface{}) *MockService_GameByToken_Call {
	return &MockService_GameByToken_Call{Call: _e.mock.On("GameByToken", ctx, token)}
}

func (_c *MockService_GameByToken_Call) Run(run func(ctx context.Context, token string)) *MockService_GameByToken_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockService_GameByToken_Call) Return(_a0 *wording.Game, _a1 error) *MockService_GameByToken_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// NewPlayerToken provides a mock function with given fields: ctx
func (_m *MockService) NewPlayerToken(ctx context.Context) string {
	ret := _m.Called(ctx)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context) string); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockService_NewPlayerToken_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'NewPlayerToken'
type MockService_NewPlayerToken_Call struct {
	*mock.Call
}

// NewPlayerToken is a helper method to define mock.On call
//  - ctx context.Context
func (_e *MockService_Expecter) NewPlayerToken(ctx interface{}) *MockService_NewPlayerToken_Call {
	return &MockService_NewPlayerToken_Call{Call: _e.mock.On("NewPlayerToken", ctx)}
}

func (_c *MockService_NewPlayerToken_Call) Run(run func(ctx context.Context)) *MockService_NewPlayerToken_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockService_NewPlayerToken_Call) Return(_a0 string) *MockService_NewPlayerToken_Call {
	_c.Call.Return(_a0)
	return _c
}

// Plays provides a mock function with given fields: ctx, gameToken, playerToken
func (_m *MockService) Plays(ctx context.Context, gameToken string, playerToken string) (*wording.Plays, error) {
	ret := _m.Called(ctx, gameToken, playerToken)

	var r0 *wording.Plays
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *wording.Plays); ok {
		r0 = rf(ctx, gameToken, playerToken)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*wording.Plays)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, gameToken, playerToken)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockService_Plays_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Plays'
type MockService_Plays_Call struct {
	*mock.Call
}

// Plays is a helper method to define mock.On call
//  - ctx context.Context
//  - gameToken string
//  - playerToken string
func (_e *MockService_Expecter) Plays(ctx interface{}, gameToken interface{}, playerToken interface{}) *MockService_Plays_Call {
	return &MockService_Plays_Call{Call: _e.mock.On("Plays", ctx, gameToken, playerToken)}
}

func (_c *MockService_Plays_Call) Run(run func(ctx context.Context, gameToken string, playerToken string)) *MockService_Plays_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockService_Plays_Call) Return(_a0 *wording.Plays, _a1 error) *MockService_Plays_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// SubmitGuess provides a mock function with given fields: ctx, gameToken, playerToken, guess
func (_m *MockService) SubmitGuess(ctx context.Context, gameToken string, playerToken string, guess string) error {
	ret := _m.Called(ctx, gameToken, playerToken, guess)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) error); ok {
		r0 = rf(ctx, gameToken, playerToken, guess)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockService_SubmitGuess_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SubmitGuess'
type MockService_SubmitGuess_Call struct {
	*mock.Call
}

// SubmitGuess is a helper method to define mock.On call
//  - ctx context.Context
//  - gameToken string
//  - playerToken string
//  - guess string
func (_e *MockService_Expecter) SubmitGuess(ctx interface{}, gameToken interface{}, playerToken interface{}, guess interface{}) *MockService_SubmitGuess_Call {
	return &MockService_SubmitGuess_Call{Call: _e.mock.On("SubmitGuess", ctx, gameToken, playerToken, guess)}
}

func (_c *MockService_SubmitGuess_Call) Run(run func(ctx context.Context, gameToken string, playerToken string, guess string)) *MockService_SubmitGuess_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(string))
	})
	return _c
}

func (_c *MockService_SubmitGuess_Call) Return(_a0 error) *MockService_SubmitGuess_Call {
	_c.Call.Return(_a0)
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
