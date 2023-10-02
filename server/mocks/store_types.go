// Code generated by MockGen. DO NOT EDIT.
// Source: ./core/store_types.go
//
// Generated by this command:
//
//	mockgen -source=./core/store_types.go -package=mocks
//
// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	core "main/core"
	reflect "reflect"
	time "time"

	gomock "go.uber.org/mock/gomock"
)

// MockPushLogger is a mock of PushLogger interface.
type MockPushLogger struct {
	ctrl     *gomock.Controller
	recorder *MockPushLoggerMockRecorder
}

// MockPushLoggerMockRecorder is the mock recorder for MockPushLogger.
type MockPushLoggerMockRecorder struct {
	mock *MockPushLogger
}

// NewMockPushLogger creates a new mock instance.
func NewMockPushLogger(ctrl *gomock.Controller) *MockPushLogger {
	mock := &MockPushLogger{ctrl: ctrl}
	mock.recorder = &MockPushLoggerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPushLogger) EXPECT() *MockPushLoggerMockRecorder {
	return m.recorder
}

// LogPush mocks base method.
func (m *MockPushLogger) LogPush(ctx context.Context, clientId string, push *core.PushMessage) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LogPush", ctx, clientId, push)
	ret0, _ := ret[0].(error)
	return ret0
}

// LogPush indicates an expected call of LogPush.
func (mr *MockPushLoggerMockRecorder) LogPush(ctx, clientId, push any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogPush", reflect.TypeOf((*MockPushLogger)(nil).LogPush), ctx, clientId, push)
}

// MockIntentEmbeddingStore is a mock of IntentEmbeddingStore interface.
type MockIntentEmbeddingStore struct {
	ctrl     *gomock.Controller
	recorder *MockIntentEmbeddingStoreMockRecorder
}

// MockIntentEmbeddingStoreMockRecorder is the mock recorder for MockIntentEmbeddingStore.
type MockIntentEmbeddingStoreMockRecorder struct {
	mock *MockIntentEmbeddingStore
}

// NewMockIntentEmbeddingStore creates a new mock instance.
func NewMockIntentEmbeddingStore(ctrl *gomock.Controller) *MockIntentEmbeddingStore {
	mock := &MockIntentEmbeddingStore{ctrl: ctrl}
	mock.recorder = &MockIntentEmbeddingStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIntentEmbeddingStore) EXPECT() *MockIntentEmbeddingStoreMockRecorder {
	return m.recorder
}

// ClosestMatchingIntent mocks base method.
func (m *MockIntentEmbeddingStore) ClosestMatchingIntent(ctx context.Context, embedding []float32) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClosestMatchingIntent", ctx, embedding)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ClosestMatchingIntent indicates an expected call of ClosestMatchingIntent.
func (mr *MockIntentEmbeddingStoreMockRecorder) ClosestMatchingIntent(ctx, embedding any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClosestMatchingIntent", reflect.TypeOf((*MockIntentEmbeddingStore)(nil).ClosestMatchingIntent), ctx, embedding)
}

