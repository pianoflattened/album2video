// TODO: set or remove app icon
const { app, BrowserWindow, dialog, ipcMain } = require("electron");
const getTags = require('./src/get_tags.js');
const path = require('path');
const ProgressBar = require('./src/progress_bar.js');
const serialize = require('serialize-javascript');
var progressBar;

function createWindow () {  
	const win = new BrowserWindow({
		useContentSize: true,
		width: 800,    
		height: 600,
		webPreferences: {
			preload: path.join(__dirname, "src/preload.js"),
			nodeIntegration: true,
			contextIsolation: false
		},
		//resizable: false,
		autoHideMenuBar: true
	});
	win.loadFile('src/index.html');
}

app.whenReady().then(() => {
	createWindow();
	
	app.on('activate', function () { 
		if (BrowserWindow.getAllWindows().length === 0) createWindow();
	});
});

app.on('window-all-closed', function () {
	if (process.platform !== 'darwin') app.quit()
});

ipcMain.on("auto-resize", function(event, width, height) {
	let browserWindow = BrowserWindow.fromWebContents(event.sender);
	browserWindow.setContentSize(width, height);
});

ipcMain.on("browse-album", function(event) {
	dialog.showOpenDialog({properties: ['openDirectory', 'showHiddenFiles']}).then(result => {
		if (!result.canceled) {
			event.reply("browse-album-successful", result.filePaths[0]);
			
			console.log("album filePaths: " + result.filePaths);
		}
	}).catch(err => {
		console.log(err);
	});
});

ipcMain.on("browse-cover", function(event) {
	dialog.showOpenDialog({properties: ['openFile', 'showHiddenFiles']}).then(result => {
		if (!result.canceled) {
			event.reply("browse-cover-successful", result.filePaths[0]);
			console.log("cover filePaths: " + result.filePaths);
		}
	}).catch(err => {
		console.log(err);
	});
});

ipcMain.on("browse-output-path", function(event) {
	dialog.showOpenDialog({properties: ['openFile', 'showHiddenFiles']}).then(result => {
		if (!result.canceled) {
			event.reply("browse-output-successful", result.filePaths[0]);
			console.log("output filePaths: " + result.filePaths);
		}
	}).catch(err => {
		console.log(err);
	});
});

ipcMain.on("browse-output-directory", function(event) {
	dialog.showOpenDialog({properties: ['openDirectory', 'showHiddenFiles']}).then(result => {
		if (!result.canceled) {
			event.reply("browse-output-successful", result.filePaths[0]);
			console.log("output filePaths: " + result.filePaths);
		}
	}).catch(err => {
		console.log(err);
	});
});

ipcMain.on("make-video", function(event, args) {
	progressBar = new ProgressBar({
		title: 'album2video', 
		text: 'collecting files', 
		detail: 'validating album path..'
	}, app);
	console.log(args);
	console.log(getTags(args, progressBar));
});