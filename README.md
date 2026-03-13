# AudioBridge

Cross-platform tool for sharing audio between computers on the same local network.

## Features

- Share microphone audio between PCs
- Share system audio (what you hear) between PCs
- Low latency UDP streaming (< 80ms on local network)
- Opus codec for efficient audio compression
- Works on Windows, macOS, and Linux

## Quick Start

### Windows

1. Download the latest release from [GitHub Releases](https://github.com/jasontiw/audiobridge/releases)
2. Extract the ZIP file
3. Run `audiobridge.exe`

### Send Audio (from PC A)

```bash
audiobridge send --target <PC_B_IP>
```

Example:
```bash
audiobridge send --target 192.168.1.100
```

### Receive Audio (on PC B)

```bash
audiobridge receive
```

## Command Line Options

### Send Command

```bash
audiobridge send [flags]

  -t, --target string    Target IP address (required)
  -p, --port int        Target UDP port (default 9876)
  -s, --source string   Audio source: microphone or system (default "microphone")
  -b, --bitrate int     Opus bitrate in bits per second (default 64000)
  -r, --rate int        Audio sample rate (default 48000)
  -c, --channels int    Number of channels: 1=mono, 2=stereo (default 1)
```

### Receive Command

```bash
audiobridge receive [flags]

  -p, --port int          Listen UDP port (default 9876)
  -j, --jitter-buffer int Jitter buffer latency in ms (default 50)
  -d, --device int        Output device index (-1 for default)
```

### Devices Command

```bash
audiobridge devices
```

Lists all available audio input and output devices.

## Configuration File

Create `audiobridge.toml` in the same folder as the executable:

```toml
[general]
log_level = "info"
port = 9876

[send]
source = "microphone"
target = "192.168.1.100"
codec = "opus"
bitrate = 64000
channels = 1
sample_rate = 48000

[receive]
output_device = "default"
jitter_buffer_ms = 50
```

## Building from Source

### Prerequisites

- Go 1.22 or later
- PortAudio library
- Opus library
- GCC (for CGO)

### Linux (Ubuntu/Debian)

```bash
sudo apt-get install -y portaudio libasound2-dev libopus0 libopusfile0 pkg-config
CGO_ENABLED=1 go build -o audiobridge
```

### Windows

1. Install [MSYS2](https://www.msys2.org/)
2. In MSYS2 UCRT64 terminal:
```bash
pacman -S mingw-w64-ucrt-x86_64-gcc
pacman -S mingw-w64-x86_64-portaudio
pacman -S mingw-w64-x86_64-opus

export CGO_ENABLED=1
go build -o audiobridge.exe
```

3. Copy required DLLs:
```bash
cp /mingw64/bin/libportaudio.dll ./
cp /mingw64/bin/libportaudiocpp.dll ./
```

### macOS

```bash
brew install portaudio opus
CGO_ENABLED=1 go build -o audiobridge
```

## Architecture

```
┌─────────────┐     UDP      ┌─────────────┐
│   Sender    │ ──────────► │  Receiver   │
│             │              │             │
│ ┌─────────┐ │              │ ┌─────────┐ │
│ │ Capture │ │              │ │  Jitter │ │
│ └────┬────┘ │              │ │ Buffer  │ │
│      │      │              │ └────┬────┘ │
│ ┌────▼────┐ │              │ ┌────▼────┐ │
│ │  Opus   │ │              │ │  Opus   │ │
│ │ Encode  │ │              │ │ Decode  │ │
│ └────┬────┘ │              │ └────┬────┘ │
│      │      │              │      │      │
│ ┌────▼────┐ │              │ ┌────▼────┐ │
│ │   UDP    │ │              │ │  Play   │ │
│ └─────────┘ │              │ └─────────┘ │
└─────────────┘              └─────────────┘
```

## Troubleshooting

### "PortAudio not available" on Windows

Ensure you have the PortAudio DLLs in the same folder as the executable:
- `libportaudio.dll`
- `libportaudiocpp.dll`

### High latency

Try reducing the jitter buffer:
```bash
audiobridge receive --jitter-buffer 30
```

### No audio device found

Check available devices:
```bash
audiobridge devices
```

## License

MIT License - see LICENSE file for details.
