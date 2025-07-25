package main

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

// import "C"

func main() {
	// Define command-line flags
	inputFile := flag.String("input", "", "Input PDF file")
	outputFile := flag.String("output", "", "Output PDF file")
	startPage := flag.Int("start-page", 1, "Page number to start pagination with")
	fromPage := flag.Int("from-page", 1, "Page of the PDF to start stamping from")
	color := flag.String("color", "0 0 0", "Color for the page numbers (e.g., '1 0 0' for red)")
	size := flag.Int("size", 12, "Font size for the page numbers")
	position := flag.String("position", "bc", "Position of the page numbers (e.g., 'tl', 'tc', 'tr', 'l', 'c', 'r', 'bl', 'bc', 'br')")
	verbose := flag.Bool("verbose", false, "Whether to log stuff")

	flag.Parse()

	pageCount, err := paginate(inputFile, outputFile, startPage, fromPage, position, size, color)
	if err != nil {
		return
	}
	if *verbose {
		fmt.Printf("Successfully stamped %d pages in %s and saved to %s\n", pageCount-*fromPage+1, *inputFile, *outputFile)
	}
}

func paginate(inputFile *string, outputFile *string, startPage *int, fromPage *int, position *string, size *int, color *string) (int, error) {

	if *inputFile == "" || *outputFile == "" {
		return 0, fmt.Errorf("input and output files must be specified")
	}

	// Get the total number of pages in the PDF
	pageCount, err := api.PageCountFile(*inputFile)
	if err != nil {
		return 0, fmt.Errorf("error getting page count: %w", err)
	}

	// Create a map of watermarks, one for each page to be stamped
	watermarks := make(map[int]*model.Watermark)

	for i := *fromPage; i <= pageCount; i++ {
		pageNum := *startPage + (i - *fromPage)
		wm, err := createWatermark(strconv.Itoa(pageNum), *position, *size, *color)
		if err != nil {
			return 0, fmt.Errorf("error creating watermark: %w", err)
		}
		watermarks[i] = wm
	}

	// Add watermarks to the PDF
	err = api.AddWatermarksMapFile(*inputFile, *outputFile, watermarks, nil)
	if err != nil {
		return 0, fmt.Errorf("error stamping PDF: %w", err)
	}
	return pageCount, nil
}

func createWatermark(text string, position string, size int, color string) (*model.Watermark, error) {
	// Description format: "text:..., pos:..., sc:..., rot:..., op:..., color:..."
	desc := fmt.Sprintf(`pos:%s, scale:1.0 abs, rot:0, offset:-2 -2, color:%s, font:Helvetica, points:%d`, position, color, size)
	return api.TextWatermark(text, desc, true, false, types.POINTS)
}
