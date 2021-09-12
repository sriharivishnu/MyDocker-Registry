variable "environment" {
    default = "development"
}

variable "aws_region" {
    default = "us-east-1"
}

variable "database_instance_type" {
    default = "db.t2.micro"
}

variable "eb_app_name" {
    type = string
}

variable "eb_app_description" {
    default = "My Shopify Challenge Application"
}

variable "eb_env_name" {
    type = string
}

variable "eb_solution_stack_name" {
    default = "64bit Amazon Linux 2 v3.4.0 running Go 1"
}

variable "eb_instance_port" {
    default = "5000"
}

variable "eb_instance_type" {
    default = "t2.micro"
}

variable "eb_asg_min_size" {
    default = 1
}

variable "eb_asg_max_size" {
    default = 2
}

# RDS

variable "rds_db_name" {
    default = "prod"
}
variable "rds_db_port" {
    default = "3306"
}
variable "rds_db_user" {
    type = string
}

# ENV Variables

variable "api_env_var_ENVIRONMENT" {
    default = "development"
}
variable "api_env_var_PORT" {
    default = "5000"
}

variable "api_env_var_S3_BUCKET_KEY" {
    type = string
}