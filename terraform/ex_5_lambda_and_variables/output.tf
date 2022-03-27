output "result" {
  description = "The lambda function return value"
  value       = resource.aws_lambda_invocation.lambda_ex.result
}
