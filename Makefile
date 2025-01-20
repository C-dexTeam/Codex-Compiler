.PHONY: help

dev.swagger.init:
	@echo "Generating swagger..."
	@~/go/bin/swag init  --parseVendor  -d . -g ./cmd/main.go 

help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Targets:"
	@echo "  dev.swagger.init    Generate Swagger"
	@echo "  help                Show this help"