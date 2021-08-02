const ffprobePath = require('@ffprobe-installer/ffprobe').path;
const ffmpeg = require('fluent-ffmpeg');
const fs = require('fs');
const path = require('path');
const ProgressBar = require('./progress_bar.js');

ffmpeg.setFfprobePath(ffprobePath);

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
		progressBar.error('albumDirectory is not a file or directory');
	}

	return true;
	// loop through form.albumDirectory


		// if mimetype is audio/* then get its tags + store in dictionary with path as a key
		// if mimetype is image/* and form.detectCover is on then add it to the list of cover art candidates
		// if form.albumDirectory has no sound files return an error
		// if form.detectCover is on and there are no image files then set the cover to a 1920 x 1080 all black image
	// if form.detectCover is off then check if form.coverPath exists
		// if it doesn't exist return an error
	// if form.detectCover is on then pick from the list of cover art candidates as follows:
		// folder.png
		// cover.png
		// folder.jpg
		// cover.jpg
		// largest image that has "cover" in its name
		// largest image that has "art" in its name
		// largest image in directory
	// return object with relevant tags + location of cover art
};