package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"os"
	"runtime"
	"sync"
	"time"
)

func loadImageAndNormalize(path string) ([]float32, image.Rectangle, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, image.Rectangle{}, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, image.Rectangle{}, err
	}

	grayImg := image.NewGray(img.Bounds())
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			grayImg.Set(x, y, color.GrayModel.Convert(img.At(x, y)))
		}
	}

	bounds := grayImg.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	normalized := make([]float32, width*height)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			normalized[y*width+x] = float32(grayImg.GrayAt(x, y).Y) / 255.0
		}
	}

	return normalized, bounds, nil
}

func multiplyMatrices(A, B []float32, width int, start, end int) []float32 {
	result := make([]float32, width*(end-start))
	for i := start; i < end; i++ {
		for j := 0; j < width; j++ {
			result[(i-start)*width+j] = A[i*width+j] * B[i*width+j]
		}
	}
	return result
}

func main() {
	startTime := time.Now()

	A, boundsA, err := loadImageAndNormalize("../img1.jpg")
	if err != nil {
		fmt.Println("Error loading image 1:", err)
		return
	}

	B, boundsB, err := loadImageAndNormalize("../img2.jpg")
	if err != nil {
		fmt.Println("Error loading image 2:", err)
		return
	}

	if boundsA != boundsB {
		fmt.Println("Image dimensions do not match!")
		return
	}

	width, height := boundsA.Max.X, boundsA.Max.Y
	numWorkers := runtime.NumCPU()
	rowsPerWorker := height / numWorkers

	var wg sync.WaitGroup
	results := make([][]float32, numWorkers)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			start, end := i*rowsPerWorker, (i+1)*rowsPerWorker
			results[i] = multiplyMatrices(A, B, width, start, end)
		}(i)
	}

	wg.Wait()

	finalResult := make([]float32, width*height)
	for i, result := range results {
		copy(finalResult[i*rowsPerWorker*width:], result)
	}

	elapsedTime := time.Since(startTime)
	fmt.Printf("Bildverarbeitung und Matrixmultiplikation GOLANG %s\n", elapsedTime)
}
