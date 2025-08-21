resource "aws_security_group" "lambda_sg" {
  name        = "${local.lambda_name}-sg"
  description = "Security group for ${local.lambda_name} Lambda function"
  vpc_id      = data.aws_vpc.alb_vpc.id

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
    description = "Allow all outbound traffic"
  }

  tags = merge(local.common_tags, {
    Name = "${local.lambda_name}-sg"
  })
}

resource "aws_security_group_rule" "alb_to_lambda" {
  type                     = "ingress"
  from_port                = local.target_group_port
  to_port                  = local.target_group_port
  protocol                 = "tcp"
  security_group_id        = aws_security_group.lambda_sg.id
  source_security_group_id = tolist(data.aws_lb.existing_alb.security_groups)[0]
  description              = "Allow traffic from ALB"
}