package cmdmosaic

import (
	"flag"

	"git.sr.ht/~ionous/tapestry/cmd/tap/internal/mosaic"
)

// collection of local flags
var mosaicFlags = struct {
	folder          mosaic.Folder
	listen, request mosaic.Port
	width, height   int
}{}

func buildFlags() (flags flag.FlagSet) {
	flags.Var(&mosaicFlags.folder, "in", "directory for if files.")
	flags.IntVar(&mosaicFlags.width, "width", 1024, "width of the application window.")
	flags.IntVar(&mosaicFlags.height, "height", 768, "height of the application window.")
	if mosaic.BuildConfig != mosaic.Prod {
		flags.Var(&mosaicFlags.listen, "listen", "end-user's port number for mosaic server. specify a port number; or, 'true' for the default (8080).")
		if mosaic.BuildConfig == mosaic.Web {
			flags.Var(&mosaicFlags.request, "www", "back-end port number for tapestry's webapps. specify a port number; or, 'true' to use the default port (3000).")
		}
	}
	return
}
