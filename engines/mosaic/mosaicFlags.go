package main

import (
	"flag"

	"git.sr.ht/~ionous/tapestry/web"
)

// collection of local command line flags (filled when the mosaic command is selected)
var mosaicFlags = struct {
	folder          Folder
	listen, request web.Port
	width, height   int
}{}

// creates a description which writes into the cfg when the mosaic command is matched
func buildFlags() (out flag.FlagSet) {
	out.Var(&mosaicFlags.folder, "in", "directory for if files.")
	out.IntVar(&mosaicFlags.width, "width", 1024, "width of the application window.")
	out.IntVar(&mosaicFlags.height, "height", 768, "height of the application window.")
	if BuildConfig != Prod {
		out.Var(&mosaicFlags.listen, "listen", "the port for your web browser. specify a port number; or, 'true' for the default (8080).")
		if BuildConfig == Web {
			out.Var(&mosaicFlags.request, "www", "local vite server where tapestry can find its webapps. specify a port number; or, 'true' to use the default port (3000).")
		}
	}
	return
}
