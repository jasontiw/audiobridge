@echo off
REM AudioBridge Release Script for Windows
REM Creates distribution packages with DLLs automatically

echo ============================================
echo AudioBridge Release Builder (Windows)
echo ============================================

set VERSION=%1
if "%VERSION%"=="" set VERSION=0.1.0

echo Building version %VERSION%

REM Set output directory
set OUTPUT_DIR=release
if exist %OUTPUT_DIR% rmdir /s /q %OUTPUT_DIR%
mkdir %OUTPUT_DIR%

REM Find MSYS2 and DLLs
set MSYS2_PATH=
set DLL_FOUND=0

REM Check common MSYS2 installation paths
for %%d in (C:\msys64 C:\msys32 D:\msys64 D:\msys32) do (
    if exist %%d (
        set MSYS2_PATH=%%d
        echo Found MSYS2 at: %%d
    )
)

if "%MSYS2_PATH%"=="" (
    echo ERROR: MSYS2 not found. Please install MSYS2 from https://www.msys2.org/
    exit /b 1
)

REM Find DLLs in MSYS2
set DLL_PORT=
set DLL_PORTAUDIOCPP=

REM Try different bin directories
for %%b in (mingw64\bin ucrt64\bin mingw32\bin) do (
    if exist "%MSYS2_PATH%\%%b\libportaudio.dll" (
        set DLL_PORT=%MSYS2_PATH%\%%b\libportaudio.dll
    )
    if exist "%MSYS2_PATH%\%%b\libportaudiocpp.dll" (
        set DLL_PORTAUDIOCPP=%MSYS2_PATH%\%%b\libportaudiocpp.dll
    )
)

REM Copy DLLs if found
echo.
echo Looking for PortAudio DLLs...

if not "%DLL_PORT%"=="" (
    echo Found: %DLL_PORT%
    copy /y "%DLL_PORT%" "%OUTPUT_DIR%\" >nul
    set DLL_FOUND=1
) else (
    echo WARNING: libportaudio.dll not found in MSYS2
)

if not "%DLL_PORTAUDIOCPP%"=="" (
    echo Found: %DLL_PORTAUDIOCPP%
    copy /y "%DLL_PORTAUDIOCPP%" "%OUTPUT_DIR%\" >nul
) else (
    echo WARNING: libportaudiocpp.dll not found in MSYS2
)

if %DLL_FOUND%==0 (
    echo.
    echo ERROR: Could not find PortAudio DLLs.
    echo Please ensure PortAudio is installed: pacman -S mingw-w64-x86_64-portaudio
    exit /b 1
)

REM Set PATH for Go build
set PATH=%MSYS2_PATH%\mingw64\bin;%MSYS2_PATH%\ucrt64\bin;%PATH%
set CGO_ENABLED=1

REM Build
echo.
echo Building audiobridge.exe...
go build -ldflags="-s -w -X main.version=v%VERSION%" -o "%OUTPUT_DIR%\audiobridge.exe"

if %ERRORLEVEL% NEQ 0 (
    echo ERROR: Build failed!
    exit /b 1
)

REM Copy README
if exist "..\README.md" (
    copy /y "..\README.md" "%OUTPUT_DIR%\" >nul
)

echo.
echo Build successful!
echo.

REM Find 7zip or zip
set ZIP_FOUND=0
where 7z >nul 2>&1
if %ERRORLEVEL%==0 (
    set ZIP_CMD=7z a -tzip
    set ZIP_FOUND=1
)

where zip >nul 2>&1
if %ERRORLEVEL%==0 (
    set ZIP_CMD=zip
    set ZIP_FOUND=1
)

REM Create ZIP inside release folder
cd %OUTPUT_DIR%

if %ZIP_FOUND%==1 (
    %ZIP_CMD% audiobridge-%VERSION%-windows-amd64.zip *
    echo.
    echo ============================================
    echo Release created: release\audiobridge-%VERSION%-windows-amd64.zip
    echo ============================================
    
    REM Delete individual files, keep only zip
    del /q audiobridge.exe 2>nul
    del /q libportaudio.dll 2>nul
    del /q libportaudiocpp.dll 2>nul
    del /q README.md 2>nul
    
    echo.
    echo Final contents:
    dir
) else (
    echo.
    echo WARNING: No zip tool found. Contents are in:
    echo %OUTPUT_DIR%
    echo.
    dir
)

cd ..
