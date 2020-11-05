import "core-js/stable";
const runtime = require("@wailsapp/runtime");
const JSONEditor = require("jsoneditor");
import "promise-polyfill/src/polyfill";

// Main entry point
function start() {
  if (!Element.prototype.matches) {
    Element.prototype.matches =
      Element.prototype.matchesSelector ||
      Element.prototype.mozMatchesSelector ||
      Element.prototype.msMatchesSelector ||
      Element.prototype.oMatchesSelector ||
      Element.prototype.webkitMatchesSelector ||
      function (s) {
        var matches = (this.document || this.ownerDocument).querySelectorAll(s),
          i = matches.length;
        while (--i >= 0 && matches.item(i) !== this) {}
        return i > -1;
      };
  }
  
  var mystore = runtime.Store.New("Json");
  var JSONIsSaved = true;
  var prevState = "";
  // Inject html
  var app = document.getElementById("app");
  app.innerHTML = `
    <div id="jsoneditor" style="width: 100%; height: 100%;"></div>
	`;
  // create the editor
  const container = document.getElementById("jsoneditor");
  const options = {
    mode: "code",
    modes: ["text", "code", "view", "tree", "preview"],
    onChangeText: function (jsonString) {
      prevState = jsonString;
      mystore.set(jsonString);
      JSONIsSaved = false;
    },
    onModeChange: function () {
      document.removeEventListener("click", newButtonListener);
      document.removeEventListener("click", openButtonListener);
      document.removeEventListener("click", saveButtonListener);
      addNewButton();
      addOpenButton();
      addSaveButton();
    },
  };
  const editor = new JSONEditor(container, options);
  addNewButton();
  addOpenButton();
  addSaveButton();

  mystore.subscribe(function (state) {
    try {
      console.log(state, prevState !== state);
      if (prevState !== state) {
        prevState = state;
        editor.set(JSON.parse(state));
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
      JSONIsSaved = true;
      return false;
    }
  };

  function addNewButton() {
    var openButton = document.getElementById("newDocument");
    if (!openButton) {
      var menu = document.getElementsByClassName("jsoneditor-menu");
      openButton = document.createElement("button");
      openButton.setAttribute("class", "jsoneditor-new");
      openButton.setAttribute("id", "newDocument");
      openButton.innerHTML = "New";
      menu[0].appendChild(openButton);

      document.addEventListener("click", newButtonListener);
    }
  }

  function addOpenButton() {
    var openButton = document.getElementById("openDocument");
    if (!openButton) {
      var menu = document.getElementsByClassName("jsoneditor-menu");
      openButton = document.createElement("button");
      openButton.setAttribute("class", "jsoneditor-open");
      openButton.setAttribute("id", "openDocument");
      openButton.innerHTML = "Open";
      menu[0].appendChild(openButton);

      document.addEventListener("click", openButtonListener);
    }
  }

  function addSaveButton() {
    var saveButton = document.getElementById("saveDocument");
    if (!saveButton) {
      var menu = document.getElementsByClassName("jsoneditor-menu");
      saveButton = document.createElement("button");
      saveButton.setAttribute("class", "jsoneditor-save");
      saveButton.setAttribute("id", "saveDocument");
      saveButton.innerHTML = "Save";
      menu[0].appendChild(saveButton);

      document.addEventListener("click", saveButtonListener);
    }
  }

  function newButtonListener(event) {
    if (event.target.matches("#newDocument")) {
      event.preventDefault();
      if (!JSONIsSaved) {
        var confirmation = confirm(
          "You have unsaved changes, are you sure to open a new file?"
        );
        if (!confirmation) {
          e.stopPropagation();
        }
      }
      window.backend.JSONFile.New();
      prevState = "";
      JSONIsSaved = true;
      return false;
    }
    return;
  }

  function openButtonListener(event) {
    if (event.target.matches("#openDocument")) {
      event.preventDefault();
      if (!JSONIsSaved) {
        var confirmation = confirm(
          "You have unsaved changes, are you sure to open a new file?"
        );
        if (!confirmation) {
          e.stopPropagation();
        }
      }
      prevState = "";
      window.backend.JSONFile.Open();
      JSONIsSaved = true;
      return true;
    }
    return;
  }

  function saveButtonListener(event) {
    if (event.target.matches("#saveDocument")) {
      event.preventDefault();
      window.backend.JSONFile.Save().then(function (res) {
        JSONIsSaved = res;
        console.log(JSONIsSaved);
      });
      return false;
    }
    return;
  }
}

// We provide our entrypoint as a callback for runtime.Init
runtime.Init(start);
