#!/bin/bash

display=:3
xauthFile=/tmp/xephyr.auth

echo "add $display . $(mcookie)" | xauth -f $xauthFile
Xephyr $display -auth /tmp/xeph.auth -screen 900x500 &
export XAUTHORITY=$xauthFile
export DISPLAY=$display
sleep 4

gnome-terminal &

xsetroot  -solid gray -cursor_name left_ptr

sleep 1

./mdtwm

echo "XAUTHORITY=$xauthFile DISPLAY=$display"
