# AudioBridge Build Scripts

Scripts para compilar AudioBridge con CGO enabled para distribución.

## Requisitos

### Windows
- [MSYS2](https://www.msys2.org/) instalado
- En terminal MSYS2:
  ```bash
  pacman -S mingw-w64-x86_64-gcc
  pacman -S mingw-w64-x86_64-portaudio
  pacman -S mingw-w64-x86_64-opus
  ```

### Linux (Ubuntu/Debian)
```bash
sudo apt-get install -y \
    build-essential \
    libportaudio2 \
    libopus0 \
    libopusfile0 \
    pkg-config
```

### macOS
```bash
brew install portaudio opus
```

## Compilación

### Windows (MSYS2 MinGW)
```bash
export CGO_ENABLED=1
export CC=gcc
go build -o audiobridge.exe
```

### Linux
```bash
CGO_ENABLED=1 go build -o audiobridge-linux-amd64
```

### macOS
```bash
CGO_ENABLED=1 go build -o audiobridge-darwin
```

## Distribución

El binario resultante incluye las librerías linkeadas estáticamente (donde es posible).
Para Windows, el binario requiere las DLLs en la misma carpeta o en el PATH:

- `libportaudio-2.dll` (si no está linkeado estáticamente)
- `libopus.dll`

## Docker Build (Alternativa)

```bash
docker build -f Dockerfile.build --target artifacts --output ./dist .
```

Esto produce binarios para todas las plataformas en `./dist/`
