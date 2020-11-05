import 'core-js/stable';
const runtime = require('@wailsapp/runtime');
const JSONEditor = require("jsoneditor");

// Main entry point
function start() {

	var mystore = runtime.Store.New('Json');

	// Inject html
	var app = document.getElementById('app');
	app.innerHTML = `
    <div id="jsoneditor" style="width: 100%; height: 100%;"></div>
	`;
	// create the editor
	const container = document.getElementById("jsoneditor")
	const options = {
		mode: 'text',
		modes: ['text', 'code', 'view'],
		onChangeText: function (jsonString) {
			console.log(jsonString)
			mystore.set(jsonString);
		}
	}
	const editor = new JSONEditor(container, options)

	var jsonStore = {};
	mystore.subscribe(function (state) {
		jsonStore = JSON.parse(state);
		editor.set(jsonStore);
	});
	window.backend.JSONFile.Start()
	// // get json
	// const updatedJson = editor.get()

};

// We provide our entrypoint as a callback for runtime.Init
runtime.Init(start);