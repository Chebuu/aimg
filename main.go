package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/Knorkebrot/ansirgb"
	"github.com/monochromegane/terminal"
	flag "github.com/ogier/pflag"
)

type Block struct {
	nocolor bool
	top     *ansirgb.Color
	bottom  *ansirgb.Color
}

func (b *Block) String() string {
	ret := ""
	if !b.nocolor {
		ret += b.bottom.Bg()
	}
	if b.top != nil {
		if !b.nocolor {
			ret += b.top.Fg()
		}
		// If it's not a UTF-8 terminal, fall back to '#'
		if strings.Contains(os.Getenv("LC_ALL"), "UTF-8") ||
			strings.Contains(os.Getenv("LANG"), "UTF-8") {
			ret += "\u2580"
		} else {
			ret += "#"
		}
	} else {
		ret += " "
	}
	return ret
}

func main() {
	var widthstr string
	flag.StringVarP(&widthstr, "width", "w", "100%", "Output width. Supports column count and percentage.")
	flag.Usage = usage
	flag.Parse()

	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	width := getColumns(widthstr) - 1 // -1 for the reset column
	for _, fpath := range flag.Args() {
		fh, err := os.Open(fpath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(10)
		}

		img, _, err := image.Decode(fh)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			fh.Close()
			os.Exit(20)
		}
		fh.Close()

		imgWidth := img.Bounds().Dx()
		imgHeight := img.Bounds().Dy()

		cols := width
		if imgWidth < cols {
			cols = imgWidth
		}

		ratio := float64(imgWidth) / float64(cols)
		rows := int(float64(imgHeight) / ratio)

		if terminal.IsTerminal(os.Stdout) {
			for i := 1; i < rows; i += 2 {
				fmt.Println("")
			}
			cursorUp(rows / 2)
		}

		for i := 1; i < rows; i += 2 {
			var bb *Block
			for j := 0; j < cols; j++ {
				// TODO: get average color of the area instead
				// of one pixel?
				x := int(ratio * float64(j))
				yTop := int(ratio * float64(i-1))
				yBottom := int(ratio * float64(i))

				top := ansirgb.Convert(img.At(x, yTop))
				bottom := ansirgb.Convert(img.At(x, yBottom))

				b := &Block{}
				b.bottom = bottom

				// Foreground colors are lighter in some terminals.
				// Ignore top (FG) if it's the same color anyway
				if top.Code != bottom.Code {
					b.top = top
				}

				if bb != nil && b.bottom.Code == bb.bottom.Code &&
					((b.top == nil && bb.top == nil) || b.top != nil && bb.top != nil && b.top.Code == bb.top.Code) {
					b.nocolor = true
				}

				bb = b

				fmt.Printf("%s", b)
			}
			reset()
			fmt.Printf("\n")
		}
		fmt.Println("File:", filepath.Base(fpath), "size:", imgWidth, "x", imgHeight)
	}
}
