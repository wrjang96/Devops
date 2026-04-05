# LocalStack 환경을 위한 AWS Provider 설정
provider "aws" {
  region                      = "us-east-1" # LocalStack은 보통 이 리전을 기본값으로 사용합니다.
  access_key                  = "test"      # LocalStack 기본 access key
  secret_key                  = "test"      # LocalStack 기본 secret key
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true
  s3_use_path_style           = true        # S3 엔드포인트 형식을 path-style로 강제

  # LocalStack 엔드포인트 오버라이드
  endpoints {
    apigateway = "http://localhost:4566"
    acm        = "http://localhost:4566"
    # 필요한 다른 서비스 엔드포인트 추가 가능
  }
}

# 변수 정의 (필요에 따라 수정)
variable "api_name" {
  description = "API Gateway 이름"
  type        = string
  default     = "my-local-api"
}

variable "stage_name" {
  description = "배포 스테이지 이름"
  type        = string
  default     = "dev"
}

variable "custom_domain_name" {
  description = "사용할 사용자 지정 도메인 이름"
  type        = string
  default     = "myapi.localstack" # 로컬 테스트용 예시 도메인
}

variable "resource_path" {
  description = "API 리소스 경로"
  type        = string
  default     = "echo"
}


# 1. API Gateway REST API 생성
resource "aws_api_gateway_rest_api" "my_api" {
  name        = var.api_name
  description = "Terraform으로 LocalStack에 배포된 API Gateway"

  # LocalStack에서 CORS 문제를 피하기 위한 기본 설정 (필요시 활성화)
  # endpoint_configuration {
  #   types = ["REGIONAL"]
  # }
}

# 수정 부분
# 2. API 리소스 생성 (/echo)
resource "aws_api_gateway_resource" "alpha_resource" {
  rest_api_id = aws_api_gateway_rest_api.my_api.id
  parent_id   = aws_api_gateway_rest_api.my_api.root_resource_id
  path_part   = "alpha"
}
# 수정 부분
resource "aws_api_gateway_resource" "beta_chris_resource" {
  rest_api_id = aws_api_gateway_rest_api.my_api.id
  parent_id   = aws_api_gateway_resource.alpha_resource.id # alpha 리소스 ID 참조
  path_part   = "beta.chris"                              # 경로의 마지막 부분
}

# 3. GET 메서드 생성
resource "aws_api_gateway_method" "get_method" {
  rest_api_id   = aws_api_gateway_rest_api.my_api.id
  resource_id   = aws_api_gateway_resource.beta_chris_resource.id
  http_method   = "GET"
  authorization = "NONE" # 인증 없음
}

# 4. Mock 통합 설정
resource "aws_api_gateway_integration" "mock_integration" {
  rest_api_id = aws_api_gateway_rest_api.my_api.id
  resource_id = aws_api_gateway_resource.beta_chris_resource.id
  http_method = aws_api_gateway_method.get_method.http_method
  type        = "MOCK" # Mock 통합 타입 지정

  # Mock 통합 요청 템플릿: API Gateway가 즉시 200 OK 응답을 생성하도록 지시
  request_templates = {
    "application/json" = jsonencode({
      statusCode = 200
    })
  }
}

# 5. 메서드 응답 설정 (200 OK)
resource "aws_api_gateway_method_response" "response_200" {
  rest_api_id = aws_api_gateway_rest_api.my_api.id
  resource_id = aws_api_gateway_resource.beta_chris_resource.id
  http_method = aws_api_gateway_method.get_method.http_method
  status_code = "200"

  # 응답 헤더 (CORS 등 필요시 추가)
  response_parameters = {
    "method.response.header.Access-Control-Allow-Origin" = true
  }

  # 응답 바디 모델 (비어있음)
  response_models = {
    "application/json" = "Empty"
  }
}

# 6. 통합 응답 설정 (Mock 응답 -> 메서드 응답 매핑)
resource "aws_api_gateway_integration_response" "integration_response" {
  rest_api_id = aws_api_gateway_rest_api.my_api.id
  resource_id = aws_api_gateway_resource.beta_chris_resource.id
  http_method = aws_api_gateway_method.get_method.http_method
  status_code = aws_api_gateway_method_response.response_200.status_code

  # 응답 헤더 매핑 (CORS 등 필요시 추가)
  response_parameters = {
    "method.response.header.Access-Control-Allow-Origin" = "'*'" # 모든 출처 허용 (테스트용)
  }

  # 응답 매핑 템플릿: 요청 파라미터($input.params)를 JSON으로 반환
  response_templates = {
    "application/json" = <<EOF
{
#   "message": "Received parameters via API Gateway Mock Integration",
  "v": "$input.params('V')",
  "name": null,
#   "headerParameters": $util.toJson($allParams.header)
}
EOF
  }

  depends_on = [aws_api_gateway_integration.mock_integration]
}

