@echo off
REM AudioBridge Build Script for Windows
REM Requires MSYS2 installed

echo ============================================
echo AudioBridge Build Script (Windows)
echo ============================================

REM Check if MSYS2 is installed
where gcc >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: GCC not found. Please install MSYS2:
    echo   1. Download from https://www.msys2.org/
    echo   2. Run: pacman -S mingw-w64-x86_64-gcc
    echo   3. Run: pacman -S mingw-w64-x86_64-portaudio
    echo   4. Run: pacman -S mingw-w64-x86_64-opus
    exit /b 1
)

echo Building with CGO enabled...
set CGO_ENABLED=1
go build -o audiobridge.exe -ldflags="-s -w"

if %ERRORLEVEL% EQU 0 (
    echo.
    echo ============================================
    echo SUCCESS: audiobridge.exe created!
    echo ============================================
) else (
    echo.
    echo ERROR: Build failed
    exit /b 1
)
