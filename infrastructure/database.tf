resource "aws_security_group" "rds-sg" {
  name = "rds-sg"

  description = "RDS (terraform-managed)"
  vpc_id      = aws_vpc.main-vpc.id

  # Only MySQL in
  ingress {
    from_port   = tonumber(var.rds_db_port)
    to_port     = tonumber(var.rds_db_port)
    protocol    = "tcp"
    cidr_blocks = ["10.0.0.0/24"]
  }

  # Allow all outbound traffic.
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
  }
}

resource "aws_db_subnet_group" "rds_subnet_group" {
  name       = "db_subnet_group"
  subnet_ids = [aws_subnet.subnet-private-1.id, aws_subnet.subnet-private-2.id]

  tags = {
    Name = "Shopify Challenge Subnet group"
  }
}
resource "random_string" "db_password" {
  length  = 16
  special = false
}
# RDS MariaDB database
resource "aws_db_instance" "database" {
  allocated_storage    = 10
  engine               = "mariadb"
  engine_version       = "10.5"
  instance_class       = var.database_instance_type
  name                 = var.rds_db_name
  username             = var.rds_db_user
  password             = random_string.db_password.result
  parameter_group_name = "default.mariadb10.5"
  port                 = tonumber(var.rds_db_port)
  multi_az             = false
  skip_final_snapshot  = true
  apply_immediately    = true
  db_subnet_group_name = aws_db_subnet_group.rds_subnet_group.name
  vpc_security_group_ids = [aws_security_group.rds-sg.id]
  tags = {
    Name = "Shopify Challenge RDS Instance"
  }
}