// TODO: set or remove app icon
const { app, BrowserWindow, dialog, ipcMain } = require("electron");
const BezierEasing = require('bezier-easing');
const ffmpegPath = require('@ffmpeg-installer/ffmpeg').path;
const ffprobePath = require('@ffprobe-installer/ffprobe').path;
const IPC = require('ipc-node-go')
const path = require('path');
var win;
var baseWidth = 0;
var baseHeight = 0;
var ipc;
var progressBar;

if (process.platform == "win32") {
	ipc = new IPC('./bin/xX_FFMP3G_BL4CKB0X_Xx.exe')
} else { // (process.platform == "linux")
	ipc = new IPC('./bin/xX_FFMP3G_BL4CKB0X_Xx');
}

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
	console.log(width, height);
	baseWidth = width+16;
	baseHeight = height+16;
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
var times = 0;
ipcMain.on("make-video", function(event, jsonData) {
    let args = JSON.parse(jsonData);
    ipc.init(args.concat([ffprobePath, ffmpegPath]));
	
	if (times === 0) {
		times += 1
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
			ipc.kill();
		});

		ipc.on('close', (code) => {
			console.log("child process closed with " + code);
		});
	}
});


// DONT LOOK DOWN HERE THIS PART IS REALLY BAD





















































































































var easeCurve = BezierEasing(0.25, 0.1, 0.25, 1.0);
ipcMain.on("resize-window", function(event, data) {
	console.log(data.btnclass); // 365 409   365 436
	console.log(baseWidth, baseHeight);
	console.log(baseWidth, baseHeight+data.offsetWidth);
	win.setResizable(true); // make sure you make this mac compatible (dont do all these stupid animation timeout lines on there)
	let outalg = (n) => baseHeight+(n*data.offsetWidth);
	let inalg = (n) => baseHeight+data.offsetWidth-(n*data.offsetWidth);
	if (data.btnclass == "btn btn-primary") {
		for (i = 1; i <= 21; i++) {
			setTimeout(_ => win.setSize(baseWidth, Math.round(outalg(easeCurve(i/21))), true), (350*i)/21);
		} 
	} else {
		for (i = 1; i <= 21; i++) {
			setTimeout(_ => win.setSize(baseWidth, Math.round(inalg(easeCurve(i/21))), true), (350*i)/21);
		} 
	}
	setTimeout(_ => win.setResizable(false), 350);
});

/*setTimeout(_ => win.setSize(baseWidth, baseHeight+1, true), 17);
setTimeout(_ => win.setSize(baseWidth, baseHeight+2, true), 33);
setTimeout(_ => win.setSize(baseWidth, baseHeight+5, true), 50);
setTimeout(_ => win.setSize(baseWidth, baseHeight+7, true), 67);
setTimeout(_ => win.setSize(baseWidth, baseHeight+10, true), 83);
setTimeout(_ => win.setSize(baseWidth, baseHeight+13, true), 100);
setTimeout(_ => win.setSize(baseWidth, baseHeight+16, true), 117);
setTimeout(_ => win.setSize(baseWidth, baseHeight+18, true), 133);
setTimeout(_ => win.setSize(baseWidth, baseHeight+19, true), 150);
setTimeout(_ => win.setSize(baseWidth, baseHeight+21, true), 167);
setTimeout(_ => win.setSize(baseWidth, baseHeight+22, true), 183);
setTimeout(_ => win.setSize(baseWidth, baseHeight+23, true), 200);
setTimeout(_ => win.setSize(baseWidth, baseHeight+24, true), 217);
setTimeout(_ => win.setSize(baseWidth, baseHeight+25, true), 233);
setTimeout(_ => win.setSize(baseWidth, baseHeight+26, true), 250); 
setTimeout(_ => win.setSize(baseWidth, baseHeight+27, true), 300);// baseheight was 27 here */

/*setTimeout(_ => win.setSize(baseWidth, baseHeight+26 435, true), 17);
setTimeout(_ => win.setSize(baseWidth, baseHeight+25 434, true), 33);
setTimeout(_ => win.setSize(baseWidth, baseHeight+22 431, true), 50);
setTimeout(_ => win.setSize(baseWidth, baseHeight+20 429, true), 67);
setTimeout(_ => win.setSize(baseWidth, baseHeight+17 426, true), 83);
setTimeout(_ => win.setSize(baseWidth, baseHeight+14 423, true), 100);
setTimeout(_ => win.setSize(baseWidth, baseHeight+11 420, true), 117);
setTimeout(_ => win.setSize(baseWidth, baseHeight+9 418, true), 133);
setTimeout(_ => win.setSize(baseWidth, baseHeight+8 417, true), 150);
setTimeout(_ => win.setSize(baseWidth, baseHeight+6 415, true), 167);
setTimeout(_ => win.setSize(baseWidth, baseHeight+5 414, true), 183);
setTimeout(_ => win.setSize(baseWidth, baseHeight+4 413, true), 200);
setTimeout(_ => win.setSize(baseWidth, baseHeight+3 412, true), 217);
setTimeout(_ => win.setSize(baseWidth, baseHeight+2 411, true), 233);
setTimeout(_ => win.setSize(baseWidth, baseHeight+1 410, true), 250);
setTimeout(_ => win.setSize(baseWidth, baseHeight, true), 300);*/