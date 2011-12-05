Xephyr -ac -br -noreset -screen 800x600 :1 &
export DISPLAY=:1
sleep 1
./mdtwm
