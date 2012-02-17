#!/bin/bash

# Don't run it becouse of changes in xproto.go
exit

if [ -f go_client.py ]; then
	git clone git://anongit.freedesktop.org/git/xcb/proto
	python go_client.py -p proto/ proto/src/xproto.xml
	gofmt -w xproto.go
fi	
