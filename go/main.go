//go:build example
// +build example

package main

// import (
// 	"syscall/js"
// )

// var (
// 	width      float64
// 	height     float64
// 	isPainting bool
// 	x          float64
// 	y          float64
// 	ctx        js.Value
// 	color      string
// )

// func main() {
// 	doc := js.Global().Get("document")
// 	canvasEl := doc.Call("getElementById", "canvas")
// 	bodyW := doc.Get("body").Get("clientWidth").Float()
// 	bodyH := doc.Get("body").Get("clientHeight").Float()

// 	canvasEl.Set("width", bodyW)
// 	canvasEl.Set("height", bodyH)
// 	ctx = canvasEl.Call("getContext", "2d")

// 	startPaint := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
// 		e := args[0]
// 		isPainting = true

// 		x = e.Get("pageX").Float() - canvasEl.Get("offsetLeft").Float()
// 		y = e.Get("pageY").Float() - canvasEl.Get("offsetTop").Float()

// 		print(x)
// 		print(y)

// 		return nil
// 	})

// 	paint := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
// 		if isPainting {
// 			e := args[0]
// 			nx := e.Get("pageX").Float() - canvasEl.Get("offsetLeft").Float()
// 			ny := e.Get("pageY").Float() - canvasEl.Get("offsetTop").Float()

// 			ctx.Set("strokeStyle", color)
// 			ctx.Set("lineJoin", "round")
// 			ctx.Set("lineWidth", 5)

// 			ctx.Call("beginPath")
// 			ctx.Call("moveTo", nx, ny)
// 			ctx.Call("lineTo", x, y)
// 			ctx.Call("closePath")

// 			// actually draw the path*
// 			ctx.Call("stroke")

// 			// Set x and y to our new coordinates*
// 			x = nx
// 			y = ny
// 		}
// 		return nil
// 	})

// 	exit := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
// 		isPainting = false
// 		return nil
// 	})

// 	canvasEl.Call("addEventListener", "mousedown", startPaint)
// 	canvasEl.Call("addEventListener", "mousemove", paint)
// 	canvasEl.Call("addEventListener", "mouseup", exit)

// 	divEl := doc.Call("getElementById", "colors")

// 	colors := [6]string{"#F4908E", "#F2F097", "#88B0DC", "#F7B5D1", "#53C4AF", "#FDE38C"}

// 	for _, x := range colors {
// 		node := doc.Call("createElement", "div")
// 		node.Call("setAttribute", "class", "color")
// 		node.Call("setAttribute", "id", x)
// 		node.Call("setAttribute", "style", "background-color: "+x)

// 		setColor := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
// 			e := args[0]
// 			color = e.Get("target").Get("id").String()
// 			return nil
// 		})

// 		node.Call("addEventListener", "click", setColor)

// 		divEl.Call("appendChild", node)

// 	}

// }

// Copyright 2014 Hajime Hoshi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"syscall/js"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// const (
// 	screenWidth  = 640
// 	screenHeight = 480
// )
var (
	screenWidth  int
	screenHeight int
)

var (
	brushImage *ebiten.Image
)

func init() {
	const (
		a0 = 0x40
		a1 = 0xc0
		a2 = 0xff
	)
	pixels := []uint8{
		a0, a1, a1, a0,
		a1, a2, a2, a1,
		a1, a2, a2, a1,
		a0, a1, a1, a0,
	}
	brushImage = ebiten.NewImageFromImage(&image.Alpha{
		Pix:    pixels,
		Stride: 4,
		Rect:   image.Rect(0, 0, 4, 4),
	})
}

type Game struct {
	touches []ebiten.TouchID
	count   int

	canvasImage *ebiten.Image
}

func NewGame() *Game {
	g := &Game{
		canvasImage: ebiten.NewImage(screenWidth, screenHeight),
	}
	g.canvasImage.Fill(color.White)
	return g
}

func (g *Game) Update() error {
	drawn := false

	// Paint the brush by mouse dragging
	mx, my := ebiten.CursorPosition()
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.paint(g.canvasImage, mx, my)
		drawn = true
	}

	// Paint the brush by touches
	g.touches = ebiten.AppendTouchIDs(g.touches[:0])
	for _, t := range g.touches {
		x, y := ebiten.TouchPosition(t)
		g.paint(g.canvasImage, x, y)
		drawn = true
	}
	if drawn {
		g.count++
	}
	return nil
}

// paint draws the brush on the given canvas image at the position (x, y).
func (g *Game) paint(canvas *ebiten.Image, x, y int) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	// Scale the color and rotate the hue so that colors vary on each frame.
	op.ColorM.Scale(1.0, 0.50, 0.125, 1.0)
	tps := ebiten.MaxTPS()
	theta := 2.0 * math.Pi * float64(g.count%tps) / float64(tps)
	op.ColorM.RotateHue(theta)
	canvas.DrawImage(brushImage, op)
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.canvasImage, nil)

	mx, my := ebiten.CursorPosition()
	msg := fmt.Sprintf("(%d, %d)", mx, my)
	for _, t := range g.touches {
		x, y := ebiten.TouchPosition(t)
		msg += fmt.Sprintf("\n(%d, %d) touch %d", x, y, t)
	}
	ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	bgColor := js.Global().Get("bg_color").String()                  //	String 	Optional. Background color in the #RRGGBB format.	Also available as the CSS variable var(--tg-theme-bg-color).
	textColor := js.Global().Get("text_color").String()              //	String 	Optional. Main text color in the #RRGGBB format.	Also available as the CSS variable var(--tg-theme-text-color).
	hintColor := js.Global().Get("hint_color").String()              //	String 	Optional. Hint text color in the #RRGGBB format.	Also available as the CSS variable var(--tg-theme-hint-color).
	linkColor := js.Global().Get("link_color").String()              //	String 	Optional. Link color in the #RRGGBB format.	Also available as the CSS variable var(--tg-theme-link-color).
	buttonColor := js.Global().Get("button_color").String()          //	String 	Optional. Button color in the #RRGGBB format.	Also available as the CSS variable var(--tg-theme-button-color).
	buttonTextColor := js.Global().Get("button_text_color").String() //	String 	Optional. Button text color in the #RRGGBB format.	Also available as the CSS variable var(--tg-theme-button-text-color).

	fmt.Println("bgColor", bgColor)
	fmt.Println("textColor", textColor)
	fmt.Println("hintColor", hintColor)
	fmt.Println("linkColor", linkColor)
	fmt.Println("buttonColor", buttonColor)
	fmt.Println("buttonTextColor", buttonTextColor)

	windowTelegramWebApp := js.Global().Get("window.Telegram.WebApp") //	String 	Optional. Button text color in the #RRGGBB format.	Also available as the CSS variable var(--tg-theme-button-text-color).
	fmt.Println("windowTelegramWebApp", windowTelegramWebApp.String())

	screenHeight = int(windowTelegramWebApp.Get("viewportHeight").Float())
	screenWidth = 1080

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Paint (Ebiten Demo)")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
