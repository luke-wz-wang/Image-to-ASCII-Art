package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"os"
	"proj2/ConversionWorker"
	"proj2/WorkBalancer"
	"runtime"
	"strconv"
)

// generate a series of conversion request and send it to channel
func generateAndRequest(req chan ConversionWorker.Request, allFinished *bool) {
	reqDone := make(chan bool)
	// spawn a fix amount of conversion requests
	for i := 0; i < 10; i++{
		source := "./dataset/input/b"
		source += strconv.Itoa(i) + ".png"
		target := "./dataset/output/result" + strconv.Itoa(i) + ".txt"

		allDone := false
		// creator's work is done
		if i == 9 {
			allDone = true
		}

		flag := make(chan bool)
		req <- ConversionWorker.Request{source, target, allDone,reqDone, flag}
		<-reqDone
	}
}


func printUsage() {
	fmt.Printf("Usage: editor [-p=[number of threads]]\n")
	fmt.Printf("\t-p=[number of workers] = An optional flag to run main.go in its parallel version.\nCall and pass the value to runtime.GOMAXPROCS(...) \nDenote the number of workers in work balancing mode.")
	fmt.Printf("\t-m=[magnitude of conversion requests] = A required flag that represents \n the magnitude of the image conversion requests. \nA total of 5*m image conversion request will be generated.")
}

func main() {

	var p int
	flag.IntVar(&p, "p", 0, "An optional flag to run main.go in its parallel version.\nCall and pass the value to runtime.GOMAXPROCS(...) \nDenote the number of workers in work balancing mode.")

	var m int
	flag.IntVar(&m, "m", 10, "A required flag that represents the magnitude of the image conversion requests.\n A total of 5*m image conversion request will be generated.")
	flag.Parse()

	if p == 0{
		// calculate actual number of requests to imitate for sequential version based on the magnitude
		totalRequests := m*10
		sequential(totalRequests)
	}else{
		numOfWorkers := p
		magnitude := m

		runtime.GOMAXPROCS(p)
		fmt.Println("Work distribution status among workers:")
		for i :=0; i < p; i++{
			fmt.Printf("%d ", i)
		}
		fmt.Printf("\n")
		fmt.Println("---------------------")

		work := make(chan ConversionWorker.Request)
		allFinished := false
		for i := 0; i < magnitude; i++ {
			// spawn a goroutine to generate conversion request
			go generateAndRequest(work, &allFinished)
		}
		WorkBalancer.Init(&allFinished, magnitude, numOfWorkers).BalanceWork(work, &allFinished)
	}
}

// sequential version
func sequential(totalRequests int){
	// spawn #totalRequests conversion requests
	for i := 0; i < totalRequests; i++{
		source := "./dataset/input/b"
		source += strconv.Itoa(i%5) + ".png"
		target := "./dataset/output/result" + strconv.Itoa(i%5) + ".txt"
		seqImageToAscii(source, target)
	}
}

func resize(img image.Image, newX int) *image.RGBA {

	bounds := img.Bounds()
	curX := bounds.Dx()
	curY := bounds.Dy()
	newImg := image.NewRGBA(image.Rect(0, 0, newX, newX*curY/curX))
	return newImg
}

func seqImageToAscii(source string, target string) {

	// symbol lists
	symList := []string{"Q", "A", "D", "E", "@", "C", ">", "!", "3", "-", "?", ";", ".", ":", "-"}

	inReader, err := os.Open(source)
	if err != nil {
		return
	}
	defer inReader.Close()

	img, _ := png.Decode(inReader)

	if img.Bounds().Dx() > 300 {
		img = resize(img, 300)
	}
	bounds := img.Bounds()
	X := bounds.Dx()
	Y := bounds.Dy()

	fileName := target

	outputTxt, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err.Error())
		return
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

}