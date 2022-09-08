# README

## About

This is a simple json enditor made with [wails](https://wails.io) and [jsoneditor](https://github.com/josdejong/jsoneditor)

jsoneditor has nice features, so this is a wrapper around it, can be used to open json files as default appication.


## Desktop file & default application

Change the json-editor.desktop file and copy it to `~/.local/share/applications`.


If you would like to use it as default for *.json, add a new line `application/json=json-editor.desktop` to  `~/.local/share/applications/defaults.list`

Sould be able to register it on first start, maybe this could be added later.

## Live Development

To run in live development mode, run `wails dev` in the project directory. The frontend dev server will run
on http://localhost:34115. Open this in your browser to connect to your application.

## Building

For a production build, use `wails build`.

