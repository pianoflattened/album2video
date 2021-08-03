const { app, ipcRenderer } = require("electron");
const getTags = require('./get_tags.js');
const ProgressBar = require('./progress_bar.js');

var progressBar;
const form = {
	albumDirectory: document.getElementById("album-dir"),
	coverPath: document.getElementById("cover-path"),
	detectCover: document.getElementById("detect-cover"),
	separateVideos: document.getElementById("separate-videos"),
	outputPath: document.getElementById("output-path")
}
const submitBtn = document.getElementById('submit');
function updateSubmitBtn(f) {
	submitBtn.disabled = !(f.albumDirectory.value && (f.coverPath.value || f.detectCover.checked) && f.outputPath.value);
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


form.detectCover.addEventListener('change', function() {
	if (this.checked) {
		console.log("checked");
		form.coverPath.setAttribute("disabled", "");
	} else {
		console.log("unchecked");
		form.coverPath.removeAttribute("disabled");
	}
});

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
						   // button because the event fire off in an order i cant control
});

for (const key in form) {
	form[key].addEventListener('input', function() {
		console.log(form.outputPath.value);
		updateSubmitBtn(form);
	});
}

submitBtn.addEventListener('click', function() {
	let formData = {};
	for (const key in form) {
		if (form[key].type == "checkbox") {
			formData[key] = form[key].checked;
		} else {
			formData[key] = form[key].value;
		}
	}
	
	let progressBar = new ProgressBar(document.querySelector(".progress-container"));
	(async () => {getTags(formData, progressBar)})();
});
