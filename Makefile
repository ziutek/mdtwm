include $(GOROOT)/src/Make.inc

TARG=mdtwm
GOFILES=\
	main.go\
	windows.go\
	config.go \
	event_handlers.go\
	manage.go\

include $(GOROOT)/src/Make.cmd
