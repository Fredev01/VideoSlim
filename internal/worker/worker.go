package worker

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	"videoslim/internal/compressor"
)

type WorkerPool struct {
	maxWorkers      int
	semaphore       chan struct{}
	totalVideos     int
	processedVideos int32 // Contador atómico
}

type ProcessResult struct {
	OriginalPath   string
	CompressedPath string
	OriginalSize   int64
	CompressedSize int64
	Error          error
}

func NewWorkerPool(workers int) *WorkerPool {
	return &WorkerPool{
		maxWorkers: workers,
		semaphore:  make(chan struct{}, workers),
	}
}

func (wp *WorkerPool) ProcessVideos(videos []string, config compressor.Config) []ProcessResult {
	wp.totalVideos = len(videos)
	var wg sync.WaitGroup
	results := make([]ProcessResult, wp.totalVideos)
	atomic.StoreInt32(&wp.processedVideos, 0)

	fmt.Printf("📝 Iniciando compresión de %d videos\n", wp.totalVideos)
	fmt.Printf("🔧 Usando %d workers concurrentes\n", wp.maxWorkers)
	fmt.Println("-----------------------------")

	for i, video := range videos {
		wg.Add(1)
		go func(index int, videoPath string) {
			defer wg.Done()
			wp.semaphore <- struct{}{}
			defer func() { <-wp.semaphore }()

			results[index] = wp.processVideo(videoPath, config)

			// Incremento atómico y actualización de progreso
			processed := atomic.AddInt32(&wp.processedVideos, 1)
			wp.updateProgress(processed)
		}(i, video)
	}

	wg.Wait()
	fmt.Println("\n✅ Proceso de compresión completado")
	return results
}

func (wp *WorkerPool) updateProgress(processed int32) {
	percentage := float32(processed) / float32(wp.totalVideos) * 100
	fmt.Printf("\rProgreso: [%-50s] %.2f%% (%d/%d)",
		generateProgressBar(percentage),
		percentage,
		processed,
		wp.totalVideos)
}

func generateProgressBar(percentage float32) string {
	const width = 50
	filled := int(percentage / 2) // La barra es de 50 caracteres
	return fmt.Sprintf("%s%s",
		strings.Repeat("█", filled),
		strings.Repeat("░", width-filled))
}

func (wp *WorkerPool) processVideo(input string, config compressor.Config) ProcessResult {
	originalInfo, _ := os.Stat(input)
	outputDir := config.OutputPath
	errCreateDir := os.MkdirAll(outputDir, 0755)
	if errCreateDir != nil {
		log.Fatal(errCreateDir)
	}

	outputFilename := fmt.Sprintf("compressed_%s", filepath.Base(input))
	output := filepath.Join(outputDir, outputFilename)

	fmt.Printf("\n🎬 Comprimiendo: %s", filepath.Base(input))

	err := compressor.CompressVideo(input, output, config)

	result := ProcessResult{
		OriginalPath:   input,
		CompressedPath: output,
		OriginalSize:   originalInfo.Size(),
	}

	if err != nil {
		fmt.Printf("\n❌ Error comprimiendo %s: %v", filepath.Base(input), err)
		result.Error = err
		return result
	}

	compressedInfo, _ := os.Stat(output)
	result.CompressedSize = compressedInfo.Size()

	return result
}