# 7. API 배포
resource "aws_api_gateway_deployment" "my_deployment" {
  rest_api_id = aws_api_gateway_rest_api.my_api.id

  # API 구성 변경 시 자동으로 재배포하기 위한 트리거 설정
  triggers = {
    redeployment = sha1(jsonencode([
      aws_api_gateway_resource.beta_chris_resource.id,
      aws_api_gateway_method.get_method.id,
      aws_api_gateway_integration.mock_integration.id,
      # 다른 리소스 ID 추가 가능
    ]))
  }

  lifecycle {
    create_before_destroy = true
  }

  depends_on = [
    aws_api_gateway_integration.mock_integration,
    aws_api_gateway_integration_response.integration_response
  ]
}

# 8. 스테이지 생성
resource "aws_api_gateway_stage" "my_stage" {
  deployment_id = aws_api_gateway_deployment.my_deployment.id
  rest_api_id   = aws_api_gateway_rest_api.my_api.id
  stage_name    = var.stage_name
}

# --- 사용자 지정 도메인 설정 ---
# 참고: LocalStack에서 사용자 지정 도메인과 ACM 인증서 처리는 불안정할 수 있습니다.
#       실제 AWS 환경과 동작이 다를 수 있으며, 추가적인 LocalStack 설정이 필요할 수 있습니다.

# LocalStack은 실제 ACM 검증을 수행하지 않으므로, 더미 ARN 또는 자체 서명 인증서를 사용할 수 있습니다.
# 여기서는 더미 ARN을 사용한다고 가정합니다. (실제로는 LocalStack 설정에 따라 다름)
# 예시: LocalStack 시작 시 DEBUG=1 환경 변수를 설정하면 ACM 관련 로그를 볼 수 있습니다.
# 또는, 로컬에서 OpenSSL 등으로 자체 서명 인증서를 생성하여 certificate_body, certificate_private_key 등을 사용할 수도 있습니다.

resource "aws_api_gateway_domain_name" "custom_domain" {
  domain_name              = var.custom_domain_name
  regional_certificate_arn = "arn:aws:acm:us-east-1:000000000000:certificate/fake-cert-uuid" # LocalStack에서 작동하는 더미 ARN (실제값 아님)

  endpoint_configuration {
    types = ["REGIONAL"]
  }

  security_policy = "TLS_1_0" # TLS 보안 정책

  # LocalStack 버그나 특성으로 인해 depends_on이 필요할 수 있음
  depends_on = [aws_api_gateway_stage.my_stage]
}

# 9. 기본 경로 매핑 (Custom Domain -> API Stage)
resource "aws_api_gateway_base_path_mapping" "my_mapping" {
  domain_name = aws_api_gateway_domain_name.custom_domain.domain_name
  api_id = aws_api_gateway_rest_api.my_api.id
  stage_name  = aws_api_gateway_stage.my_stage.stage_name
  # base_path   = "(none)" # 루트 경로(/)에 매핑하려면 생략하거나 (none) 사용
}

# 출력 값
output "base_url" {
  description = "API Gateway 호출 URL (기본)"
  value       = aws_api_gateway_stage.my_stage.invoke_url
}

output "api_endpoint" {
  description = "생성된 리소스의 전체 경로"
  value       = "${aws_api_gateway_stage.my_stage.invoke_url}/${var.resource_path}"
}

output "custom_domain_url" {
  description = "사용자 지정 도메인 URL (HTTPS)"
  value       = "https://${aws_api_gateway_domain_name.custom_domain.regional_domain_name}" # REGIONAL 엔드포인트 사용 시 regional_domain_name 사용
}

output "custom_domain_api_endpoint" {
  description = "사용자 지정 도메인을 사용한 API 엔드포인트"
  # 기본 경로 매핑이 루트(/)가 아닐 경우 base_path_mapping.base_path 를 추가해야 할 수 있음
  value       = "https://${aws_api_gateway_domain_name.custom_domain.regional_domain_name}/${var.resource_path}"
}

# curl -X GET "http://localhost:4566/restapis/{}/dev/alpha/beta.chris?V=1.0"