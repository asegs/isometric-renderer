package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"image"
	"image/draw"
	"image/png"
	"os"
	"strconv"
	"strings"
)

var width = 1920
var height = 1080

var lineWidth = 100
var fps = 3

func getImageFromFilePath(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, err := png.Decode(f)
	return img, err
}


func permutePoint (point image.Point) image.Point {
	offsetX := -100
	offsetY := -200

	return image.Point{
		X: -1 * (point.X * 40) + (point.X - point.Y) * (40 * 0.5) + offsetX,
		Y: -1 * (point.Y * 20) + (point.X + point.Y) * (20 * 0.5 ) + offsetY,
	}
}
func main () {
	topLeft := image.Point{0, 0}
	bottomRight := image.Point{width, height}

	rgbaImage := image.NewRGBA(image.Rectangle{topLeft, bottomRight})

	mainApp := app.New()
	imageWindow := mainApp.NewWindow("Images")

	txt,_ := ReadToString("assets/lookup.txt")
	table := strings.Split(txt,"\n")



	grass,err := getImageFromFilePath("assets/grass.png")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	water,err := getImageFromFilePath("assets/water.png")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	trees,err := getImageFromFilePath("assets/trees.png")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	sand,err := getImageFromFilePath("assets/sand.png")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	r := image.Rectangle{
		Min: image.Point{0,0},
		Max: image.Point{width,height},
	}

	textures := make([]image.Image, 4)
	textures[0] = grass
	textures[1] = water
	textures[2] = trees
	textures[3] = sand

	for col := 19 ; col >= 0 ; col -- {
		for row := 0 ; row < 20 ; row ++ {
			fromTable,_ := strconv.Atoi(string(table[row][col]))
			draw.Draw(rgbaImage,r,textures[fromTable],permutePoint(image.Point{col,row}),draw.Over)
		}
	}
	canvasToWrite := canvas.NewRasterFromImage(rgbaImage)

	imageWindow.SetContent(canvasToWrite)
	imageWindow.Resize(fyne.NewSize(float32(width), float32(height)))

	//go func() {
	//	pos := 0
	//	for true {
	//		drawImageRed(pos)
	//		pos ++
	//		canvas.Refresh(canvasToWrite)
	//		time.Sleep(time.Duration(1000 / fps) * time.Millisecond)
	//	}
	//}()

	imageWindow.ShowAndRun()
}