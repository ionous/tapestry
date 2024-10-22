package print

import (
	"bytes"
	"io"
	"testing"
)

func TestBracket(t *testing.T) {
	span := Parens()
	w := span.ChunkOutput()
	io.WriteString(w, "hello")
	io.WriteString(w, "you")
	w.Close()
	if str := span.String(); str != "(hello you)" {
		t.Fatal("mismatched", str)
	}
}

func TestManualBracket(t *testing.T) {
	span := NewSpanner()
	w := span.ChunkOutput()
	io.WriteString(w, "hey")
	io.WriteString(w, "(you)")
	io.WriteString(w, "guys")
	if str := span.String(); str != "hey (you) guys" {
		t.Fatal("mismatched", str)
	}
}

func TestCapitalize(t *testing.T) {
	span := NewSpanner()
	w := Capitalize(span.ChunkOutput())
	io.WriteString(w, "hello")
	io.WriteString(w, "you")
	if str := span.String(); str != "Hello you" {
		t.Fatal("mismatched", str)
	}
}

func TestLowercase(t *testing.T) {
	span := NewSpanner()
	w := Lowercase(span.ChunkOutput())
	io.WriteString(w, "Hello")
	io.WriteString(w, "Hugh")
	if str := span.String(); str != "hello hugh" {
		t.Fatal("mismatched", str)
	}
}

func TestTitlecase(t *testing.T) {
	span := NewSpanner()
	w := TitleCase(span.ChunkOutput())
	io.WriteString(w, "hello")
	io.WriteString(w, "you")
	if str := span.String(); str != "Hello You" {
		t.Fatal("mismatched", str)
	}
}

func TestTag(t *testing.T) {
	var buf bytes.Buffer
	w := Tag(&buf, "tag")
	io.WriteString(w, "hello")
	io.WriteString(w, "goodbye")
	w.Close()
	if str := buf.String(); str != "<tag>hellogoodbye</tag>" {
		t.Fatal("mismatched", str)
	}
}

func TestTagSpan(t *testing.T) {
	span := Brackets("<tag>", "</tag>")
	w := span.ChunkOutput()
	io.WriteString(w, "hello")
	io.WriteString(w, "goodbye")
	w.Close()
	if str := span.String(); str != "<tag>hello goodbye</tag>" {
		t.Fatal("mismatched", str)
	}
}
