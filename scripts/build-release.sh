# AudioBridge Release Script
# This script creates a distribution package with all necessary files

set -e

echo "========================================"
echo "AudioBridge Release Builder"
echo "========================================"

# Check for dependencies
if ! command -v zip &> /dev/null; then
    echo "Error: zip not found. Install with: apt install zip (Linux)"
    exit 1
fi

# Set variables
VERSION=${1:-"0.1.0"}
DATE=$(date +%Y-%m-%d)
OUTPUT_DIR="release"

echo "Building version $VERSION"

# Create output directory
mkdir -p $OUTPUT_DIR

# Build Linux
echo "Building Linux..."
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X main.version=v$VERSION" -o $OUTPUT_DIR/audiobridge-linux-amd64

# Build Windows  
echo "Building Windows..."
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -X main.version=v$VERSION" -o $OUTPUT_DIR/audiobridge-windows-amd64.exe

# Copy Windows DLLs (you need to provide these)
echo "Note: For Windows, copy the following DLLs to the release folder:"
echo "  - libportaudio.dll"
echo "  - libportaudiocpp.dll"

# Create README for distribution
cat > $OUTPUT_DIR/README.txt << 'EOF'
AudioBridge vVERSION

Quick Start:
1. Run audiobridge-windows-amd64.exe (Windows) or ./audiobridge-linux-amd64 (Linux)
2. On sending PC: audiobridge send --target <RECEIVER_IP>
3. On receiving PC: audiobridge receive

For full documentation, visit: https://github.com/jasontiw/audiobridge

EOF

sed -i "s/VERSION/$VERSION/" $OUTPUT_DIR/README.txt

# Create ZIP files
echo "Creating distribution packages..."
cd $OUTPUT_DIR

# Windows
zip -r audiobridge-$VERSION-windows-amd64.zip audiobridge-windows-amd64.exe README.txt
# Note: Add DLLs manually before distributing

# Linux  
zip -r audiobridge-$VERSION-linux-amd64.zip audiobridge-linux-amd64 README.txt

cd ..

echo "========================================"
echo "Release files created in ./release/"
echo "========================================"
ls -la $OUTPUT_DIR/
