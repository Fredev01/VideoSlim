package cli

import (
	"flag"
	"log"
	"path/filepath"
	"runtime"
)

type Config struct {
	InputPath     string
	OutputPath    string
	Concurrency   int
	Quality       int
	MaxSizeTarget float64
	VideoFormat   string
}

func ParseFlags() Config {
	inputPath := flag.String("input", "", "Ruta de archivo o directorio de videos")
	outputPath := flag.String("output", "", "Directorio de salida (opcional)")
	concurrency := flag.Int("concurrency", getCPUCores(), "Número de workers concurrentes")
	quality := flag.Int("quality", 23, "Nivel de calidad (0-51, menor = mejor calidad)")
	maxSize := flag.Float64("max-size", 500, "Tamaño máximo de video en MB")
	format := flag.String("format", "mp4", "Formato de salida")

	flag.Parse()

	// Validaciones
	if *inputPath == "" {
		flag.Usage()
		log.Fatal("Debe especificar un directorio o archivo de entrada")
	}

	// Establecer ruta de salida por defecto
	if *outputPath == "" {
		*outputPath = filepath.Join(filepath.Dir(*inputPath), "compressed")
	}

	return Config{
		InputPath:     *inputPath,
		OutputPath:    *outputPath,
		Concurrency:   *concurrency,
		Quality:       *quality,
		MaxSizeTarget: *maxSize,
		VideoFormat:   *format,
	}
}

func getCPUCores() int {
	return runtime.NumCPU() - 1
}
