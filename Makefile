include $(GOROOT)/src/Make.inc

TARG=aurgo
GOFILES=\
	main.go\
	aur.go

include $(GOROOT)/src/Make.cmd

