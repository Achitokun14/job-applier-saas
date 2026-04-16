@echo off
echo Installing Python Dependencies...
cd /d "%~dp0..\python-service"
pip install -r requirements.txt
echo Dependencies installed.
