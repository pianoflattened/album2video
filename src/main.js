// TODO: set or remove app icon
// TODO: download fonts for consistency across platforms
const { app, BrowserWindow, dialog, ipcMain } = require("electron");
const BezierEasing = require('bezier-easing');
const child_process = require('child_process');
// const ffmpegPath = require('ffmpeg-static');
// const ffprobePath = require('ffprobe-static').path;
const fs = require('fs');
const IPC = require('ipc-node-go')
const path = require('path');
import { platform } from 'os';
import { rootPath } from 'electron-root-path';

function getPlatform() {
  switch (platform()) {
    case 'aix':
    case 'freebsd':
    case 'linux':
    case 'openbsd':
    case 'android':
      return 'linux';
    case 'darwin':
    case 'sunos':
      return 'mac';
    case 'win32':
      return 'win';
  }
};

// debug function
// var _getAllFilesFromFolder = function(dir) {
//     var filesystem = require("fs");
//     var results = [];
//     filesystem.readdirSync(dir).forEach(function(file) {
//         file = dir+'/'+file;
//         var stat = filesystem.statSync(file);
//         if (stat && stat.isDirectory()) {
//             results.push(file)
//             // results = results.concat(_getAllFilesFromFolder(file))
//         } else results.push(file);
//     });
//     return results;
// };
//
// console.log(_getAllFilesFromFolder("."));

const root = rootPath;console.log("30");
const binPath = process.mainModule.filename.indexOf('app.asar') !== -1 ? path.join(path.dirname(app.getAppPath()), '..', './resources', './bin') : path.join(root, './bin');

function getBin(p) {
    return path.resolve(path.join(binPath, p)) + (getPlatform() == "win" ? ".exe" : "");
}

const blackboxPath = getBin('xX_FFMP3G_BL4CKB0X_Xx');
const trackfmtPath = getBin('trackfmt');
const ffmpegPath = getBin('ffmpeg');
const ffprobePath = getBin('ffprobe');

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
			preload: path.join(__dirname, "src", "preload_launcher.js"),
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
    xX_FFMP3G_BL4CKB0X_Xx = new IPC(blackboxPath);
	// if (process.platform == "win32") {
	// 	xX_FFMP3G_BL4CKB0X_Xx = new IPC('./bin/xX_FFMP3G_BL4CKB0X_Xx.exe');
	// } else { // (process.platform == "linux")
	// 	xX_FFMP3G_BL4CKB0X_Xx = new IPC('./bin/xX_FFMP3G_BL4CKB0X_Xx');
	// }
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
	// let trackfmtPath;
    //
	// if (process.platform == "win32") {
	// 	trackfmtPath = "./bin/trackfmt.exe";
	// 	// trackfmt = new IPC('./bin/trackfmt.exe');
	// } else { // (process.platform == "linux")
	// 	// trackfmt = new IPC('./bin/trackfmt');
	// 	trackfmtPath = "./bin/trackfmt";
	// }
	var trackfmt = child_process.spawn(trackfmtPath, [format, timestamps]);
	trackfmt.stdout.on('data', data => {
		event.reply("timestamps", data.toString());
	});

	trackfmt.stderr.on('data', data => {
		console.log(data.toString());
	});

	trackfmt.on("close", code => {
		console.log("trackfmt closed with " + code);
	});
});


//     ____  ____     _   ______  ______   __    ____  ____  __ __    __   __   __   __   __   __
//    / __ \/ __ \   / | / / __ \/_  __/  / /   / __ \/ __ \/ //_/   / /  / /  / /  / /  / /  / /
//   / / / / / / /  /  |/ / / / / / /    / /   / / / / / / / ,<     / /  / /  / /  / /  / /  / /
//  / /_/ / /_/ /  / /|  / /_/ / / /    / /___/ /_/ / /_/ / /| |   /_/  /_/  /_/  /_/  /_/  /_/
// /_____/\____/  /_/ |_/\____/ /_/    /_____/\____/\____/_/ |_|  (_)  (_)  (_)  (_)  (_)  (_)


var easeCurve = BezierEasing(0.25, 0.1, 0.25, 1.0);
ipcMain.on("resize-window", function(event, data) {
	let browserWindow = BrowserWindow.fromWebContents(event.sender);
	console.log(data.btnclass); // 365 409   365 436
	console.log(baseWidth, baseHeight);
	console.log(baseWidth, baseHeight+data.offsetHeight);
	win.setResizable(true); // make sure you make this mac compatible (dont do all these stupid animation timeout lines on there)
	let outalg = (n) => baseHeight+(n*data.offsetHeight);
	let inalg = (n) => baseHeight+data.offsetHeight-(n*data.offsetHeight);
	// i have these lines hidden in my editor and you should too          1000 of below
	// setTimeout(_ => browserWindow.setContentSize(baseWidth, Math.round(outalg(easeCurve(1/21))), true), 17);
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
});
