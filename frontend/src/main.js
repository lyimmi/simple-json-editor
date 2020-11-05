import "core-js/stable";
const runtime = require("@wailsapp/runtime");
const JSONEditor = require("jsoneditor");
import "promise-polyfill/src/polyfill";

// Main entry point
function start() {
  var mystore = runtime.Store.New("Json");

  // Inject html
  var app = document.getElementById("app");
  app.innerHTML = `
    <div id="jsoneditor" style="width: 100%; height: 100%;"></div>
	`;
  // create the editor
  var changed = "";
  const container = document.getElementById("jsoneditor");
  const options = {
    mode: "code",
    modes: ["text", "code", "view", "tree", "preview"],
    onChangeText: function (jsonString) {
      changed = jsonString;
      mystore.set(jsonString);
    },
    onModeChange: function() {
      addSaveButton();
    }
  };
  const editor = new JSONEditor(container, options);
  addSaveButton();

  var jsonStore = {};
  mystore.subscribe(function (state) {
    try {
      if (changed !== state) {
        jsonStore = JSON.parse(state);
        editor.set(jsonStore);
      }
    } catch (e) {
      console.warn("invalid json");
    }
  });
  window.backend.JSONFile.Start();

  //Save event
  document.onkeydown = function (e) {
    if (e.key == "s" && e.ctrlKey == true) {
      e.preventDefault();
      window.backend.JSONFile.Save();
      return false;
    }
  };
}

// We provide our entrypoint as a callback for runtime.Init
runtime.Init(start);

function addSaveButton() {
  var saveButton = document.getElementById("saveDocument");
  if (!saveButton) {
    var menu = document.getElementsByClassName("jsoneditor-menu");
    saveButton = document.createElement("button");
    saveButton.setAttribute("class", "jsoneditor-save");
    saveButton.setAttribute("id", "saveDocument");
    saveButton.innerHTML = "Save";
    //<a href="#" class="jsoneditor-save" id="saveDocument">Save</a>
    menu[0].appendChild(saveButton);

    document.addEventListener("click", function (event) {
      if (event.target.matches("#saveDocument")) {
        event.preventDefault();
        window.backend.JSONFile.Save();
        return false;
      }
      return;
    });
  }
}
