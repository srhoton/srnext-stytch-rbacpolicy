package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stytchauth/stytch-management-go/v2/pkg/models/rbacpolicy"
	"go.uber.org/zap"
)

type RBACPolicyClient interface {
	Get(ctx context.Context, body rbacpolicy.GetRequest) (*rbacpolicy.GetResponse, error)
	Set(ctx context.Context, body rbacpolicy.SetRequest) (*rbacpolicy.SetResponse, error)
}

type Handler struct {
	client    RBACPolicyClient
	projectID string
	logger    *zap.Logger
}

func NewHandler(client RBACPolicyClient, projectID string, logger *zap.Logger) *Handler {
	return &Handler{
		client:    client,
		projectID: projectID,
		logger:    logger,
	}
}

func (h *Handler) HandleRequest(ctx context.Context, request events.ALBTargetGroupRequest) (events.ALBTargetGroupResponse, error) {
	h.logger.Info("Processing request",
		zap.String("method", request.HTTPMethod),
		zap.String("path", request.Path),
	)

	// Handle health check endpoint
	if (request.Path == "/health" || request.Path == "/rbacpolicy/health") && request.HTTPMethod == http.MethodGet {
		return h.handleHealthCheck()
	}

	switch request.HTTPMethod {
	case http.MethodGet:
		return h.handleGet(ctx)
	case http.MethodPut:
		return h.handlePut(ctx, request.Body)
	case http.MethodPost:
		return h.handlePut(ctx, request.Body)
	case http.MethodDelete:
		return h.handleDelete(ctx)
	default:
		return h.errorResponse(http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *Handler) handleGet(ctx context.Context) (events.ALBTargetGroupResponse, error) {
	req := rbacpolicy.GetRequest{
		ProjectID: h.projectID,
	}

	resp, err := h.client.Get(ctx, req)
	if err != nil {
		h.logger.Error("Failed to get RBAC policy", zap.Error(err))
		return h.errorResponse(http.StatusInternalServerError, fmt.Sprintf("Failed to get RBAC policy: %v", err))
	}

	body, err := json.Marshal(resp.Policy)
	if err != nil {
		h.logger.Error("Failed to marshal response", zap.Error(err))
		return h.errorResponse(http.StatusInternalServerError, "Failed to marshal response")
	}

	return events.ALBTargetGroupResponse{
		StatusCode:        http.StatusOK,
		StatusDescription: http.StatusText(http.StatusOK),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body:            string(body),
		IsBase64Encoded: false,
	}, nil
}

func (h *Handler) handlePut(ctx context.Context, body string) (events.ALBTargetGroupResponse, error) {
	var policy rbacpolicy.Policy
	if err := json.Unmarshal([]byte(body), &policy); err != nil {
		h.logger.Error("Failed to unmarshal request body", zap.Error(err))
		return h.errorResponse(http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
	}

	req := rbacpolicy.SetRequest{
		ProjectID: h.projectID,
		Policy:    policy,
	}

	resp, err := h.client.Set(ctx, req)
	if err != nil {
		h.logger.Error("Failed to set RBAC policy", zap.Error(err))
		return h.errorResponse(http.StatusInternalServerError, fmt.Sprintf("Failed to set RBAC policy: %v", err))
	}

	responseBody, err := json.Marshal(resp.Policy)
	if err != nil {
		h.logger.Error("Failed to marshal response", zap.Error(err))
		return h.errorResponse(http.StatusInternalServerError, "Failed to marshal response")
	}

	return events.ALBTargetGroupResponse{
		StatusCode:        http.StatusOK,
		StatusDescription: http.StatusText(http.StatusOK),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body:            string(responseBody),
		IsBase64Encoded: false,
	}, nil
}

func (h *Handler) handleDelete(ctx context.Context) (events.ALBTargetGroupResponse, error) {
	// First get the current policy to preserve stytch roles
	getReq := rbacpolicy.GetRequest{
		ProjectID: h.projectID,
	}
	
	getResp, err := h.client.Get(ctx, getReq)
	if err != nil {
		h.logger.Error("Failed to get current RBAC policy", zap.Error(err))
		return h.errorResponse(http.StatusInternalServerError, fmt.Sprintf("Failed to get current RBAC policy: %v", err))
	}
	
	// Clear only custom roles and resources, preserve stytch roles
	clearedPolicy := rbacpolicy.Policy{
		StytchMember:    getResp.Policy.StytchMember,
		StytchAdmin:     getResp.Policy.StytchAdmin,
		StytchResources: getResp.Policy.StytchResources,
		CustomRoles:     []rbacpolicy.Role{},
		CustomResources: []rbacpolicy.Resource{},
	}

	req := rbacpolicy.SetRequest{
		ProjectID: h.projectID,
		Policy:    clearedPolicy,
	}

	_, err = h.client.Set(ctx, req)
	if err != nil {
		h.logger.Error("Failed to clear RBAC policy", zap.Error(err))
		return h.errorResponse(http.StatusInternalServerError, fmt.Sprintf("Failed to clear RBAC policy: %v", err))
	}

	return events.ALBTargetGroupResponse{
		StatusCode:        http.StatusNoContent,
		StatusDescription: http.StatusText(http.StatusNoContent),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body:            "",
		IsBase64Encoded: false,
	}, nil
}

func (h *Handler) handleHealthCheck() (events.ALBTargetGroupResponse, error) {
	healthResponse := map[string]string{
		"status": "healthy",
	}

	body, _ := json.Marshal(healthResponse)

	return events.ALBTargetGroupResponse{
		StatusCode:        http.StatusOK,
		StatusDescription: http.StatusText(http.StatusOK),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body:            string(body),
		IsBase64Encoded: false,
	}, nil
}

func (h *Handler) errorResponse(statusCode int, message string) (events.ALBTargetGroupResponse, error) {
	errorBody := map[string]string{
		"error": message,
	}

	body, _ := json.Marshal(errorBody)

	return events.ALBTargetGroupResponse{
		StatusCode:        statusCode,
		StatusDescription: http.StatusText(statusCode),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body:            string(body),
		IsBase64Encoded: false,
	}, nil
}
