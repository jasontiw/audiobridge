# AudioBridge — Product Requirements Document

**Versión:** 0.2-draft  
**Fecha:** 2026-03-13  
**Estado:** En revisión  
**Autor:** Por definir

---

## 1. Resumen ejecutivo

AudioBridge es una herramienta multiplataforma escrita en Go que permite compartir audio del sistema entre computadoras en la misma red local. Funciona en Windows, macOS y Linux sin necesidad de drivers adicionales ni herramientas de terceros.

La experiencia de usuario se centra en una **TUI interactiva** (interfaz de terminal tipo htop) complementada por un **ícono en la bandeja del sistema (systray)** para acceso rápido desde el escritorio. El usuario puede operar la herramienta completamente desde la terminal o simplemente dejarla corriendo en segundo plano y gestionarla desde la barra de tareas.

El proyecto nace como alternativa open source a soluciones como VoiceMeeter + VBAN, orientada a uso general en cualquier contexto multi-PC.

---

## 2. Problema

Los setups multi-PC son comunes en streaming, gaming y trabajo en casa, pero compartir audio entre computadoras sigue siendo sorprendentemente difícil:

- **VoiceMeeter** solo funciona en Windows y requiere aprender una interfaz compleja.
- **Soluciones de hardware** (cables, switches) son costosas y rígidas.
- **No existe una herramienta CLI multiplataforma** que haga esto de forma simple y confiable.

El usuario técnico que trabaja con múltiples PCs no debería necesitar instalar software propietario de terceros para una tarea tan básica como escuchar el audio de otra máquina.

---

## 3. Objetivos del producto

### 3.1 Objetivos para el MVP

| # | Objetivo | Métrica de éxito |
|---|----------|-----------------|
| O1 | Compartir audio del sistema entre dos PCs | Latencia < 80ms en red local |
| O2 | Compartir micrófono entre dos PCs | Latencia < 80ms en red local |
| O3 | Funcionar en Windows, macOS y Linux | Binarios independientes para las 3 plataformas |
| O4 | Configuración simple via TOML | Setup funcional en < 5 minutos |
| O5 | Distribución como binario único sin dependencias de runtime | `go build` produce un ejecutable portable |

### 3.2 Fuera de alcance (MVP)

- Interfaz gráfica (GUI)
- Encriptación del stream de audio
- Soporte para más de 2 PCs simultáneas
- Mezcla de múltiples fuentes de audio
- Transmisión por internet (WAN)

---

## 4. Usuarios objetivo

### Perfil primario: El desarrollador / técnico multi-PC

- Tiene 2 o más computadoras en su escritorio o red local
- Usa Linux, macOS o Windows (o una combinación)
- Se siente cómodo usando la terminal
- Prefiere herramientas configurables sobre interfaces visuales
- Casos de uso típicos: streaming, trabajo remoto, gaming, desarrollo

### Perfil secundario: El entusiasta de audio / streamer técnico

- Setup de streaming con PC de captura y PC de juego separadas
- Quiere enrutar audio entre máquinas sin depender de soluciones Windows-only
- Valora la estabilidad y el control sobre la experiencia visual

---

## 5. Casos de uso principales

### UC-01: Escuchar el audio de la PC secundaria en los parlantes de la principal

**Actor:** Usuario con dos PCs conectadas a la misma red.  
**Flujo:**
1. En la PC secundaria: `audiobridge send --audio-system`
2. En la PC principal: `audiobridge receive --output default`
3. El audio del sistema de la PC secundaria suena en los parlantes de la PC principal.

**Criterio de aceptación:** Latencia percibida < 80ms, sin cortes en condiciones normales de red local.

---

### UC-02: Compartir micrófono de la PC principal a la secundaria

**Actor:** Usuario cuyo micrófono está conectado a la PC principal pero necesita usarlo en la secundaria.  
**Flujo:**
1. En la PC principal: `audiobridge send --microphone`
2. En la PC secundaria: `audiobridge receive --create-virtual-device`
3. Las apps en la PC secundaria ven un dispositivo de audio virtual con el micrófono de la PC principal.

