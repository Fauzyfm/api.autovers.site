#!/bin/bash
# Generate Swagger documentation
swag init -g main.go
echo "✅ Swagger documentation generated!"
echo "Visit: http://localhost:8080/swagger/index.html"
