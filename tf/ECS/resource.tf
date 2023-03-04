resource "aws_ecs_task_definition" "myecs" {
  family                   = join("-", [var.base_name, "task", "definition"])
  requires_compatibilities = ["FARGATE"]

  network_mode = "awsvpc"
  cpu          = 256
  memory       = 512

  container_definitions = jsonencode([
    {
      name      = "gRPC-server"
      image     = "${data.aws_ecr_repository.myecs.repository_url}:${var.image_tag}"
      essential = true
      portMappings = [
        {
          containerPort = 8080
          hostPort      = 8080
        }
      ]
      logConfiguration = {
        logDriver = "awsfirelens"
        options = {
          Name              = "cloudwatch"
          region            = var.region
          log_group_name    = join("/", ["ecs", var.base_name])
          log_stream_prefix = "grpc"
        }
      }
    },
    {
      name      = "log-router"
      image     = "public.ecr.aws/aws-observability/aws-for-fluent-bit:stable"
      essential = true

      firelensConfiguration = {
        type = "fluentbit"
        options = {
          enable-ecs-log-metadata = "true"
          config-file-type        = "file"
          config-file-value       = "/fluent-bit/configs/parse-json.conf"
        }
      }
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          awslogs-region        = var.region
          awslogs-group         = join("/", ["ecs", var.base_name])
          awslogs-stream-prefix = "logger"
        }
      }
    }
  ])

  execution_role_arn = aws_iam_role.myecs_task_execution_role.arn
  task_role_arn      = aws_iam_role.myecs_task_role.arn
}

resource "aws_ecs_service" "myecs" {
  name    = join("-", [var.base_name, "service"])
  cluster = aws_ecs_cluster.myecs.id

  task_definition = aws_ecs_task_definition.myecs.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  depends_on = [aws_lb_listener.myecs]

  load_balancer {
    target_group_arn = aws_lb_target_group.myecs.arn
    container_name   = "gRPC-server"
    container_port   = 8080
  }

  network_configuration {
    subnets          = data.aws_subnets.myecs_private.ids
    security_groups  = [aws_security_group.myecs_service.id]
    assign_public_ip = false
  }
}

resource "aws_security_group" "myecs_service" {
  name   = join("-", [var.base_name, "service", "sg"])
  vpc_id = data.aws_vpc.myecs.id

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port       = 8080
    to_port         = 8080
    protocol        = "tcp"
    security_groups = [aws_security_group.myecs_alb.id]
  }
}

resource "aws_ecs_cluster" "myecs" {
  name = join("-", [var.base_name, "cluster"])
}

resource "aws_ecs_cluster_capacity_providers" "myecs" {
  cluster_name = aws_ecs_cluster.myecs.name

  capacity_providers = ["FARGATE"]

  default_capacity_provider_strategy {
    base              = 1
    weight            = 100
    capacity_provider = "FARGATE"
  }
}
