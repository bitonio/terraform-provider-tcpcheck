#!/bin/bash
set -e

echo "Building tcpcheck provider..."

# Build for current platform
go build -o terraform-provider-tcpcheck

echo "Installing provider locally..."

# Determine OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
esac

# Provider installation directory
PLUGIN_DIR="${HOME}/.terraform.d/plugins/registry.terraform.io/bitonio/tcpcheck/1.0.0/${OS}_${ARCH}"

mkdir -p "$PLUGIN_DIR"
cp terraform-provider-tcpcheck "$PLUGIN_DIR/"

echo "Provider installed to: $PLUGIN_DIR"
echo "You can now use the provider in your Terraform configuration"
