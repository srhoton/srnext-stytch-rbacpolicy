locals {
  # Resource naming (shortened for AWS limits)
  resource_prefix = "stytch-rbac-${var.environment}"
  lambda_name     = "${local.resource_prefix}-lambda"

  # Common tags
  common_tags = merge(
    {
      Environment = var.environment
      Project     = var.project_name
      ManagedBy   = "Terraform"
      Owner       = "srnext"
      CreatedDate = timestamp()
    },
    var.tags
  )

  # Lambda configuration
  lambda_source_dir  = "${path.module}/../lambda"
  lambda_binary_path = "${path.module}/builds/bootstrap"
  lambda_zip_path    = "${path.module}/builds/lambda-deployment.zip"

  # ALB configuration
  alb_listener_port = 443
  target_group_port = 80
  health_check_path = "/health"
  path_pattern      = "/rbacpolicy/*"
}