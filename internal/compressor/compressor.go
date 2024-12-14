package compressor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Definir Config aquí o importarla si está en otro paquete
type Config struct {
	InputPath     string
	OutputPath    string
	Concurrency   int
	Quality       int
	MaxSizeTarget float64
	VideoFormat   string
}

func ListVideoFiles(dir string) ([]string, error) {
	videoExtensions := []string{".mp4", ".avi", ".mov", ".mkv", ".flv"}
	var videos []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			ext := strings.ToLower(filepath.Ext(path))
			for _, validExt := range videoExtensions {
				if ext == validExt {
					videos = append(videos, path)
					break
				}
			}
		}
		return nil
	})

	return videos, err
}

func CompressVideo(input, output string, config Config) error {
	// Comandos FFmpeg para compresión
	cmd := exec.Command("ffmpeg",
		"-i", input,
		"-vf", fmt.Sprintf("scale=%d:-2", calculateResolution(input)),
		"-crf", fmt.Sprintf("%d", config.Quality),
		"-preset", "medium",
		"-y", // Sobrescribir sin preguntar
		output,
	)

	// Capturar la salida de error
	outputBytes, err := cmd.CombinedOutput()
	if err != nil {
		// Convertir bytes a string para el mensaje de error
		return fmt.Errorf("error comprimiendo video: %v - %s", err, string(outputBytes))
	}

	return nil
}

func ValidateFFmpeg() error {
	cmd := exec.Command("ffmpeg", "-version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("FFmpeg no está instalado")
	}
	return nil
}

func calculateResolution(input string) int {
	// Lógica para calcular nueva resolución basada en tamaño objetivo
	return 1280 // Por defecto 720p
}
