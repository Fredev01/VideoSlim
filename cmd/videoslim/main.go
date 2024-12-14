package main

import (
	"fmt"
	"log"
	"os"

	"videoslim/internal/cli"
	"videoslim/internal/compressor"
	"videoslim/internal/worker"
)

func main() {
	// Parsear argumentos de línea de comandos
	config := cli.ParseFlags()

	compressorConfig := compressor.Config{
		OutputPath:    config.OutputPath,
		Quality:       config.Quality,
		MaxSizeTarget: config.MaxSizeTarget,
		VideoFormat:   config.VideoFormat,
	}

	// Validar dependencias
	if err := compressor.ValidateFFmpeg(); err != nil {
		log.Fatalf("Error: FFmpeg no encontrado. %v", err)
	}

	// Obtener lista de archivos de video
	videoPaths, err := getVideoPaths(config.InputPath)
	if err != nil {
		log.Fatalf("Error obteniendo videos: %v", err)
	}

	// Crear pool de workers
	pool := worker.NewWorkerPool(config.Concurrency)

	// Procesar videos
	results := pool.ProcessVideos(videoPaths, compressorConfig)

	// Generar reporte de resultados
	generateReport(results)
}

func getVideoPaths(inputPath string) ([]string, error) {
	// Lógica para obtener rutas de videos (archivo único o múltiples en directorio)
	stat, err := os.Stat(inputPath)
	if err != nil {
		return nil, err
	}

	var videoPaths []string
	if stat.IsDir() {
		videoPaths, err = compressor.ListVideoFiles(inputPath)
	} else {
		videoPaths = []string{inputPath}
	}

	return videoPaths, err
}

func generateReport(results []worker.ProcessResult) {
	fmt.Println("\n--- Reporte de Compresión ---")
	for _, result := range results {
		if result.Error != nil {
			fmt.Printf("❌ %s: Compresión fallida - %v\n", result.OriginalPath, result.Error)
		} else {
			compressionRate := (1 - float64(result.CompressedSize)/float64(result.OriginalSize)) * 100
			fmt.Printf("✅ %s: Comprimido (%.2f%%) - Original: %dMB, Comprimido: %dMB\n",
				result.OriginalPath, compressionRate,
				result.OriginalSize/(1024*1024),
				result.CompressedSize/(1024*1024))
		}
	}
}
