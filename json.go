package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails"
)

// JSONFile is the json handler
type JSONFile struct {
	r        *wails.Runtime
	store    *wails.Store
	path     string
	fileName string
}

// WailsInit is called when the component is being initialised
func (c *JSONFile) WailsInit(runtime *wails.Runtime) error {
	c.r = runtime
	c.store = runtime.Store.New("Json", `{}`)
	args := os.Args[1:]
	if len(args) > 0 {
		c.path = args[0]
		c.fileName = filepath.Base(c.path)
		c.r.Window.SetTitle("Simple JSON editor - " + c.fileName)
	}
	return nil
}

//Start json editor
func (c *JSONFile) Start() {
	args := os.Args[1:]
	if c.path != "" {
		dat, err := ioutil.ReadFile(args[0])
		check(err)
		c.store.Set(string(dat))
	}
}

//Save a file
func (c *JSONFile) Save() bool {
	if c.path == "" {
		c.path = c.r.Dialog.SelectSaveFile("Select a file", "*.json", "*.txt")
		c.r.Window.SetTitle("Simple JSON editor - " + c.fileName)
	}
	s, ok := c.store.Get().(string)
	if ok {
		err := ioutil.WriteFile(c.path, []byte(s), 0644)
		if err != nil {
			return false
		}
		return true
	}
	return false
}

// SetJSON reads a file if exists
func (c *JSONFile) SetJSON(data string) {
	c.store.Set(data)
}
