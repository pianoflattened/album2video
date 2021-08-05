// TODO: set or remove app icon
const { app, BrowserWindow, dialog, ipcMain } = require("electron");
//const BezierEasing = require('bezier-easing');
const ffmpegPath = require('@ffmpeg-installer/ffmpeg').path;
const ffprobePath = require('@ffprobe-installer/ffprobe').path;
const IPC = require('ipc-node-go')
const path = require('path');
var win;

var progressBar;
const ipc = new IPC('./bin/xX_FFMP3G_BL4CKB0X_Xx')

function createWindow () {  
	win = new BrowserWindow({
		useContentSize: true,
		width: 366,    
		height: 411,
		webPreferences: {
			preload: path.join(__dirname, "src/preload.js"),
			nodeIntegration: true,
			contextIsolation: false
		},
		resizable: false,
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
	dialog.showOpenDialog({
        properties: ['openDirectory', 'showHiddenFiles'], 
        title: 'choose album directory'
    }).then(result => {
		if (!result.canceled) {
			event.reply("browse-album-successful", result.filePaths[0]);
			
			console.log("album filePaths: " + result.filePaths);
		}
	}).catch(err => {
		console.log(err);
	});
});

ipcMain.on("browse-cover", function(event) {
	dialog.showOpenDialog({
        properties: ['openFile', 'showHiddenFiles'],
        title: 'choose cover path'
    }).then(result => {
		if (!result.canceled) {
			event.reply("browse-cover-successful", result.filePaths[0]);
			console.log("cover filePaths: " + result.filePaths);
		}
	}).catch(err => {
		console.log(err);
	});
});

ipcMain.on("browse-output-path", function(event) {
	dialog.showOpenDialog({
        properties: ['openFile', 'showHiddenFiles'],
        title: 'choose output path'
    }).then(result => {
		if (!result.canceled) {
			event.reply("browse-output-successful", result.filePaths[0]);
			console.log("output filePaths: " + result.filePaths);
		}
	}).catch(err => {
		console.log(err);
	});
});

ipcMain.on("browse-output-directory", function(event) {
	dialog.showOpenDialog({
        properties: ['openDirectory', 'showHiddenFiles'],
        title: 'choose output directory'
    }).then(result => {
		if (!result.canceled) {
			event.reply("browse-output-successful", result.filePaths[0]);
			console.log("output filePaths: " + result.filePaths);
		}
	}).catch(err => {
		console.log(err);
	});
});

// ipc is so easy :D
ipcMain.on("make-video", function(event, jsonData) {
    let args = JSON.parse(jsonData);
    ipc.init(args.concat([ffprobePath, ffmpegPath]));

    ipc.on("log", console.log);
    ipc.on("error", console.error);
    ipc.on("progress-label", data => {
        event.reply("progress-label", data);
    });
    
    ipc.on("make-determinate", data => {
    	event.reply("make-determinate");
    });
    
    ipc.on("set-progress", data => {
    	event.reply("set-progress", data);
    });
    
    ipc.on("timestamps", timestamps => {
    	event.reply("timestamps", timestamps);
    });
    
    ipc.on("set-complete", data => {
    	event.reply("set-complete", data);
    });

    ipc.on('close', (code) => {
        console.log("child process closed with " + code);
    });
});


// DONT LOOK DOWN HERE THIS PART IS REALLY BAD





















































































































//var easeCurve = BezierEasing(0.25, 0.1, 0.25, 1.0);
ipcMain.on("resize-window", function(event, btnclass) {
	console.log(btnclass); // 365 409   365 436
	win.setResizable(true);
	if (btnclass == "btn btn-primary") { // faster if i calculate by hand lol
		/*setTimeout(_ => win.setSize(365, 410, true), 17);
		setTimeout(_ => win.setSize(365, 411, true), 33);
		setTimeout(_ => win.setSize(365, 414, true), 50);
		setTimeout(_ => win.setSize(365, 416, true), 67);
		setTimeout(_ => win.setSize(365, 419, true), 83);
		setTimeout(_ => win.setSize(365, 422, true), 100);
		setTimeout(_ => win.setSize(365, 425, true), 117);
		setTimeout(_ => win.setSize(365, 427, true), 133);
		setTimeout(_ => win.setSize(365, 428, true), 150);
		setTimeout(_ => win.setSize(365, 430, true), 167);
		setTimeout(_ => win.setSize(365, 431, true), 183);
		setTimeout(_ => win.setSize(365, 432, true), 200);
		setTimeout(_ => win.setSize(365, 433, true), 217);
		setTimeout(_ => win.setSize(365, 434, true), 233);
		setTimeout(_ => win.setSize(365, 435, true), 250);
		setTimeout(_ => win.setSize(365, 435, true), 283);*/
		setTimeout(_ => win.setSize(365, 436, true), 300);
	} else {
		/*setTimeout(_ => win.setSize(365, 435, true), 17);
		setTimeout(_ => win.setSize(365, 434, true), 33);
		setTimeout(_ => win.setSize(365, 431, true), 50);
		setTimeout(_ => win.setSize(365, 429, true), 67);
		setTimeout(_ => win.setSize(365, 426, true), 83);
		setTimeout(_ => win.setSize(365, 423, true), 100);
		setTimeout(_ => win.setSize(365, 420, true), 117);
		setTimeout(_ => win.setSize(365, 418, true), 133);
		setTimeout(_ => win.setSize(365, 417, true), 150);
		setTimeout(_ => win.setSize(365, 415, true), 167);
		setTimeout(_ => win.setSize(365, 414, true), 183);
		setTimeout(_ => win.setSize(365, 413, true), 200);
		setTimeout(_ => win.setSize(365, 412, true), 217);
		setTimeout(_ => win.setSize(365, 411, true), 233);
		setTimeout(_ => win.setSize(365, 410, true), 250);*/
		setTimeout(_ => win.setSize(365, 409, true), 300);
	}
	setTimeout(_ => win.setResizable(false), 350);
});
