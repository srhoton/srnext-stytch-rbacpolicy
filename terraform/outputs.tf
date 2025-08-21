output "lambda_function_arn" {
  description = "ARN of the Lambda function"
  value       = aws_lambda_function.rbac_policy_lambda.arn
}

output "lambda_function_name" {
  description = "Name of the Lambda function"
  value       = aws_lambda_function.rbac_policy_lambda.function_name
}

output "target_group_arn" {
  description = "ARN of the ALB target group"
  value       = aws_lb_target_group.lambda_target_group.arn
}

output "cloudwatch_log_group_name" {
  description = "Name of the CloudWatch log group"
  value       = aws_cloudwatch_log_group.lambda_logs.name
}

output "endpoint_url" {
  description = "HTTPS endpoint URL for the RBAC policy API"
  value       = "https://${var.domain_name}${local.path_pattern}"
}

output "security_group_id" {
  description = "ID of the Lambda security group"
  value       = aws_security_group.lambda_sg.id
}

output "iam_role_arn" {
  description = "ARN of the Lambda execution IAM role"
  value       = aws_iam_role.lambda_execution_role.arn
}