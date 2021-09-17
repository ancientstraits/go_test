package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"encoding/hex"

	"github.com/golang/freetype"
	"golang.org/x/image/font"
)

var (
	dpi      = float64(72)
	fontfile = "ttf/FiraCode-Regular.ttf"
	hinting  = "full"
	size     = float64(12)
	spacing  = float64(1.5)
	wonb     = true
)

func rectColor(code string) *image.Uniform {
	dat, err := hex.DecodeString(code)
	if err != nil {
		panic(err)
	}
	return &image.Uniform{color.RGBA{dat[0], dat[1], dat[2], 0xff}}
}

func rectColorAlpha(code string) *image.Uniform {
	dat, err := hex.DecodeString(code)
	if err != nil {
		panic(err)
	}
	return &image.Uniform{color.RGBA{dat[0], dat[1], dat[2], dat[3]}}
}

func read(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scan := bufio.NewScanner(file)
	content := []string{}
	for scan.Scan() {
		content = append(content, scan.Text())
	}

	return content
}

func genIMG(str string) {
	flag.Parse()

	// Read the font data.
	fontBytes, err := ioutil.ReadFile(fontfile)
	if err != nil {
		log.Println(err)
		return
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		return
	}

	// Initialize the context.
	ruler := color.RGBA{0x00, 0x00, 0x00, 0xff}
	if wonb {
		ruler = color.RGBA{0xff, 0xff, 0xff, 0xff}
	}
	lnNumColor := rectColor("ffffbb")
	txtColor   := rectColor("ffffff")
	rgba := image.NewRGBA(image.Rect(0, 0, 640, 480))
	draw.Draw(
		rgba,
		image.Rect(0, 0, 37, 840),
		// &image.Uniform{color.RGBA{0x33, 0x33, 0x77, 0xff}},
		rectColor("333377"),
		image.Point{0, 0},
		draw.Src,
	)
	draw.Draw(
		rgba,
		image.Rect(37, 0, 640 - 37, 840),
		&image.Uniform{color.RGBA{0x77, 0x77, 0x33, 0xff}},
		image.Point{0, 0},
		draw.Src,
	)
	c := freetype.NewContext()
	c.SetDPI(dpi)
	c.SetFont(f)
	c.SetFontSize(size)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	switch hinting {
	default:
		c.SetHinting(font.HintingNone)
	case "full":
		c.SetHinting(font.HintingFull)
	}

	// Draw the guidelines.
	for i := 0; i < 480; i++ {
		rgba.Set(37, i, ruler)
	}

	// Draw the text.
	pt := freetype.Pt(10, 10+int(c.PointToFixed(size)>>6))
	txt := read(str)
	for ln, str := range txt {
		c.SetSrc(lnNumColor)
		_, err = c.DrawString(fmt.Sprintf("%3d  ", ln + 1), pt)
		if err != nil {
			log.Println(err)
			return
		}
		c.SetSrc(txtColor)
		_, err = c.DrawString(
			"     " + strings.Replace(str, "\t", "    ", -1),
		pt)
		if err != nil {
			log.Println(err)
			return
		}
		pt.Y += c.PointToFixed(size * spacing)
	}

	// Save that RGBA image to disk.
	outFile, err := os.Create("out.png")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer outFile.Close()
	b := bufio.NewWriter(outFile)
	err = png.Encode(b, rgba)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	err = b.Flush()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	fmt.Println("Wrote out.png OK.")
}

func main() {
	genIMG("main.go")
}
