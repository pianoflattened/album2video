const ffprobePath = require('@ffprobe-installer/ffprobe').path;
const ffmpeg = require('fluent-ffmpeg');
const fs = require('fs');
const { FileMagic, MagicFlags } = require('@npcz/magic');
const path = require('path');
const ProgressBar = require('./progress_bar.js');

ffmpeg.setFfprobePath(ffprobePath);

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

// this magic library sucks and setup is ridiculous and cryptic shout out to this guy abdes good
// job making ur shit impenetrable i spent like a whole hour squinting at ur readme and trying to
// figure out where the require was supposed to go to and when you would get to the meat of how a
// normal person is supposed to use this. i just want a library that prints out a little string
// after reading th contents of a file - all this setup shit seems 2 be boilerplate and i dont
// understand why it has to be done on th user-side. i am not the only one who finds ur shit to be
// absolutely preposterous the literal only other use of what should be the most popular mimetype
// detection library fr node js (because it's the only maintained one, mmmagic broke on me bc it
// was built for a versino of node that was too old and i didnt want to spend time fixing it lol)
// involved literal copy-pasting of your example because nobody wants to fucking wrap their head
// around an entire god damned zalgathor froghorn's anthologie of eldtrich magicks obskura anyway 
// heres the asinine boilerplate i had to write to use filemagic and get the audio/mp3 string or
// whatever maybe make ur library usable and u wont send any more immature 19yr olds into a burning
// fury and get thusly chewed out in some dark corner of their github code thanks
FileMagic.magicFile = path.normalize(
	path.join(__dirname, 'node_modules', '@npcz', 'magic', 'dist', 'magic.mgc')
);

if (process.platform === 'darwin' || process.platform === 'linux') {
	FileMagic.defaulFlags = MagicFlags.MAGIC_PRESERVE_ATIME;
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
		progressBar.error('albumDirectory is not a file or directory');
	}

	// loop through form.albumDirectory
	progressBar.detail = 'reading ' + form.albumDirectory;
	const albumDir = fs.opendirSync(form.albumDirectory);
	for (const f of albumDir) {
		fullpath = path.join(form.albumDirectory, f.name);
		
	}

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