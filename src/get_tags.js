const ffprobePath = require('@ffprobe-installer/ffprobe').path;
const ffmpeg = require('fluent-ffmpeg');
const fs = require('fs').promises;
const mmm = require('mmmagic');
const path = require('path');

ffmpeg.setFfprobePath(ffprobePath);

module.exports = {
	get_tags: function(args) {
		// indeterminate progress bar
		// check if args.albumPath exists + is a directory
			// if it isn't a directory, check if it's a file
				// if it's a file, set args.albumPath to its parent directory and continue
				// if not then return an error
			// if it doesn't exist return an error
		// loop through args.albumPath
			// if mimetype is audio/* then get its tags + store in dictionary with path as a key
			// if mimetype is image/* and args.detectCover is on then add it to the list of cover art candidates
			// if args.albumPath has no sound files return an error
			// if args.detectCover is on and there are no image files then set the cover to a 1920 x 1080 all black image
		// if args.detectCover is off then check if args.coverPath exists
			// if it doesn't exist return an error
		// if args.detectCover is on then pick from the list of cover art candidates as follows:
			// folder.png
			// cover.png
			// folder.jpg
			// cover.jpg
			// largest image that has "cover" in its name
			// largest image that has "art" in its name
			// largest image in directory
	}
};