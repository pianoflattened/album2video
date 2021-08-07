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
var xX_FFMP3G_BL4CKB0X_Xx, trackfmt;
var progressBar;

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
		resizable: true,
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
	baseWidth = width;
	baseHeight = height;
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
var blackboxtimes = 0;
var trackfmttimes = 0;
ipcMain.on("make-video", function(event, jsonData) {
    let args = JSON.parse(jsonData);
	if (process.platform == "win32") {
		xX_FFMP3G_BL4CKB0X_Xx = new IPC('./bin/xX_FFMP3G_BL4CKB0X_Xx.exe');
	} else { // (process.platform == "linux")
		xX_FFMP3G_BL4CKB0X_Xx = new IPC('./bin/xX_FFMP3G_BL4CKB0X_Xx');
	}
    xX_FFMP3G_BL4CKB0X_Xx.init(args.concat([ffprobePath, ffmpegPath]));
	xX_FFMP3G_BL4CKB0X_Xx.on("log", console.log);
	xX_FFMP3G_BL4CKB0X_Xx.on("error", console.error);
	xX_FFMP3G_BL4CKB0X_Xx.on("close", (code) => console.log("xX_FFMP3G_BL4CKB0X_Xx closed with " + code));
	xX_FFMP3G_BL4CKB0X_Xx.on("progress-label", data => {
		event.reply("progress-label", data);
	});

	xX_FFMP3G_BL4CKB0X_Xx.on("make-determinate", data => {
		event.reply("make-determinate");
	});

	xX_FFMP3G_BL4CKB0X_Xx.on("set-progress", data => {
		event.reply("set-progress", data);
	});

	xX_FFMP3G_BL4CKB0X_Xx.on("timestamps", timestamps => {
		event.reply("get-timestamp-format", JSON.stringify(timestamps));
	});

	xX_FFMP3G_BL4CKB0X_Xx.on("set-complete", data => {
		event.reply("set-complete", data);
		xX_FFMP3G_BL4CKB0X_Xx.kill();
	});
});

ipcMain.on("timestamp-format", function(event, data) {
	let format = data.format;
	let timestamps = data.timestamps;

	if (process.platform == "win32") {
		trackfmt = new IPC('./bin/trackfmt.exe');
	} else { // (process.platform == "linux")
		trackfmt = new IPC('./bin/trackfmt');
	}
	trackfmt.init([format, timestamps]);
	trackfmt.on("log", console.log);
	trackfmt.on("error", console.error);
	trackfmt.on("close", (code) => console.log("trackfmt closed with " + code));
	trackfmt.on("result", data => {
		console.log(data);
		event.reply("timestamps", data);
		trackfmt.kill();
	});
});

// DONT LOOK DOWN HERE THIS PART IS REALLY BAD

