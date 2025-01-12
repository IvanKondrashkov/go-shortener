// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/service.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	url "net/url"
	reflect "reflect"

	models "github.com/IvanKondrashkov/go-shortener/internal/models"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockRepository) Close(ctx context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close", ctx)
}

// Close indicates an expected call of Close.
func (mr *MockRepositoryMockRecorder) Close(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockRepository)(nil).Close), ctx)
}

// GetByID mocks base method.
func (m *MockRepository) GetByID(ctx context.Context, id uuid.UUID) (*url.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*url.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockRepositoryMockRecorder) GetByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockRepository)(nil).GetByID), ctx, id)
}

// Load mocks base method.
func (m *MockRepository) Load(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Load", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Load indicates an expected call of Load.
func (mr *MockRepositoryMockRecorder) Load(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Load", reflect.TypeOf((*MockRepository)(nil).Load), ctx)
}

// Ping mocks base method.
func (m *MockRepository) Ping(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockRepositoryMockRecorder) Ping(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockRepository)(nil).Ping), ctx)
}

// Save mocks base method.
func (m *MockRepository) Save(ctx context.Context, id uuid.UUID, url *url.URL) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, id, url)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Save indicates an expected call of Save.
func (mr *MockRepositoryMockRecorder) Save(ctx, id, url interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockRepository)(nil).Save), ctx, id, url)
}

// SaveBatch mocks base method.
func (m *MockRepository) SaveBatch(ctx context.Context, batch []*models.RequestShortenAPIBatch) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveBatch", ctx, batch)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveBatch indicates an expected call of SaveBatch.
func (mr *MockRepositoryMockRecorder) SaveBatch(ctx, batch interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveBatch", reflect.TypeOf((*MockRepository)(nil).SaveBatch), ctx, batch)
}
