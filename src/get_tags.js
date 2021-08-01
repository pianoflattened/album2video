const ffprobePath = require('@ffprobe-installer/ffprobe').path;
const ffmpeg = require('fluent-ffmpeg');
const fs = require('fs');
const mmm = require('mmmagic');
const path = require('path');
const ProgressBar = require('electron-progressbar');

ffmpeg.setFfprobePath(ffprobePath);

module.exports = {
	get_tags: function(form) {
		// indeterminate progress bar
		// check if form.albumPath exists + is a directory
			// if it isn't a directory, check if it's a file
				// if it's a file, set form.albumPath to its parent directory and continue
				// if not then return an error
			// if it doesn't exist return an error
		// loop through form.albumPath
			// if mimetype is audio/* then get its tags + store in dictionary with path as a key
			// if mimetype is image/* and form.detectCover is on then add it to the list of cover art candidates
			// if form.albumPath has no sound files return an error
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
	}
};