package main

import (
	"io/ioutil"
	"os"

	"github.com/wailsapp/wails"
)

// JSONFile is the json handler
type JSONFile struct {
	r     *wails.Runtime
	store *wails.Store
}

// WailsInit is called when the component is being initialised
func (c *JSONFile) WailsInit(runtime *wails.Runtime) error {
	c.r = runtime
	c.store = runtime.Store.New("Json", `{}`)
	// c.store.Subscribe(func(data string) {
	// 	println("New Value:", data)
	// })
	return nil
}

func (c *JSONFile) Start() {
	//uEnc := b64.URLEncoding.EncodeToString(res)
	args := os.Args[1:]
	if len(args) > 0 {
		dat, err := ioutil.ReadFile(args[0])
		check(err)
		c.store.Set(string(dat))
	}
}

// SetJson reads a file if exists
func (c *JSONFile) SetJson(data string) {
	c.store.Set(data)
}