var easeCurve = BezierEasing(0.25, 0.1, 0.25, 1.0);
ipcMain.on("resize-window", function(event, data) {
	let browserWindow = BrowserWindow.fromWebContents(event.sender);
	console.log(data.btnclass); // 365 409   365 436
	console.log(baseWidth, baseHeight);
	console.log(baseWidth, baseHeight+data.offsetHeight);
	win.setResizable(true); // make sure you make this mac compatible (dont do all these stupid animation timeout lines on there)
	let outalg = (n) => baseHeight+(n*data.offsetHeight);
	let inalg = (n) => baseHeight+data.offsetHeight-(n*data.offsetHeight);
	if (data.btnclass == "btn btn-primary") {
		setTimeout(_ => browserWindow.setContentSize(baseWidth, Math.round(outalg(easeCurve(1/21))), true), 17);
		setTimeout(_ => browserWindow.setContentSize(baseWidth, Math.round(outalg(easeCurve(2/21))), true), 33);
		setTimeout(_ => browserWindow.setContentSize(baseWidth, Math.round(outalg(easeCurve(3/21))), true), 50);
		setTimeout(_ => browserWindow.setContentSize(baseWidth, Math.round(outalg(easeCurve(4/21))), true), 67);
		setTimeout(_ => browserWindow.setContentSize(baseWidth, Math.round(outalg(easeCurve(5/21))), true), 83);
		setTimeout(_ => browserWindow.setContentSize(baseWidth, Math.round(outalg(easeCurve(6/21))), true), 100);
		setTimeout(_ => browserWindow.setContentSize(baseWidth, Math.round(outalg(easeCurve(7/21))), true), 117);
		setTimeout(_ => browserWindow.setContentSize(baseWidth, Math.round(outalg(easeCurve(8/21))), true), 133);
		setTimeout(_ => browserWindow.setContentSize(baseWidth, Math.round(outalg(easeCurve(9/21))), true), 150);
		setTimeout(_ => browserWindow.setContentSize(baseWidth, Math.round(outalg(easeCurve(10/21))), true), 167);
		setTimeout(_ => browserWindow.setContentSize(baseWidth, Math.round(outalg(easeCurve(11/21))), true), 183);
		setTimeout(_ => browserWindow.setContentSize(baseWidth, Math.round(outalg(easeCurve(12/21))), true), 200);
		setTimeout(_ => browserWindow.setContentSize(baseWidth, Math.round(outalg(easeCurve(13/21))), true), 217);
		setTimeout(_ => browserWindow.setContentSize(baseWidth, Math.round(outalg(easeCurve(14/21))), true), 233);
		setTimeout(_ => browserWindow.setContentSize(baseWidth, Math.round(outalg(easeCurve(15/21))), true), 250);
		setTimeout(_ => browserWindow.setContentSize(baseWidth, Math.round(outalg(easeCurve(16/21))), true), 267);
		setTimeout(_ => browserWindow.setContentSize(baseWidth, Math.round(outalg(easeCurve(17/21))), true), 283);
		setTimeout(_ => browserWindow.setContentSize(baseWidth, Math.round(outalg(easeCurve(18/21))), true), 300);
		setTimeout(_ => browserWindow.setContentSize(baseWidth, Math.round(outalg(easeCurve(19/21))), true), 317);
		setTimeout(_ => browserWindow.setContentSize(baseWidth, Math.round(outalg(easeCurve(20/21))), true), 333);
		setTimeout(_ => browserWindow.setContentSize(baseWidth, Math.round(outalg(easeCurve(21/21))), true), 350);
	} else {
		for (i = 1; i <= 21; i++) {
			setTimeout(_ => browserWindow.setContentSize(baseWidth, Math.round(inalg(easeCurve(2/21))), true), 33);
			setTimeout(_ => browserWindow.setContentSize(baseWidth, Math.round(inalg(easeCurve(1/21))), true), 17);
			setTimeout(_ => browserWindow.setContentSize(baseWidth, Math.round(inalg(easeCurve(3/21))), true), 50);
			setTimeout(_ => browserWindow.setContentSize(baseWidth, Math.round(inalg(easeCurve(4/21))), true), 67);
			setTimeout(_ => browserWindow.setContentSize(baseWidth, Math.round(inalg(easeCurve(5/21))), true), 83);
			setTimeout(_ => browserWindow.setContentSize(baseWidth, Math.round(inalg(easeCurve(6/21))), true), 100);
			setTimeout(_ => browserWindow.setContentSize(baseWidth, Math.round(inalg(easeCurve(7/21))), true), 117);
			setTimeout(_ => browserWindow.setContentSize(baseWidth, Math.round(inalg(easeCurve(8/21))), true), 133);
			setTimeout(_ => browserWindow.setContentSize(baseWidth, Math.round(inalg(easeCurve(9/21))), true), 150);
			setTimeout(_ => browserWindow.setContentSize(baseWidth, Math.round(inalg(easeCurve(10/21))), true), 167);
			setTimeout(_ => browserWindow.setContentSize(baseWidth, Math.round(inalg(easeCurve(11/21))), true), 183);
			setTimeout(_ => browserWindow.setContentSize(baseWidth, Math.round(inalg(easeCurve(12/21))), true), 200);
			setTimeout(_ => browserWindow.setContentSize(baseWidth, Math.round(inalg(easeCurve(13/21))), true), 217);
			setTimeout(_ => browserWindow.setContentSize(baseWidth, Math.round(inalg(easeCurve(14/21))), true), 233);
			setTimeout(_ => browserWindow.setContentSize(baseWidth, Math.round(inalg(easeCurve(15/21))), true), 250);
			setTimeout(_ => browserWindow.setContentSize(baseWidth, Math.round(inalg(easeCurve(16/21))), true), 267);
			setTimeout(_ => browserWindow.setContentSize(baseWidth, Math.round(inalg(easeCurve(17/21))), true), 283);
			setTimeout(_ => browserWindow.setContentSize(baseWidth, Math.round(inalg(easeCurve(18/21))), true), 300);
			setTimeout(_ => browserWindow.setContentSize(baseWidth, Math.round(inalg(easeCurve(19/21))), true), 317);
			setTimeout(_ => browserWindow.setContentSize(baseWidth, Math.round(inalg(easeCurve(20/21))), true), 333);
			setTimeout(_ => browserWindow.setContentSize(baseWidth, Math.round(inalg(easeCurve(21/21))), true), 350);
		}
	}
});

/*
17
33
50
67
83
100
117
133
150
167
183
200
217
233
250
267
283
300
317
333
350
*/

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
