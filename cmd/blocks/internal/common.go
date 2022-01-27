package blocks

import (
	"strconv"
	"strings"
)

// ex. "message0": "%1%2%3%4%5%6%7%8",
func writeMsg(out *Js, cnt int) {
	for i := 1; i <= cnt; i++ {
		out.R(percent).S(strconv.Itoa(i))
	}
}

type Input struct {
	Name  string // optional name for the input, will be capitalized
	Label string // optional label
	Check string // optional type check ( what's allowed to connect )
	Type  string // one of input_statement, etc. dummy if not specified
}

// write a number of fields, followed by the input that they merge down into
// using blockly json's message interpolation syntax.
func writeInput(args *Js, in Input, fields ...func(*Js)) (ret int) {
	// write the label, if any
	if n := in.Label; len(n) > 0 {
		args.Brace(obj, func(lab *Js) {
			lab.
				Kv("type", FieldLabel).R(comma).
				Kv("text", n)
		}).R(comma)
		ret++
	}
	// write the fields
	for _, field := range fields {
		args.Brace(obj, field).R(comma)
		ret++
	}
	// write the input they are a part of
	args.Brace(obj, func(tail *Js) {
		if n := in.Name; len(n) > 0 {
			tail.Kv("name", strings.ToUpper(n)).R(comma)
		}
		if c := in.Check; len(c) > 0 {
			tail.Kv("check", c).R(comma)
		}
		if t := in.Type; len(t) > 0 {
			tail.Kv("type", t)
		} else {
			tail.Kv("type", InputDummy)
		}
		ret++
	})
	return
}
