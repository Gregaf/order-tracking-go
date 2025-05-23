locals {
  bucket_name = "${lower(replace(var.project_name, "_", "-"))}-tfstate-${var.environment}"
}
