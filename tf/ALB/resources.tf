### ALB for routing gRPC ###

### ALB
resource "aws_lb" "myecs" {
  name               = join("-", [var.base_name, "alb"])
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.myecs_alb.id]
  subnets            = data.aws_subnets.myecs_public.ids
}

### ALB Listener
resource "aws_lb_listener" "myecs" {
  load_balancer_arn = aws_lb.myecs.arn
  port              = "443"
  protocol          = "HTTPS"
  certificate_arn   = data.aws_acm_certificate.myecs.arn

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.myecs.arn
  }
}

# ALB Target Group
resource "aws_lb_target_group" "myecs" {
  name = join("-", [var.base_name, "tg"])

  protocol         = "HTTP"
  protocol_version = "GRPC"
  port             = 8080

  vpc_id      = data.aws_vpc.myecs.id
  target_type = "ip"

  lifecycle {
    create_before_destroy = true
  }
}

### ALB Security Group
resource "aws_security_group" "myecs_alb" {
  name   = join("-", [var.base_name, "alb", "sg"])
  vpc_id = data.aws_vpc.myecs.id

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
