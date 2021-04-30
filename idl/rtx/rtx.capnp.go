// Code generated by capnpc-go. DO NOT EDIT.

package rtx

import (
	capnp "zombiezen.com/go/capnproto2"
	text "zombiezen.com/go/capnproto2/encoding/text"
	schemas "zombiezen.com/go/capnproto2/schemas"
)

type Assignment struct{ capnp.Struct }

// Assignment_TypeID is the unique identifier for the type Assignment.
const Assignment_TypeID = 0x91fcee5c442c25a7

func NewAssignment(s *capnp.Segment) (Assignment, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return Assignment{st}, err
}

func NewRootAssignment(s *capnp.Segment) (Assignment, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return Assignment{st}, err
}

func ReadRootAssignment(msg *capnp.Message) (Assignment, error) {
	root, err := msg.RootPtr()
	return Assignment{root.Struct()}, err
}

func (s Assignment) String() string {
	str, _ := text.Marshal(0x91fcee5c442c25a7, s.Struct)
	return str
}

func (s Assignment) Eval() (capnp.Pointer, error) {
	return s.Struct.Pointer(0)
}

func (s Assignment) HasEval() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Assignment) EvalPtr() (capnp.Ptr, error) {
	return s.Struct.Ptr(0)
}

func (s Assignment) SetEval(v capnp.Pointer) error {
	return s.Struct.SetPointer(0, v)
}

func (s Assignment) SetEvalPtr(v capnp.Ptr) error {
	return s.Struct.SetPtr(0, v)
}

// Assignment_List is a list of Assignment.
type Assignment_List struct{ capnp.List }

// NewAssignment creates a new list of Assignment.
func NewAssignment_List(s *capnp.Segment, sz int32) (Assignment_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	return Assignment_List{l}, err
}

func (s Assignment_List) At(i int) Assignment { return Assignment{s.List.Struct(i)} }

func (s Assignment_List) Set(i int, v Assignment) error { return s.List.SetStruct(i, v.Struct) }

func (s Assignment_List) String() string {
	str, _ := text.MarshalList(0x91fcee5c442c25a7, s.List)
	return str
}

// Assignment_Promise is a wrapper for a Assignment promised by a client call.
type Assignment_Promise struct{ *capnp.Pipeline }

func (p Assignment_Promise) Struct() (Assignment, error) {
	s, err := p.Pipeline.Struct()
	return Assignment{s}, err
}

func (p Assignment_Promise) Eval() *capnp.Pipeline {
	return p.Pipeline.GetPipeline(0)
}

type BoolEval struct{ capnp.Struct }

// BoolEval_TypeID is the unique identifier for the type BoolEval.
const BoolEval_TypeID = 0xd7932881bef41f81

func NewBoolEval(s *capnp.Segment) (BoolEval, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return BoolEval{st}, err
}

func NewRootBoolEval(s *capnp.Segment) (BoolEval, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return BoolEval{st}, err
}

func ReadRootBoolEval(msg *capnp.Message) (BoolEval, error) {
	root, err := msg.RootPtr()
	return BoolEval{root.Struct()}, err
}

func (s BoolEval) String() string {
	str, _ := text.Marshal(0xd7932881bef41f81, s.Struct)
	return str
}

func (s BoolEval) Eval() (capnp.Pointer, error) {
	return s.Struct.Pointer(0)
}

func (s BoolEval) HasEval() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s BoolEval) EvalPtr() (capnp.Ptr, error) {
	return s.Struct.Ptr(0)
}

func (s BoolEval) SetEval(v capnp.Pointer) error {
	return s.Struct.SetPointer(0, v)
}

func (s BoolEval) SetEvalPtr(v capnp.Ptr) error {
	return s.Struct.SetPtr(0, v)
}

// BoolEval_List is a list of BoolEval.
type BoolEval_List struct{ capnp.List }

// NewBoolEval creates a new list of BoolEval.
func NewBoolEval_List(s *capnp.Segment, sz int32) (BoolEval_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	return BoolEval_List{l}, err
}

