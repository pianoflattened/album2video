## album2video
converts directory of mp3s -> video + tracklist
i have only used this on linux but i tried to make it cross-platform. sorry if it doesnt work

### building
have git and go installed
```
git clone https://github.com/pianoflattened/album2video
cd album2video
mkdir bin
go build -o bin
```

### using
just run th executable that th command spits out (it will be in bin and called album2video). it has a help menu

## THINGS I NEED TO ADD/FIX
\- substantial slowdown if the image is large - resize image to be 720px tall if it is bigger + option to disable <br>
\- if the image dimensions arent both even numbers it stops (yuv480p thing cant help it). going to crop out a row/col of pixels of the image in memory to fix <br>
\- guessing artist/title from filenames

### track detection from filenames (for when ur stuffs not tagged)
the regex is currently as follows:
```regex
^([0-9]+|[A-Za-z]{1,2}|[0-9]+[A-Za-z]|)([-_ ]| - |)([0-9]+|[A-Za-z])[ _.]
```
slap it into a site like https://regexr.com/ and type in track names to see if yours work. they probably will but if they dont let me know and i'll try to fix it

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
