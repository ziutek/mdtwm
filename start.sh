#!/bin/bash

display=:3
xauthFile=/tmp/xephyr.auth

echo "add $display . $(mcookie)" | xauth -f $xauthFile
Xephyr $display -auth /tmp/xeph.auth -screen 640x400 &
export XAUTHORITY=$xauthFile
export DISPLAY=$display
sleep 2

xterm &
./mdtwm
