package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed frontend/src
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:             "Json editor",
		Width:             1200,
		Height:            800,
		MinWidth:          1024,
		MinHeight:         768,
		MaxWidth:          2000,
		MaxHeight:         2000,
		DisableResize:     false,
		Fullscreen:        false,
		Frameless:         false,
		StartHidden:       false,
		HideWindowOnClose: false,
		RGBA:              &options.RGBA{R: 255, G: 255, B: 255, A: 255},
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
				menu.Radio("en", true, nil, func(cd *menu.CallbackData) {
					runtime.EventsEmit(app.ctx, "change-lang", "en")
				}),
				menu.Radio("hu", false, nil, func(cd *menu.CallbackData) {
					runtime.EventsEmit(app.ctx, "change-lang", "hu")
				}),
			)),
		),
		Logger:           nil,
		LogLevel:         logger.DEBUG,
		OnStartup:        app.startup,
		OnDomReady:       app.domReady,
		OnBeforeClose:    app.beforeClose,
		OnShutdown:       app.shutdown,
		WindowStartState: options.Normal,
		Bind: []interface{}{
			app,
		},
		// Windows platform specific options
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
			// DisableFramelessWindowDecorations: false,
			WebviewUserDataPath: "",
		},
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: true,
				HideTitle:                  false,
				HideTitleBar:               false,
				FullSizeContent:            false,
				UseToolbar:                 false,
				HideToolbarSeparator:       true,
			},
			Appearance:           mac.NSAppearanceNameDarkAqua,
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			About: &mac.AboutInfo{
				Title:   "Json editor",
				Message: "",
				Icon:    icon,
			},
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}
