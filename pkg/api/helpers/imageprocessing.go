package helpers

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"os"

	"github.com/fogleman/gg"
)

func cropToRatio(img image.Image) image.Image {
	// Open the image
	//imgPath := "your_image.jpg" // Change this to the path of your image
	//img, err := gg.LoadImage(imgPath)
	//if err != nil {
	//	panic(err)
	//}

	// Calculate the dimensions for the cropped image
	targetWidth := 1.5 * float64(img.Bounds().Dx())
	targetHeight := float64(img.Bounds().Dy())

	// Initialize a new context with the desired dimensions
	ctx := gg.NewContext(int(targetWidth), int(targetHeight))

	// Calculate the cropping area
	cropX := (float64(img.Bounds().Dx()) - targetWidth) / 2
	cropY := 0.0

	// Draw the cropped portion to the new context
	ctx.DrawImage(img, int(-cropX), int(-cropY))

	// Save the cropped image to a new file
	outputFilePath := "/home/kabeer/Documents/READON/img/cropped_image001.jpg" // Change this to the desired output file path
	ctx.SavePNG(outputFilePath)
	if err := ctx.SavePNG(outputFilePath); err != nil {
		panic(err)
	}

	// Print a success message
	//println("Image cropped and saved to", outputFilePath)
	return img
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
	err = jpeg.Encode(&buffer, croppedImage, nil) // Use "image/png" for PNG encoding
	if err != nil {
		fmt.Println("Error encoding the cropped image:", err)
		return nil, err
	}

	// Get the bytes from the buffer
	croppedImageData := buffer.Bytes()

	outputFilePath := "/home/kabeer/Documents/READON/img/cropped_image002.jpg"

	// Save the cropped image
	err = SaveCroppedImage(croppedImageData, outputFilePath)
	if err != nil {
		return nil, err
	}

	return croppedImageData, nil
}

func SaveCroppedImage(croppedImageData []byte, outputFilePath string) error {
	// Create or open the output file for writing
	outputFile, err := os.Create(outputFilePath)
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
