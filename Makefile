include $(GOROOT)/src/Make.inc

GC = $Og -N

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
	tile.go\

include $(GOROOT)/src/Make.cmd
