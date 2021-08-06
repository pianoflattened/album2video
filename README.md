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
^([0-9]+|[A-Za-z]{1,2}|[0-9]+[A-Za-z]|)(-| - |_| |)([0-9]+|[A-Za-z])(. | |_)
```
slap it into a site like https://regexr.com/ and type in track names to see if yours work. they probably will but if they dont submit a pr or otherwise let me know and i'll try to fix it

### tracklist formatting
reference for myself dw about it yet although if youve ever used printf in your life most of this should make sense
```%t title
%s timestamp
%r artist (indiscriminate)
%a artist (discriminate)
%d disc
%n track number (overall)
%w track number (within disc)
%% percent
%{ left brace
%} right brace

%{ }v only include characters inside {} if %v exists 
	v is an example value. %v is rendered at the %. i do not plan on making nesting work unless
	someone somehow comes up w a practical use case - this means that the rules listed here do not 
	apply to inside curly braces
c lowercase
C title case
(number) pad zeroes
	cannot be less than 3 since thats the minimum for a valid yt timestamp (0:00)
---

# SORRY FOR MAKING AN ELECTRON APP !!!!!!!!!!!!!!!
