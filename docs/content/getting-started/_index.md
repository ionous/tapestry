---
weight: 1 # smaller is higher
---

The Tapestry tool -- called `tap` -- is a one-stop shop turning story files into playable games. The stories themselves can be written in any editor of your choice. ( I tend to use [Sublime](https://www.sublimetext.com/). )

# Installing Tapestry 

At this time, there are no pre-built releases. To run Tapestry, you will need to build Tapestry. This probably requires some coding expertise.

1. **Install Go**. To install Go, visit the Go website and follow its [installation instructions](https://go.dev/doc/install).

2. **Download the source code**: Use git to clone `https://git.sr.ht/~ionous/tapestry` to a local directory. Or, you can either fork the [{{< icon "gdoc_github" >}} GitHub mirror](https://github.com/ionous/tapestry), or download a zip of the source from that GitHub link.

3. **Build the source**: In that directory, `cd cmd/tap`, and either run `go install` to install the tap command globally; or use `go build` to create the `tap` executable in that directory.

4. **Verify Installation**. At the command line, run `tap help`. If everything worked successfully, you should see a help message describing the available commands.

{{< hint type=caution >}}
On Windows: Because Tapestry uses [go-sqlite3](https://github.com/mattn/go-sqlite3) ( which uses cgo ) you will need gcc. [TDM-GCC](https://jmeubank.github.io/tdm-gcc/) is a good option. Follow their [installation instructions](https://jmeubank.github.io/tdm-gcc/download/), and then Tapestry should build successfully.
{{< /hint >}}

## Alternative Installation

This method doesn't get tested much, so it may be unreliable. After installing Go, instead of manually downloading and building the source, at the command line enter: `go install git.sr.ht/~ionous/tapestry/cmd/tap@latest`. That will get the source and compile it for you. Assuming it completes successfully, you can run the tap tool to verify the installation.
    
