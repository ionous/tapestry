package serve

import "io"

// helper to close pipes ( esp for example on error )
type pipes []io.Closer

func (ps *pipes) Close() {
	for _, p := range *ps {
		p.Close()
	}
}
func (p *pipes) AddReader(i io.ReadCloser, e error) (io.ReadCloser, error) {
	if e == nil {
		*p = append(*p, i)
	}
	return i, e
}
func (p *pipes) AddWriter(i io.WriteCloser, e error) (io.WriteCloser, error) {
	if e == nil {
		*p = append(*p, i)
	}
	return i, e
}
