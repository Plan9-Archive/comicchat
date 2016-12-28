package main

import (
	"bufio"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"

	"github.com/golang/freetype/truetype"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"golang.org/x/image/draw"
	"golang.org/x/image/math/fixed"
)

var (
	fontname = "gooddog-plain.regular"

	comicfontdata = draw2d.FontData{
		Name:   fontname,
		Family: draw2d.FontFamilySans,
		Style:  draw2d.FontStyleNormal,
	}

	comicfont *truetype.Font

	// username -> facename
	userface = map[string]string{}

	// facename -> image
	faces = map[string]image.Image{}
	// slice of facenames
	facekeys []string
)

func init() {
	// seed random
	rand.Seed(time.Now().Unix())

	// load font
	fnt, err := ioutil.ReadFile(fontname + ".ttf")
	if err != nil {
		panic(err)
	}
	comicfont, err = truetype.Parse(fnt)
	if err != nil {
		panic(err)
	}

	draw2d.RegisterFont(comicfontdata, comicfont)

	// load faces
	if fis, err := ioutil.ReadDir("faces/"); err != nil {
		panic(err)
	} else {
		for _, fi := range fis {
			if i := loadFromPngFile("faces/" + fi.Name()); i != nil {
				base := strings.TrimSuffix(fi.Name(), path.Ext(fi.Name()))
				log.Printf("loaded face %s", base)

				faces[base] = i
				facekeys = append(facekeys, base)
			}
		}
	}
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
		return nil
	}
	//log.Printf("Read %s OK.\n", filePath)
	return i
}

func makeusercomic(who, text string) image.Image {
	var facename string
	if fname, ok := userface[who]; !ok {
		// assign random face
		userface[who] = facekeys[rand.Intn(len(facekeys))]
		facename = userface[who]
	} else {
		facename = fname
	}

	i, _ := makecomic(text, facename)
	return i
}

func GetFontHeight(fd draw2d.FontData, size float64) (height float64) {
	font := draw2d.GetFont(fd)
	fupe := font.FUnitsPerEm()
	bounds := font.Bounds(fixed.Int26_6(fupe))
	height = float64(bounds.Max.Y-bounds.Min.Y) * size / float64(fupe)
	return
}
func RenderString(text string, fd draw2d.FontData, size float64, fill color.Color) (buffer image.Image) {

	const stretchFactor = 1.2

	height := GetFontHeight(fd, size) * stretchFactor
	widthMax := float64(len(text)) * size

	buf := image.NewRGBA(image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{int(widthMax + 1), int(height + 1)},
	})

	gc := draw2dimg.NewGraphicContext(buf)
	gc.Translate(0, height/stretchFactor)
	gc.SetFontData(fd)
	gc.SetFontSize(size)
	gc.SetStrokeColor(color.Black)
	gc.SetFillColor(fill)
	width := gc.FillStringAt(text, 1, 1)
	gc.StrokeStringAt(text, 1, 1)

	buffer = buf.SubImage(image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{int(width + 1), int(height + 1)},
	})

	return
}

// turn a string into a slice of strings based on a width limit
// TODO: make utf8-safe
func wrap(str string, width int) []string {
	var out []string

	lines := strings.Split(str, "\n")

	for _, l := range lines {

		words := strings.Split(l, " ")
		if len(words) == 0 {
			return out
		}

		// current line we are making
		current := words[0]

		// # spaces left before the end
		remaining := width - len(current)

		for _, word := range words[1:] {
			if len(word)+1 > remaining {
				out = append(out, current)
				current = word
				remaining = width - len(word)
			} else {
				current += " " + word
				remaining -= 1 + len(word)
			}
		}

		out = append(out, current)
	}

	return out
}

var iRect = image.Rect(0, 0, 300, 300)

func makecomic(text, face string) (image.Image, error) {
	facei := faces[face]
	dest := image.NewRGBA(iRect)
	b := iRect

	draw2dimg.DrawImage(facei, dest, draw2d.NewIdentityMatrix(), draw.Over, draw2dimg.BilinearFilter)

	size := 24.0
	lines := wrap(text, 25)
	height := GetFontHeight(comicfontdata, size)
	var ilines []image.Image
	for _, l := range lines {
		ilines = append(ilines, RenderString(l, comicfontdata, size, color.White))
	}

	tr := draw2d.NewIdentityMatrix()
	tr.Translate(25.0, float64(b.Dy())-10.0)
	tr.Translate(0, -float64(len(ilines))*height)
	for _, i := range ilines {
		draw2dimg.DrawImage(i, dest, tr, draw.Over, draw2dimg.BilinearFilter)
		tr.Translate(0, height)
	}

	//	gc.StrokeStringAt(text, float64(20), float64((b.Dy())-20))
	//	gc.FillStringAt(text, float64(20), float64((b.Dy())-20))

	return dest, nil
}
