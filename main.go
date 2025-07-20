package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

func main() {
	// Define command-line flags
	inputFile := flag.String("input", "", "Input PDF file")
	outputFile := flag.String("output", "", "Output PDF file")
	fromPage := flag.Int("from-page", 1, "Page of the PDF to start stamping from")
	color := flag.String("color", "0 0 0", "Color for the page numbers (e.g., '1 0 0' for red)")
	size := flag.Int("size", 12, "Font size for the page numbers")
	position := flag.String("position", "bc", "Position of the page numbers (e.g., 'tl', 'tc', 'tr', 'l', 'c', 'r', 'bl', 'bc', 'br')")
	verbose := flag.Bool("verbose", false, "Whether to log stuff")

	flag.Parse()

	err := paginate(inputFile, outputFile, fromPage, position, size, color, verbose)
	if err != nil {
		log.Fatal(err)
	}

}

// paginate
func paginate(inputFile *string, outputFile *string, fromPage *int, position *string, size *int, color *string, verbose *bool) error {

	if *inputFile == "" || *outputFile == "" {
		fmt.Printf("%s %s\n", *inputFile, *outputFile)
		return fmt.Errorf("Input and output files must be specified")
	}

	// Create a single watermark object to be used as a stamp.
	// We use the "%p" placeholder, which pdfcpu replaces with the page number during processing.
	// This is the key to mimicking the `stamp` command's behavior.
	stamp, err := createStamp("%p", *position, *size, *color)
	if err != nil {
		return fmt.Errorf("Error creating stamp: %v", err)
	}

	// Use AddWatermarksFile to apply the single stamp object to all selected pages.
	// This function correctly handles shared page resources when used with a placeholder,
	// preventing the "jumbled" text issue.
	if *fromPage > 1 {
		err = api.AddWatermarksFile(*inputFile, *outputFile, []string{fmt.Sprintf("%d-", *fromPage)}, stamp, nil)
	} else {

		err = api.AddWatermarksFile(*inputFile, *outputFile, nil, stamp, nil)
	}
	if err != nil {
		return fmt.Errorf("Error stamping PDF: %v", err)
	}

	if *verbose {
		fmt.Printf("Successfully stamped pages in %s and saved to %s\n", *inputFile, *outputFile)
	}
	return nil
}

// createStamp creates a watermark object using the documented api.TextWatermark function.
func createStamp(text, position string, size int, color string) (*model.Watermark, error) {
	// The description string configures the appearance of the stamp.
	// `onTop` is true to ensure it's a stamp (overlay) rather than a watermark (underlay).
	desc := fmt.Sprintf(`pos:%s, scale:1.0 abs, rot:0, offset: -2 -2, color:%s, font:Helvetica, points:%d`, position, color, size)
	return api.TextWatermark(text, desc, true, false, types.POINTS)
}
