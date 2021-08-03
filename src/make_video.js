const ffmpegPath = require('@ffmpeg-installer/ffmpeg').path;
const { FFmpeg } = require("kiss-ffmpeg");
const fs = require("fs").promises;
const path = require('path');
const tmp = require('tmp');

FFmpeg.command = ffmpegPath;
let stage1 = new FFmpeg();
let stage1Proc;

// wanted something more robust than a substring of iso which breaks when u get past 99:99:99
function secondsToTimestamp(seconds) { 
    let hours = Math.floor(seconds/3600);
    let minutes = Math.floor(seconds/60) - (hours*60);
    let secs = seconds - (hours*3600) - (minutes*60);
    let t = [
        hours.toString(), 
        minutes.toString().padStart(2, "0"), 
        secs.toString().padStart(2, "0")
    ].join(":");
    while ((t.startsWith("0") || t.startsWith(":")) && t.length > 4) {
        t = t.substr(1);    
    }
    return t;
}

module.exports = async function makeVideo(data, textArea) {
    if (data.form.separateFiles) {
    
    } else {
        await makeSingleVideo(data, textArea);
    }
};

async function makeSingleVideo(data, textArea) {
    concatWavPath = await makeConcatWav(data);
    console.log(concatWavPath);
    // run ffmpeg command that creates video with concat.wav and cover (2/2)
}

async function makeConcatWav(data) {
    let length = 0;
    let fileList = "";
    let timestamps = [];
    let concatWavPath;
    for (const f of data.audioFiles) { // calculate timestamps
        timestamps.push({
            title: f.title,
            time: secondsToTimestamp(length)
        });
        // https://www.ffmpeg.org/ffmpeg-utils.html#Quoting-and-escaping
        fileList += "file '" + f.filename.replace("'", "'\\''") + "'\n";
        length += f.time;
    }
    
    // run ffmpeg command that concatenates all the sound files (1/2)
    tmp.file({postfix: '.txt'}, function(err, path1, fd, cleanupCallback) {
        if (err) throw err;
        stage1.inputs = [{
            url: path1,
            options: "-f concat -safe 0"
        }];
        fs.writeFile(path1, fileList)
        .then(_ => tmp.tmpName({postfix: '.wav'}, function(err, path2) {
            let path3 = path.join(path.dirname(data.form.outputPath), path.basename(path2))
            console.log(path3);
            if (err) throw err;
            stage1.outputs = [{
                url: path2,
                options: "-c copy"          
            }];
            stage1Proc = new Promise(function(resolve, reject) {
                let proc = stage1.run();
                proc.on('close', function(code) {
                    resolve(code);
                });
                proc.on('error', function(err) {
                    reject(err);
                });
            });
        }));
        cleanupCallback();
    });
    await stage1Proc;
    return {
        concatWavPath: concatWavPath,
        timestamps: timestamps,
        totalLength: length
    };
}
