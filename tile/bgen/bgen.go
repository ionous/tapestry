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
// 	DetailedExt = ".ifx"
// 	CompactExt  = ".if"
// )

// type Output struct {
// }

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

// // starting a new flow; which represents a fresh block that has no influence on parent blocks
// // "blocks" might be the topLevel array of blocks, or the value of a block or shadow key.
// func newBlock(blocks *js.Builder, typeName string) jsn.State {
// 	var fields, inputs, next Js
// 	var nextKey string
// 	return &StateMix{
// 		// one of every extant member of the flow ( skipping optional elements lacking a value )
// 		OnKey: func(key string, _ string) error {
// 			nextKey = key
// 			// this might be a field or input
// 			// we might write to next when the block is *followed* by another in a repeat.
// 			// therefore we cant close the block in Commit --
// 			// but we might close child blocks
// 		},
// 		// an embedded flow
// 		OnMap: func(string, jsn.FlowBlock) (okay bool) { return },
// 		// a value that fills a slot
// 		OnSlot: func(string, jsn.SlotBlock) (okay bool) { return },
// 		// a member that's a swap
// 		OnSwap: func(string, jsn.SwapBlock) (okay bool) { return },
// 		// a member that repeats
// 		OnRepeat: func(string, jsn.SliceBlock) (okay bool) { return },
// 		// a single value
// 		OnValue: func(_ string, pv interface{}) (err error) {
// 			if b, e := json.Marshal(pv); e != nil {
// 				err = e
// 			} else {
// 				fields.Q(nextKey).R(js.Colon).Write(b)
// 			}
// 			return
// 		},

// 		OnEnd: func(interface{}) {
// 			// blocks.Brace(js.Obj, func(blk *js.Builder) {
// 			blk.Kv("type", typeName)
// 			blk.Q("extraState").R(js.Colon).Brace(js.Obj, func(out *js.Builder) {
// 				// it seems the extra state blk is required for loading to work
// 				// even if there's nothing to put in it.
// 			})
// 			writeContents(blk, "fields", &fields)
// 			writeContents(blk, "inputs", &inputs)
// 			// writeContents(blk, "next", &next)
// 			// })
// 		},
// 	}
// }

// func writeContents(out *js.Builder, key string, contents *js.Builder) {
// 	if contents.Len() > 0 {
// 		out.R(js.Comma).Q("fields").R(js.Colon).Brace(js.Obj, func(out *js.Builder) {
// 			out.S(contents.String())
// 		})
// 	}
// }
