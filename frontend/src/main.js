import "jsoneditor/dist/jsoneditor.min.css"
import "./main.css"
import JSONEditor from "jsoneditor"

import "./darktheme.css"

import { EventsEmit, EventsOn } from '../wailsjs/runtime';
import { GetCurrentFile, New, Alert, GetDarkMode, GetLocale } from "../wailsjs/go/main/App";
import langHu from "./languages/hu";

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


document.addEventListener("dragover", (e) => {
  e.preventDefault();
});

window.addEventListener(
  "drop",
  (e) => {
    e = e || event;
    e.preventDefault();
    if (e.dataTransfer.items) {
      // Use DataTransferItemList interface to access the file(s)
      for (let i = 0; i < e.dataTransfer.items.length; i++) {
        // If dropped items aren't files, reject them
        if (e.dataTransfer.items[i].kind === "file") {
          try {
            let file = e.dataTransfer.items[i].getAsFile();
            if (file.type !== "application/json" && file.type !== "application/text") {
              alert("only json file is allowed");
              return;
            }
            if (file.size > 10000000) {
              alert("file is greater than 10Mb");
              return;
            }
            handleDroppedFile(file);
            break;
          } catch (error) {
            console.error(error)
          }
        }
      }
    } else {
      // todo: finish
      // Use DataTransfer interface to access the file(s)
      for (let i = 0; i < e.dataTransfer.files.length; i++) {
        try {
          let file = e.dataTransfer.files[i]
          if (file.type !== "application/json" && file.type !== "application/text") {
            alert("only json file is allowed");
            return;
          }
          if (file.size > 10000000) {
            alert("file is greater than 10Mb");
            return;
          }
          handleDroppedFile(file);
          break;
        } catch (error) {
          console.error(error)
        }
      }
    }
  },
  false
);

function handleDroppedFile(file) {
  console.log(file.webkitRelativePath);
  let reader = file.stream().getReader(),
    result = new Uint8Array();

  reader.read().then(function processText({ done, value }) {
    if (done) {

      New(btoa(new TextDecoder().decode(result))).then(data => {
        console.log(data)
      });
      return;
    }
    if (typeof value !== "undefined") {
      let tmp = new Uint8Array(result.length + value.length);
      tmp.set(result);
      tmp.set(value, result.length);
      result = tmp;
    }
    // Read some more, and call this function again
    return reader.read().then(processText);
  });
}