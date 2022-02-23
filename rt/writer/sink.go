package writer

import "io"

// Sink implements a container for Output.
// it just happens to help provide simple output handling for rt.Runtime implementations.
type Sink struct {
	Output io.Writer
}

func (k *Sink) Writer() io.Writer {
	return k.Output
}

func (k *Sink) SetWriter(out io.Writer) (ret io.Writer) {
	ret, k.Output = k.Output, out
	return
}
