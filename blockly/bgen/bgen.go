package bgen

// import (
// 	"encoding/json"
// 	"log"

// 	"git.sr.ht/~ionous/tapestry"
// 	"git.sr.ht/~ionous/tapestry/dl/story"
// 	"git.sr.ht/~ionous/tapestry/jsn"
// 	"git.sr.ht/~ionous/tapestry/jsn/chart"
// 	"git.sr.ht/~ionous/tapestry/web/files"
// 	"git.sr.ht/~ionous/tapestry/web/js"
// 	"github.com/ionous/errutil"
// )

// const (
// 	CompactExt  = ".if"
// )

// func (out *Output) importStory(tgt jsn.Marshalee) error {
// 	enc := chart.MakeEncoder()
// 	return enc.Marshal(tgt, story.Map(&enc, story.BlockMap{
// 		story.OtherBlocks: story.KeyMap{
// 			story.BlockStart: func(b jsn.Block, _ interface{}) (err error) {
// 				// if flow, ok := b.(jsn.FlowBlock); ok {
// 				// }
// 				return
// 			},
// 		},
// 	}))
// }

// // read a comma-separated list of files and directories
// func importStoryFiles(k *Output, srcPath string) (err error) {
// 	if e := files.ReadPaths(srcPath,
// 		[]string{CompactExt, DetailedExt}, func(p string) error {
// 			return readOne(k, p)
// 		}); e != nil {
// 		err = errutil.New("couldn't read file", srcPath, e)
// 	} else {
// 		k.Flush()
// 	}
// 	return
// }

// func readOne(k *story.Importer, path string) (err error) {
// 	log.Println("reading", path)
// 	if b, e := files.ReadFile(path); e != nil {
// 		err = e
// 	} else if script, e := decodeStory(path, b); e != nil {
// 		err = errutil.New("couldn't decode", path, "b/c", e)
// 	} else if e := k.ImportStory(path, script); e != nil {
// 		err = errutil.New("couldn't import", path, "b/c", e)
// 	}
// 	return
// }

// func decodeStory(path string, b []byte) (ret story.StoryLines, err error) {
// 	var curr story.Story
// 	if e := story.Decode(&curr, b, tapestry.AllSignatures); e != nil {
// 		err = e
// 	} else {
// 		ret = curr.reformat()
// 	}
// 	return
// }
