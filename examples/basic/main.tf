# Example usage of the tcpcheck module

terraform {
  required_providers {
    tcpcheck = {
      source = "bitonio/tcpcheck"
      version = "1.0.0"
    }
  }
}

module "web_tcpcheck" {
  source = "../"

  host     = "example.com"
  port     = 443
  interval = 5
  timeout  = 30
}

output "tcpcheck_result" {
  value = {
    id   = module.web_tcpcheck.check_id
    host = module.web_tcpcheck.host
    port = module.web_tcpcheck.port
  }
}
