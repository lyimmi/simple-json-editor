package main

import "github.com/wailsapp/wails/v2/pkg/runtime"

func initListeners(a *App) {

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
	runtime.EventsOn(a.ctx, "change-lang", func(optionalData ...interface{}) {
		loc := "en"
		if len(optionalData) > 0 {
			loc = optionalData[0].(string)
		}
		a.UserLocale = loc
		a.SaveSettings()
	})

	runtime.EventsOn(a.ctx, "toggle-dark-mode", func(optionalData ...interface{}) {
		a.DarkMode = !a.DarkMode
		a.SaveSettings()
	})

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
	runtime.EventsOn(a.ctx, "new-json", func(optionalData ...interface{}) {
		a.New(nil)
	})
	runtime.EventsOn(a.ctx, "open-json", func(optionalData ...interface{}) {
		a.Open()
	})
	runtime.EventsOn(a.ctx, "save-json", func(optionalData ...interface{}) {
		a.Save()
	})
}
