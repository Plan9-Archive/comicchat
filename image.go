package main

import (
	"bufio"
	"code.google.com/p/draw2d/draw2d"
	"code.google.com/p/freetype-go/freetype/truetype"
	"image"
	_ "image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"
)

var (
	fontname = "gooddog-plain.regular"
	pic      = "guy.jpg"

	comicfontdata = draw2d.FontData{
		Name:   fontname,
		Family: draw2d.FontFamilySans,
		Style:  draw2d.FontStyleNormal,
	}
)

var comicfont *truetype.Font

func init() {
	fnt, err := ioutil.ReadFile(fontname + ".ttf")
	if err != nil {
		panic(err)
	}
	comicfont, err = truetype.Parse(fnt)
	if err != nil {
		panic(err)
	}

	draw2d.RegisterFont(comicfontdata, comicfont)
}

func saveToPngFile(filePath string, m image.Image) {
	f, err := os.Create(filePath)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()
	b := bufio.NewWriter(f)
	err = png.Encode(b, m)
	if err != nil {
		log.Println(err)
		return
	}
	err = b.Flush()
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("Wrote %s OK.\n", filePath)
}

func loadFromPngFile(filePath string) image.Image {
	f, err := os.OpenFile(filePath, 0, 0)
	if f == nil {
		log.Printf("can't open file; err=%s\n", err)
		return nil
	}
	defer f.Close()
	b := bufio.NewReader(f)
	i, _, err := image.Decode(b)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Printf("Read %s OK.\n", filePath)
	return i
}

func makecomic(text string) (image.Image, error) {
	source := loadFromPngFile(pic)
	b := source.Bounds()
	dest := image.NewRGBA(b)
	gc := draw2d.NewGraphicContext(dest)

	gc.SetFontData(comicfontdata)
	gc.SetFontSize(18.0)

	gc.DrawImage(source)

	gc.StrokeStringAt(text, float64(20), float64((b.Dy())-20))
	gc.FillStringAt(text, float64(20), float64((b.Dy())-20))

	return dest, nil
}
