---
weight: 1 # smaller is higher
---

The Tapestry tool -- called `tap` -- is a one-stop shop turning story files into playable games. The story files themselves can be written in the editor of your choice. ( I tend to use Sublime[https://www.sublimetext.com/] )

# Installing Tapestry 

At this time, there are no pre-built releases. So to run Tapestry, you will need to build Tapestry. 

1. **Install Go**. To install Go, visit the Go website and follow its [installation instructions](https://go.dev/doc/install).

1. **Install Tapestry**. After installing Go, you can either 1) use Go to install Tapestry, or 2) get the Tapestry source code and build it yourself.

    1. **Use Go to install Tapestry.** Open a terminal or command line prompt and type `go install git.sr.ht/~ionous/tapestry/cmd/tap@latest`. Assuming it completes successfully, you can then run the tap tool.
    
    1. **Build Tapestry from source**:

        1. Use git to clone `https://git.sr.ht/~ionous/tapestry` to a local directory;
        2. In that directory, `cd cmd/tap`, and either run `go install` to install the tap command globally; or,run `go build` to create the `tap` executable in that directory.


3. **Verify Installation**. At the command line, run `tap help`. If everything worked successfully, you should see a help message describing the available commands.


{{< hint type=caution >}}
On Windows you'll need gcc ( Tapestry uses [go-sqlite3](https://github.com/mattn/go-sqlite3), and it requires cgo. ) [TDM-GCC](https://jmeubank.github.io/tdm-gcc/) is a good option. Follow their [installation instructions](https://jmeubank.github.io/tdm-gcc/download/) and then install, or clone and build Tapestry.
{{< /hint >}}
