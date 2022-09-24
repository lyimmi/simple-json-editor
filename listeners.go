package main

import "github.com/wailsapp/wails/v2/pkg/runtime"

func initListeners(a *App) {

	// Create a message dialog on error event.
	runtime.EventsOn(a.ctx, "error", func(optionalData ...interface{}) {
		msg := "An error occurred..."
		if len(optionalData) > 0 {
			msg = optionalData[0].(string)
		}
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Error",
			Message: msg,
		})
	})

	// Change the language on change-lang event and save it.
	runtime.EventsOn(a.ctx, "change-lang", func(optionalData ...interface{}) {
		loc := "en"
		if len(optionalData) > 0 {
			loc = optionalData[0].(string)
		}
		a.UserLocale = loc
		a.SaveSettings()
	})

	// Toggle dark mode and save it.
	runtime.EventsOn(a.ctx, "toggle-dark-mode", func(optionalData ...interface{}) {
		a.DarkMode = !a.DarkMode
		a.SaveSettings()
	})

	// Change the window title and store the new json data.
	runtime.EventsOn(a.ctx, "json-edited", func(optionalData ...interface{}) {
		if a.jsonFileSaved {
			runtime.WindowSetTitle(a.ctx, a.jsonFileName+"*")
			a.jsonFileSaved = false
		}
		if len(optionalData) > 0 {
			d, ok := optionalData[0].(string)
			if ok && len(d) > 0 {
				a.jsonFile = []byte(d)
			}
		}
	})

	// Open a new empty editor.
	runtime.EventsOn(a.ctx, "new-json", func(optionalData ...interface{}) {
		a.New(nil)
	})

	// Open a file dialog and load a file form disk.
	runtime.EventsOn(a.ctx, "open-json", func(optionalData ...interface{}) {
		a.Open()
	})

	// Save the json data to disk.
	runtime.EventsOn(a.ctx, "save-json", func(optionalData ...interface{}) {
		a.Save()
	})
}
