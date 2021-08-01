const { ipcRenderer } = require("electron");

const browseAlbumBtn = document.getElementById('browse-album');
browseAlbumBtn.addEventListener('click', function() {
	ipcRenderer.send("browse-album");
});
ipcRenderer.on("browse-album-successful", function(event, filePath) {
	document.getElementById("album-dir").value = filePath;
});

const browseCoverBtn = document.getElementById('browse-cover');
browseCoverBtn.addEventListener('click', function() {
	ipcRenderer.send("browse-cover");
});
ipcRenderer.on("browse-cover-successful", function(event, filePath) {
	document.getElementById("cover-path").value = filePath;
});