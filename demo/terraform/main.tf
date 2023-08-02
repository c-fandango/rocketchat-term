provider "aws" {
  region = "eu-west-1"
}

resource "aws_vpc" "demo_vpc" {
  cidr_block = "10.0.0.0/16"
  tags = {
    Name = "demo vpc"
  }
}

resource "aws_subnet" "demo_subnet" {
  vpc_id            = aws_vpc.demo_vpc.id
  cidr_block        = "10.0.0.0/16"
  availability_zone = "eu-west-1a"
  tags = {
    Name = "demo subnet"
  }
}

resource "aws_security_group" "demo_sg" {
  name_prefix = "demo-sg"
  vpc_id      = aws_vpc.demo_vpc.id

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["81.101.7.36/32"]
  }

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["81.101.7.36/32"]
  }

  egress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 53
    to_port     = 53
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 53
    to_port     = 53
    protocol    = "udp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_internet_gateway" "demo_igw" {
  vpc_id = aws_vpc.demo_vpc.id
}

resource "aws_route_table" "demo_route_table" {
  vpc_id = aws_vpc.demo_vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.demo_igw.id
  }
}

resource "aws_route_table_association" "demo_association" {
  subnet_id      = aws_subnet.demo_subnet.id
  route_table_id = aws_route_table.demo_route_table.id
}

resource "aws_instance" "demo_instance" {
  ami                         = "ami-01dd271720c1ba44f"
  instance_type               = "t2.micro"
  subnet_id                   = aws_subnet.demo_subnet.id
  vpc_security_group_ids      = [aws_security_group.demo_sg.id]
  key_name                    = "demo"
  associate_public_ip_address = true
  tags = {
    Name = "demo"
  }
}

resource "aws_ec2_instance_state" "demo_state" {
  instance_id = aws_instance.demo_instance.id
  state       = "stopped"
}

