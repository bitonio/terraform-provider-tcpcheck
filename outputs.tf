output "check_id" {
  description = "The unique identifier for the health check"
  value       = tcpcheck_tcp_check.this.id
}

output "host" {
  description = "The host that was checked"
  value       = tcpcheck_tcp_check.this.host
}

output "port" {
  description = "The port that was checked"
  value       = tcpcheck_tcp_check.this.port
}
