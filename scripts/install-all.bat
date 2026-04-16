@echo off
echo ============================================
echo   Installing All Dependencies
echo ============================================

echo.
echo [1/3] Installing Backend Dependencies...
call "%~dp0backend\install.bat"

echo.
echo [2/3] Installing Frontend Dependencies...
call "%~dp0frontend\install.bat"

echo.
echo [3/3] Installing Python Dependencies...
call "%~dp0python-service\install.bat"

echo.
echo All dependencies installed!
pause
