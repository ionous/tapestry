module git.sr.ht/~ionous/tapestry/cmd/tap

go 1.18

require (
	git.sr.ht/~ionous/tapestry v0.0.0-00010101000000-000000000000
	github.com/ionous/errutil v0.0.0-20210301225645-f05ff1857722
	github.com/wailsapp/wails/v2 v2.0.0-beta.42
	golang.org/x/sys v0.0.0-20220114195835-da31bd327af9
)

require (
	github.com/bep/debounce v1.2.1 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/google/uuid v1.1.2 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/ionous/inflect v0.0.0-20211113032332-adfb17b87f92 // indirect
	github.com/ionous/num2words v0.0.0-20210224003458-c9a432ced842 // indirect
	github.com/jchv/go-winloader v0.0.0-20210711035445-715c2860da7e // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/labstack/echo/v4 v4.7.2 // indirect
	github.com/labstack/gommon v0.3.1 // indirect
	github.com/leaanthony/go-ansi-parser v1.0.1 // indirect
	github.com/leaanthony/gosod v1.0.3 // indirect
	github.com/leaanthony/slicer v1.5.0 // indirect
	github.com/mattn/go-colorable v0.1.11 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mattn/go-sqlite3 v1.14.6 // indirect
	github.com/pkg/browser v0.0.0-20210706143420-7d21f8c997e2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rogpeppe/go-internal v1.6.1 // indirect
	github.com/tkrajina/go-reflector v0.5.5 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.1 // indirect
	github.com/wailsapp/mimetype v1.4.1 // indirect
	golang.org/x/crypto v0.0.0-20210817164053-32db794688a5 // indirect
	golang.org/x/net v0.0.0-20211015210444-4f30a5c0130f // indirect
	golang.org/x/text v0.3.7 // indirect
)

// so that it will use the local source tree, and not the currently released version.
replace git.sr.ht/~ionous/tapestry => ../../../tapestry
