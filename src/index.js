const { ipcRenderer } = require("electron");
const browseAlbumBtn = document.getElementById('browse-album');

browseAlbumBtn.addEventListener('click', function() {
	ipcRenderer.send("browse-album");
});

ipcRenderer.on("browse-album-successful", function(event, filePath) {
	document.getElementById("album-dir").value = filePath;
});