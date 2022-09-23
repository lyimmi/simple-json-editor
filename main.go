package main

import (
	"changeme/internal/lang"
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
			menu.SubMenu(lang.Text(app.UserLocale, "file"), menu.NewMenuFromItems(
				menu.Text(lang.Text(app.UserLocale, "file.new"), keys.CmdOrCtrl("n"), func(cd *menu.CallbackData) {
					runtime.EventsEmit(app.ctx, "new-json")
				}),
				menu.Text(lang.Text(app.UserLocale, "file.open"), keys.CmdOrCtrl("o"), func(cd *menu.CallbackData) {
					runtime.EventsEmit(app.ctx, "open-json")
				}),
				menu.Text(lang.Text(app.UserLocale, "file.save"), keys.CmdOrCtrl("s"), func(cd *menu.CallbackData) {
					runtime.EventsEmit(app.ctx, "save-json")
				}),
				menu.Separator(),
				menu.Text(lang.Text(app.UserLocale, "file.quit"), keys.CmdOrCtrl("q"), func(cd *menu.CallbackData) {
					runtime.Quit(app.ctx)
				}),
			)),
			menu.SubMenu(lang.Text(app.UserLocale, "language"), menu.NewMenuFromItems(
				menu.Radio(lang.Text(app.UserLocale, "language.en"), isLocaleSelected("en"), nil, func(cd *menu.CallbackData) {
					runtime.EventsEmit(app.ctx, "change-lang", "en")
				}),
				menu.Radio(lang.Text(app.UserLocale, "language.hu"), isLocaleSelected("hu"), nil, func(cd *menu.CallbackData) {
					runtime.EventsEmit(app.ctx, "change-lang", "hu")
				}),
			)),
			menu.SubMenu(lang.Text(app.UserLocale, "view"), menu.NewMenuFromItems(
				menu.Checkbox(lang.Text(app.UserLocale, "view.darkmode"), app.DarkMode, nil, func(cd *menu.CallbackData) {
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
