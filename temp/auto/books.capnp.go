// Code generated by capnpc-go. DO NOT EDIT.

package books

import (
	strconv "strconv"
	capnp "zombiezen.com/go/capnproto2"
	text "zombiezen.com/go/capnproto2/encoding/text"
	schemas "zombiezen.com/go/capnproto2/schemas"
)

const Baz = uint64(0xf0007e64b3ac5312)

type Circle struct{ capnp.Struct }

// Circle_TypeID is the unique identifier for the type Circle.
const Circle_TypeID = 0xe94487436412f670

func NewCircle(s *capnp.Segment) (Circle, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return Circle{st}, err
}

func NewRootCircle(s *capnp.Segment) (Circle, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return Circle{st}, err
}

func ReadRootCircle(msg *capnp.Message) (Circle, error) {
	root, err := msg.RootPtr()
	return Circle{root.Struct()}, err
}

func (s Circle) String() string {
	str, _ := text.Marshal(0xe94487436412f670, s.Struct)
	return str
}

func (s Circle) Title() (string, error) {
	p, err := s.Struct.Ptr(0)
	return p.Text(), err
}

func (s Circle) HasTitle() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Circle) TitleBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	return p.TextBytes(), err
}

func (s Circle) SetTitle(v string) error {
	return s.Struct.SetText(0, v)
}

// Circle_List is a list of Circle.
type Circle_List struct{ capnp.List }

// NewCircle creates a new list of Circle.
func NewCircle_List(s *capnp.Segment, sz int32) (Circle_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	return Circle_List{l}, err
}

func (s Circle_List) At(i int) Circle { return Circle{s.List.Struct(i)} }

func (s Circle_List) Set(i int, v Circle) error { return s.List.SetStruct(i, v.Struct) }

func (s Circle_List) String() string {
	str, _ := text.MarshalList(0xe94487436412f670, s.List)
	return str
}

// Circle_Promise is a wrapper for a Circle promised by a client call.
type Circle_Promise struct{ *capnp.Pipeline }

func (p Circle_Promise) Struct() (Circle, error) {
	s, err := p.Pipeline.Struct()
	return Circle{s}, err
}

type Squeak struct{ capnp.Struct }

// Squeak_TypeID is the unique identifier for the type Squeak.
const Squeak_TypeID = 0x9f8475e4999f967e

func NewSqueak(s *capnp.Segment) (Squeak, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0})
	return Squeak{st}, err
}

func NewRootSqueak(s *capnp.Segment) (Squeak, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0})
	return Squeak{st}, err
}

func ReadRootSqueak(msg *capnp.Message) (Squeak, error) {
	root, err := msg.RootPtr()
	return Squeak{root.Struct()}, err
}

func (s Squeak) String() string {
	str, _ := text.Marshal(0x9f8475e4999f967e, s.Struct)
	return str
}

func (s Squeak) Width() int32 {
	return int32(s.Struct.Uint32(0))
}

func (s Squeak) SetWidth(v int32) {
	s.Struct.SetUint32(0, uint32(v))
}

// Squeak_List is a list of Squeak.
type Squeak_List struct{ capnp.List }

// NewSqueak creates a new list of Squeak.
func NewSqueak_List(s *capnp.Segment, sz int32) (Squeak_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0}, sz)
	return Squeak_List{l}, err
}

func (s Squeak_List) At(i int) Squeak { return Squeak{s.List.Struct(i)} }

func (s Squeak_List) Set(i int, v Squeak) error { return s.List.SetStruct(i, v.Struct) }

func (s Squeak_List) String() string {
	str, _ := text.MarshalList(0x9f8475e4999f967e, s.List)
	return str
}

// Squeak_Promise is a wrapper for a Squeak promised by a client call.
type Squeak_Promise struct{ *capnp.Pipeline }

func (p Squeak_Promise) Struct() (Squeak, error) {
	s, err := p.Pipeline.Struct()
	return Squeak{s}, err
}

type Book struct{ capnp.Struct }
type Book_Which uint16

const (
	Book_Which_circle Book_Which = 0
	Book_Which_square Book_Which = 1
)

