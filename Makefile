include $(GOROOT)/src/Make.inc

GC = $Og -N

TARG=mdtwm
GOFILES=\
	main.go\
	geometry.go\
	window.go\
	box.go\
	window_box.go\
	panel_box.go\
	box_list.go \
	config.go \
	event_handlers.go\
	manage.go\
	atoms.go\
	utils.go\

include $(GOROOT)/src/Make.cmd
