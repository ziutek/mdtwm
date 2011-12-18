include $(GOROOT)/src/Make.inc

GC = $Og -N

TARG=mdtwm
GOFILES=\
	main.go\
	geometry.go\
	window.go\
	box.go\
	box_list.go\
	root_panel.go\
	panel.go\
	boxed_window.go \
	config.go \
	event_handlers.go\
	manage.go\
	atoms.go\
	utils.go\

include $(GOROOT)/src/Make.cmd
