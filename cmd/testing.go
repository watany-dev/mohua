package cmd

import (
	"context"
	"mohua/internal/sagemaker"
	"github.com/stretchr/testify/mock"
)

// MockSageMakerClient is a mock implementation of the sagemaker.Client interface
type MockSageMakerClient struct {
	mock.Mock
}

func (m *MockSageMakerClient) ValidateConfiguration(ctx context.Context) (bool, error) {
	args := m.Called(ctx)
	return args.Bool(0), args.Error(1)
}

func (m *MockSageMakerClient) ListEndpoints(ctx context.Context) ([]sagemaker.ResourceInfo, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]sagemaker.ResourceInfo), args.Error(1)
}

func (m *MockSageMakerClient) ListNotebooks(ctx context.Context) ([]sagemaker.ResourceInfo, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]sagemaker.ResourceInfo), args.Error(1)
}

func (m *MockSageMakerClient) ListStudioApps(ctx context.Context) ([]sagemaker.ResourceInfo, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]sagemaker.ResourceInfo), args.Error(1)
}

func (m *MockSageMakerClient) GetRegion() string {
	args := m.Called()
	return args.String(0)
}