func (s BoolEval_List) At(i int) BoolEval { return BoolEval{s.List.Struct(i)} }

func (s BoolEval_List) Set(i int, v BoolEval) error { return s.List.SetStruct(i, v.Struct) }

func (s BoolEval_List) String() string {
	str, _ := text.MarshalList(0xd7932881bef41f81, s.List)
	return str
}

// BoolEval_Promise is a wrapper for a BoolEval promised by a client call.
type BoolEval_Promise struct{ *capnp.Pipeline }

func (p BoolEval_Promise) Struct() (BoolEval, error) {
	s, err := p.Pipeline.Struct()
	return BoolEval{s}, err
}

func (p BoolEval_Promise) Eval() *capnp.Pipeline {
	return p.Pipeline.GetPipeline(0)
}

type Execute struct{ capnp.Struct }

// Execute_TypeID is the unique identifier for the type Execute.
const Execute_TypeID = 0xf1fe2e1f686cf364

func NewExecute(s *capnp.Segment) (Execute, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return Execute{st}, err
}

func NewRootExecute(s *capnp.Segment) (Execute, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return Execute{st}, err
}

func ReadRootExecute(msg *capnp.Message) (Execute, error) {
	root, err := msg.RootPtr()
	return Execute{root.Struct()}, err
}

func (s Execute) String() string {
	str, _ := text.Marshal(0xf1fe2e1f686cf364, s.Struct)
	return str
}

func (s Execute) Eval() (capnp.Pointer, error) {
	return s.Struct.Pointer(0)
}

func (s Execute) HasEval() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Execute) EvalPtr() (capnp.Ptr, error) {
	return s.Struct.Ptr(0)
}

func (s Execute) SetEval(v capnp.Pointer) error {
	return s.Struct.SetPointer(0, v)
}

func (s Execute) SetEvalPtr(v capnp.Ptr) error {
	return s.Struct.SetPtr(0, v)
}

// Execute_List is a list of Execute.
type Execute_List struct{ capnp.List }

// NewExecute creates a new list of Execute.
func NewExecute_List(s *capnp.Segment, sz int32) (Execute_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	return Execute_List{l}, err
}

func (s Execute_List) At(i int) Execute { return Execute{s.List.Struct(i)} }

func (s Execute_List) Set(i int, v Execute) error { return s.List.SetStruct(i, v.Struct) }

func (s Execute_List) String() string {
	str, _ := text.MarshalList(0xf1fe2e1f686cf364, s.List)
	return str
}

// Execute_Promise is a wrapper for a Execute promised by a client call.
type Execute_Promise struct{ *capnp.Pipeline }

func (p Execute_Promise) Struct() (Execute, error) {
	s, err := p.Pipeline.Struct()
	return Execute{s}, err
}

func (p Execute_Promise) Eval() *capnp.Pipeline {
	return p.Pipeline.GetPipeline(0)
}

type NumListEval struct{ capnp.Struct }

// NumListEval_TypeID is the unique identifier for the type NumListEval.
const NumListEval_TypeID = 0xb45129a42f682d82

func NewNumListEval(s *capnp.Segment) (NumListEval, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return NumListEval{st}, err
}

func NewRootNumListEval(s *capnp.Segment) (NumListEval, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return NumListEval{st}, err
}

func ReadRootNumListEval(msg *capnp.Message) (NumListEval, error) {
	root, err := msg.RootPtr()
	return NumListEval{root.Struct()}, err
}

func (s NumListEval) String() string {
	str, _ := text.Marshal(0xb45129a42f682d82, s.Struct)
	return str
}

func (s NumListEval) Eval() (capnp.Pointer, error) {
	return s.Struct.Pointer(0)
}

func (s NumListEval) HasEval() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s NumListEval) EvalPtr() (capnp.Ptr, error) {
	return s.Struct.Ptr(0)
}

func (s NumListEval) SetEval(v capnp.Pointer) error {
	return s.Struct.SetPointer(0, v)
}

