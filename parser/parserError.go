package parser

import "errors"

type AlwaysError struct{}

func (n AlwaysError) Scan(ctx Context, bounds Bounds, cs Cursor) (ret Result, err error) {
	err = errors.New("nothing matched")
	return
}
