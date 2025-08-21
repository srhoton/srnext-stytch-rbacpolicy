resource "aws_cloudwatch_log_group" "lambda_logs" {
  name              = "/aws/lambda/${local.lambda_name}"
  retention_in_days = var.cloudwatch_logs_retention_days

  tags = merge(local.common_tags, {
    Name = "${local.lambda_name}-logs"
  })
}

resource "aws_lambda_function" "rbac_policy_lambda" {
  filename         = data.archive_file.lambda_zip.output_path
  function_name    = local.lambda_name
  role             = aws_iam_role.lambda_execution_role.arn
  handler          = "bootstrap"
  runtime          = "provided.al2023"
  architectures    = ["arm64"]
  source_code_hash = data.archive_file.lambda_zip.output_base64sha256
  memory_size      = var.lambda_memory_size
  timeout          = var.lambda_timeout

  environment {
    variables = {
      ENVIRONMENT                 = var.environment
      STYTCH_WORKSPACE_KEY_ID     = jsondecode(data.aws_secretsmanager_secret_version.stytch_credentials_current.secret_string)["workspace_key"]
      STYTCH_WORKSPACE_KEY_SECRET = jsondecode(data.aws_secretsmanager_secret_version.stytch_credentials_current.secret_string)["workspace_secret"]
      STYTCH_PROJECT_ID           = "project-test-478debed-30da-42a5-9216-97240a34bd1f"
    }
  }

  vpc_config {
    subnet_ids         = data.aws_subnets.private_subnets.ids
    security_group_ids = [aws_security_group.lambda_sg.id]
  }

  depends_on = [
    aws_cloudwatch_log_group.lambda_logs,
    aws_iam_role_policy_attachment.lambda_basic_execution,
    aws_iam_role_policy_attachment.lambda_vpc_execution,
    null_resource.lambda_build
  ]

  tags = merge(local.common_tags, {
    Name = local.lambda_name
  })
}

resource "aws_lambda_permission" "alb_invoke" {
  statement_id  = "AllowExecutionFromALB"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.rbac_policy_lambda.function_name
  principal     = "elasticloadbalancing.amazonaws.com"
  source_arn    = aws_lb_target_group.lambda_target_group.arn
}