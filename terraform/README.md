# Terraform Infrastructure for Stytch RBAC Policy Lambda

This Terraform configuration deploys the Stytch RBAC Policy Lambda function and associated infrastructure to AWS.

## Infrastructure Components

- **Lambda Function**: ARM64-based Lambda running on `provided.al2023` runtime
- **ALB Integration**: Target group and listener rules for existing ALB
- **VPC Configuration**: Lambda deployed in private subnets with appropriate security groups
- **IAM Roles**: Execution role with permissions for Secrets Manager and VPC access
- **CloudWatch Logs**: Log group with configurable retention
- **Route53**: DNS record for the RBAC policy endpoint

## Prerequisites

- Terraform >= 1.5.0
- AWS CLI configured with appropriate credentials
- Go 1.23+ for building Lambda
- Existing ALB with HTTPS listener on port 443
- Existing Secrets Manager secret with Stytch credentials
- Existing Route53 hosted zone

## Configuration

### Backend Configuration

The Terraform state is stored in S3:
- Bucket: `steve-rhoton-tfstate`
- Key: `srnext-stytch-rbacpolicy/terraform.tfstate`
- Region: `us-west-2`

### Important Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `aws_region` | AWS region for deployment | `us-west-2` |
| `environment` | Environment name | `sandbox` |
| `alb_arn` | ARN of existing ALB | Configured |
| `domain_name` | FQDN for the endpoint | `srnext-stytch-rbac-policy.sb.int.fullbayapi.com` |
| `stytch_credentials_secret_arn` | Secrets Manager ARN | Configured |

## Deployment

1. Initialize Terraform:
```bash
terraform init
```

2. Review the plan:
```bash
terraform plan
```

3. Apply the configuration:
```bash
terraform apply
```

## Endpoints

After deployment, the Lambda will be accessible at:
- Base URL: `https://srnext-stytch-rbac-policy.sb.int.fullbayapi.com`
- Health Check: `https://srnext-stytch-rbac-policy.sb.int.fullbayapi.com/health`
- RBAC Policy: `https://srnext-stytch-rbac-policy.sb.int.fullbayapi.com/rbacpolicy/*`

## API Operations

- `GET /rbacpolicy` - Retrieve RBAC policy
- `PUT /rbacpolicy` - Create/Update RBAC policy
- `POST /rbacpolicy` - Create/Update RBAC policy
- `DELETE /rbacpolicy` - Clear RBAC policy
- `GET /health` - Health check endpoint

## Security

- Lambda runs in private subnets
- Credentials stored in AWS Secrets Manager
- Security group restricts ingress to ALB only
- IAM role follows least privilege principle

## Monitoring

- CloudWatch Logs: `/aws/lambda/srnext-stytch-rbacpolicy-sandbox-lambda`
- Log retention: 7 days (configurable)

## Outputs

After deployment, Terraform will output:
- Lambda function ARN
- Target group ARN
- CloudWatch log group name
- Endpoint URL
- Security group ID
- IAM role ARN