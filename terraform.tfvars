#--------------------------------------------------------------
# General
#--------------------------------------------------------------

name   = "app"
region = "us-west-2"

#--------------------------------------------------------------
# Network
#--------------------------------------------------------------

vpc_cidr = "10.0.0.0/16"

private_subnet_cidrs = {
  "a" = "10.0.1.0/24",
  "c" = "10.0.2.0/24",
}

#--------------------------------------------------------------
# RDS
#--------------------------------------------------------------

db             = "app-db"
db_name        = "my_db"
db_username    = "user"
db_port        = "3306"
engine         = "mysql"
engine_version = "8.0.39"
db_instance    = "db.t3.micro"

#--------------------------------------------------------------
# App Runner
#--------------------------------------------------------------
image_identifier_url    = "<IMAGE_IDENTIFIER_URL>"
base_url                = "<BASE_URL>"