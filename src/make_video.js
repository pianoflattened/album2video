const ffmpegPath = require('@ffmpeg-installer/ffmpeg').path;
const { FFmpeg } = require("kiss-ffmpeg");

FFmpeg.command = ffmpegPath;
let ffmpeg = new FFmpeg();

module.exports = async function makeVideo(data) {
    let paths = [];
    for (const f of data.audioFiles) {
        paths.push(f.filename)
    }
    console.log(paths);
    // run ffmpeg command that concatenates all the sound files (1/2)
    ffmpeg.inputs = paths;

    // run ffmpeg command that creates video with concat.wav and cover (2/2)
};
