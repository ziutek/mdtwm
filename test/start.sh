#!/bin/bash

display=:3
xauthFile=/tmp/xephyr.auth
xephyrCmd=$(which Xephyr 2>/dev/null)
termEmu="gnome-terminal"

while getopts ":t:" opt
do
    case $opt in
        t) termEmu=$OPTARG;;
    esac
done

termEmuCmd=$(which $termEmu 2>/dev/null)
[ "${termEmuCmd:-undef}" == "undef" ] && {
    echo "can't find $termEmu, specify a terminal emulator with -t"
    exit 1;
}

[ "${xephyrCmd:-undef}" == "undef" ] && {
    echo "you need to install Xephyr to run this script!"
    exit 1 ;
}

echo "add $display . $(mcookie)" | xauth -f $xauthFile
Xephyr $display -auth /tmp/xeph.auth -screen 900x500 &
export XAUTHORITY=$xauthFile
export DISPLAY=$display
sleep 3

$terEmuCmd &

#xsetroot  -solid gray -cursor_name left_ptr

sleep 2

./mdtwm

echo "XAUTHORITY=$xauthFile DISPLAY=$display"
