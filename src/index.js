const { ipcRenderer } = require("electron");

const browseAlbumBtn = document.getElementById('browse-album');
const albumDirInput = document.getElementById("album-dir");
browseAlbumBtn.addEventListener('click', function() {
	ipcRenderer.send("browse-album");
});
ipcRenderer.on("browse-album-successful", function(event, filePath) {
	albumDirInput.value = filePath;
});

const browseCoverBtn = document.getElementById('browse-cover');
const coverPathInput = document.getElementById('cover-path');
browseCoverBtn.addEventListener('click', function() {
	ipcRenderer.send("browse-cover");
});
ipcRenderer.on("browse-cover-successful", function(event, filePath) {
	coverPathInput.value = filePath;
});

const browseOutputBtn = document.getElementById('browse-output');
browseOutputBtn.addEventListener('click', function() {
	if (this.getAttribute("browsetype") == "directory") {
		ipcRenderer.send("browse-output-directory");
	} else if (this.getAttribute("browsetype" == "path") {
		ipcRenderer.send("browse-output-path");
	}
});
ipcRenderer.on("browse-output-successful", function(event, filePath) {
	document.getElementById("output-path").value = filePath;
});

const detectCoverCheckbox = document.getElementById('detect-cover');
detectCoverCheckbox.addEventListener('change', function() {
	if (this.checked) {
		console.log("checked");
		coverPathInput.setAttribute("class", "grayout");
		coverPathInput.setAttribute("readonly", "");
	} else {
		console.log("unchecked");
		coverPathInput.removeAttribute("class");
		coverPathInput.removeAttribute("readonly");
	}
});

const separateVideosCheckbox = document.getElementById('separate-videos');
const outputPathLabel = document.querySelector('label[for="output-path"]');
const outputPath = document.getElementById("output-path");
let outpath = "";
let outdir = "";
separateVideosCheckbox.addEventListener('change', function() {
	if (this.checked) {
		console.log("checked");
		outpath = outputPath.value;
		outputPathLabel.textContent = "output directory";
		browseOutputBtn.setAttribute("browsetype", "directory");
		outputPath.value = outdir;
	} else {
		console.log("unchecked");
		outdir = outputPath.value;
		outputPathLabel.textContent = "output path";
		browseOutputBtn.setAttribute("browsetype", "path");
		outputPath.value = outpath;
	}
});

const submitBtn = document.getElementById('submit');
submitBtn.addEventListener('click', function() {
	let a2vCallObj = {};
	a2vCallObj.albumDir = albumDirInput.value;
	a2vCallObj.coverPath = coverPathInput.value;
	a2vCallObj.detectCover = detectCoverCheckbox.checked;
	a2vCallObj.separateVideos = separateVideosCheckbox.checked;
	a2vCallObj.outputPath = outputPath.value;
	
	ipcRenderer.send("make-video", a2vCallObj);
});