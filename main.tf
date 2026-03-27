terraform {
  required_providers {
    tcpcheck = {
      source = "bitonio/tcpcheck"
      version = "1.0.0"
    }
  }
}

resource "tcpcheck_tcp_check" "this" {
  host     = var.host
  port     = var.port
  interval = var.interval
  timeout  = var.timeout
}
