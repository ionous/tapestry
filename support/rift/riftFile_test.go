package rift_test

import (
	"bufio"
	"embed"
	"encoding/json"
	"io"
	"path"
	"testing"

	"git.sr.ht/~ionous/tapestry/support/rift"
	"git.sr.ht/~ionous/tapestry/support/rift/stdmap"
	"github.com/kr/pretty"
)

//go:embed testdata/*.rift
var riftData embed.FS

//go:embed testdata/*.json
var jsonData embed.FS

const testFolder = "testdata"

func TestFiles(t *testing.T) {

	if files, e := riftData.ReadDir(testFolder); e != nil {
		t.Fatal(e)
	} else {
		for _, info := range files {
			riftName := path.Join(testFolder, info.Name())
			jsonName := riftName[:len(riftName)-4] + "json"
			//
			if got, e := readRift(riftName); e != nil {
				t.Fatal(e)
			} else if want, e := readJson(jsonName); e != nil {
				t.Fatal(e)
			} else {
				// reflect.DeepEqual
				if diff := pretty.Diff(got, want); len(diff) > 0 {
					t.Fatal(riftName, diff)
				} else {
					t.Log("ok: ", riftName)
				}
			}
		}
	}
}

func readRift(filePath string) (ret any, err error) {
	if fp, e := riftData.Open(filePath); e != nil {
		err = e
	} else {
		src := bufio.NewReader(fp)
		doc := rift.Document{MakeMap: stdmap.Build}
		if e := doc.ReadLines(src, doc.NewEntry()); e != nil {
			err = e
		} else {
			ret = doc.Value
		}
	}
	return

}

func readJson(filePath string) (ret any, err error) {
	if fp, e := jsonData.Open(filePath); e != nil {
		err = e
	} else if b, e := io.ReadAll(fp); e != nil {
		err = e
	} else if e := json.Unmarshal(b, &ret); e != nil {
		err = e
	}
	return

}
