package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stytchauth/stytch-management-go/v2/pkg/models/rbacpolicy"
	"go.uber.org/zap"
)

type mockRBACPolicyClient struct {
	getFunc func(ctx context.Context, body rbacpolicy.GetRequest) (*rbacpolicy.GetResponse, error)
	setFunc func(ctx context.Context, body rbacpolicy.SetRequest) (*rbacpolicy.SetResponse, error)
}

func (m *mockRBACPolicyClient) Get(ctx context.Context, body rbacpolicy.GetRequest) (*rbacpolicy.GetResponse, error) {
	if m.getFunc != nil {
		return m.getFunc(ctx, body)
	}
	return nil, errors.New("get not implemented")
}

func (m *mockRBACPolicyClient) Set(ctx context.Context, body rbacpolicy.SetRequest) (*rbacpolicy.SetResponse, error) {
	if m.setFunc != nil {
		return m.setFunc(ctx, body)
	}
	return nil, errors.New("set not implemented")
}

// mockUnmarshalableType is a type that causes json.Marshal to fail
type mockUnmarshalableType struct {
	Channel chan int `json:"channel"` // channels cannot be marshaled to JSON
}

func TestHandleGet(t *testing.T) {
	logger := zap.NewNop()
	projectID := "test-project-id"

	tests := []struct {
		name           string
		mockClient     *mockRBACPolicyClient
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "Successful GET",
			mockClient: &mockRBACPolicyClient{
				getFunc: func(ctx context.Context, body rbacpolicy.GetRequest) (*rbacpolicy.GetResponse, error) {
					return &rbacpolicy.GetResponse{
						StatusCode: 200,
						RequestID:  "req-123",
						Policy: rbacpolicy.Policy{
							CustomRoles: []rbacpolicy.Role{
								{
									RoleID:      "admin",
									Description: "Administrator role",
									Permissions: []rbacpolicy.Permission{
										{
											ResourceID: "documents",
											Actions:    []string{"read", "write", "delete"},
										},
									},
								},
							},
							CustomResources: []rbacpolicy.Resource{
								{
									ResourceID:       "documents",
									Description:      "Document resources",
									AvailableActions: []string{"read", "write", "delete"},
								},
							},
						},
					}, nil
				},
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name: "Failed GET - Client error",
			mockClient: &mockRBACPolicyClient{
				getFunc: func(ctx context.Context, body rbacpolicy.GetRequest) (*rbacpolicy.GetResponse, error) {
					return nil, errors.New("client error")
				},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewHandler(tt.mockClient, projectID, logger)

			request := events.ALBTargetGroupRequest{
				HTTPMethod: http.MethodGet,
				Path:       "/rbacpolicy",
			}

			response, err := handler.HandleRequest(context.Background(), request)

			if tt.expectedError && err == nil {
				t.Errorf("Expected error but got none")
			}

			if !tt.expectedError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if response.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, response.StatusCode)
			}

			if tt.expectedStatus == http.StatusOK {
				var policy rbacpolicy.Policy
				if err := json.Unmarshal([]byte(response.Body), &policy); err != nil {
					t.Errorf("Failed to unmarshal response body: %v", err)
				}
			}
		})
	}
}

func TestHandlePut(t *testing.T) {
	logger := zap.NewNop()
	projectID := "test-project-id"

	validPolicy := rbacpolicy.Policy{
		CustomRoles: []rbacpolicy.Role{
			{
				RoleID:      "editor",
				Description: "Editor role",
				Permissions: []rbacpolicy.Permission{
					{
						ResourceID: "documents",
						Actions:    []string{"read", "write"},
					},
				},
			},
		},
		CustomResources: []rbacpolicy.Resource{
			{
				ResourceID:       "documents",
				Description:      "Document resources",
				AvailableActions: []string{"read", "write", "delete"},
			},
		},
	}

	validPolicyJSON, _ := json.Marshal(validPolicy)

	tests := []struct {
		name           string
		body           string
		mockClient     *mockRBACPolicyClient
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "Successful PUT",
			body: string(validPolicyJSON),
			mockClient: &mockRBACPolicyClient{
				setFunc: func(ctx context.Context, body rbacpolicy.SetRequest) (*rbacpolicy.SetResponse, error) {
					return &rbacpolicy.SetResponse{
						StatusCode: 200,
						RequestID:  "req-456",
						Policy:     body.Policy,
					}, nil
				},
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "Invalid JSON body",
			body:           "invalid json",
			mockClient:     &mockRBACPolicyClient{},
			expectedStatus: http.StatusBadRequest,
			expectedError:  false,
		},
		{
			name: "Failed PUT - Client error",
			body: string(validPolicyJSON),
			mockClient: &mockRBACPolicyClient{
				setFunc: func(ctx context.Context, body rbacpolicy.SetRequest) (*rbacpolicy.SetResponse, error) {
					return nil, errors.New("client error")
				},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewHandler(tt.mockClient, projectID, logger)

			request := events.ALBTargetGroupRequest{
				HTTPMethod: http.MethodPut,
				Path:       "/rbacpolicy",
				Body:       tt.body,
			}

			response, err := handler.HandleRequest(context.Background(), request)

			if tt.expectedError && err == nil {
				t.Errorf("Expected error but got none")
			}

			if !tt.expectedError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if response.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, response.StatusCode)
			}
		})
	}
}

