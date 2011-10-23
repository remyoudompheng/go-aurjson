include $(GOROOT)/src/Make.inc

DEPS=aurjson

TARG=aurgo
GOFILES=\
	main.go\

include $(GOROOT)/src/Make.cmd

