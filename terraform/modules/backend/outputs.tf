output "s3_bucket_name" {
  description = "The name of the S3 bucket to use as a Terraform backend."
  value       = aws_s3_bucket.this.id
}
