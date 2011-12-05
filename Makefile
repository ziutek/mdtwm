include $(GOROOT)/src/Make.inc

TARG=mdtwm
GOFILES=\
	main.go\
	windows.go\
	config.go \

include $(GOROOT)/src/Make.cmd
