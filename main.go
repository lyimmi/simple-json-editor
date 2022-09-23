package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
)

//go:embed frontend/dist
var assets embed.FS

func main() {

	// Create an instance of the app structure
	app := NewApp()

	isLocaleSelected := func(l string) bool {
		return app.UserLocale == l
	}

	// Create application with options
	err := wails.Run(&options.App{
		Title:             "Json editor",
		Width:             1200,
		Height:            800,
		MinWidth:          200,
		MinHeight:         200,
		DisableResize:     false,
		Fullscreen:        false,
		Frameless:         false,
		StartHidden:       false,
		HideWindowOnClose: false,
		BackgroundColour:  &options.RGBA{R: 255, G: 255, B: 255, A: 255},
		Assets:            assets,
		Menu: menu.NewMenuFromItems(
			menu.SubMenu("File", menu.NewMenuFromItems(
				menu.Text("New", keys.CmdOrCtrl("n"), func(cd *menu.CallbackData) {
					runtime.EventsEmit(app.ctx, "new-json")
				}),
				menu.Text("Open", keys.CmdOrCtrl("o"), func(cd *menu.CallbackData) {
					runtime.EventsEmit(app.ctx, "open-json")
				}),
				menu.Text("Save", keys.CmdOrCtrl("s"), func(cd *menu.CallbackData) {
					runtime.EventsEmit(app.ctx, "save-json")
				}),
				menu.Separator(),
				menu.Text("Quit", keys.CmdOrCtrl("q"), func(cd *menu.CallbackData) {
					runtime.Quit(app.ctx)
				}),
			)),
			menu.SubMenu("Language", menu.NewMenuFromItems(
				menu.Radio("en", isLocaleSelected("en"), nil, func(cd *menu.CallbackData) {
					runtime.EventsEmit(app.ctx, "change-lang", "en")
				}),
				menu.Radio("hu", isLocaleSelected("hu"), nil, func(cd *menu.CallbackData) {
					runtime.EventsEmit(app.ctx, "change-lang", "hu")
				}),
			)),
			menu.SubMenu("View", menu.NewMenuFromItems(
				menu.Checkbox("Darkmode", app.DarkMode, nil, func(cd *menu.CallbackData) {
					runtime.EventsEmit(app.ctx, "toggle-dark-mode")
				}),
			)),
		),
		OnStartup:        app.startup,
		OnDomReady:       app.domReady,
		OnBeforeClose:    app.beforeClose,
		OnShutdown:       app.shutdown,
		WindowStartState: options.Normal,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}