// UpdateIntentEmbedding mocks base method.
func (m *MockIntentEmbeddingStore) UpdateIntentEmbedding(ctx context.Context, intent string, embedding []float32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateIntentEmbedding", ctx, intent, embedding)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateIntentEmbedding indicates an expected call of UpdateIntentEmbedding.
func (mr *MockIntentEmbeddingStoreMockRecorder) UpdateIntentEmbedding(ctx, intent, embedding any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateIntentEmbedding", reflect.TypeOf((*MockIntentEmbeddingStore)(nil).UpdateIntentEmbedding), ctx, intent, embedding)
}

// MockScheduler is a mock of Scheduler interface.
type MockScheduler struct {
	ctrl     *gomock.Controller
	recorder *MockSchedulerMockRecorder
}

// MockSchedulerMockRecorder is the mock recorder for MockScheduler.
type MockSchedulerMockRecorder struct {
	mock *MockScheduler
}

// NewMockScheduler creates a new mock instance.
func NewMockScheduler(ctrl *gomock.Controller) *MockScheduler {
	mock := &MockScheduler{ctrl: ctrl}
	mock.recorder = &MockSchedulerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockScheduler) EXPECT() *MockSchedulerMockRecorder {
	return m.recorder
}

// ClearScheduledEvent mocks base method.
func (m *MockScheduler) ClearScheduledEvent(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClearScheduledEvent", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// ClearScheduledEvent indicates an expected call of ClearScheduledEvent.
func (mr *MockSchedulerMockRecorder) ClearScheduledEvent(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClearScheduledEvent", reflect.TypeOf((*MockScheduler)(nil).ClearScheduledEvent), ctx, id)
}

// ClearScheduledRecurringEvent mocks base method.
func (m *MockScheduler) ClearScheduledRecurringEvent(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClearScheduledRecurringEvent", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// ClearScheduledRecurringEvent indicates an expected call of ClearScheduledRecurringEvent.
func (mr *MockSchedulerMockRecorder) ClearScheduledRecurringEvent(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClearScheduledRecurringEvent", reflect.TypeOf((*MockScheduler)(nil).ClearScheduledRecurringEvent), ctx, id)
}

// ReadyEvents mocks base method.
func (m *MockScheduler) ReadyEvents(ctx context.Context, frontier time.Time, eventType string, infoParser func(*core.ScheduledEvent, string) error) ([]core.ScheduledEvent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadyEvents", ctx, frontier, eventType, infoParser)
	ret0, _ := ret[0].([]core.ScheduledEvent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadyEvents indicates an expected call of ReadyEvents.
func (mr *MockSchedulerMockRecorder) ReadyEvents(ctx, frontier, eventType, infoParser any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadyEvents", reflect.TypeOf((*MockScheduler)(nil).ReadyEvents), ctx, frontier, eventType, infoParser)
}

// RecurringEvents mocks base method.
func (m *MockScheduler) RecurringEvents(ctx context.Context) ([]core.ScheduledRecurringEvent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecurringEvents", ctx)
	ret0, _ := ret[0].([]core.ScheduledRecurringEvent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RecurringEvents indicates an expected call of RecurringEvents.
func (mr *MockSchedulerMockRecorder) RecurringEvents(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecurringEvents", reflect.TypeOf((*MockScheduler)(nil).RecurringEvents), ctx)
}

// ScheduleEvent mocks base method.
func (m *MockScheduler) ScheduleEvent(ctx context.Context, event *core.ScheduledEvent) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ScheduleEvent", ctx, event)
	ret0, _ := ret[0].(error)
	return ret0
}

// ScheduleEvent indicates an expected call of ScheduleEvent.
func (mr *MockSchedulerMockRecorder) ScheduleEvent(ctx, event any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ScheduleEvent", reflect.TypeOf((*MockScheduler)(nil).ScheduleEvent), ctx, event)
}

// ScheduleRecurringEvent mocks base method.
func (m *MockScheduler) ScheduleRecurringEvent(ctx context.Context, event *core.ScheduledRecurringEvent) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ScheduleRecurringEvent", ctx, event)
	ret0, _ := ret[0].(error)
	return ret0
}

// ScheduleRecurringEvent indicates an expected call of ScheduleRecurringEvent.
func (mr *MockSchedulerMockRecorder) ScheduleRecurringEvent(ctx, event any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ScheduleRecurringEvent", reflect.TypeOf((*MockScheduler)(nil).ScheduleRecurringEvent), ctx, event)
}

// UpdateRecurringEventTimestamp mocks base method.
func (m *MockScheduler) UpdateRecurringEventTimestamp(ctx context.Context, id string, stamp time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateRecurringEventTimestamp", ctx, id, stamp)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateRecurringEventTimestamp indicates an expected call of UpdateRecurringEventTimestamp.
func (mr *MockSchedulerMockRecorder) UpdateRecurringEventTimestamp(ctx, id, stamp any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateRecurringEventTimestamp", reflect.TypeOf((*MockScheduler)(nil).UpdateRecurringEventTimestamp), ctx, id, stamp)
}

// MockClientRegistry is a mock of ClientRegistry interface.
type MockClientRegistry struct {
	ctrl     *gomock.Controller
	recorder *MockClientRegistryMockRecorder
}

// MockClientRegistryMockRecorder is the mock recorder for MockClientRegistry.
type MockClientRegistryMockRecorder struct {
	mock *MockClientRegistry
}

// NewMockClientRegistry creates a new mock instance.
func NewMockClientRegistry(ctrl *gomock.Controller) *MockClientRegistry {
	mock := &MockClientRegistry{ctrl: ctrl}
	mock.recorder = &MockClientRegistryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClientRegistry) EXPECT() *MockClientRegistryMockRecorder {
	return m.recorder
}

// Client mocks base method.
func (m *MockClientRegistry) Client(ctx context.Context, source, clientId string, infoParser func(*core.Client, string) error) (core.Client, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Client", ctx, source, clientId, infoParser)
	ret0, _ := ret[0].(core.Client)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Client indicates an expected call of Client.
func (mr *MockClientRegistryMockRecorder) Client(ctx, source, clientId, infoParser any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Client", reflect.TypeOf((*MockClientRegistry)(nil).Client), ctx, source, clientId, infoParser)
}

// ClientsForUser mocks base method.
func (m *MockClientRegistry) ClientsForUser(ctx context.Context, userId string, infoParser func(*core.Client, string) error) ([]core.Client, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClientsForUser", ctx, userId, infoParser)
	ret0, _ := ret[0].([]core.Client)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ClientsForUser indicates an expected call of ClientsForUser.
func (mr *MockClientRegistryMockRecorder) ClientsForUser(ctx, userId, infoParser any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClientsForUser", reflect.TypeOf((*MockClientRegistry)(nil).ClientsForUser), ctx, userId, infoParser)
}

// UpsertClient mocks base method.
func (m *MockClientRegistry) UpsertClient(ctx context.Context, source, clientId string, info any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertClient", ctx, source, clientId, info)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertClient indicates an expected call of UpsertClient.
func (mr *MockClientRegistryMockRecorder) UpsertClient(ctx, source, clientId, info any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertClient", reflect.TypeOf((*MockClientRegistry)(nil).UpsertClient), ctx, source, clientId, info)
}

// UserForClient mocks base method.
func (m *MockClientRegistry) UserForClient(ctx context.Context, source, clientId string) (core.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserForClient", ctx, source, clientId)
	ret0, _ := ret[0].(core.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserForClient indicates an expected call of UserForClient.
func (mr *MockClientRegistryMockRecorder) UserForClient(ctx, source, clientId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserForClient", reflect.TypeOf((*MockClientRegistry)(nil).UserForClient), ctx, source, clientId)
}
