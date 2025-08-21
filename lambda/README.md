# Stytch RBAC Policy Lambda

AWS Lambda function for managing Stytch RBAC policies through CRUD operations.

## Features

- Full CRUD operations for Stytch RBAC policies
- ALB event support
- ARM64 architecture optimized
- Comprehensive test coverage (>80%)
- Production-ready logging with Zap

## Environment Variables

Required environment variables:

- `STYTCH_WORKSPACE_KEY_ID`: Stytch workspace key ID
- `STYTCH_WORKSPACE_KEY_SECRET`: Stytch workspace key secret
- `STYTCH_PROJECT_ID`: Stytch project ID
- `ENVIRONMENT`: (Optional) Set to "production" for production logging

## API Endpoints

### GET /rbacpolicy
Retrieve the current RBAC policy.

**Response:**
```json
{
  "stytch_member": {...},
  "stytch_admin": {...},
  "stytch_resources": [...],
  "custom_roles": [...],
  "custom_resources": [...]
}
```

### PUT/POST /rbacpolicy
Create or update the RBAC policy.

**Request Body:**
```json
{
  "custom_roles": [
    {
      "role_id": "admin",
      "description": "Administrator role",
      "permissions": [
        {
          "resource_id": "documents",
          "actions": ["read", "write", "delete"]
        }
      ]
    }
  ],
  "custom_resources": [
    {
      "resource_id": "documents",
      "description": "Document resources",
      "available_actions": ["read", "write", "delete"]
    }
  ]
}
```

### DELETE /rbacpolicy
Clear the RBAC policy (sets an empty policy).

## Development

### Prerequisites

- Go 1.24
- golangci-lint
- Make

### Building

```bash
# Build for ARM64 Lambda
make build

# Run all checks and build
make all
```

### Testing

```bash
# Run tests
make test

# Run tests with coverage
make coverage

# Check coverage meets threshold (80%)
make coverage
```

### Code Quality

```bash
# Format code
make fmt

# Run linter
make lint

# Run go vet
make vet
```

### Deployment

```bash
# Create deployment package
make package
```

The deployment package will be created at `bin/lambda-deployment.zip`.

## Architecture

```
lambda/
├── cmd/
│   └── lambda/       # Main Lambda entry point
├── internal/
│   ├── config/       # Configuration management
│   └── handler/      # Request handlers
├── Makefile          # Build and test automation
└── go.mod            # Go module definition
```

## Testing

The project maintains >80% test coverage across all packages:

- Configuration validation tests
- Handler unit tests for all CRUD operations
- Mock client implementations for testing
- Error handling verification

## Security

- Uses environment variables for sensitive credentials
- No hardcoded secrets
- Comprehensive input validation
- Structured error responses without exposing internal details

## Performance

- Built for ARM64 architecture (Graviton2)
- Minimal dependencies
- Efficient JSON marshaling/unmarshaling
- Proper context handling for cancellation