func (s NumListEval) SetEvalPtr(v capnp.Ptr) error {
	return s.Struct.SetPtr(0, v)
}

// NumListEval_List is a list of NumListEval.
type NumListEval_List struct{ capnp.List }

// NewNumListEval creates a new list of NumListEval.
func NewNumListEval_List(s *capnp.Segment, sz int32) (NumListEval_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	return NumListEval_List{l}, err
}

func (s NumListEval_List) At(i int) NumListEval { return NumListEval{s.List.Struct(i)} }

func (s NumListEval_List) Set(i int, v NumListEval) error { return s.List.SetStruct(i, v.Struct) }

func (s NumListEval_List) String() string {
	str, _ := text.MarshalList(0xb45129a42f682d82, s.List)
	return str
}

// NumListEval_Promise is a wrapper for a NumListEval promised by a client call.
type NumListEval_Promise struct{ *capnp.Pipeline }

func (p NumListEval_Promise) Struct() (NumListEval, error) {
	s, err := p.Pipeline.Struct()
	return NumListEval{s}, err
}

func (p NumListEval_Promise) Eval() *capnp.Pipeline {
	return p.Pipeline.GetPipeline(0)
}

type NumberEval struct{ capnp.Struct }

// NumberEval_TypeID is the unique identifier for the type NumberEval.
const NumberEval_TypeID = 0x92a24bdedf5f7511

func NewNumberEval(s *capnp.Segment) (NumberEval, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return NumberEval{st}, err
}

func NewRootNumberEval(s *capnp.Segment) (NumberEval, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return NumberEval{st}, err
}

func ReadRootNumberEval(msg *capnp.Message) (NumberEval, error) {
	root, err := msg.RootPtr()
	return NumberEval{root.Struct()}, err
}

func (s NumberEval) String() string {
	str, _ := text.Marshal(0x92a24bdedf5f7511, s.Struct)
	return str
}

func (s NumberEval) Eval() (capnp.Pointer, error) {
	return s.Struct.Pointer(0)
}

func (s NumberEval) HasEval() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s NumberEval) EvalPtr() (capnp.Ptr, error) {
	return s.Struct.Ptr(0)
}

func (s NumberEval) SetEval(v capnp.Pointer) error {
	return s.Struct.SetPointer(0, v)
}

func (s NumberEval) SetEvalPtr(v capnp.Ptr) error {
	return s.Struct.SetPtr(0, v)
}

// NumberEval_List is a list of NumberEval.
type NumberEval_List struct{ capnp.List }

// NewNumberEval creates a new list of NumberEval.
func NewNumberEval_List(s *capnp.Segment, sz int32) (NumberEval_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	return NumberEval_List{l}, err
}

func (s NumberEval_List) At(i int) NumberEval { return NumberEval{s.List.Struct(i)} }

func (s NumberEval_List) Set(i int, v NumberEval) error { return s.List.SetStruct(i, v.Struct) }

func (s NumberEval_List) String() string {
	str, _ := text.MarshalList(0x92a24bdedf5f7511, s.List)
	return str
}

// NumberEval_Promise is a wrapper for a NumberEval promised by a client call.
type NumberEval_Promise struct{ *capnp.Pipeline }

func (p NumberEval_Promise) Struct() (NumberEval, error) {
	s, err := p.Pipeline.Struct()
	return NumberEval{s}, err
}

func (p NumberEval_Promise) Eval() *capnp.Pipeline {
	return p.Pipeline.GetPipeline(0)
}

type RecordEval struct{ capnp.Struct }

// RecordEval_TypeID is the unique identifier for the type RecordEval.
const RecordEval_TypeID = 0x9ae3b5ecaf74323a

func NewRecordEval(s *capnp.Segment) (RecordEval, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return RecordEval{st}, err
}

func NewRootRecordEval(s *capnp.Segment) (RecordEval, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return RecordEval{st}, err
}

