package rift_test

import (
	"bufio"
	"embed"
	"encoding/json"
	"io"
	"log"
	"path"
	"strings"
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
	var focus string
	// focus = "headerFooterComment3"
	if files, e := riftData.ReadDir(testFolder); e != nil {
		t.Fatal(e)
	} else {
		for _, info := range files {
			shortName := info.Name()
			riftName := path.Join(testFolder, shortName)
			jsonName := riftName[:len(riftName)-4] + "json"
			if (len(focus) > 0 && !strings.Contains(riftName, focus)) ||
				strings.HasPrefix(shortName, "x_") {
				t.Log("skipping", riftName)
				continue
			}
			//
			t.Log("trying", riftName)
			if got, e := readRift(riftName); e != nil {
				t.Log("error", e)
				t.Fail()
			} else if want, e := readJson(jsonName); e != nil {
				if a, e := json.MarshalIndent(got, "", " "); e != nil {
					t.Fatal(e)
				} else {
					t.Log(string(a))
				}
				t.Fail()

			} else {
				// reflect.DeepEqual
				if diff := pretty.Diff(got, want); len(diff) > 0 {
					log.Println("ng: ", riftName)
					log.Println("got: ", pretty.Sprint(got))
					t.Log(diff)
					t.Fail()
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
		keepComments := strings.Contains(strings.ToLower(filePath), "comment")
		comments := rift.DiscardCommentWriter
		if keepComments {
			comments = rift.KeepCommentWriter
		}
		doc := rift.NewDocument(stdmap.Build, comments)
		if res, e := doc.ReadDoc(bufio.NewReader(fp)); e != nil {
			err = e
		} else if len(res.Comment) > 0 {
			ret = map[string]any{
				"content": res.Content,
				"comment": res.Comment,
			}
		} else {
			ret = map[string]any{
				"content": res.Content,
			}
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
