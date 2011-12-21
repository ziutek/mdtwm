I like tiling window managers but I like full mouse control to. So this is my
attempt to write some mouse-driven tiling wm in Go.

I have never written any wm before, so this code shouldn't be treated as example
how to write wm at all.

This program isn't mature window manager yet, but I started using it for
everyday work (it replaced xmodad).

mdtwm uses only right button as follows:

- single-click-move: move a window,
- single-slow-click: send right click to window,
- double-click-move: tile/untile window (not implemented yet),
- triple-click: close window.

Default layout (see config.go)

- two desks
- first desk contains two panels
- second desk contains one (fullscreen) panel

Default keybindings (see config.go)

- Mod+Enter: new terminal (gnome-terminal)
- Mod+q: exit
- Mod+1: first desk
- Mod+2: second desk

mdtwm doesn't contain its own status bar yet, but there is support for dzen2. If
cfg.StatusLogger is set to Dzen2Logger (as in default config.go) you can use
mdtwm together with dzen2 as follows:

    mdtwm 2>~/mdtwm.log |dzen2 -e '' -ta l -fg '#ddddcc' -bg '#555588'