func ReadRootRecordEval(msg *capnp.Message) (RecordEval, error) {
	root, err := msg.RootPtr()
	return RecordEval{root.Struct()}, err
}

func (s RecordEval) String() string {
	str, _ := text.Marshal(0x9ae3b5ecaf74323a, s.Struct)
	return str
}

func (s RecordEval) Eval() (capnp.Pointer, error) {
	return s.Struct.Pointer(0)
}

func (s RecordEval) HasEval() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s RecordEval) EvalPtr() (capnp.Ptr, error) {
	return s.Struct.Ptr(0)
}

func (s RecordEval) SetEval(v capnp.Pointer) error {
	return s.Struct.SetPointer(0, v)
}

func (s RecordEval) SetEvalPtr(v capnp.Ptr) error {
	return s.Struct.SetPtr(0, v)
}

// RecordEval_List is a list of RecordEval.
type RecordEval_List struct{ capnp.List }

// NewRecordEval creates a new list of RecordEval.
func NewRecordEval_List(s *capnp.Segment, sz int32) (RecordEval_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	return RecordEval_List{l}, err
}

func (s RecordEval_List) At(i int) RecordEval { return RecordEval{s.List.Struct(i)} }

func (s RecordEval_List) Set(i int, v RecordEval) error { return s.List.SetStruct(i, v.Struct) }

func (s RecordEval_List) String() string {
	str, _ := text.MarshalList(0x9ae3b5ecaf74323a, s.List)
	return str
}

// RecordEval_Promise is a wrapper for a RecordEval promised by a client call.
type RecordEval_Promise struct{ *capnp.Pipeline }

func (p RecordEval_Promise) Struct() (RecordEval, error) {
	s, err := p.Pipeline.Struct()
	return RecordEval{s}, err
}

func (p RecordEval_Promise) Eval() *capnp.Pipeline {
	return p.Pipeline.GetPipeline(0)
}

type RecordListEval struct{ capnp.Struct }

// RecordListEval_TypeID is the unique identifier for the type RecordListEval.
const RecordListEval_TypeID = 0xa309c167368bfe10

func NewRecordListEval(s *capnp.Segment) (RecordListEval, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return RecordListEval{st}, err
}

func NewRootRecordListEval(s *capnp.Segment) (RecordListEval, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return RecordListEval{st}, err
}

func ReadRootRecordListEval(msg *capnp.Message) (RecordListEval, error) {
	root, err := msg.RootPtr()
	return RecordListEval{root.Struct()}, err
}

func (s RecordListEval) String() string {
	str, _ := text.Marshal(0xa309c167368bfe10, s.Struct)
	return str
}

func (s RecordListEval) Eval() (capnp.Pointer, error) {
	return s.Struct.Pointer(0)
}

func (s RecordListEval) HasEval() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s RecordListEval) EvalPtr() (capnp.Ptr, error) {
	return s.Struct.Ptr(0)
}

func (s RecordListEval) SetEval(v capnp.Pointer) error {
	return s.Struct.SetPointer(0, v)
}

func (s RecordListEval) SetEvalPtr(v capnp.Ptr) error {
	return s.Struct.SetPtr(0, v)
}

// RecordListEval_List is a list of RecordListEval.
type RecordListEval_List struct{ capnp.List }

// NewRecordListEval creates a new list of RecordListEval.
func NewRecordListEval_List(s *capnp.Segment, sz int32) (RecordListEval_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	return RecordListEval_List{l}, err
}

func (s RecordListEval_List) At(i int) RecordListEval { return RecordListEval{s.List.Struct(i)} }

func (s RecordListEval_List) Set(i int, v RecordListEval) error { return s.List.SetStruct(i, v.Struct) }

func (s RecordListEval_List) String() string {
	str, _ := text.MarshalList(0xa309c167368bfe10, s.List)
	return str
}

// RecordListEval_Promise is a wrapper for a RecordListEval promised by a client call.
type RecordListEval_Promise struct{ *capnp.Pipeline }

