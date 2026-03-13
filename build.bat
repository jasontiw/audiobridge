@echo off
REM AudioBridge Build Script for Windows
REM Requires MSYS2 installed

echo ============================================
echo AudioBridge Build Script (Windows)
echo ============================================

REM Check if pkg-config is available
where pkg-config >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: pkg-config not found. Please install MSYS2:
    echo   1. Download from https://www.msys2.org/
    echo   2. Install and run MSYS2
    echo   3. In MSYS2 terminal run:
    echo      pacman -S mingw-w64-x86_64-gcc
    echo      pacman -S mingw-w64-x86_64-portaudio
    echo      pacman -S mingw-w64-x86_64-opus
    echo      pacman -S mingw-w64-x86_64-pkg-config
    echo   4. Add MSYS2\mingw64\bin to your PATH
    exit /b 1
)

REM Check if GCC is available
where gcc >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: GCC not found. Please install MSYS2 with mingw-w64
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