**Criterio de aceptación:** El dispositivo virtual aparece en la lista de micrófonos del OS secundario.

---

### UC-03: Modo bidireccional (audio + micrófono simultáneos)

**Actor:** Usuario que quiere audio completo entre dos PCs.  
**Flujo:**
1. Se configura `audiobridge.toml` en ambas máquinas con los roles definidos.
2. `audiobridge start` en cada máquina establece los dos streams simultáneamente.

**Criterio de aceptación:** Ambos streams operan en paralelo sin interferencia.

---

### UC-04: Verificar estado de la conexión

**Actor:** Usuario que quiere diagnosticar problemas.  
**Flujo:**
1. `audiobridge status` muestra streams activos, latencia medida, paquetes perdidos y dispositivos en uso.

---

## 6. Requerimientos funcionales

### 6.1 Captura y reproducción de audio

| ID | Requerimiento |
|----|---------------|
| RF-01 | Capturar audio del micrófono por defecto del sistema |
| RF-02 | Capturar audio del sistema (loopback) |
| RF-03 | Reproducir audio en el dispositivo de salida por defecto |
| RF-04 | Permitir seleccionar dispositivo de entrada/salida por nombre o índice |
| RF-05 | Soportar sample rate de 44100 Hz y 48000 Hz |
| RF-06 | Soportar 1 canal (mono) y 2 canales (estéreo) |

### 6.2 Transporte de red

| ID | Requerimiento |
|----|---------------|
| RF-07 | Transmitir audio via UDP en red local |
| RF-08 | Implementar jitter buffer configurable (default: 50ms) |
| RF-09 | Numerar paquetes para detección de pérdida y reordenamiento |
| RF-10 | Puerto configurable (default: 9876) |
| RF-11 | Modo de descubrimiento automático en la red local (mDNS / UDP broadcast) |

### 6.3 Compresión (Opus)

| ID | Requerimiento |
|----|---------------|
| RF-12 | Comprimir audio con codec Opus antes de transmitir |
| RF-13 | Bit rate configurable (default: 64 kbps para voz, 128 kbps para música) |
| RF-14 | Modo de transmisión sin compresión (PCM crudo) disponible como opción |

### 6.4 Configuración

| ID | Requerimiento |
|----|---------------|
| RF-15 | Aceptar archivo de configuración `audiobridge.toml` |
| RF-16 | Flags de CLI que sobreescriben los valores del archivo de configuración |
| RF-17 | Comando `audiobridge init` que genera un archivo de configuración de ejemplo |
| RF-18 | Validar el archivo de configuración al arrancar y reportar errores claros |

### 6.5 Observabilidad

| ID | Requerimiento |
|----|---------------|
| RF-19 | Log de eventos en stdout con niveles: `info`, `warn`, `error` |
| RF-20 | Comando `audiobridge status` con métricas en tiempo real |
| RF-21 | Métricas: latencia promedio, paquetes enviados/recibidos/perdidos, dispositivos activos |

---

## 7. Requerimientos no funcionales

| ID | Categoría | Requerimiento |
|----|-----------|---------------|
| RNF-01 | Latencia | End-to-end < 80ms en red local con Ethernet |
| RNF-02 | Latencia | End-to-end < 120ms en red local con WiFi 5GHz |
| RNF-03 | Estabilidad | Sin cortes en sesiones de hasta 8 horas continuas |
| RNF-04 | CPU | Uso de CPU < 5% en hardware de referencia (CPU de 4 núcleos, 2020+) |
| RNF-05 | Memoria | Uso de RAM < 50 MB en operación normal |
| RNF-06 | Compatibilidad | Windows 10+, macOS 12+, Ubuntu 20.04+ |
| RNF-07 | Distribución | Binario único sin dependencias de runtime para el usuario final |
| RNF-08 | Licencia | Licencia MIT, código abierto en GitHub |

---

## 8. Arquitectura técnica

### 8.1 Stack tecnológico

