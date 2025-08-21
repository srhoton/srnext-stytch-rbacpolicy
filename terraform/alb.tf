resource "aws_lb_target_group" "lambda_target_group" {
  name        = "${local.resource_prefix}-tg"
  target_type = "lambda"
  vpc_id      = data.aws_vpc.alb_vpc.id

  health_check {
    enabled             = true
    path                = local.health_check_path
    interval            = 30
    timeout             = 5
    healthy_threshold   = 2
    unhealthy_threshold = 2
    matcher             = "200"
  }

  tags = merge(local.common_tags, {
    Name = "${local.resource_prefix}-tg"
  })
}

resource "aws_lb_target_group_attachment" "lambda_attachment" {
  target_group_arn = aws_lb_target_group.lambda_target_group.arn
  target_id        = aws_lambda_function.rbac_policy_lambda.arn

  depends_on = [aws_lambda_permission.alb_invoke]
}

resource "aws_lb_listener_rule" "rbac_policy_rule" {
  listener_arn = data.aws_lb_listener.https_listener.arn
  priority     = 200

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.lambda_target_group.arn
  }

  condition {
    host_header {
      values = [var.domain_name]
    }
  }

  condition {
    path_pattern {
      values = [local.path_pattern]
    }
  }

  tags = merge(local.common_tags, {
    Name = "${local.resource_prefix}-listener-rule"
  })
}