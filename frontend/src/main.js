
window.alert = (msg, title) => {
  msg = typeof msg !== "string" ? "" : msg
  title = typeof title !== "string" ? "" : title
  window.go.main.App.Alert(msg, title);
}

// create the editor
const container = document.getElementById("jsoneditor");

let saved = true;

const options = {
  languages: {
    hu: {
      array: 'Tömb',
      auto: 'Auto',
      appendText: 'Hozzáad',
      appendTitle: 'Adjon hozzá egy új \'auto\' típusú mezőt ezután a mező után (Ctrl+Shift+Ins)',
      appendSubmenuTitle: 'Válassza ki a hozzáadandó mező típusát',
      appendTitleAuto: 'Adjon hozzá egy új \'auto\' típusú mezőt (Ctrl+Shift+Ins)',
      ascending: 'Növekvő',
      ascendingTitle: 'Rendezze a(z) ${type} típus gyerekeit növekvő sorrendbe',
      actionsMenu: 'Kattintson azt actions menü megnyitásához (Ctrl+M)',
      cannotParseFieldError: 'Nem lehet a mezőt JSON-né alakítani',
      cannotParseValueError: 'Nem lehet az értéket JSON-né alakítani',
      collapseAll: 'Az összes mező összecsukása',
      compactTitle: 'JSON adatok sűrítése, fölösleges karakterek eltávolítása (Ctrl+Shift+I)',
      descending: 'Csökkenő',
      descendingTitle: 'Rendezze a(z) ${type} típus gyerekeit csökkenő sorrendbe',
      drag: 'Húzza a mezőt a mozgatáshoz (Alt+Shift+Arrows)',
      duplicateKey: 'duplikált kulcs',
      duplicateText: 'Duplikáció',
      duplicateTitle: 'Kiválasztott mezők duplikálása (Ctrl+D)',
      duplicateField: 'Duplikálja ezt a mezőt (Ctrl+D)',
      duplicateFieldError: 'Mező nevének duplikálása',
      empty: 'üres',
      expandAll: 'Összes mező kinyitása',
      expandTitle: 'Kattintson a mező ki/be csukásához (Ctrl+E). \n' +
        'Ctrl+Kattintás az összes hozzátartozó elem ki/be csukásához.',
      formatTitle: 'JSON adat formázása, megfelelő bekezdésekkel és soremelésekkel (Ctrl+I)',
      insert: 'Beilleszt',
      insertTitle: 'Illesszen be egy új \'auto\' típusú mezőt ez elé a mező elé (Ctrl+Ins)',
      insertSub: 'Válassza ki a beillesztendő mező típusát',
      object: 'Object',
      ok: 'Ok',
      redo: 'Vissza (Ctrl+Shift+Z)',
      removeText: 'Eltávolítás',
      removeTitle: 'Kiválasztott mezők eltávolítása (Ctrl+Del)',
      removeField: 'Távolítsa el ezt a mezőt (Ctrl+Del)',
      repairTitle: 'JSON javítása: idézőjelek és escape karakterek javítása, megjegyzések és JSONP jelölések eltávolítása, JavaScript objektumok JSON-né alakítása.',
      searchTitle: 'Mezők és értékek keresése',
      searchNextResultTitle: 'Következő eredmény (Enter)',
      searchPreviousResultTitle: 'Előző eredmény (Shift + Enter)',
      selectNode: 'Válasszon egy csomópontot...',
      showAll: 'mind mutatása',
      showMore: 'több mutatása',
      showMoreStatus: 'megjelenítve ${visibleChilds} / ${totalChilds}.',
      sort: 'Rendezés',
      sortTitle: 'A(z) ${type} gyerekeinek rendezése',
      sortTitleShort: 'Tartalom rendezése',
      sortFieldLabel: 'Mező:',
      sortDirectionLabel: 'Irány:',
      sortFieldTitle: 'Válassza ki azt a beágyazott mezőt, amely alapján a tömböt vagy objektumot rendezni kívánja.',
      sortAscending: 'Növekvő',
      sortAscendingTitle: 'A kiválasztott mező növekvő sorrendbe rendezése',
      sortDescending: 'Csökkenő',
      sortDescendingTitle: 'A kiválasztott mező rendezése csökkenő sorrendben',
      string: 'String',
      transform: 'Átalakítás',
      transformTitle: 'Szűrje, rendezze vagy alakítsa át ennek a(z) ${type} a gyermekeit.',
      transformTitleShort: 'Szűrjön, rendezzen, vagy alakítsa át a tartalmat',
      extract: 'Kivonat',
      extractTitle: 'Vegye ki a következő típust ${type}',
      transformQueryTitle: 'Adjon meg egy JMESPath-lekérdezést',
      transformWizardLabel: 'Varázsló',
      transformWizardFilter: 'Szűrő',
      transformWizardSortBy: 'Rendezés',
      transformWizardSelectFields: 'Mezők kiválasztása',
      transformQueryLabel: 'Lekérdezés',
      transformPreviewLabel: 'Előnézet',
      type: 'Típus',
      typeTitle: 'A mező típusának módosítása',
      openUrl: 'Ctrl+Click, vagy Ctrl+Enter az url új ablakban történő megnyitásához',
      undo: 'Utolsó művelet visszavonása (Ctrl+Z)',
      validationCannotMove: 'Egy mezőt nem lehet áthelyezni a saját gyermekébe',
      autoType: 'Mező típus "auto". ' +
        'A mező típusa automatikusan meghatározásra kerül az értékből, ' +
        'és lehet karakterlánc, szám, boolean vagy null.',
      objectType: 'Field type "object". ' +
        'Egy objektum kulcs/érték párok rendezetlen halmazát tartalmazza.',
      arrayType: 'Mező típus "array". ' +
        'Egy tömb értékek rendezett gyűjteményét tartalmazza.',
      stringType: 'Mező típus "string". ' +
        'A mező típusa nem az értékből kerül meghatározásra, hanem mindig stringként adódik vissza.',
      modeEditorTitle: 'Szerkesztő mód váltása',
      modeCodeText: 'Kód',
      modeCodeTitle: 'Váltson kódkiemelőre',
      modeFormText: 'Űrlap',
      modeFormTitle: 'Váltás űrlapszerkesztőre',
      modeTextText: 'Szöveg',
      modeTextTitle: 'Váltás egyszerű szövegszerkesztőre',
      modeTreeText: 'Fa',
      modeTreeTitle: 'Faszerkesztőre váltás',
      modeViewText: 'Megtekintés',
      modeViewTitle: 'Váltás fa nézetre',
      modePreviewText: 'Előnézet',
      modePreviewTitle: 'Váltás előnézeti módra',
      examples: 'Minták',
      default: 'Alapértelmezett',
      containsInvalidProperties: 'Érvénytelen tulajdonságokat tartalmaz',
      containsInvalidItems: 'Érvénytelen elemeket tartalmaz'
    },
  },
  mode: 'code',
  modes: ['code', 'form', 'text', 'tree', 'view', 'preview'], // allowed modes
  onModeChange: function (newMode, oldMode) {
    console.log('Mode switched from', oldMode, 'to', newMode)
  },
  onChangeText: function (node, event) {
    if (saved) {
      saved = false;
    }
    try {
      window.runtime.EventsEmit("json-edited", editor.getText());
    } catch (e) {
      console.warn(e)
    }
  }
};

let editor = new JSONEditor(container, options);

// set json
const initialJson = {};
editor.set(initialJson);

getCurrentFile();

window.runtime.EventsOn("json-data", data => {
  try {
    let dJson = JSON.parse(data)
    editor.set(dJson);
  } catch (e) {
    console.warn(e);
    editor.set(data);
  }
});

window.runtime.EventsOn("change-lang", data => {
  options.language = data;
  editor.destroy();
  editor = new JSONEditor(container, options);
  getCurrentFile();
});


function getCurrentFile() {
  window.go.main.App.GetCurrentFile().then(data => {
    try {
      let dJson = JSON.parse(data);
      editor.set(dJson);
    } catch (e) {
      console.warn(e);
    }
  });
}

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

      window.go.main.App.New(btoa(new TextDecoder().decode(result))).then(data => {
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