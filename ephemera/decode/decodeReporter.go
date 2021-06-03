package decode

import (
	"git.sr.ht/~ionous/iffy/ephemera/reader"
)

type IssueReport func(reader.Position, error)

func NewDecoderReporter(report IssueReport) *Decoder {
	return &Decoder{source: "decoder", cmds: make(map[string]cmdRec), issueFn: report}
}
func (m *Decoder) SetSource(source string) *Decoder {
	m.source = source
	return m
}
func (m *Decoder) report(ofs string, err error) {
	m.issueFn(reader.Position{Source: m.source, Offset: ofs}, err)
	m.IssueCount++
}
