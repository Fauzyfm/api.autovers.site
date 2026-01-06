@echo off
REM Generate Swagger documentation for Windows
swag init -g main.go
echo.
echo ✅ Swagger documentation generated!
echo.
echo Visit: http://localhost:8080/swagger/index.html
pause
