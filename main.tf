# Variables
variable "region" {}

variable "name" {}

variable "vpc_cidr" {}

variable "private_subnet_cidrs" {
  type = map(string)
}

variable "db_name" {}

variable "db_username" {}

variable "db_port" {}

variable "engine" {}

variable "engine_version" {}

variable "db_instance" {}

variable "db_password" {
  description = "The password for the RDS instance"
  type        = string
  sensitive   = true
}

variable "image_identifier_url" {}

# Terraform configuration
terraform {
  required_version = "=1.9.7"
}

provider "aws" {
  region = var.region
}


resource "aws_vpc" "main" {
  cidr_block = var.vpc_cidr

  tags = {
    Name = "${var.name}-vpc"
  }
}


resource "aws_subnet" "private" {
  for_each = var.private_subnet_cidrs

  vpc_id            = aws_vpc.main.id
  cidr_block        = each.value
  availability_zone = "${var.region}${each.key}"

  tags = {
    Name = "${var.name}-private-subnet-${each.key}"
  }
}

resource "aws_security_group" "rds_security_group" {
  name        = "${var.name}-rds-sg"
  description = "Security group for RDS"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port   = var.db_port
    to_port     = var.db_port
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "${var.name}-rds-sg"
  }
}

resource "aws_db_subnet_group" "main" {
  name       = "${var.name}-subnet-group"
  subnet_ids = [for id in aws_subnet.private : id.id]

  tags = {
    Name = "${var.name}-subnet-group"
  }
}

resource "aws_db_instance" "rds" {
  allocated_storage           = 10
  storage_type                = "gp2"
  engine                      = var.engine
  engine_version              = var.engine_version
  instance_class              = var.db_instance
  identifier                  = var.db_name
  username                    = var.db_username
  password                    = var.db_password
  port                        = var.db_port

  vpc_security_group_ids      = [aws_security_group.rds_security_group.id]
  db_subnet_group_name        = aws_db_subnet_group.main.name
  skip_final_snapshot         = true  

  # Multi-AZ を無効化
  multi_az                    = false 

  tags = {
    Name = "${var.name}-rds-instance"
  }
}

output "rds_address" {
  description = "The address of the RDS instance"
  value       = aws_db_instance.rds.address
}

resource "aws_apprunner_vpc_connector" "main" {
  vpc_connector_name = "${var.name}-vpc-connector"
  subnets            = [
    aws_subnet.private["a"].id,
    aws_subnet.private["c"].id
  ]
  security_groups    = [aws_security_group.rds_security_group.id]
}

resource "aws_iam_role" "app_runner" {
  name = "${var.name}-app-runner-role"

  assume_role_policy = jsonencode({
    Statement = [
      {
        Action = "sts:AssumeRole",
        Effect = "Allow",
        Principal = {
          Service = "build.apprunner.amazonaws.com"
        }
      }
    ],
    Version = "2012-10-17"
  })
}

resource "aws_iam_role_policy" "app_runner_policy" {
  name = "${var.name}-app-runner-policy"
  role = aws_iam_role.app_runner.id

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Action = [
          "ecr:GetDownloadUrlForLayer",
          "ecr:BatchGetImage",
          "ecr:BatchCheckLayerAvailability"
        ],
        Resource = "*"
      },
      {
        Effect = "Allow",
        Action = "ecr:GetAuthorizationToken",
        Resource = "*"
      }
    ]
  })
}

resource "aws_apprunner_auto_scaling_configuration_version" "main" {                            
  auto_scaling_configuration_name = "main"
  min_size = 1
  max_size = 5
  max_concurrency = 50

  tags = {
    Name = "${var.name}-auto-scaling-configuration"
  }
}

resource "aws_apprunner_service" "main" {
  service_name = "${var.name}-service"
  auto_scaling_configuration_arn = aws_apprunner_auto_scaling_configuration_version.main.arn

  source_configuration {
    authentication_configuration {
      access_role_arn = aws_iam_role.app_runner.arn
    }
    auto_deployments_enabled = false

    image_repository {
      image_identifier      = var.image_identifier_url
      image_repository_type = "ECR"
      image_configuration {
        port = "8080"
        runtime_environment_variables = {
          DB_HOST     = aws_db_instance.rds.address
          DB_NAME     = "${var.db_name}"
          DB_USER     = var.db_username
          DB_PASSWORD = var.db_password
          DB_PORT     = var.db_port
        }
      }
    }
  }

  network_configuration {
    egress_configuration {
      egress_type       = "VPC"
      vpc_connector_arn = aws_apprunner_vpc_connector.main.arn
    }
  }

  instance_configuration {
    cpu    = "256"
    memory = "512"
  }

  tags = {
    Name = "${var.name}-service"
  }
}

output "service_url" {
  value = aws_apprunner_service.main.service_url
}
