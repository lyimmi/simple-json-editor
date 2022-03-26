package main

// JSONFile is the json handler
type JsonFile struct {
	// r        *wails.Runtime
	// store    *wails.Store
	// path     string
	// fileName string
}

// //New a file
// func (c *JSONFile) New() bool {
// 	c.path = ""
// 	c.fileName = ""
// 	c.r.Window.SetTitle("Simple JSON editor")
// 	c.store.Set("{}")
// 	return true
// }

// //Open a file
// func (c *JSONFile) Open() bool {
// 	f := c.r.Dialog.SelectFile()
// 	if f == "" {
// 		return false
// 	}
// 	dat, err := ioutil.ReadFile(f)
// 	if err != nil {
// 		return false
// 	}
// 	c.path = f
// 	c.fileName = filepath.Base(c.path)
// 	c.r.Window.SetTitle("Simple JSON editor - " + c.fileName)
// 	c.store.Set(string(dat))
// 	return true
// }

// //Save a file
// func (c *JSONFile) Save() bool {
// 	if c.path == "" {
// 		c.path = c.r.Dialog.SelectSaveFile("Select a file", "*.json", "*.txt")
// 		c.r.Window.SetTitle("Simple JSON editor - " + c.fileName)
// 	}
// 	s, ok := c.store.Get().(string)
// 	if ok {
// 		err := ioutil.WriteFile(c.path, []byte(s), 0644)
// 		if err != nil {
// 			return false
// 		}
// 		return true
// 	}
// 	return false
// }

// // SetJSON reads a file if exists
// func (c *JSONFile) SetJSON(data string) {
// 	c.store.Set(data)
// }
