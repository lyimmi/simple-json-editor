package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx           context.Context
	userDir       string
	UserLocale    string `json:"userLocale"`
	jsonFile      []byte
	jsonFileName  string
	jsonFilePath  string
	jsonFileSaved bool
	DarkMode      bool `json:"darkMode"`
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

	args := os.Args[1:]
	if len(args) > 0 {
		a.jsonFilePath = args[0]
		a.jsonFileName = filepath.Base(a.jsonFilePath)
		a.jsonFile, _ = os.ReadFile(a.jsonFilePath)
	} else {
		a.jsonFileName = "untitled file"
	}

	return a
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	initListeners(a)
	runtime.WindowSetTitle(a.ctx, a.jsonFileName)

	runtime.OnFileDrop(a.ctx, func(x, y int, paths []string) {
		if len(paths) == 0 {
			return
		}

		errs := make([]string, 0)
		for _, p := range paths {
			d, err := os.ReadFile(p)
			if err != nil {
				a.Alert(fmt.Sprintf("failed to open file %s", p), "Error")
			}
			if err = validateJSONFile(p, d); err != nil {
				errs = append(errs, err.Error())
				continue
			}

			a.New(d)
			return
		}

		a.Alert(fmt.Sprintf("No json file found errors: %s", strings.Join(errs, "\n")), "ERROR")
	})
}

func validateJSONFile(jsonFilePath string, fileData []byte) error {
	fi, err := os.Stat(jsonFilePath)
	if err != nil {
		return err
	}
	if fi.IsDir() {
		return errors.New("file is a directory")
	}

	if fi.Size() > 30_000_000 {
		return fmt.Errorf("file size too big: %d", fi.Size())
	}

	if json.Valid(fileData) {
		return nil
	}

	parts := strings.Split(jsonFilePath, ".")
	if parts[len(parts)-1] == "json" && ((fileData[0] == '{' && fileData[len(fileData)-1] == '}') || (fileData[0] == '[' && fileData[len(fileData)-1] == ']')) {
		return nil
	}

	return errors.New("not a json file")
}

// domReady is called after front-end resources have been loaded
func (a *App) domReady(ctx context.Context) {
	runtime.EventsEmit(a.ctx, "change-lang", a.UserLocale)
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
	fileData, err := os.ReadFile(a.jsonFilePath)
	if err != nil {
		runtime.EventsEmit(a.ctx, "error", err.Error())
		return false
	}

	if err = validateJSONFile(a.jsonFilePath, fileData); err != nil {
		runtime.EventsEmit(a.ctx, "error", err.Error())
		return false
	}

	a.jsonFile = fileData
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

	_, err = f.Write(d)
	if err != nil {
		log.Fatal(err)
	}
}

// GetLocale returns the app's currently selected locale.
func (a *App) GetLocale() string {
	return a.UserLocale
}

// GetDarkMode returns if the app is in dark mode or not.
func (a *App) GetDarkMode() bool {
	return a.DarkMode
}
