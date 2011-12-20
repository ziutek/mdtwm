include $(GOROOT)/src/Make.inc

GC = $Og -N
ALL = mdtwm test

all : $(ALL)

mdtwm :	\
	main.go\
	globals.go\
	geometry.go\
	window.go\
	box.go\
	box_list.go\
	root_panel.go\
	panel.go\
	boxed_window.go \
	config.go \
	manage.go\
	atoms.go\
	utils.go\
	events.go\
	input_events.go\

	$(GC) -o $@.$O $^
	$(LD) -o $@ $@.$O

test : \
	test.go\
	window.go\
	geometry.go\
	atoms.go\

	$(GC) -o $@.$O $^
	$(LD) -o $@ $@.$O

clean:
	rm -rf *.[68] $(ALL)

.PHONY : clean