func TestHandlePost(t *testing.T) {
	logger := zap.NewNop()
	projectID := "test-project-id"

	validPolicy := rbacpolicy.Policy{
		CustomRoles: []rbacpolicy.Role{
			{
				RoleID:      "viewer",
				Description: "Viewer role",
				Permissions: []rbacpolicy.Permission{
					{
						ResourceID: "documents",
						Actions:    []string{"read"},
					},
				},
			},
		},
	}

	validPolicyJSON, _ := json.Marshal(validPolicy)

	mockClient := &mockRBACPolicyClient{
		setFunc: func(ctx context.Context, body rbacpolicy.SetRequest) (*rbacpolicy.SetResponse, error) {
			return &rbacpolicy.SetResponse{
				StatusCode: 200,
				RequestID:  "req-789",
				Policy:     body.Policy,
			}, nil
		},
	}

	handler := NewHandler(mockClient, projectID, logger)

	request := events.ALBTargetGroupRequest{
		HTTPMethod: http.MethodPost,
		Path:       "/rbacpolicy",
		Body:       string(validPolicyJSON),
	}

	response, err := handler.HandleRequest(context.Background(), request)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.StatusCode)
	}
}

func TestHandleDelete(t *testing.T) {
	logger := zap.NewNop()
	projectID := "test-project-id"

	tests := []struct {
		name           string
		mockClient     *mockRBACPolicyClient
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "Successful DELETE",
			mockClient: &mockRBACPolicyClient{
				setFunc: func(ctx context.Context, body rbacpolicy.SetRequest) (*rbacpolicy.SetResponse, error) {
					// Verify that an empty policy is being set
					if len(body.Policy.CustomRoles) != 0 || len(body.Policy.CustomResources) != 0 {
						t.Errorf("Expected empty policy for delete operation")
					}
					return &rbacpolicy.SetResponse{
						StatusCode: 200,
						RequestID:  "req-delete",
						Policy:     body.Policy,
					}, nil
				},
			},
			expectedStatus: http.StatusNoContent,
			expectedError:  false,
		},
		{
			name: "Failed DELETE - Client error",
			mockClient: &mockRBACPolicyClient{
				setFunc: func(ctx context.Context, body rbacpolicy.SetRequest) (*rbacpolicy.SetResponse, error) {
					return nil, errors.New("client error")
				},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewHandler(tt.mockClient, projectID, logger)

			request := events.ALBTargetGroupRequest{
				HTTPMethod: http.MethodDelete,
				Path:       "/rbacpolicy",
			}

			response, err := handler.HandleRequest(context.Background(), request)

			if tt.expectedError && err == nil {
				t.Errorf("Expected error but got none")
			}

			if !tt.expectedError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if response.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, response.StatusCode)
			}
		})
	}
}

func TestHandleUnsupportedMethod(t *testing.T) {
	logger := zap.NewNop()
	projectID := "test-project-id"
	mockClient := &mockRBACPolicyClient{}

	handler := NewHandler(mockClient, projectID, logger)

	request := events.ALBTargetGroupRequest{
		HTTPMethod: http.MethodPatch,
		Path:       "/rbacpolicy",
	}

	response, err := handler.HandleRequest(context.Background(), request)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if response.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected status code %d, got %d", http.StatusMethodNotAllowed, response.StatusCode)
	}
}

func TestErrorResponse(t *testing.T) {
	logger := zap.NewNop()
	projectID := "test-project-id"
	mockClient := &mockRBACPolicyClient{}

	handler := NewHandler(mockClient, projectID, logger)

	tests := []struct {
		name           string
		statusCode     int
		message        string
		expectedStatus int
	}{
		{
			name:           "Bad Request Error",
			statusCode:     http.StatusBadRequest,
			message:        "Invalid input",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Internal Server Error",
			statusCode:     http.StatusInternalServerError,
			message:        "Server error",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Not Found Error",
			statusCode:     http.StatusNotFound,
			message:        "Resource not found",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := handler.errorResponse(tt.statusCode, tt.message)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if response.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, response.StatusCode)
			}

			var errorBody map[string]string
			if err := json.Unmarshal([]byte(response.Body), &errorBody); err != nil {
				t.Errorf("Failed to unmarshal error response body: %v", err)
			}

			if errorBody["error"] != tt.message {
				t.Errorf("Expected error message '%s', got '%s'", tt.message, errorBody["error"])
			}
		})
	}
}
