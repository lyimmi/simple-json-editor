package main

import (
	"context"
	"encoding/json"
	"os"
	"path"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx           context.Context `json:"-"`
	userDir       string          `json:"-"`
	UserLocale    string          `json:"userLocale"`
	jsonFile      []byte          `json:"-"`
	jsonFileName  string          `json:"-"`
	jsonFilePath  string          `json:"-"`
	jsonFileSaved bool            `json:"-"`
	DarkMode      bool            `json:"darkMode"`
}

// NewApp creates a new App application struct
func NewApp() *App {
	a := &App{
		jsonFileSaved: true,
		UserLocale:    "en",
		DarkMode:      false,
	}
	var err error
	a.userDir, err = os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	a.LoadSettings()
	return a
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	var err error
	a.ctx = ctx

	initListeners(a)

	runtime.EventsEmit(a.ctx, "change-lang", a.UserLocale)

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

// alertBeforeQuit opens an alert / confirm dialog before quitting if the editor has unsaved changes.
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

// GetCurrentFile returns the currently opened json file's contents.
func (a *App) GetCurrentFile() string {
	return string(a.jsonFile)
}

// New creates an empty editor.
func (a *App) New(data []byte) bool {
	res, err := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Title:   "Alert, closing file!",
		Message: "Are you sure to close the current file and create a new one?",
		Type:    runtime.QuestionDialog,
		Buttons: []string{
			"Yes",
			"No",
		},
	})

	if err != nil {
		panic(err)
	}

	if res == "Yes" {
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

	return false
}

// Open loads a file from disk.
func (a *App) Open() bool {

	prevent := a.alertBeforeQuit()
	if prevent {
		return false
	}

	var err error
	newPath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
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

	if newPath != "" {
		a.jsonFilePath = newPath
	} else {
		return false
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

// Save stores the currently opened json file on disk.
func (a *App) Save() bool {
	var (
		jPath string
		err   error
	)
	if a.jsonFilePath == "" {
		jPath, err = runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
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

		if jPath == "" {
			return false
		}
		a.jsonFilePath = jPath
	}

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

// Alert shows an alert dialog.
func (a *App) Alert(message string, title string) {
	runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:    runtime.WarningDialog,
		Title:   title,
		Message: message,
	})
}

// LoadSettings loads the app's configs from disk.
func (a *App) LoadSettings() {
	fp := path.Join(a.userDir, ".jsoneditor-settings")
	if _, err := os.Stat(fp); err != nil {
		return
	}

	d, err := os.ReadFile(fp)
	if err != nil {
		panic(err)
	}

	nA := App{
		ctx:           a.ctx,
		userDir:       a.userDir,
		jsonFile:      a.jsonFile,
		jsonFileName:  a.jsonFileName,
		jsonFilePath:  a.jsonFilePath,
		jsonFileSaved: a.jsonFileSaved,
		UserLocale:    a.UserLocale,
		DarkMode:      a.DarkMode,
	}

	err = json.Unmarshal(d, &nA)
	if err != nil {
		panic(err)
	}

	*a = nA
}

// SaveSettings stores the app's settings on disk.
func (a *App) SaveSettings() {
	d, err := json.Marshal(a)
	if err != nil {
		panic(err)
	}
	f, err := os.Create(path.Join(a.userDir, ".jsoneditor-settings"))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.Write(d)
}

// GetLocale returns the app's currently selected locale.
func (a *App) GetLocale() string {
	return a.UserLocale
}

// GetDarkMode returns if the app is in dark mode or not.
func (a *App) GetDarkMode() bool {
	return a.DarkMode
}
