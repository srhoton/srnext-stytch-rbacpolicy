resource "aws_route53_record" "rbac_policy_dns" {
  zone_id = data.aws_route53_zone.base_zone.zone_id
  name    = var.domain_name
  type    = "A"

  alias {
    name                   = data.aws_lb.existing_alb.dns_name
    zone_id                = data.aws_lb.existing_alb.zone_id
    evaluate_target_health = true
  }
}