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
	"math/rand"
	"os"
	"path"
	"strings"
	"time"
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

func makecomic(text, face string) (image.Image, error) {
	facei := faces[face]
	dest := image.NewRGBA(image.Rect(0, 0, 300, 300))
	b := dest.Bounds()
	gc := draw2d.NewGraphicContext(dest)

	gc.SetFontData(comicfontdata)
	gc.SetFontSize(18.0)

	gc.DrawImage(facei)

	gc.StrokeStringAt(text, float64(20), float64((b.Dy())-20))
	gc.FillStringAt(text, float64(20), float64((b.Dy())-20))

	return dest, nil
}
