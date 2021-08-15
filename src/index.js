const { app, ipcRenderer } = require("electron");
const path = require('path');
const ProgressBar = require('./progress_bar.js');
var allowRunningTimestamps = false
var timestamps

var progressBar;
const form = {
	albumDirectory: document.getElementById("album-dir"),
	coverPath: document.getElementById("cover-path"),
	extractCover: document.getElementById("extract-cover"),
	separateVideos: document.getElementById("separate-videos"),
	outputPath: document.getElementById("output-path")
}
const submitBtn = document.getElementById('submit');
function updateSubmitBtn(f) {
	submitBtn.disabled = !(f.albumDirectory.value && (f.coverPath.value || f.extractCover.checked) && f.outputPath.value);
}

// file browse events
const browseAlbumBtn = document.getElementById('browse-album');
browseAlbumBtn.addEventListener('click', function() {
	ipcRenderer.send("browse-album");
});
ipcRenderer.on("browse-album-successful", function(event, filePath) {
	form.albumDirectory.value = filePath;
	updateSubmitBtn(form);
});

const browseCoverBtn = document.getElementById('browse-cover');
browseCoverBtn.addEventListener('click', function() {
	ipcRenderer.send("browse-cover");
});
ipcRenderer.on("browse-cover-successful", function(event, filePath) {
	form.coverPath.value = filePath;
	updateSubmitBtn(form);
});

const browseOutputBtn = document.getElementById('browse-output');
browseOutputBtn.addEventListener('click', function() {
	if (this.getAttribute("browsetype") == "directory") {
		ipcRenderer.send("browse-output-directory");
	}
	if (this.getAttribute("browsetype") == "path") {
		ipcRenderer.send("browse-output-path");
	}
});
ipcRenderer.on("browse-output-successful", function(event, filePath) {
	document.getElementById("output-path").value = filePath;
	updateSubmitBtn(form);
});

document.getElementById("collapse-formatting").addEventListener('click', function() {
	let t = document.getElementById("timestamps");
	if (t.classList.contains("btn-collapsed")) {
		t.classList.toggle("btn-collapsed");
	} else {
		setTimeout(_ => t.classList.toggle("btn-collapsed"), 300);
	}
});

let trackFormatting = document.getElementById("track-formatting");
trackFormatting.addEventListener('input', function() {
	if (allowRunningTimestamps) ipcRenderer.send("timestamp-format", {
		format: trackFormatting.value,
		timestamps: timestamps
	});
});

form.extractCover.addEventListener('change', function() {
	if (this.checked) {
		console.log("checked");
		document.getElementById("browse-cover").setAttribute("disabled", "");
	} else {
		console.log("unchecked");
		document.getElementById("browse-cover").removeAttribute("disabled");
	}
});

let browsetype;
let outdir ="";
let outpath = "";
form.separateVideos.addEventListener('change', function() {
	if (this.checked) {
		outpath = form.outputPath.value;
		form.outputPath.placeholder = "output directory";
		browseOutputBtn.setAttribute("browsetype", "directory");
		browsetype = "directory";
		form.outputPath.value = outdir;
	} else {
		outdir = form.outputPath.value;
		form.outputPath.placeholder = "output path";
		browseOutputBtn.setAttribute("browsetype", "path");
		browsetype = "path";
		form.outputPath.value = outpath;
	}
	updateSubmitBtn(form); // fixes a bug where checking the separate videos box breaks the submit
						   // button because the event fires off in an order i cant control
});

for (const key in form) {
	form[key].addEventListener('input', function() {
		console.log(form.outputPath.value);
		updateSubmitBtn(form);
	});
}

submitBtn.addEventListener('click', function() {
    submitBtn.setAttribute("disabled","");
	let ack = form.albumDirectory.value;
	if (!ack.endsWith(path.sep)) {
		ack += path.sep
	}

	let formData = [
        ack,
        form.coverPath.value,
        form.extractCover.checked,
        form.separateVideos.checked,
        form.outputPath.value
    ];

	progressBar = new ProgressBar(document.querySelector(".progress-container"));
    progressBar.makeIndeterminate();
    progressBar.setLabel('starting subprocess..');
    ipcRenderer.send("make-video", JSON.stringify(formData));
});

const collapseFormatting = document.getElementById('collapse-formatting');
collapseFormatting.addEventListener('click', function() {
	ipcRenderer.send('resize-window', {
		btnclass: collapseFormatting.className,
		offsetHeight: process.platform == "win32" ? 28 : 27 // LOL
	});
});

ipcRenderer.on("progress-label", function(event, msg) {
    progressBar.setLabel(msg);
});

ipcRenderer.on("make-determinate", function(event) {
	progressBar.makeDeterminate();
});

ipcRenderer.on("set-progress", function(event, progress) {
	progressBar.setProgress(progress);
});

ipcRenderer.on("get-timestamp-format", function(event, _timestamps) {
	timestamps = _timestamps
	ipcRenderer.send("timestamp-format", {
		format: trackFormatting.value,
		timestamps: _timestamps
	});
});

ipcRenderer.on("set-complete", function(event) {
	progressBar.setComplete().then(_ => submitBtn.removeAttribute("disabled"));
});

ipcRenderer.on("set-error", function(event, code) {
	progressBar.error(code).then(_ => submitBtn.removeAttribute("disabled"));
});

ipcRenderer.on("timestamps", function(event, timestamps) {
	document.getElementById("timestamps").value = timestamps;
	document.getElementById("timestamps").removeAttribute("readonly");
	allowRunningTimestamps = true;
});
