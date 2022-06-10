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

func drawAtCoord (onto  *image.RGBA, from image.Image, x int, y int, bounds image.Rectangle) {
	draw.Draw(onto,bounds,from,permutePoint(image.Point{x,y}),draw.Over)
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

	red,err := getImageFromFilePath("assets/red.png")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	x := 0
	y := 0



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
			drawAtCoord(rgbaImage,textures[fromTable],col,row,r)
		}
	}
	canvasToWrite := canvas.NewRasterFromImage(rgbaImage)

	imageWindow.SetContent(canvasToWrite)
	imageWindow.Resize(fyne.NewSize(float32(width), float32(height)))

	imageWindow.Canvas().SetOnTypedKey(func(event *fyne.KeyEvent) {
		fmt.Println(event.Name)
		newX := x
		newY := y
		switch event.Name {
		case "Left":
			newX--
			break
		case "Right":
			newX++
			break
		case "Up":
			newY--
			break
		case "Down":
			newY++
			break
		}
		fromTable,_ := strconv.Atoi(string(table[y][x]))
		drawAtCoord(rgbaImage,textures[fromTable],x,y,r)
		drawAtCoord(rgbaImage,red,newX,newY,r)
		x = newX
		y = newY
		canvas.Refresh(canvasToWrite)
	})

	imageWindow.ShowAndRun()
}