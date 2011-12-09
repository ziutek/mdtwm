Xephyr -ac -br -noreset -screen 640x400 :1 &
export DISPLAY=:1
sleep 1
./mdtwm
