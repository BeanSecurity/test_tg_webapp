package main

import (
	"fmt"
	"log"
	"syscall/js"

	"github.com/BeanSecurity/test_tg_webapp/go/game"
	"github.com/hajimehoshi/ebiten/v2"
)

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

	screen := js.Global().Get("screen")

	if windowTelegramWebApp.IsUndefined() {
		game.ScreenHeight = 1080
	} else {
		game.ScreenHeight = windowTelegramWebApp.Get("viewportHeight").Int()
	}
	game.ScreenWidth = screen.Get("width").Int()

	fmt.Printf("game.ScreenHeight: %+#v\n", game.ScreenHeight) // DEBUG: dump var
	fmt.Printf("game.ScreenWidth: %+#v\n", game.ScreenWidth)   // DEBUG: dump var

	ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)
	ebiten.SetWindowTitle("Paint (Ebiten Demo)")
	if err := ebiten.RunGame(game.NewGame()); err != nil {
		log.Fatal(err)
	}
}
