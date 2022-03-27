terraform {
    backend "s3" {
        bucket = "terraform-bucket-ex2"
        key = "./terraform.tfstate" 
        region = "us-east-1"
    }
}
