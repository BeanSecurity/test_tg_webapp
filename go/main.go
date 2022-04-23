package main

import (
	"fmt"
	"log"
	"syscall/js"

	"github.com/BeanSecurity/test_tg_webapp/go/game"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// window.Telegram.WebApp
	window := js.Global().Get("window")
	if window.IsUndefined() {
		return
	}
	webapp := window.Get("Telegram").Get("WebApp")

	bgColor := webapp.Get("bg_color").String()                  //	String 	Optional. Background color in the #RRGGBB format.	Also available as the CSS variable var(--tg-theme-bg-color).
	textColor := webapp.Get("text_color").String()              //	String 	Optional. Main text color in the #RRGGBB format.	Also available as the CSS variable var(--tg-theme-text-color).
	hintColor := webapp.Get("hint_color").String()              //	String 	Optional. Hint text color in the #RRGGBB format.	Also available as the CSS variable var(--tg-theme-hint-color).
	linkColor := webapp.Get("link_color").String()              //	String 	Optional. Link color in the #RRGGBB format.	Also available as the CSS variable var(--tg-theme-link-color).
	buttonColor := webapp.Get("button_color").String()          //	String 	Optional. Button color in the #RRGGBB format.	Also available as the CSS variable var(--tg-theme-button-color).
	buttonTextColor := webapp.Get("button_text_color").String() //	String 	Optional. Button text color in the #RRGGBB format.	Also available as the CSS variable var(--tg-theme-button-text-color).

	fmt.Println("bgColor", bgColor)
	fmt.Println("textColor", textColor)
	fmt.Println("hintColor", hintColor)
	fmt.Println("linkColor", linkColor)
	fmt.Println("buttonColor", buttonColor)
	fmt.Println("buttonTextColor", buttonTextColor)

	screen := js.Global().Get("screen")

	if webapp.IsUndefined() {
		game.ScreenHeight = 1080
	} else {
		game.ScreenHeight = webapp.Get("viewportHeight").Int()
	}
	game.ScreenWidth = 2 * screen.Get("width").Int()

	fmt.Printf("game.ScreenHeight: %+#v\n", game.ScreenHeight) // DEBUG: dump var
	fmt.Printf("game.ScreenWidth: %+#v\n", game.ScreenWidth)   // DEBUG: dump var

	game.ScreenHeight = 800
	game.ScreenWidth = 400

	// ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)
	ebiten.SetWindowSize(400, 800)
	ebiten.SetWindowTitle("Paint (Ebiten Demo)")
	if err := ebiten.RunGame(game.NewGame()); err != nil {
		log.Fatal(err)
	}
}
