resource "aws_s3_bucket" "my_local_bucket" {
  bucket = "my-local-test-bucket"

  tags = {
    Environment = "local"
    Terraform   = "true"
  }
}