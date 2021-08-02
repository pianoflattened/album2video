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



// like he says "oh it's good practice to do this it's good practice to make sure you set 
// MAGIC_PRESERVE_ATIME on linux-based platforms" yet not only he doesn't make this a default
// feature DUDE MAKES HIS EXAMPLES FUCKING INDECIPHERABLE im so frustrated i spent 6 hours of my
// day figuring out how to make this work because dude wouldn't spell out exactly how this worked
// this shit should be as minimal as possible i should be importing one function that spits out a
// file's detected mimetype and instead i have to dedicate all this time to figuring out how to get
// like 12 goddamn characters into memory so i can know what the hell guys r using my code on and
// adjust behavior accordingly like !!!!!!!!!!!



// AND THEN I TEST TH SHIT AND IT GIVES ME AN AMBIGUOUS "FAILED TO INITIALIZE" ERROR WHAT AM I
// SUPPOSED TO DO W THIS SHIT IM JUST GOING TO CHECK FILE EXTENSIONS MAYBE SOMEWHERE ALONG THE LINE
// ILL PACKAGE THIS WITH SOME BLACK BOX BINARY EXECUTABLE THAT JUST DOES IT PROBABLY BASED ON
// GOLANG OR PYTHON OR LITERALLY ANYTHING BUT PREFERRABLY SOMETHING FAST IDK ANYWAY THATS WHAT
// WE'RE DOING NOW :D

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
	FileMagic.getInstance().then((magic) => { // https://files.catbox.moe/law8p4.png
		for (const f of albumDir) {
			fullpath = path.join(form.albumDirectory, f.name);
			let file_mimetype = magic.detect(file, magic.flags | MagicFlags.MAGIC_MIME_TYPE);
			console.log(file_mimetype);
		}
		FileMagic.close();
	}).catch((err) => {
		console.log(err);
		FileMagic.close();
	});

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