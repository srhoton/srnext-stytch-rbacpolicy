resource "null_resource" "lambda_build" {
  triggers = {
    source_hash = filebase64sha256("${local.lambda_source_dir}/go.mod")
    build_hash  = filebase64sha256("${local.lambda_source_dir}/cmd/lambda/main.go")
  }

  provisioner "local-exec" {
    command = <<-EOT
      mkdir -p ${path.module}/builds
      cd ${local.lambda_source_dir}
      GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags lambda.norpc -o ${abspath(local.lambda_binary_path)} ./cmd/lambda
    EOT
  }
}

data "archive_file" "lambda_zip" {
  type        = "zip"
  source_file = local.lambda_binary_path
  output_path = local.lambda_zip_path

  depends_on = [null_resource.lambda_build]
}