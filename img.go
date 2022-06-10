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
	"time"
)

var width = 1920
var height = 1080

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

func lookupToNumber ( lookup []string) [][]int {
	result := make([][]int, len(lookup))
	for row := 0 ; row < len(lookup) ; row ++ {
		list := make([]int,len(lookup[0]))
		for col := 0 ; col < len(lookup[0]) ; col ++ {
			list[col],_ = strconv.Atoi(string(lookup[row][col]))
		}
		result[row] = list
	}
	return result
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

	fire1,err := getImageFromFilePath("assets/fire1.png")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fire2, err := getImageFromFilePath("assets/fire2.png")
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

	lookupTable := lookupToNumber(table)

	for col := 19 ; col >= 0 ; col -- {
		for row := 0 ; row < 20 ; row ++ {
			fromTable := lookupTable[row][col]
			drawAtCoord(rgbaImage,textures[fromTable % len(textures)],col,row,r)
		}
	}

	fires := make([][]bool,len(lookupTable))
	for i := 0 ; i < len(lookupTable) ; i ++ {
		fires[i] = make([]bool, len(lookupTable[0]))
	}


	canvasToWrite := canvas.NewRasterFromImage(rgbaImage)

	go func() {
		playFirst := true
		for true {
			for row := 0 ; row < len(lookupTable) ; row ++ {
				for col := 0 ; col < len(lookupTable[0]) ; col ++ {
					if fires[row][col] {
						drawAtCoord(rgbaImage,textures[lookupTable[row][col] % len(textures)],col,row,r)
						if playFirst {
							drawAtCoord(rgbaImage,fire1,col,row,r)
						}else {
							drawAtCoord(rgbaImage,fire2,col,row,r)
						}
					}
				}
			}
			playFirst = !playFirst
			canvas.Refresh(canvasToWrite)
			time.Sleep(500 * time.Millisecond)
		}
	}()

	imageWindow.SetContent(canvasToWrite)
	imageWindow.Resize(fyne.NewSize(float32(width), float32(height)))

	imageWindow.Canvas().SetOnTypedKey(func(event *fyne.KeyEvent) {
		fmt.Println(event.Name)
		newX := x
		newY := y
		drawNew := false
		switch event.Name {
		case "Left":
			if x > 0 {
				newX--
			}
			break
		case "Right":
			if x < len(lookupTable[0]) - 1 {
				newX++
			}
			break
		case "Up":
			if y > 0 {
				newY--
			}
			break
		case "Down":
			if y < len(lookupTable) - 1 {
				newY++
			}
			break
		case "Space":
			drawNew = true
			break
		}
		if event.Name == "Q" {
			fires[y][x] = !fires[y][x]
			if fires[y][x] {
				drawAtCoord(rgbaImage,fire1,newX,newY,r)
			}else {
				drawAtCoord(rgbaImage,textures[lookupTable[newY][newX] % len(textures)],newX,newY,r)
			}
		} else {
			drawAtCoord(rgbaImage, textures[lookupTable[y][x]%len(textures)], x, y, r)
			if fires[y][x] {
				drawAtCoord(rgbaImage,fire1,x,y,r)
			}
			if drawNew {
				lookupTable[newY][newX]++
				drawAtCoord(rgbaImage, textures[lookupTable[newY][newX]%len(textures)], newX, newY, r)
			}
			drawAtCoord(rgbaImage, red, newX, newY, r)
			x = newX
			y = newY
		}
		canvas.Refresh(canvasToWrite)
	})

	imageWindow.ShowAndRun()
}