func (w Book_Which) String() string {
	const s = "circlesquare"
	switch w {
	case Book_Which_circle:
		return s[0:6]
	case Book_Which_square:
		return s[6:12]

	}
	return "Book_Which(" + strconv.FormatUint(uint64(w), 10) + ")"
}

// Book_TypeID is the unique identifier for the type Book.
const Book_TypeID = 0x8100cc88d7d4d47c

func NewBook(s *capnp.Segment) (Book, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 2})
	return Book{st}, err
}

func NewRootBook(s *capnp.Segment) (Book, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 2})
	return Book{st}, err
}

func ReadRootBook(msg *capnp.Message) (Book, error) {
	root, err := msg.RootPtr()
	return Book{root.Struct()}, err
}

func (s Book) String() string {
	str, _ := text.Marshal(0x8100cc88d7d4d47c, s.Struct)
	return str
}

func (s Book) Which() Book_Which {
	return Book_Which(s.Struct.Uint16(4))
}
func (s Book) Title() (string, error) {
	p, err := s.Struct.Ptr(0)
	return p.Text(), err
}

func (s Book) HasTitle() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Book) TitleBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	return p.TextBytes(), err
}

func (s Book) SetTitle(v string) error {
	return s.Struct.SetText(0, v)
}

func (s Book) PageCount() int32 {
	return int32(s.Struct.Uint32(0))
}

func (s Book) SetPageCount(v int32) {
	s.Struct.SetUint32(0, uint32(v))
}

func (s Book) Circle() (Circle, error) {
	if s.Struct.Uint16(4) != 0 {
		panic("Which() != circle")
	}
	p, err := s.Struct.Ptr(1)
	return Circle{Struct: p.Struct()}, err
}

func (s Book) HasCircle() bool {
	if s.Struct.Uint16(4) != 0 {
		return false
	}
	p, err := s.Struct.Ptr(1)
	return p.IsValid() || err != nil
}

func (s Book) SetCircle(v Circle) error {
	s.Struct.SetUint16(4, 0)
	return s.Struct.SetPtr(1, v.Struct.ToPtr())
}

// NewCircle sets the circle field to a newly
// allocated Circle struct, preferring placement in s's segment.
func (s Book) NewCircle() (Circle, error) {
	s.Struct.SetUint16(4, 0)
	ss, err := NewCircle(s.Struct.Segment())
	if err != nil {
		return Circle{}, err
	}
	err = s.Struct.SetPtr(1, ss.Struct.ToPtr())
	return ss, err
}

func (s Book) Square() (Squeak, error) {
	if s.Struct.Uint16(4) != 1 {
		panic("Which() != square")
	}
	p, err := s.Struct.Ptr(1)
	return Squeak{Struct: p.Struct()}, err
}

func (s Book) HasSquare() bool {
	if s.Struct.Uint16(4) != 1 {
		return false
	}
	p, err := s.Struct.Ptr(1)
	return p.IsValid() || err != nil
}

func (s Book) SetSquare(v Squeak) error {
	s.Struct.SetUint16(4, 1)
	return s.Struct.SetPtr(1, v.Struct.ToPtr())
}

// NewSquare sets the square field to a newly
// allocated Squeak struct, preferring placement in s's segment.
func (s Book) NewSquare() (Squeak, error) {
	s.Struct.SetUint16(4, 1)
	ss, err := NewSqueak(s.Struct.Segment())
	if err != nil {
		return Squeak{}, err
	}
	err = s.Struct.SetPtr(1, ss.Struct.ToPtr())
	return ss, err
}

// Book_List is a list of Book.
type Book_List struct{ capnp.List }

// NewBook creates a new list of Book.
func NewBook_List(s *capnp.Segment, sz int32) (Book_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 2}, sz)
	return Book_List{l}, err
}

func (s Book_List) At(i int) Book { return Book{s.List.Struct(i)} }

func (s Book_List) Set(i int, v Book) error { return s.List.SetStruct(i, v.Struct) }

func (s Book_List) String() string {
	str, _ := text.MarshalList(0x8100cc88d7d4d47c, s.List)
	return str
}

