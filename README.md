This program isn't yet ready for use.

I like tiling window managers but I like full mouse control to. So this is my
attempt to write some mouse-driven tiling wm in Go.

I have never written any wm before, so this code shouldn't be treated as example
how to write wm at all.

mdtwm uses only right button as follows:

- single-click-move: move a window,
- single-slow-click: send right click to window (not implemented yet),
- double-click-move: tile/untile window (not implemented yet),
- triple-click: close window.