func (p RecordListEval_Promise) Struct() (RecordListEval, error) {
	s, err := p.Pipeline.Struct()
	return RecordListEval{s}, err
}

func (p RecordListEval_Promise) Eval() *capnp.Pipeline {
	return p.Pipeline.GetPipeline(0)
}

type TextEval struct{ capnp.Struct }

// TextEval_TypeID is the unique identifier for the type TextEval.
const TextEval_TypeID = 0xa2c8e611e5b84332

func NewTextEval(s *capnp.Segment) (TextEval, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return TextEval{st}, err
}

func NewRootTextEval(s *capnp.Segment) (TextEval, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return TextEval{st}, err
}

func ReadRootTextEval(msg *capnp.Message) (TextEval, error) {
	root, err := msg.RootPtr()
	return TextEval{root.Struct()}, err
}

func (s TextEval) String() string {
	str, _ := text.Marshal(0xa2c8e611e5b84332, s.Struct)
	return str
}

func (s TextEval) Eval() (capnp.Pointer, error) {
	return s.Struct.Pointer(0)
}

func (s TextEval) HasEval() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s TextEval) EvalPtr() (capnp.Ptr, error) {
	return s.Struct.Ptr(0)
}

func (s TextEval) SetEval(v capnp.Pointer) error {
	return s.Struct.SetPointer(0, v)
}

func (s TextEval) SetEvalPtr(v capnp.Ptr) error {
	return s.Struct.SetPtr(0, v)
}

// TextEval_List is a list of TextEval.
type TextEval_List struct{ capnp.List }

// NewTextEval creates a new list of TextEval.
func NewTextEval_List(s *capnp.Segment, sz int32) (TextEval_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	return TextEval_List{l}, err
}

func (s TextEval_List) At(i int) TextEval { return TextEval{s.List.Struct(i)} }

func (s TextEval_List) Set(i int, v TextEval) error { return s.List.SetStruct(i, v.Struct) }

func (s TextEval_List) String() string {
	str, _ := text.MarshalList(0xa2c8e611e5b84332, s.List)
	return str
}

// TextEval_Promise is a wrapper for a TextEval promised by a client call.
type TextEval_Promise struct{ *capnp.Pipeline }

func (p TextEval_Promise) Struct() (TextEval, error) {
	s, err := p.Pipeline.Struct()
	return TextEval{s}, err
}

func (p TextEval_Promise) Eval() *capnp.Pipeline {
	return p.Pipeline.GetPipeline(0)
}

type TextListEval struct{ capnp.Struct }

// TextListEval_TypeID is the unique identifier for the type TextListEval.
const TextListEval_TypeID = 0xbe170ad5982c7478

func NewTextListEval(s *capnp.Segment) (TextListEval, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return TextListEval{st}, err
}

func NewRootTextListEval(s *capnp.Segment) (TextListEval, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return TextListEval{st}, err
}

func ReadRootTextListEval(msg *capnp.Message) (TextListEval, error) {
	root, err := msg.RootPtr()
	return TextListEval{root.Struct()}, err
}

func (s TextListEval) String() string {
	str, _ := text.Marshal(0xbe170ad5982c7478, s.Struct)
	return str
}

func (s TextListEval) Eval() (capnp.Pointer, error) {
	return s.Struct.Pointer(0)
}

func (s TextListEval) HasEval() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s TextListEval) EvalPtr() (capnp.Ptr, error) {
	return s.Struct.Ptr(0)
}

func (s TextListEval) SetEval(v capnp.Pointer) error {
	return s.Struct.SetPointer(0, v)
}

func (s TextListEval) SetEvalPtr(v capnp.Ptr) error {
	return s.Struct.SetPtr(0, v)
}

// TextListEval_List is a list of TextListEval.
type TextListEval_List struct{ capnp.List }

// NewTextListEval creates a new list of TextListEval.
func NewTextListEval_List(s *capnp.Segment, sz int32) (TextListEval_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	return TextListEval_List{l}, err
}

