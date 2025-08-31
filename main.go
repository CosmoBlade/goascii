package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strconv"
	"strings"
)

func main() {
	args := os.Args[1:]

	if len(args) < 2 {
		fmt.Println("Usage: goascii <input.png/jpg> <output.txt> [--width=100] [--height=50]")
		fmt.Println("Example: goascii image.png output.txt --width=100 --height=50")
		os.Exit(1)
	}

	inputFile := args[0]
	outputFile := args[1]

	width := 0
	height := 0

	for i := 2; i < len(args); i++ {
		arg := args[i]

		if strings.HasPrefix(arg, "--width=") {
			if val, err := strconv.Atoi(arg[8:]); err == nil {
				width = val
			}
		} else if strings.HasPrefix(arg, "--height=") {
			if val, err := strconv.Atoi(arg[9:]); err == nil {
				height = val
			}
		} else {
			fmt.Printf("Unknown argument: %s\n", arg)
			os.Exit(1)
		}
	}

	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Printf("Error decoding image: %v\n", err)
		os.Exit(1)
	}

	bounds := img.Bounds()
	origWidth, origHeight := bounds.Max.X, bounds.Max.Y

	targetWidth := width
	targetHeight := height

	if targetWidth == 0 && targetHeight == 0 {
		targetWidth = origWidth
		targetHeight = origHeight
	} else if targetWidth == 0 {
		targetWidth = int(float64(targetHeight) * float64(origWidth) / float64(origHeight))
	} else if targetHeight == 0 {
		targetHeight = int(float64(targetWidth) * float64(origHeight) / float64(origWidth))
	}

	fmt.Printf("Original size: %dx%d\n", origWidth, origHeight)
	fmt.Printf("Target size: %dx%d\n", targetWidth, targetHeight)

	asciiChars := []string{"@", "#", "S", "%", "?", "*", "+", ";", ":", ",", "."}

	outFile, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		os.Exit(1)
	}
	defer outFile.Close()

	var asciiArt strings.Builder

	xStep := float64(origWidth) / float64(targetWidth)
	yStep := float64(origHeight) / float64(targetHeight)

	for y := 0; y < targetHeight; y++ {
		for x := 0; x < targetWidth; x++ {
			startX := int(float64(x) * xStep)
			startY := int(float64(y) * yStep)
			endX := int(float64(x+1) * xStep)
			endY := int(float64(y+1) * yStep)

			if endX > origWidth {
				endX = origWidth
			}
			if endY > origHeight {
				endY = origHeight
			}

			var totalR, totalG, totalB, count float64
			for origY := startY; origY < endY; origY++ {
				for origX := startX; origX < endX; origX++ {
					color := img.At(origX, origY)
					r, g, b, _ := color.RGBA()
					totalR += float64(r / 256)
					totalG += float64(g / 256)
					totalB += float64(b / 256)
					count++
				}
			}

			avgR := totalR / count
			avgG := totalG / count
			avgB := totalB / count

			gray := 0.299*avgR + 0.587*avgG + 0.114*avgB

			charIndex := int(gray * float64(len(asciiChars)-1) / 255)
			if charIndex < 0 {
				charIndex = 0
			}
			if charIndex >= len(asciiChars) {
				charIndex = len(asciiChars) - 1
			}

			asciiArt.WriteString(asciiChars[charIndex])
		}
		asciiArt.WriteString("\n")
	}

	_, err = outFile.WriteString(asciiArt.String())
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Success! ASCII art saved to %s\n", outputFile)
	fmt.Printf("ASCII size: %dx%d characters\n", targetWidth, targetHeight)
}
