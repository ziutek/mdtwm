#!/bin/bash

display=:3
xauthFile=/tmp/xephyr.auth

echo "add $display . $(mcookie)" | xauth -f $xauthFile
Xephyr $display -auth /tmp/xeph.auth -screen 1000x400 &
export XAUTHORITY=$xauthFile
export DISPLAY=$display
sleep 2

xsetroot -solid white
xterm &
xterm &
xterm &

sleep 2

./mdtwm

echo "XAUTHORITY=$xauthFile"
echo "DISPLAY=$display"
