const ffprobePath = require('@ffprobe-installer/ffprobe').path;
const { ffprobe } = require('kiss-ffmpeg');
const fs = require('fs');
const mime = require('mime-types');
const path = require('path');
var sizeOf = require('image-size');

ffprobe.command = ffprobePath;

function sleepSync(ms) {
    let start = new Date().getTime(), expire = start + ms;
    while (new Date().getTime() < expire) { }
    return;
}

/*// makes directories synchronously iterable. thankz stack over flow :D
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
});*/

function overallTrackNumber(track, disc, discTracks) {
	let n = track;
	for (let i = 1; i <= disc-1; i++) {
		n += discTracks[i];
	}
	return n;
}

module.exports = async function getTags(form, progressBar) {
	// indeterminate progress bar
	await progressBar.makeIndeterminate();
	await progressBar.setLabel('validating album path..')

	let stats;
	try { // check if form.albumDirectory exists + is a directory
		stats = fs.statSync(form.albumDirectory);
	} catch (err) {
		await progressBar.error(err.toString()); 
		return false; 
	}
	if (stats.isDirectory()) { /* pass */ } else if (stats.isFile()) {
		form.albumDirectory = path.dirname(form.albumDirectory);
	} else {
		await progressBar.error(form.albumDirectory + ' is not a file or directory');
		return false;
	}

	// loop through form.albumDirectory
	await progressBar.setLabel('reading ' + form.albumDirectory + '..');
	const albumDir = fs.opendirSync(form.albumDirectory);
	var audioFiles = [];
	let imageFiles = {};
	let discTracks = {};
	let disc;
	let track;
	let fileMimetype;
	let metadata;
    let promises = [];
	for await (const f of albumDir) { // need to somehow ".then()" this for loop :(
		if (f.name != "concat.wav") {
			fullpath = path.join(form.albumDirectory, f.name);
			await progressBar.setLabel('reading ' + fullpath + '..');

			fileMimetype = mime.lookup(fullpath);
            console.log(fullpath);
			switch (fileMimetype.split("/")[0]) {
				case "audio": // if mimetype is audio/* then get its tags + store in dictionary with path as key
                    promises.push(ffprobe(fullpath).then(function(info) {
                        metadata = info.format;
                        if (metadata.tags.disc) {
						    disc = parseInt(metadata.tags.disc.split("/")[0]);
					    } else { disc = 1; }
					    track = parseInt(metadata.tags.track.split("/")[0]);
					    audioFiles.push({
						    filename: metadata.filename,
						    artist: metadata.tags.artist,
						    albumArtist: metadata.tags.albumArtist || "",
						    title: metadata.tags.title,
						    track: track,
						    disc: disc,
						    time: parseFloat(metadata.duration) // - parseFloat(metadata.start_time)
					    });
					    if (discTracks[disc]) {
						    if (discTracks[disc] < track) {
							    discTracks[disc] = track;
						    }
					    } else {
						    discTracks[disc] = track;
					    }
                    }).catch(function(err) {
                        console.log(err);
                        (async () => {await progressBar.error(err.toString())})();
                    }));
					break;
				case "image": // if mimetype is image/* and form.detectCover add it to cover art candidates list
					if (form.detectCover) imageFiles[f.name] = fullpath;
					break;
				default:
					break;
			}
		}
	}

	await progressBar.setLabel('checking if ' + form.albumDirectory + 'has no sound files..');
	if (!audioFiles) { // if form.albumDirectory has no sound files return an error
		await progressBar.error(form.albumDirectory + ' does not contain any sound files');
		return false;
	}
    
    await progressBar.setLabel('ordering audio files..');
    Promise.all(promises).then(function() { // FURIOUS that i have to put this stupid while loop in here
        let limit = 300;
        let c = 0;
        while (promises.length != audioFiles.length) {
            sleepSync(16);
            console.log("slept");
            c += 1;
            if (c >= limit) {
                console.log("broken");
                break;
            }
        }
        return c < limit;
    }).then(function(ok) {
        if (!ok) {
            Promise.reject('audio files miscounted');
        }
        console.log("VERIFY THAT THERE ARE " + promises.length + " TRACKS IN THE ALBUM");
        console.log(discTracks);
	    audioFiles.sort(function(a, b) {
		    aOverall = overallTrackNumber(a.track, a.disc || 1, discTracks);
		    bOverall = overallTrackNumber(b.track, b.disc || 1, discTracks);
            // console.log(a.track, b.track);
            // console.log(a.disc, b.disc);
            // console.log(aOverall, bOverall);
            // console.log("---------");
            if (aOverall > bOverall) return 1;
            if (aOverall < bOverall) return -1;
            if (aOverall == bOverall) return 0;
	    });
    });

	await progressBar.setLabel('checking cover art..');
	// if form.detectCover is on + no image files set the cover to a 1920 x 1080 all black image
	if (!imageFiles && form.detectCover) { 
		form.coverPath = '../assets/black.png'
	}

	if (!form.detectCover) { // if form.detectCover is off then check if form.coverPath exists
		let stats;
		try { // check if form.coverPath exists + is a file
			stats = fs.statSync(form.coverPath);
		} catch (err) { // if it doesn't exist return an error
			await progressBar.error(err.toString()); 
			return false; 
		}
		if (stats.isFile()) { /* pass */ } else { // if it's not a file return an error
			await progressBar.error(form.coverPath + ' is not a file');
			return false;
		}
	} else {
		await progressBar.setLabel('choosing cover art..');
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
	
    await progressBar.setLabel('wait..');
	return { // return object with relevant tags + location of cover art
		form: form,
		audioFiles: audioFiles,
        progressBar: progressBar,
        discTracks: discTracks
	}
};
