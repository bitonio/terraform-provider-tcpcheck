variable "host" {
  description = "The hostname or IP address to check"
  type        = string
}

variable "port" {
  description = "The TCP port to check"
  type        = number
  
  validation {
    condition     = var.port > 0 && var.port <= 65535
    error_message = "Port must be between 1 and 65535."
  }
}

variable "interval" {
  description = "The interval in seconds between health check attempts"
  type        = number
  default     = 5
  
  validation {
    condition     = var.interval > 0
    error_message = "Interval must be greater than 0."
  }
}

variable "timeout" {
  description = "The total timeout in seconds for all health check attempts"
  type        = number
  default     = 60
  
  validation {
    condition     = var.timeout > 0
    error_message = "Timeout must be greater than 0."
  }
}
