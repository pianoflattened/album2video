const { ipcRenderer } = require("electron");

const browseAlbumBtn = document.getElementById('browse-album');
browseAlbumBtn.addEventListener('click', function() {
	ipcRenderer.send("browse-album");
});
ipcRenderer.on("browse-album-successful", function(event, filePath) {
	document.getElementById("album-dir").value = filePath;
});

const browseCoverBtn = document.getElementById('browse-cover');
const coverPathInput = document.getElementById('cover-path');
browseCoverBtn.addEventListener('click', function() {
	ipcRenderer.send("browse-cover");
});
ipcRenderer.on("browse-cover-successful", function(event, filePath) {
	coverPathInput.value = filePath;
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
separateVideosCheckbox.addEventListener('change', function() {
	if (this.checked) {
		console.log("checked");
		outputPathLabel.textContent = "output directory";
	} else {
		console.log("unchecked");
		outputPathLabel.textContent = "output path";
	}
});


const submitBtn = document.getElementById('submit');
submitBtn.addEventListener('click', function() {
	let albumDir = document.getElementById("album-dir").value;
	let coverPath = document.getElementById("cover-path").value;
	ipcRenderer.send("make-video", );
});