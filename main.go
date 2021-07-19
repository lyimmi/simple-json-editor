package main

import (
	_ "embed"

	"github.com/wailsapp/wails"
)

//go:embed frontend/build/main.js
var js string

//go:embed frontend/build/main.css
var css string

//go:embed frontend/build/jsoneditor.min.css
var jecss string

func main() {

	app := wails.CreateApp(&wails.AppConfig{
		Width:     1024,
		Height:    768,
		Title:     "Simple JSON editor",
		JS:        js,
		CSS:       jecss + css,
		Colour:    "#131313",
		Resizable: true,
	})

	app.Bind(&JSONFile{})
	app.Run()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
