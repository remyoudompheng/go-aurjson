include $(GOROOT)/src/Make.inc

DEPS=aurjson

TARG=aurgo
GOFILES=\
	main.go\

include $(GOROOT)/src/Make.cmd

clean: recursive-clean

recursive-clean:
	$(MAKE) -C aurjson clean
