package main

import (
	"github.com/leaanthony/mewn"
	"github.com/wailsapp/wails"
)

func main() {

	js := mewn.String("./frontend/build/main.js")
	css := mewn.String("./frontend/build/main.css")
	jsonEditorCSS := mewn.String("./frontend/build/jsoneditor.min.css")

	app := wails.CreateApp(&wails.AppConfig{
		Width:     1024,
		Height:    768,
		Title:     "Simple JSON editor",
		JS:        js,
		CSS:       jsonEditorCSS + css,
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
