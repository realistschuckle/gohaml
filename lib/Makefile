# Copyright 2009 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

CLEANFILES += lang.go

include $(GOROOT)/src/Make.inc

TARG=gohaml
GOFILES=\
	haml.go\
	parser.go\
	tree.go\
	lang.go\

include $(GOROOT)/src/Make.pkg

lang.go: lang.y
	$(QUOTED_GOBIN)/goll1e lang.y lang.go