| Componente | Tecnología | Justificación |
|------------|-----------|---------------|
| Lenguaje | Go 1.22+ | Multiplataforma, binario único, buena concurrencia |
| Audio I/O | `portaudio-go` (binding a PortAudio) | Soporte WASAPI / CoreAudio / ALSA / PipeWire |
| Codec | `hraban/opus` (binding a libopus via CGO) | Estándar de facto para audio de baja latencia |
| Transporte | `net.UDPConn` (stdlib) | Latencia mínima, no requiere librerías externas |
| Configuración | `BurntSushi/toml` | Simple, legible, ampliamente usado en Go |
| Descubrimiento | `grandcat/zeroconf` (mDNS) | Autodescubrimiento en red local sin configuración manual |
| CLI | `spf13/cobra` | Estándar de facto para CLIs en Go |

### 8.2 Estructura de módulos

```
audiobridge/
├── cmd/
│   ├── root.go          # Comando raíz y flags globales
│   ├── send.go          # Subcomando: audiobridge send
│   ├── receive.go       # Subcomando: audiobridge receive
│   ├── start.go         # Subcomando: audiobridge start (modo bidireccional)
│   ├── status.go        # Subcomando: audiobridge status
│   └── init.go          # Subcomando: audiobridge init (genera config)
├── audio/
│   ├── capture.go       # Captura de micrófono y loopback
│   ├── playback.go      # Reproducción en dispositivo de salida
│   └── devices.go       # Listado y selección de dispositivos
├── transport/
│   ├── sender.go        # Envío de paquetes UDP
│   ├── receiver.go      # Recepción de paquetes UDP
│   └── packet.go        # Estructura del paquete (header + payload)
├── codec/
│   ├── opus.go          # Encode/decode Opus
│   └── pcm.go           # Pass-through PCM sin compresión
├── jitter/
│   └── buffer.go        # Jitter buffer con reordenamiento por número de secuencia
├── discovery/
│   └── mdns.go          # Anuncio y descubrimiento via mDNS
├── config/
│   └── config.go        # Parsing y validación del archivo TOML
└── main.go
```

### 8.3 Formato del paquete UDP

```
┌─────────────────────────────────────────────────────────┐
│  Header (12 bytes)                                       │
│  ┌──────────┬──────────┬──────────┬────────────────┐    │
│  │ Magic    │ Version  │ Seq num  │ Timestamp      │    │
│  │ 2 bytes  │ 1 byte   │ 4 bytes  │ 5 bytes        │    │
│  └──────────┴──────────┴──────────┴────────────────┘    │
├─────────────────────────────────────────────────────────┤
│  Payload (variable)                                      │
│  Audio comprimido en Opus o PCM crudo                    │
└─────────────────────────────────────────────────────────┘
```

---

## 9. Interfaz de usuario (CLI)

### 9.1 Comandos principales

```bash
# Enviar audio del sistema al receiver en 192.168.1.10
audiobridge send --audio-system --target 192.168.1.10

# Enviar micrófono al receiver en 192.168.1.10
audiobridge send --microphone --target 192.168.1.10

# Recibir y reproducir en parlantes por defecto
audiobridge receive

# Modo bidireccional usando archivo de configuración
audiobridge start --config audiobridge.toml

# Ver estado y métricas
audiobridge status

# Listar dispositivos de audio disponibles
audiobridge devices

# Generar archivo de configuración de ejemplo
audiobridge init
```

### 9.2 Archivo de configuración (`audiobridge.toml`)

```toml
[general]
log_level = "info"        # debug | info | warn | error
port = 9876

[send]
source = "microphone"     # microphone | system | both
target = "192.168.1.10"
codec = "opus"            # opus | pcm
bitrate = 64000           # bits por segundo (solo para opus)
channels = 1              # 1 = mono, 2 = stereo
sample_rate = 44100

[receive]
output_device = "default"
jitter_buffer_ms = 50

[discovery]
enabled = true
announce_name = "mi-pc-principal"
```

### 9.3 Salida de `audiobridge status`

