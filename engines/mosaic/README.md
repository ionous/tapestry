# README

## About

This is the tapestry GUI. It uses [Wails](https://wails.io/) to provide a native webbrowser in a window managed by go.

## Details

The app supports three modes: `ext`, `dev`, and `production`. The first two rely on [Vite](https://vitejs.dev/) to serve assets.

Ext
: Use your local web browser to [interact](http://localhost:8080/mosaic/) with the app. For simplicity's sake, while built into the application, Wails isn't used. The app can be built and run with no special tags.

Dev
: Uses wails' "embedded" browser to interact with the app. The app must be built with the "dev" tag. ( ie. `go build -tags dev` ) Wails normally has its own web server in this mode, but it's disabled so everything can route through the same code used by "web".

Production
: Uses wail's browser to interact with pre-built content. all assets get built into the go app, except for story files which are read/written to by the app. The "production" and "desktop" tags are required, as are several linker options. ( -w -s -H windowsgui ). Wails would normally serve the embedded assets, however Tapestry routes everything through the same multiplexer used for web and dev returning the embedded assets manually.

Tapestry doesn't currently expose a javascript api for the wails window.

### running web and dev

after building the app, in the `www` directory:

> npm run dev

### building production

in the www directory:

> npm run build

in the cmd/tap directory, either:

> go build -tags desktop,production -ldflags "-w -s -H windowsgui"

or

> wails build -s -noPackage

( ^ requires the wails build system to be installed: 
ex. go install github.com/wailsapp/wails/v2/cmd/wails@v2.4.0 )


## Compiler options

build tags
: control which version of the source code gets built:

* "dev": activates the wails browser and launches its dev server.
* "debug": activates the console for production builds.
* "bindings": appears to create an app capable of building the go <-> js bindings.
* "production": switches off the wails and tapestry servers; embeds wails and tapestry assets.
* "desktop": embeds a prebuilt, minified version of the [wails runtime](https://github.com/wailsapp/wails/blob/master/v2/internal/frontend/runtime/runtime_prod_desktop.js) for logging, events, native window controls, etc.

gcflags
: controls the go compiler. the specified options make debugging with delve ( ex. in visual studio code ) easier:

* all: apply the options to all packages.
* -l: no inlined functions.
* -N: disable optimizations.

ldflags
: controls the linker.

* -w: Omit the DWARF symbol table.
* -s: Omit the symbol table and debug information.
* -H windowsgui writes a "GUI binary" ( instead of a 'console binary' )

# Overview

```
================================== DEV MODE ===============================================

USER +-> [wails webkit](http://wails.localhost)
               +
               |
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ WEB & DEV MODES ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
               +
               |
USER +-> [tapestry server](http://localhost:8080)
               |
             <mux> +--> (unknown url?) +--> npm+vite:3000 +--------->+
               |                                                     |
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ ALL MODES  ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
               +                                                     +
               |                                                     |
     { tapestry backend }                                  { tapestry frontend  }
    ............................................    ....................................
    .  /stories/ ----> [ Documents: .if/db ]   .    .   /www/index.html                .
    .  /ramble/  ----> take turns              .    .   /www/mosaic/.vue, .js, .css    .
    .  /blocks/  ----> blocks from .if file    .    .   /www/ramble/.vue, .js, .css    .
    .  /shapes/  ----> blockly definitions     .    .   /www/assets/.png, .etc         .
    ............................................    ....................................
               |                                                     |
               +                                                     +
================================== PRODUCTION MODE ========================================
               |                                                     |
             <mux> +-----> (unknown url?) +----->        package tapestry/www
               |                                    embedded assets built by vite cli
               +                                          served by http.FileServer
          AssetsHandler
               |
               +
            Assets +------>❌ wails treats all requests as files, then falls back to AssetsHandler if unhandled.
               |                Tapestry always returns "file not found" to delegate everything to the handler.
               +
USER +->  [wails webkit](http://wails.localhost)

==========================================================================================
```