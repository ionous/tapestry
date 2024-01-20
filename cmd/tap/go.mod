module git.sr.ht/~ionous/tapestry/cmd/tap

go 1.21

toolchain go1.21.4

require (
	git.sr.ht/~ionous/tapestry v0.0.0-20220223221837-b25bd82d20fe
	github.com/ionous/errutil v0.0.0-20231013205411-87ce252b8e2a
	github.com/wailsapp/wails/v2 v2.4.0
	golang.org/x/sys v0.6.0
)

require (
	github.com/bep/debounce v1.2.1 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/ionous/inflect v0.0.0-20211113032332-adfb17b87f92 // indirect
	github.com/ionous/num2words v0.0.0-20210224003458-c9a432ced842 // indirect
	github.com/ionous/tell v0.8.0 // indirect
	github.com/jchv/go-winloader v0.0.0-20210711035445-715c2860da7e // indirect
	github.com/labstack/echo/v4 v4.10.2 // indirect
	github.com/labstack/gommon v0.4.0 // indirect
	github.com/leaanthony/go-ansi-parser v1.6.0 // indirect
	github.com/leaanthony/gosod v1.0.3 // indirect
	github.com/leaanthony/slicer v1.6.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/mattn/go-sqlite3 v1.14.16 // indirect
	github.com/pkg/browser v0.0.0-20210911075715-681adbf594b8 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rivo/uniseg v0.4.4 // indirect
	github.com/samber/lo v1.37.0 // indirect
	github.com/tkrajina/go-reflector v0.5.6 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	github.com/wailsapp/mimetype v1.4.1 // indirect
	golang.org/x/crypto v0.7.0 // indirect
	golang.org/x/exp v0.0.0-20230810033253-352e893a4cad // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/text v0.8.0 // indirect
)

// fix: not sure how to get rid of this:
replace git.sr.ht/~ionous/tapestry => ../../../tapestry

// for local debugging:
// replace github.com/wailsapp/wails/v2 v2.0.0-beta.42 => C:\Dev\Go\pkg\mod\github.com\wailsapp\wails\v2@v2.0.0-beta.42
