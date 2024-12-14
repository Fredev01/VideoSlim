# VideoSlim 🎥

## Descripción
VideoSlim es una herramienta de línea de comandos para comprimir videos de manera eficiente utilizando FFmpeg.

## Requisitos Previos
- Go 1.21+
- FFmpeg instalado

## Instalación
```bash
git clone https://github.com/tu-usuario/videoslim.git
cd videoslim
go build cmd/videoslim/main.go
```

## Uso
```bash
# Comprimir un video individual
./videoslim -input video.mp4

# Comprimir todos los videos en un directorio
./videoslim -input /ruta/a/videos

# Opciones adicionales
./videoslim -input video.mp4 -concurrency 4 -quality 25 -max-size 500
```

## Flags Disponibles
- `-input`: Ruta de archivo o directorio (requerido)
- `-output`: Directorio de salida
- `-concurrency`: Número de workers paralelos
- `-quality`: Nivel de compresión (0-51)
- `-max-size`: Tamaño máximo objetivo en MB
- `-format`: Formato de salida

## Licencia
MIT License
```

## Características del Proyecto

1. **Modularidad**: Separación clara de responsabilidades
2. **Concurrencia**: Procesamiento paralelo de videos
3. **Flexibilidad**: Múltiples opciones de configuración
4. **Robusto manejo de errores**
5. **Reporte detallado de compresión**

## Pasos Siguientes

1. Implementar más pruebas unitarias
2. Agregar soporte para más formatos de video
3. Incluir más opciones de compresión
4. Mejorar detección de resolución original
