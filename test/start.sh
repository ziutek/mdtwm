#!/bin/bash

display=:3
xauthFile=/tmp/xephyr.auth
xephyrCmd=$(which Xephyr 2>/dev/null)

[ "${xephyrCmd:-undef}" == "undef" ] && {
    echo "you need to install Xephyr to run this script!"
    exit 1 ;
}

echo "add $display . $(mcookie)" | xauth -f $xauthFile
Xephyr $display -auth /tmp/xeph.auth -screen 900x500 &
export XAUTHORITY=$xauthFile
export DISPLAY=$display
sleep 3

gnome-terminal &

#xsetroot  -solid gray -cursor_name left_ptr

sleep 2

./mdtwm

echo "XAUTHORITY=$xauthFile DISPLAY=$display"
