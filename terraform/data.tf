data "aws_lb" "existing_alb" {
  arn = var.alb_arn
}

data "aws_lb_listener" "https_listener" {
  load_balancer_arn = var.alb_arn
  port              = local.alb_listener_port
}

data "aws_vpc" "alb_vpc" {
  id = data.aws_lb.existing_alb.vpc_id
}

data "aws_subnets" "private_subnets" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.alb_vpc.id]
  }

  filter {
    name   = "tag:Name"
    values = ["*private*"]
  }
}

data "aws_route53_zone" "base_zone" {
  name         = var.base_domain
  private_zone = false
}

data "aws_secretsmanager_secret" "stytch_credentials" {
  arn = var.stytch_credentials_secret_arn
}

data "aws_secretsmanager_secret_version" "stytch_credentials_current" {
  secret_id = data.aws_secretsmanager_secret.stytch_credentials.id
}