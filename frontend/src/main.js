import "jsoneditor/dist/jsoneditor.min.css"
import "./main.css"
import JSONEditor from "jsoneditor"

import "./darktheme.css"

import {EventsEmit, EventsOn, OnFileDrop} from '../wailsjs/runtime';
import { GetCurrentFile, New, Alert, GetDarkMode, GetLocale } from "../wailsjs/go/main/App";
import langHu from "./languages/hu";

OnFileDrop(() => {}, false)

window.alert = (msg, title) => {
  msg = typeof msg !== "string" ? "" : msg
  title = typeof title !== "string" ? "" : title
  Alert(msg, title);
}

const editorOptions= {
  languages: {
    hu: langHu,
  },
  language: "en",
  mode: 'code',
  modes: ['code', 'form', 'text', 'tree', 'view', 'preview'], // allowed modes
  onModeChange: function (newMode, oldMode) {
    console.log('Mode switched from', oldMode, 'to', newMode);
  },
  onChangeText: function (node, event) {
    if (saved) {
      saved = false;
    }
    try {
      EventsEmit("json-edited", editor.getText());
    } catch (e) {
      console.warn(e);
    }
  }
}

let saved = true;
let editor = null;

GetDarkMode().then(mode => {
  if (mode) {
    document.body.classList.add("dark-mode");
  }
});

// create the editor
const container = document.getElementById("jsoneditor");
GetLocale().then(lang => {
  editorOptions.language = lang
  editor = new JSONEditor(container, editorOptions)
  GetCurrentFile().then(data => {
    try {
      editor.set(JSON.parse(data));
    } catch (e) {
      console.warn(e);
    }
  });
});

EventsOn("json-data", data => {
  try {
    let dJson = JSON.parse(data);
    editor.set(dJson);
  } catch (e) {
    console.warn(e);
    editor.set(data);
  }
});

EventsOn("change-lang", data => {
  editorOptions.language = data;
  editor.destroy();

  editor = new JSONEditor(container, editorOptions);
  GetCurrentFile().then(data => {
    try {
      let dJson = JSON.parse(data);
      editor.set(dJson);
    } catch (e) {
      console.warn(e);
    }
  });
});

EventsOn("toggle-dark-mode", data => {
  if (document.body.classList.contains("dark-mode")) {
    document.body.classList.remove("dark-mode");
  } else {
    document.body.classList.add("dark-mode");
  }
});