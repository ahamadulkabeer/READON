package helpers

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"os"
	"time"
)

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

func cropToRatio(img image.Image) image.Image {
	ratio := 0.75
	bounds := img.Bounds()
	targetWidth := float64(bounds.Dy()) * ratio
	offsetX := (float64(bounds.Dx()) - targetWidth) / 2

	rect := image.Rect(
		int(offsetX),
		0,
		int(offsetX+targetWidth),
		bounds.Dy(),
	)

	croppedImg := image.NewRGBA(rect)
	draw.Draw(croppedImg, croppedImg.Bounds(), img, rect.Min, draw.Src)
	return croppedImg
}

// Example code for cropping an image to a specific ratio
/*func cropToRatio(inputImage image.Image, ratio float64) image.Image {

	// Calculate the dimensions for the cropped image based on the desired ratio
	// For example, you can maintain the width and adjust the height accordingly
	width := 1600 // Desired width
	height := int(float64(width) / ratio)

	// Create a new RGBA image with the desired dimensions
	croppedImage := image.NewRGBA(image.Rect(0, 0, width, height))

	// Use an image processing library to crop the input image to the new dimensions
	// and draw it onto the croppedImage

	return croppedImage
}*/

func CropImage(imagetocrop []byte) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(imagetocrop))
	if err != nil {
		fmt.Println("Error decoding the image:", err)
		return nil, err
	}
	croppedImage := cropToRatio(img)
	// Create a buffer to store the encoded image
	var buffer bytes.Buffer

	// Encode the cropped image as JPEG (or your desired format)
	err = jpeg.Encode(&buffer, croppedImage, nil)
	if err != nil {
		fmt.Println("Error encoding the cropped image:", err)
		return nil, err
	}

	// Get the bytes from the buffer
	croppedImageData := buffer.Bytes()

	// Save the cropped image
	err = SaveCroppedImage(croppedImageData)
	if err != nil {
		return nil, err
	}

	return croppedImageData, nil
}

func SaveCroppedImage(croppedImageData []byte) error {
	outputFilePath := "/home/kabeer/Documents/READON/img/cropped/"
	// Create or open the output file for writing
	filename := "IMG" + time.Now().Format("2006-01-02_03-04-05PM") + ".jpeg"
	outputFile, err := os.Create(outputFilePath + filename)
	if err != nil {
		fmt.Println("Error creating the output image file:", err)
		return err
	}
	defer outputFile.Close()

	// Write the croppedImageData to the file
	_, err = outputFile.Write(croppedImageData)
	if err != nil {
		fmt.Println("Error writing the cropped image data to the file:", err)
		return err
	}

	fmt.Println("Cropped image saved as", outputFilePath)
	return nil
}
