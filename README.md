## album2video

gui application that makes a folder of mp3 files + album art into a youtube uploadable format

### building fr development

this is pretty half-baked idk why anyone would want to know how as of right now but i need to write this in so i don't forget how to do this in the future. this has no dependencies if you're building on yr target platform because it uses this magic node module that downloads ffmpeg for you. will update w instructions for cross-compilation later (again also for myself lol) 

you will need to have a reasonably new version of nodejs / npm and golang 1.16
```bash
git clone https://github.com/sunglasseds/album2video.git album2video
cd album2video
mkdir bin
# build the scary evil go binary that does all the work
go build -o ./bin/xX_FFMP3G_BL4CKB0X_Xx ./src/xX_FFMP3G_BL4CKB0X_Xx/*.go
npm install
```
`npm start` 2 run

### track detection
the regex is currently as follows:
```regex
^([0-9]+|[A-Za-z]|[0-9]+[A-Za-z]|)(-| - |_| |)([0-9]+|[A-Za-z])(?=. | |_)
```
slap it into a site like https://regexr.com/ and type in track names to see if yours work. they probably will but if they dont submit a pr or otherwise let me know and i'll try to fix it

---

# SORRY FOR MAKING AN ELECTRON APP !!!!!!!!!!!!!!!