// Book_Promise is a wrapper for a Book promised by a client call.
type Book_Promise struct{ *capnp.Pipeline }

func (p Book_Promise) Struct() (Book, error) {
	s, err := p.Pipeline.Struct()
	return Book{s}, err
}

func (p Book_Promise) Circle() Circle_Promise {
	return Circle_Promise{Pipeline: p.Pipeline.GetPipeline(1)}
}

func (p Book_Promise) Square() Squeak_Promise {
	return Squeak_Promise{Pipeline: p.Pipeline.GetPipeline(1)}
}

const schema_85d3acc39d94e0f8 = "x\xda|\x91=\x88\xd4P\x14\x85\xcf\xb9/c,\xb2" +
	"\x98L\x06\xac\x96\xb4\x0a\xfe\xe0V2\xcd\xa8\xa3\xbd\xcf" +
	"i-\xcc&q'L\xdc\xccO\x82\"\xee\x82\xa2(" +
	"\x82\x9d\x08\xbb0,\xf6\xdbY\x09\xca\x82\xb2\x9d\x8d\xb8" +
	"`k\xa1\x88\x9d6be\xe4!3k,\xb6=\xe7" +
	"\xbc{\xefw\x9e\xfb\xea\x9c\x9ci\xac\x10\xd0n\xe3\xd0" +
	"\xce\x9d\xbd\xbd\x8f\x8f\xde\xdd\xd5\x0eY\xfd\xfa\xf4t\xfa" +
	"v\xfb\xc3\x03\\\x12[(\xde\xc6c\xef\xb9\x0dx\xd3" +
	"\xaf`\xb5\xfelk\xe3sy\x7f\x0b\xf5\xa8e\x03~" +
	"\xc9M\x7f\x8d\xb6\xbf\xc6\xc0\x9f\xb2\x03V\xc3\x9f\xcd\xb8" +
	"\xfb\xf0\xe27x\xce?\xe1\x06M\xfa57\xfd]\xda" +
	"\xfe.\x03\xff\x8bI\xef4{\xdb/\xe2\xf5\xef\xef\x9d" +
	"Fud?\x0d\xfa\xc7\xe5\x9e\x7fRl\xa0wL\x14" +
	"\xc1j9\xcf\x07\x93SQ\xc8\xe1\xea\xb0}!\xcf\x07" +
	"\xc0eR\xbb\xca\x02,\x02^\xb8\x04\xe8\xab\x8a\xba/" +
	"$[4Zr\x05\xd0\xb1\xa2\x1e\x0a\x17\xa5\xaa\xd8\xa2" +
	"\x00\xde\x8d6\xa0\xfb\x8a\xba\x10.\xaa\xdfFV\x807" +
	"2r\xa6\xa8o\x09\x83\"-\xb2\x84\x0e\x84\x8e\xa1\x0a" +
	"W\x92n^\xae\x82\x05-\x08-\xb0\x13\xa5\xe3(K" +
	"\xe8\xee#\x83t\xc1\xcedT\x86cc\xcc\x9b\xfbk" +
	"\xd4\x19z\xa32T\xe3\xc4@X\x94\xea\xc7\x93\xd3G" +
	"\x9b\xd7^\xbe\x81\xb6\x84\xe7]\xd2\x01<\xb6;\xbdQ" +
	"\x99\x84\x03@[s\xd2\x05CzXQ\xb7\x84\xc1\xcd" +
	"4.\xfa\xb3\x9b\xea\x1b\xba\xe98R\xd9l\xc3\xac\xeb" +
	"\xda\xfc\xa5 \xca\xc2\xc9\xe4\xbf\xf1\xdeB\xa0O(\xea" +
	"\xb3\xb3\x1e\x0ex\x7f=M\xb2\x18\x987U\xbb`9" +
	"\xbcm~\xc9\x98\x7f\x02\x00\x00\xff\xff\x83\xee\xa6\xee"

func init() {
	schemas.Register(schema_85d3acc39d94e0f8,
		0x8100cc88d7d4d47c,
		0x9f8475e4999f967e,
		0xe94487436412f670,
		0xf0007e64b3ac5312)
}
