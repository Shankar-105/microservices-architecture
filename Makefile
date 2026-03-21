ifeq ($(OS),Windows_NT)
include Makefile.windows
else
include Makefile.linux
endif
