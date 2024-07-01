package mosaic

import (
	"git.sr.ht/~ionous/tapestry/web"
)

func FilesApi(cfg *Config) web.Resource {
	return &web.Wrapper{
		Finds: func(name string) (ret web.Resource) {
			// by adding a trailing slash, walk will follow a symlink.
			path := cfg.Stories()
			switch name {
			case "blocks":
				// generates blockly files out of story files.
				where := storyFolder{cfg, path}
				ret = blocksRoot{blocksFolder{where}}
			case "stories":
				// serves raw story files in json format.
				where := storyFolder{cfg, path}
				ret = rootFolder{where}
			}
			return
		},
	}
}

// note: the root folder used to handle "put" to receive detailed format stories from the inline composer;
// removed in favor of the mosaic blockly editor which puts individual (compact) files to the file endpoints.
type rootFolder struct {
	storyFolder
}