```
AudioBridge v0.1.0 — Estado actual
────────────────────────────────────
Stream activo:    SEND (micrófono → 192.168.1.10:9876)
Codec:            Opus 64kbps, mono, 44100 Hz
Uptime:           00:42:17

Métricas de red:
  Paquetes enviados:    127,401
  Paquetes perdidos:    12 (0.009%)
  Latencia promedio:    18ms
  Latencia máxima:      34ms

Dispositivo de captura:  MacBook Pro Microphone (idx: 0)
```

---

## 10. Plan de desarrollo (MVP)

### Fase 1 — Fundamentos (semanas 1–2)

- [ ] Setup del repositorio: módulos Go, CI básico (GitHub Actions), linting
- [ ] Implementación de captura de audio (`portaudio-go`)
- [ ] Implementación de reproducción de audio
- [ ] Comando `audiobridge devices` funcional
- [ ] Transporte UDP básico (sin jitter buffer, sin codec)

### Fase 2 — Audio por red (semanas 3–4)

- [ ] Codec Opus (encode/decode)
- [ ] Jitter buffer con reordenamiento por número de secuencia
- [ ] Estructura de paquetes con header completo
- [ ] Comandos `send` y `receive` funcionales en red local
- [ ] Tests de latencia en los 3 sistemas operativos

### Fase 3 — UX y configuración (semanas 5–6)

- [ ] Parsing de archivo TOML (`audiobridge init` + validación)
- [ ] Comando `audiobridge status` con métricas en tiempo real
- [ ] Descubrimiento automático via mDNS
- [ ] Logging con niveles
- [ ] Modo bidireccional (`audiobridge start`)

### Fase 4 — Distribución (semana 7)

- [ ] Binarios compilados para Windows (amd64), macOS (amd64 + arm64), Linux (amd64)
- [ ] Pipeline de release en GitHub Actions
- [ ] README con instalación y guía de inicio rápido
- [ ] Documentación de todos los flags y opciones del TOML

---

## 11. Riesgos y mitigaciones

| Riesgo | Probabilidad | Impacto | Mitigación |
|--------|-------------|---------|-----------|
| CGO complica la compilación cruzada (PortAudio, Opus) | Alta | Alto | Usar Docker con toolchains por OS; documentar el proceso de build |
| Loopback de audio del sistema no disponible sin permisos en macOS | Media | Medio | Documentar requisito de permiso; considerar ScreenCaptureKit como alternativa |
| Variabilidad de latencia en WiFi 2.4GHz | Alta | Medio | Jitter buffer adaptativo; recomendar Ethernet en docs |
| Dispositivo virtual de micrófono requiere driver de kernel en Windows/Linux | Alta | Alto | Aclarar en MVP que el "dispositivo virtual" es una fase 2 post-MVP |

---

## 12. Métricas de éxito post-lanzamiento

| Métrica | Objetivo a 3 meses |
|---------|-------------------|
| Estrellas en GitHub | > 500 |
| Issues abiertos sin respuesta > 7 días | < 10% del total |
| Latencia media reportada por usuarios | < 80ms |
| Plataformas con binarios funcionales | Windows, macOS, Linux (las 3) |

---

## 13. Glosario

| Término | Definición |
|---------|-----------|
| PCM | Pulse-Code Modulation. Audio crudo sin compresión. |
| Opus | Codec de audio de baja latencia, estándar en WebRTC y VoIP. |
| Jitter buffer | Buffer que absorbe variaciones de latencia en la red para reproducción suave. |
| Loopback | Captura del audio que el sistema operativo está reproduciendo (audio del sistema). |
| mDNS | Multicast DNS. Protocolo de descubrimiento de servicios en red local (como Bonjour). |
| WASAPI | Windows Audio Session API. API de audio de baja latencia en Windows. |
| ALSA / PipeWire | Sistemas de audio en Linux. PipeWire es el moderno, ALSA el clásico. |
| CoreAudio | API de audio nativa de macOS. |
| CGO | Mecanismo de Go para llamar código C desde Go. |

---

*Documento generado el 2026-03-13. Sujeto a revisión antes de inicio de desarrollo.*
