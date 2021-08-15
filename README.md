## album2video

gui application that makes a folder of sound files (any format that ffmpeg likes) + album art into a youtube-uploadable format

## THIGNS I NEED TO FIX
- if the image dimensions arent both even numbers then it stops (yuv480p quirk sorry abt this youtube requires this for uploading lol)
- substantial slowdown if image is large - resize image to be 720px tall if it is taller
	- add an option to disable this

### setting ffmpeg/ffprobe paths
JUST KIDDING i saw you grimace it was pretty funny makes you look like a cartoon character when you do that. i use a magic node library that automatically downloads binaries depending on ur os (linked below) and i built it with those. if you want to use a different version / have some modified version of ffmpeg bc youre a sociopath then build from source i guess

### building from source
lol

### track detection from filenames (for when ur stuffs not tagged)
the regex is currently as follows:
```regex
^([0-9]+|[A-Za-z]{1,2}|[0-9]+[A-Za-z]|)([-_ ]| - |)([0-9]+|[A-Za-z])[ _.]
```
slap it into a site like https://regexr.com/ and type in track names to see if yours work. they probably will but if they dont submit a pr or otherwise let me know and i'll try to fix it

### tracklist formatting
if youve ever used printf in your life most of this should make sense
```%s song
%t timestamp
%r artist (indiscriminate)
%a artist (discriminate)
%d disc
%n track number (overall)
%w track number (within disc)
%% percent

%< ]x only include characters inside < ] or [ > if %x exists
%[ >x
	x is an example value. %x is rendered on the side that the arrow is on. i do not plan on
	making nesting work unless someone somehow comes up w a practical use case - this means that
	the rules listed here do not apply to inside these braces. \ is used to escape inside (if you
	want to use a > or ] put a \ before it and the regex will ignore it)

^ uppercase
- title case
v lowercase
(number) pad zeroes
```
example (the format i use): `%v[ - >a%vs - %3vt`

if nothing is entered the script will use `%[ - >a%s - %t`

you can write your own formatter and replace the binary if you want. the blackbox has some ipc stuff going on that makes it a bit less practical to immediately replace but the formatter is just straight input -> output stuff v straightforward

### encoding problems on windows (┬░ displaying instead of °, etc)
this superuser answer explains it:
https://superuser.com/questions/1584842/ffprobe-output-text-wrong-encoding/1588628#1588628

> Goto Control Panel\Clock and Region <br>
> Click Change date, time, or number format <br>
> In Region window, click tab Administrative and click Change system locale... <br>
> Check the checkbox Beta: Use Unicode UTF-8 <br>
> Click Ok and restart computer. <br>

i had to search "region" 2 find the setting

---
# ***SORRY FOR MAKING AN ELECTRON APP !!!!!!!!!!!!!!!***
