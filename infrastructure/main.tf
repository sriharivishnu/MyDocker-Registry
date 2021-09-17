# Elastic beanstalk application
resource "aws_elastic_beanstalk_application" "server" {
    name        = var.eb_app_name
    description = var.eb_app_description
}

resource "aws_security_group" "server-sg" {
  name = "server-sg"

  description = "Server (terraform-managed)"
  vpc_id      = aws_vpc.main-vpc.id

  # Allow all outbound traffic.
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
  }
}

resource "aws_security_group_rule" "server-to-rds" {
    type = "ingress"
    from_port   = tonumber(var.rds_db_port)
    to_port     = tonumber(var.rds_db_port)
    protocol    = "tcp"
    security_group_id = aws_security_group.server-sg.id
    source_security_group_id = aws_security_group.rds-sg.id
}

resource "random_string" "jwt_password" {
  length  = 16
  special = true
}

# Elastic beanstalk environment
resource "aws_elastic_beanstalk_environment" "serverEnvironment" {
    name                = var.eb_env_name
    application         = aws_elastic_beanstalk_application.server.name
    solution_stack_name = var.eb_solution_stack_name
    setting {
        namespace = "aws:ec2:vpc"
        name      = "VPCId"
        value     = aws_vpc.main-vpc.id
    }
    setting {
        namespace = "aws:ec2:vpc"
        name      = "Subnets"
        value     = join(",", [aws_subnet.subnet-public-1.id, aws_subnet.subnet-public-2.id])
    }
    setting {
        namespace = "aws:ec2:vpc"
        name      = "ELBScheme"
        value     = "internet facing"
    }
    setting {
        namespace = "aws:elasticbeanstalk:environment"
        name      = "LoadBalancerType"
        value     = "application"
    }
    setting {
        namespace = "aws:elasticbeanstalk:environment:process:default"
        name      = "MatcherHTTPCode"
        value     = "200"
    }
    setting {
        namespace = "aws:elasticbeanstalk:environment:process:default"
        name      = "Port"
        value     = var.eb_instance_port
    }
    setting {
        namespace = "aws:autoscaling:launchconfiguration"
        name      = "InstanceType"
        value     = var.eb_instance_type
    }
    setting {
        namespace = "aws:autoscaling:asg"
        name      = "MinSize"
        value     = var.eb_asg_min_size
    }
    setting {
        namespace = "aws:autoscaling:asg"
        name      = "MaxSize"
        value     = var.eb_asg_max_size
    }
    setting {
      namespace = "aws:autoscaling:launchconfiguration"
      name = "IamInstanceProfile"
      value = "aws-elasticbeanstalk-ec2-role"
    }
    setting {
      namespace = "aws:autoscaling:launchconfiguration"
      name = "SecurityGroups"
      value = aws_security_group.server-sg.id
    }
    setting {
        namespace = "aws:elasticbeanstalk:healthreporting:system"
        name      = "SystemType"
        value     = "enhanced"
    }
    setting {
      namespace = "aws:elasticbeanstalk:cloudwatch:logs"
      name = "StreamLogs"
      value = "true"
    }

    # ENVIRONMENT VARS
    setting {
      namespace = "aws:elasticbeanstalk:application:environment"
      name      = "ENVIRONMENT"
      value     = var.api_env_var_ENVIRONMENT
    }
    setting {
      namespace = "aws:elasticbeanstalk:application:environment"
      name      = "PORT"
      value     = var.api_env_var_PORT
    }
    setting {
      namespace = "aws:elasticbeanstalk:application:environment"
      name      = "DATABASE_NAME"
      value     = var.rds_db_name
    }
    setting {
      namespace = "aws:elasticbeanstalk:application:environment"
      name      = "DATABASE_HOST"
      value     = aws_db_instance.database.address
    }
    setting {
      namespace = "aws:elasticbeanstalk:application:environment"
      name      = "DATABASE_PORT"
      value     = var.rds_db_port
    }
    setting {
      namespace = "aws:elasticbeanstalk:application:environment"
      name      = "DATABASE_USER"
      value     = var.rds_db_user
    }
    setting {
      namespace = "aws:elasticbeanstalk:application:environment"
      name      = "DATABASE_PASSWORD"
      value     = random_string.db_password.result
    }
    setting {
      namespace = "aws:elasticbeanstalk:application:environment"
      name      = "S3_BUCKET_KEY"
      value     = var.api_env_var_S3_BUCKET_KEY
    }
    setting {
      namespace = "aws:elasticbeanstalk:application:environment"
      name      = "AWS_ACCESS_KEY_ID"
      value     = aws_iam_access_key.server_access_key.id
    }
    setting {
      namespace = "aws:elasticbeanstalk:application:environment"
      name      = "AWS_SECRET_ACCESS_KEY"
      value     = aws_iam_access_key.server_access_key.secret
    }
    setting {
      namespace = "aws:elasticbeanstalk:application:environment"
      name      = "AWS_DEFAULT_REGION"
      value     = var.aws_region
    }
    setting {
      namespace = "aws:elasticbeanstalk:application:environment"
      name      = "JWT_SECRET"
      value     = random_string.jwt_password.result
    }

}

# S3 bucket

resource "aws_s3_bucket" "b" {
  bucket = var.api_env_var_S3_BUCKET_KEY
  acl    = "private"

  tags = {
    Name        = "My bucket for Shopify Challenge"
    Environment = var.api_env_var_ENVIRONMENT
  }
  cors_rule {
        allowed_headers = [
            "*",
        ]
        allowed_methods = [
            "GET",
            "PUT",
            "POST",
            "HEAD",
        ]
        allowed_origins = [
            "*",
        ]
        expose_headers  = []
        max_age_seconds = 3000
    }
}

resource "aws_s3_bucket_public_access_block" "b_block" {
  bucket = aws_s3_bucket.b.id

  block_public_acls   = true
  block_public_policy = true
}

resource "aws_iam_user" "server_access_user" {
  name = "shopify-server-access-user"
  force_destroy = true
}

resource "aws_iam_user_policy" "server_access_role" {
  name = "shopify-server-access-policy"
  user = aws_iam_user.server_access_user.name
  
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
        {
            "Sid" = "Stmt1631409231466"
            "Action" = [
                "s3:DeleteObject",
                "s3:GetObject",
                "s3:GetObjectVersion",
                "s3:ListBucket",
                "s3:ListBucketVersions",
                "s3:PutObject"
            ]
            "Effect" = "Allow"
            "Resource" = "arn:aws:s3:::srihari-shopify-challenge/*"
        },
    ]
  })
}

resource "aws_iam_access_key" "server_access_key" {
  user = aws_iam_user.server_access_user.name
}