I like tiling window management, but all tiling window managers known to me force to use keyboard for dozen of things. So this is my attempt to write some mouse-driven tiling wm in Go.

I have never written any wm before, so this code isn't probably a good example how to write wm at all but maybe some example how to write wm in Go.

This program isn't mature window manager yet, but I started using it for everyday work (it replaced xmodad).

mdtwm uses only right mouse button as follows:

- single-click-move: move a window (you can move between desktops),
- single-slow-click: send right click to window,
- double-click-move: tile/untile window,
- triple-click: close window.

Default layout (see config.go)

- two desks
- first desk contains two panels
- second desk contains one horizontal panel
- third desk contains one vertical panel

Default keybindings (see config.go)

- Mod+Enter: new terminal (gnome-terminal)
- Mod+q: exit
- Mod+1: first desk
- Mod+2: second desk
- Mod+3: thrid desk

mdtwm doesn't contain its own status bar yet, but there is support for dzen2. If cfg.StatusLogger is set to Dzen2Logger (as in default config.go) you can use mdtwm together with dzen2 as follows:

    mdtwm 2>~/mdtwm.log |dzen2 -e '' -ta l -fg '#ddddcc' -bg '#555588'

Known issues:

- You can't change current desktops without use of keyboard.
- See *issues* directory.

Build instruction:

    go install github.com/ziutek/mdtwm/cmd/mdtwm

    The executable will be installed in $GOPATH/bin/mdtwm

Configuration:
    Edit the file $GOPATH/src/pkg/github.com/ziutek/mdtwm/config.go
    (if you do not have $GOPATH configured, try $GOROOT).

    then run:
    go install github.com/ziutek/mdtwm/cmd/mdtwm
    
Screenshot:

![Screenshot image](/ziutek/mdtwm/raw/master/screenshot.jpg)
