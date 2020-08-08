package ImageConversion

import (
	"fmt"
	"image"
	"image/png"
	"os"
)

// resize image and returns a new RGBA img
func Resize(img image.Image, newX int) *image.RGBA {
	bounds := img.Bounds()
	curX := bounds.Dx()
	curY := bounds.Dy()
	newImg := image.NewRGBA(image.Rect(0, 0, newX, newX*curY/curX))
	return newImg
}

// convert image to ascii code
func ImageToAscii(source string, outpath string, allDone bool, allFinished *bool) bool{

	// symbol lists
	symList := []string{"Q", "A", "D", "E", "@", "C", ">", "!", "3", "-", "?", ";", ".", ":", "-"}

	// load img
	inReader, err := os.Open(source)
	if err != nil {
		return false
	}
	defer inReader.Close()

	img, _ := png.Decode(inReader)
	target := outpath

	if img == nil{
		return false
	}

	// resize img if too wide
	if img.Bounds().Dx() > 250 {
		img = Resize(img, 250)
	}

	bounds := img.Bounds()
	X := bounds.Dx()
	Y := bounds.Dy()

	fileName := target

	outputTxt, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	// convert pixel to symbol
	for i := 0; i < Y; i++ {
		for j := 0; j < X; j++ {
			// calculate symbol
			color := img.At(j, i)
			_, g, _, _ := color.RGBA()
			average := uint8(g >> 8)
			idx := average / 18
			outputTxt.WriteString(symList[idx])
			// line break for each line of pixels
			if j == X-1 {
				outputTxt.WriteString("\n")
			}
		}
	}
	outputTxt.Close()

	// set allFinished flag to be true
	if allDone{
		*allFinished = true
	}
	return true
}
