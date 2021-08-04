## album2video

gui application that makes a folder of mp3 files + album art into a youtube uploadable format

### building fr development

this is pretty half-baked idk why anyone would want to know how as of right now but i need to write this in so i don't forget how to do this in the future. this has no dependencies if you're building on yr target platform because it uses this magic node module that downloads ffmpeg for you. will update w instructions for cross-compilation later (again also for myself lol) 
```bash
git clone https://github.com/sunglasseds/album2video.git album2video
cd album2video
mkdir bin
# build the scary evil go binary that does all the work
go build -o ./bin/xX_FFMP3G_BL4CKB0X_Xx ./src/*.go
npm install
```
`npm start` 2 run

# SORRY FOR MAKING AN ELECTRON APP !!!!!!!!!!!
