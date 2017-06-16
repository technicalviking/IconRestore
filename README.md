# README #

Are you a windows user frustrated with the fact that sometimes a random video driver update will randomly scatter your desktop icons to the four winds?  Do you use a laptop that seems incapable of maintaining your carefully arranged desktop icons when you plug in a second monitor?

Maybe I can help!

IconRestore (because I'm the best at naming things!) was built to somewhat solve that problem for me. Yeah, there's other software out there that does this.  I hear people say nice things about something called [DesktopOK](http://bfy.tw/CO9K).  I mostly wanted to figure out how something like that software might work and learn GoLang at the same time.

### What is this repository for? ###

* Providing the source code, and a basic exe for saving and restoring

### How do I get set up? ###

* I just wanna run it!
Just download the executable then!

* * Saving
* * * Put your desktop icons where you want em.
* * * Refresh the desktop a few times (F5)
* * * Open a command prompt
* * * type in cd path/to/IconRestore.exe
* * * type in IconRestore.exe save
* * * if you see a file called regKeyData in the same folder as IconRestore.exe, you're done!

* * Restoring
* * * Open a command prompt
* * * type in cd path/to/IconRestore.exe
* * * type in IconRestore.exe
* * * If you noticed the screen flicker, or the taskbar went away for a second, that's a sign that it worked!  Your Icons SHOULD be in the position they were saved in!

* I wanna build from source!
* * if you're computer isn't set up to develop golang applications, [go get that done](https://golang.org/doc/install)
* * download main.go from this repo into some source folder
* * run "go build" from a command prompt navigated to the src/folder wherre you put main.go
* * fromm here the instructions are the same as if you'd downloaded the exe!


* Dependencies/troubleshooting
* * Only Tested on Windows 10
* * for laptop users: windows may actually store multiple icon configurations, which could mean that the save file created when you weren't plugged into a monitor won't actually fix the icons when you are. (also vice versa).  So you may have to run 'save' from multiple configurations.  I may actually code an improvement to fix that in the future.  Sorry.

### Contribution guidelines ###

* If you feel like this could be improved, you probably won't get any argument from me.  Feel free to download this code and modify it as you wish.  Alternately fork this repo or even contribute to it!

### Who do I talk to? ###