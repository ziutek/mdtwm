include $(GOROOT)/src/Make.inc

TARG=mdtwm
GOFILES=\
	main.go\
	window.go\
	config.go \
	event_handlers.go\
	manage.go\
	atoms.go\
	utils.go\
	box.go\

include $(GOROOT)/src/Make.cmd
