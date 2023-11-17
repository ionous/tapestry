module git.sr.ht/~ionous/tapestry

go 1.21

toolchain go1.21.4

require (
	github.com/ionous/errutil v0.0.0-20230227195626-6de478250a3b
	github.com/ionous/inflect v0.0.0-20211113032332-adfb17b87f92
	github.com/ionous/num2words v0.0.0-20210224003458-c9a432ced842
	github.com/ionous/sliceOf v0.0.0-20170627065049-c4e57a86cb77
	github.com/ionous/tell v0.0.0-20231117054346-43acc5600dbf
	github.com/kr/pretty v0.3.1
	github.com/mattn/go-sqlite3 v1.14.16
	github.com/reiver/go-porterstemmer v1.0.1
)

require (
	github.com/kr/text v0.2.0 // indirect
	github.com/rogpeppe/go-internal v1.9.0 // indirect
)

// so that it will use the local source tree, and not the currently released version.
replace github.com/ionous/tell => /Users/ionous/Dev/GitHub/tell
