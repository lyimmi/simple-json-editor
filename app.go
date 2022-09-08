package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx           context.Context
	userDir       string
	userLocale    string
	jsonFile      []byte
	jsonFileName  string
	jsonFilePath  string
	jsonFileSaved bool
}

// NewApp creates a new App application struct
func NewApp(locale string) *App {
	return &App{
		jsonFileSaved: true,
		userLocale:    locale,
	}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	var err error
	a.ctx = ctx
	a.userDir, err = os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	runtime.EventsEmit(a.ctx, "change-lang", a.userLocale)

	args := os.Args[1:]

	if len(args) > 0 {
		a.jsonFilePath = args[0]
		a.jsonFileName = filepath.Base(a.jsonFilePath)
		a.jsonFile, err = os.ReadFile(a.jsonFilePath)
		if err != nil {
			runtime.EventsEmit(a.ctx, "error", err.Error())
		}
	} else {
		a.jsonFileName = "untitled file"
	}
	runtime.WindowSetTitle(a.ctx, a.jsonFileName)

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
	runtime.EventsOn(a.ctx, "json-edited", func(optionalData ...interface{}) {
		runtime.WindowSetTitle(a.ctx, a.jsonFileName+"*")
		a.jsonFileSaved = false
		if len(optionalData) > 0 {
			if d, ok := optionalData[0].([]interface{}); ok {
				if len(d) > 0 {
					a.jsonFile = []byte(d[0].(string))
				}
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

// domReady is called after front-end resources have been loaded
func (a App) domReady(ctx context.Context) {
	// Add your action here
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	return a.alertBeforeQuit()
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
}

func (a *App) alertBeforeQuit() bool {
	if !a.jsonFileSaved {
		res, err := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Alert, unsaved changes!",
			Message: "You have unsaved changes, save before quit?",
			Type:    runtime.QuestionDialog,
			Buttons: []string{
				"Yes",
				"No",
				"Cancel",
			},
		})

		if err != nil {
			panic(err)
		}

		if res == "No" {
			return false
		}
		if res == "Yes" {
			return !a.Save()
		}
		if res == "Cancel" {
			return true
		}
	}

	return false
}

func (a *App) GetCurrentFile() string {
	return string(a.jsonFile)
}

func (a *App) New(data []byte) bool {
	prevent := a.alertBeforeQuit()
	if prevent {
		return false
	}
	if data == nil {
		data = []byte("{}")
	}
	a.jsonFile = data
	a.jsonFileName = ""
	a.jsonFilePath = ""
	a.jsonFileSaved = true

	runtime.WindowSetTitle(a.ctx, a.jsonFileName)
	runtime.EventsEmit(a.ctx, "json-saved", a.jsonFileSaved)
	runtime.EventsEmit(a.ctx, "json-data", string(a.jsonFile))

	return true
}

func (a *App) Open() bool {

	prevent := a.alertBeforeQuit()
	if prevent {
		return false
	}

	var err error
	a.jsonFilePath, err = runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		DefaultDirectory: a.userDir,
		Title:            "select a json file",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "*.json",
				Pattern:     "*.json",
			},
		},
	})
	if err != nil {
		runtime.EventsEmit(a.ctx, "error", err.Error())
	}

	a.jsonFileName = filepath.Base(a.jsonFilePath)
	a.jsonFile, err = os.ReadFile(a.jsonFilePath)
	if err != nil {
		runtime.EventsEmit(a.ctx, "error", err.Error())
	}
	a.jsonFileSaved = true

	runtime.WindowSetTitle(a.ctx, a.jsonFileName)
	runtime.EventsEmit(a.ctx, "json-saved", a.jsonFileSaved)
	runtime.EventsEmit(a.ctx, "json-data", string(a.jsonFile))

	return true
}

func (a *App) Save() bool {
	var err error
	if a.jsonFilePath == "" {
		a.jsonFilePath, err = runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
			DefaultDirectory: a.userDir,
			Title:            "select a json file",
			Filters: []runtime.FileFilter{
				{
					DisplayName: "*.json",
					Pattern:     "*.json",
				},
			},
		})
		if err != nil {
			runtime.EventsEmit(a.ctx, "error", err.Error())
		}
	}

	// log.Println(a.jsonFilePath)

	err = os.WriteFile(a.jsonFilePath, a.jsonFile, 0644)
	if err != nil {
		runtime.EventsEmit(a.ctx, "error", err.Error())
		return false
	}

	a.jsonFileName = filepath.Base(a.jsonFilePath)
	a.jsonFileSaved = true

	runtime.EventsEmit(a.ctx, "json-saved", a.jsonFileSaved)
	runtime.WindowSetTitle(a.ctx, a.jsonFileName)

	return true
}

func (a *App) Alert(message string, title string) {
	runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:    runtime.WarningDialog,
		Title:   title,
		Message: message,
	})
}
