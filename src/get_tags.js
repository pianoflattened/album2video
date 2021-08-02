const ffprobePath = require('@ffprobe-installer/ffprobe').path;
const { ffprobeSync } = require('kiss-ffmpeg');
const ffmpeg = require('fluent-ffmpeg');
const fs = require('fs');
const mime = require('mime-types');
const path = require('path');
const ProgressBar = require('./progress_bar.js');
var sizeOf = require('image-size');

ffprobeSync.command = ffprobePath;

// makes directories synchronously iterable. thankz stack over flow :D
const p = fs.Dir.prototype;
if (p.hasOwnProperty(Symbol.iterator)) { return; }
const entriesSync = function* () {
try {
	let dirent;
	while ((dirent = this.readSync()) !== null) { yield dirent; }
    } finally { this.closeSync(); }
}
if (!p.hasOwnProperty(entriesSync)) { p.entriesSync = entriesSync; }
Object.defineProperty(p, Symbol.iterator, {
	configurable: true,
	enumerable: false,
	value: entriesSync,
	writable: true
});

function overallTrackNumber(track, disc, discTracks) {
	let n = track;
	for (let i = 0; i <= disc-2; i++) {
		n += discTracks[i];
	}
	return n;
}

module.exports = function getTags (form) {
	// indeterminate progress bar
	let progressBar = new ProgressBar({
		title: 'album2video',
		text: 'collecting files',
		detail: 'validating album path..'
	});
	
	console.log(form.albumDirectory);

	let stats;
	try { // check if form.albumDirectory exists + is a directory
		stats = fs.statSync(form.albumDirectory);
	} catch (err) {
		progressBar.error(err.toString()); 
		return false; 
	}
	if (stats.isDirectory()) { /* pass */ } else if (stats.isFile()) {
		form.albumDirectory = path.dirname(form.albumDirectory);
	} else {
		progressBar.error(form.albumDirectory + ' is not a file or directory');
		return false;
	}

	// loop through form.albumDirectory
	progressBar.detail = 'reading ' + form.albumDirectory + '..';
	const albumDir = fs.opendirSync(form.albumDirectory);
	var audioFiles = [];
	let imageFiles = {};
	let discTracks = {};
	let disc;
	let track;
	let fileMimetype;
	let metadata;
	for (const f of albumDir) {
		console.log(f.name);
		if (f.name != "concat.wav") {
			fullpath = path.join(form.albumDirectory, f.name);
			progressBar.detail = 'reading ' + fullpath + '..';

			fileMimetype = mime.lookup(fullpath);
			console.log(fileMimetype);
			switch (fileMimetype.split("/")[0]) {
				case "audio": // if mimetype is audio/* then get its tags + store in dictionary with path as key
					try {
						metadata = ffprobeSync(fullpath).format;
					} catch (err) {
						progressBar.error(err.toString());
					}
					if (metadata.tags.disc) {
						disc = parseInt(metadata.tags.disc.split("/")[0]);
					} else { disc = 1; }
					track = parseInt(metadata.tags.track.split("/")[0]);
					audioFiles.push({
						filename: metadata.filename, 
						fullpath: fullpath,
						artist: metadata.tags.artist,
						albumArtist: metadata.tags.albumArtist || "",
						title: metadata.tags.title,
						track: track,
						disc: disc,
						length: parseFloat(metadata.duration) - parseFloat(metadata.start_time)
					});
					console.log(audioFiles);
					if (discTracks[disc]) {
						if (discTracks[disc] < track) {
							discTracks[disc] = track;
						}
					} else {
						discTracks[disc] = track;
					}
					break;
				case "image": // if mimetype is image/* and form.detectCover add it to cover art candidates list
					if (form.detectCover) imageFiles[f.name] = fullpath;
					break;
				default:
					break;
			}

			progressBar.detail = 'checking if ' + form.albumDirectory + 'has no sound files..';
			if (!audioFiles) { // if form.albumDirectory has no sound files return an error
				progressBar.error(form.albumDirectory + ' does not contain any sound files');
				return false;
			}

			progressBar.detail = 'ordering audio files..';
			audioFiles.sort(function(a, b) {
				aOverall = overallTrackNumber(a.track, a.disc, discTracks);
				bOverall = overallTrackNumber(b.track, b.disc, discTracks);
				return aOverall - bOverall;
			});

			progressBar.detail = 'checking cover art..';
			// if form.detectCover is on + no image files set the cover to a 1920 x 1080 all black image
			if (!imageFiles && form.detectCover) { 
				form.coverPath = '../assets/black.png'
			}

			if (!form.detectCover) { // if form.detectCover is off then check if form.coverPath exists
				let stats;
				try { // check if form.coverPath exists + is a file
					stats = fs.statSync(form.coverPath);
				} catch (err) { // if it doesn't exist return an error
					progressBar.error(err.toString()); 
					return false; 
				}
				if (stats.isFile()) { /* pass */ } else { // if it's not a file return an error
					progressBar.error(form.coverPath + ' is not a file');
					return false;
				}
			} else {
				progressBar.detail = 'choosing cover art..';
				let found = false; // check for these names, highest to lowest priority
				["folder.png", "cover.png", "folder.jpg", "cover.jpg"].forEach(function(n) {
					if (n in imageFiles) {
						form.coverPath = imageFiles[n];
						found = true;
					}
				});

				// i could definitely optimize this section what with the recalculating of image sizes 3 times 
				// in a row BUT this is the easiest way for me to conceptualize this and i just want it to work
				if (!found) { // largest image that has "cover" in its name
					let withCover = {};
					let dim;
					Object.keys(imageFiles).forEach(function(key) {
						if (key.toLowerCase().includes("cover")) {
							dim = sizeOf(imageFiles[key]);
							withCover[dim.width*dim.height] = imageFiles[key];
						}
					});
					if (withCover) {
						form.coverPath = withCover[Math.max(...Object.keys(withCover))];
						found = true;
					} 
				}

				if (!found) { // largest image that has "art" in its name
					let withArt = {};
					let dim;
					Object.keys(imageFiles).forEach(function(key) {
						if (key.toLowerCase().includes("art")) {
							dim = sizeOf(imageFiles[key]);
							withCover[dim.width*dim.height] = imageFiles[key];
						}
					});
					if (withArt) {
						form.coverPath = withArt[Math.max(...Object.keys(withArt))];
						found = true;
					}
				}

				if (!found) { // largest image in directory
					let covers = {};
					let dim;
					Object.keys(imageFiles).forEach(function(key) {
						dim = sizeOf(imageFiles[key]);
						covers[dim.width*dim.height] = imageFiles[key];
					});
					form.coverPath = covers[Math.max(...Object.keys(covers))];
				}
			}
		}
	}
	
	return { // return object with relevant tags + location of cover art
		form: form,
		audioFiles: audioFiles
	}
};