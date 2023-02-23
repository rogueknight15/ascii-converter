package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

type Color struct {
	r, g, b, a uint32
}

func (c *Color) GetBrightness() uint32 {
	return (2*c.r + 3*c.g + c.b) / 6
}

func (c *Color) AsCharacter() string {
	brightness := c.GetBrightness()
	if brightness <= 25 {
		return " "
	} else if brightness <= 50 {
		return "."
	} else if brightness <= 75 {
		return ":"
	} else if brightness <= 100 {
		return "-"
	} else if brightness <= 125 {
		return "="
	} else if brightness <= 150 {
		return "+"
	} else if brightness <= 175 {
		return "*"
	} else if brightness <= 200 {
		return "#"
	} else if brightness <= 225 {
		return "%"
	} else {
		return "@"
	}
}

type AsciiMode int

const (
	ASCII_MODE_PLAIN AsciiMode = iota
	ASCII_MODE_FORE
	ASCII_MODE_BACK
)

type Image struct {
	data   []Color
	width  int
	height int
}

func (i Image) At(x, y int) color.Color {
	c := i.data[x+y*i.width]
	return color.RGBA{uint8(c.r), uint8(c.g), uint8(c.b), uint8(c.a)}
}

func (i Image) Bounds() image.Rectangle {
	return image.Rectangle{image.Point{0, 0}, image.Point{i.width, i.height}}
}

func (i Image) ColorModel() color.Model {
	return color.RGBAModel
}

func LoadImage(filepath string) Image {
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	image, err := png.Decode(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	width, height := image.Bounds().Dx(), image.Bounds().Dy()

	pictureData := make([]Color, width*height)
	for y := 0; y < image.Bounds().Dy(); y++ {
		for x := 0; x < image.Bounds().Dx(); x++ {
			r, g, b, a := image.At(x, y).RGBA()
			newR, newG, newB, newA := r/256, g/256, b/256, a/256
			pictureData[x+y*width] = Color{newR, newG, newB, newA}
		}
	}

	return Image{
		data:   pictureData,
		width:  width,
		height: height,
	}
}

func (i *Image) ToAscii(mode AsciiMode) string {
	picture := ""
	for y := 0; y < i.height; y++ {
		for x := 0; x < i.width; x++ {
			color := i.data[x+y*i.width]
			if mode == ASCII_MODE_PLAIN {
				picture += color.AsCharacter()
			} else if mode == ASCII_MODE_BACK {
				picture += fmt.Sprintf("\033[48;2;%d;%d;%dm", color.r, color.g, color.b) + color.AsCharacter() + "\033[0m"
			} else if mode == ASCII_MODE_FORE {
				picture += fmt.Sprintf("\033[38;2;%d;%d;%dm", color.r, color.g, color.b) + color.AsCharacter() + "\033[0m"
			}
		}
		picture += "\n"
	}
	return picture
}

func (i *Image) Resize(xScale float64, yScale float64) {
	newWidth := int(float64(i.width) * xScale)
	newHeight := int(float64(i.height) * yScale)
	newImage := Image{
		make([]Color, newWidth*newHeight),
		newWidth,
		newHeight,
	}
	for y := 0; y < newImage.height; y++ {
		for x := 0; x < newImage.width; x++ {
			px := float64(x) / xScale
			py := float64(y) / yScale
			newImage.data[x+y*newImage.width] = i.data[int(px)+int(py)*i.width]
		}
	}
	*i = newImage
}

func (i *Image) SavePNG(filepath string) {
	file, err := os.Create(filepath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = png.Encode(file, *i)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (i *Image) SaveAscii(filepath string, mode AsciiMode) {
	file, err := os.Create(filepath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	_, err = file.WriteString(i.ToAscii(mode))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	file.Close()
}

func main() {
	filepath := flag.String("f", "", "Please specify a path to a .png file")
	xFactor := flag.Float64("xscale", 1.0, "Please specify a scaling factor for the x axis")
	yFactor := flag.Float64("yscale", 1.0, "Please specify a scaling factor for the y axis")
	outputFormat := flag.Int64("format", 0, "Please specify the output format (0 = NO_COLOR, 1 = FOREGROUD_COLORED, 2 = BACKGROUND_COLORED, 3 = PNG)")
	outputFilepath := flag.String("o", "out.txt", "Please specify an output filepath for the converted image")
	flag.Parse()

	image := LoadImage(*filepath)
	image.Resize(*xFactor, *yFactor)

	switch *outputFormat {
	case 0:
		image.SaveAscii(*outputFilepath, ASCII_MODE_PLAIN)
	case 1:
		image.SaveAscii(*outputFilepath, ASCII_MODE_FORE)
	case 2:
		image.SaveAscii(*outputFilepath, ASCII_MODE_BACK)
	case 3:
		image.SavePNG(*outputFilepath)
	default:
		{
			fmt.Println("Invalid output format: " + "Please specify the output format (0 = NO_COLOR, 1 = FOREGROUD_COLORED, 2 = BACKGROUND_COLORED, 3 = PNG)")
			os.Exit(1)
		}
	}

	fmt.Println("File saved to: " + *outputFilepath)
}
