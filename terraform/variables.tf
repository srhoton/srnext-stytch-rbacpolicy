variable "aws_region" {
  description = "AWS region for deployment"
  type        = string
  default     = "us-west-2"
}

variable "environment" {
  description = "Environment name (e.g., sandbox, dev, staging, prod)"
  type        = string
  default     = "sandbox"

  validation {
    condition     = contains(["sandbox", "dev", "staging", "prod"], var.environment)
    error_message = "Environment must be one of: sandbox, dev, staging, prod."
  }
}

variable "project_name" {
  description = "Project name for resource naming"
  type        = string
  default     = "srnext-stytch-rbacpolicy"
}

variable "alb_arn" {
  description = "ARN of the existing Application Load Balancer"
  type        = string
  default     = "arn:aws:elasticloadbalancing:us-west-2:345594586248:loadbalancer/app/external-private-alb/720e2b5474d3d602"

  validation {
    condition     = can(regex("^arn:aws:elasticloadbalancing:", var.alb_arn))
    error_message = "ALB ARN must be a valid ELB ARN."
  }
}

variable "domain_name" {
  description = "Domain name for the RBAC policy endpoint"
  type        = string
  default     = "srnext-stytch-rbac-policy.sb.int.fullbayapi.com"
}

variable "base_domain" {
  description = "Base domain for Route53 zone lookup"
  type        = string
  default     = "sb.int.fullbayapi.com"
}

variable "stytch_credentials_secret_arn" {
  description = "ARN of the Secrets Manager secret containing Stytch credentials"
  type        = string
  default     = "arn:aws:secretsmanager:us-west-2:345594586248:secret:srnext/stytchCredentials-BpFPtL"

  validation {
    condition     = can(regex("^arn:aws:secretsmanager:", var.stytch_credentials_secret_arn))
    error_message = "Secret ARN must be a valid Secrets Manager ARN."
  }
}

variable "lambda_memory_size" {
  description = "Memory allocation for Lambda function in MB"
  type        = number
  default     = 512

  validation {
    condition     = var.lambda_memory_size >= 128 && var.lambda_memory_size <= 10240
    error_message = "Lambda memory size must be between 128 and 10240 MB."
  }
}

variable "lambda_timeout" {
  description = "Lambda function timeout in seconds"
  type        = number
  default     = 30

  validation {
    condition     = var.lambda_timeout >= 1 && var.lambda_timeout <= 900
    error_message = "Lambda timeout must be between 1 and 900 seconds."
  }
}

variable "cloudwatch_logs_retention_days" {
  description = "CloudWatch Logs retention period in days"
  type        = number
  default     = 7

  validation {
    condition     = contains([1, 3, 5, 7, 14, 30, 60, 90, 120, 150, 180, 365, 400, 545, 731, 1827, 3653], var.cloudwatch_logs_retention_days)
    error_message = "CloudWatch Logs retention must be a valid AWS retention period."
  }
}

variable "tags" {
  description = "Additional tags to apply to resources"
  type        = map(string)
  default     = {}
}