func (s TextListEval_List) At(i int) TextListEval { return TextListEval{s.List.Struct(i)} }

func (s TextListEval_List) Set(i int, v TextListEval) error { return s.List.SetStruct(i, v.Struct) }

func (s TextListEval_List) String() string {
	str, _ := text.MarshalList(0xbe170ad5982c7478, s.List)
	return str
}

// TextListEval_Promise is a wrapper for a TextListEval promised by a client call.
type TextListEval_Promise struct{ *capnp.Pipeline }

func (p TextListEval_Promise) Struct() (TextListEval, error) {
	s, err := p.Pipeline.Struct()
	return TextListEval{s}, err
}

func (p TextListEval_Promise) Eval() *capnp.Pipeline {
	return p.Pipeline.GetPipeline(0)
}

const schema_8167b37f5d971356 = "x\xda2\x90at`2d\x95gf`\x08Ta" +
	"e\xfb\xbf\\U\xc7%\xe6\xdd\x9f\x89\x0c\x82\\\x8c\xff" +
	"\xc3\x84\xa7\xc7\xd6oNod`edg`\x10." +
	"e\\%\\\x0bfU2\xda30\xfe\x17,\x8d\xbf" +
	"\x7f\xcf{\xd1$ljg2\xae\x12^\x08f\xcd\x05" +
	"\xab\xb52*Y\xfff\xeb\xe3Y\xd8\xd4\xeee\\%" +
	"|\x14\xcc:\x08Vk\xe4\xbc\xe3\xa9\xe0\xb3\x13\x8b\xb0" +
	"\xa9}\xc88K\xf8%\x98\xf5\x14\xacV\xe0_\xb7Y" +
	"\xfaA\xce\xc5\xd8\xd4\xb22\x9d\x12\x16d\x02\xb1x\x99" +
	"@j\x9bt3\xf4\x97h\x06n\xc1\xa6V\x97i\x93" +
	"\xb0)X\xad!XmE\x89\xce\x8c\xab\\\xe2\xfb\xb0" +
	"\xa9\x0dd\xda%\x1c\x09V\x1b\x0aV\xdb(\xffe_" +
	"\xa3\xc6\xe4\xebX\xc3\x8ci\x96p-Xm%Xm" +
	"\xca\xe7\x9c\x0cy\xbd\x7f\x1f\xb1\x86\x19\xd3$\xe1\x85`" +
	"\xb5s\xc1j\x8bJ*\xf4\x92\x13\x0b\xf2\x18\x0b\xac\x1c" +
	"\x8b\x8b3\xd3\xf3\xf8sS\xf3J\x02\x18\x19\x03Y\x98" +
	"Y\x18\x18X\x18\x19\x18\x04y\xb5\x18\x18\x029\x98\x19" +
	"\x03E\x98\x18\xf9S\xcb\x12s\x18\x85\x18\x98\x18\x85P" +
	"u\xfb\x95\xe6&\xa5\x16\xf1\xbb\x96%\xe6\x90\xa1;(" +
	"59\xbf(\x85\\\xdd!\xa9\x15%\xaee\xcc\xa4\xeb" +
	"e\x82\xd9\xec\x93Y\\\x02\xb2\x9c\x81\x81<\xaf\xfbd" +
	"\x16\xcb\x97P\xe0z\x9fL\xfbbr\xf5;\xe5\xe7\xe7" +
	"\x90\xe5{\xc6\x02+\xd7\x8a\xd4\xe4\xd2\x12\xc6T\"\xb5" +
	"\x02\x02\x00\x00\xff\xff\xc4\xac\xe0t"

func init() {
	schemas.Register(schema_8167b37f5d971356,
		0x91fcee5c442c25a7,
		0x92a24bdedf5f7511,
		0x9ae3b5ecaf74323a,
		0xa2c8e611e5b84332,
		0xa309c167368bfe10,
		0xb45129a42f682d82,
		0xbe170ad5982c7478,
		0xd7932881bef41f81,
		0xf1fe2e1f686cf364)
}
