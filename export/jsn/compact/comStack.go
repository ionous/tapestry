package compact

type comStack []compactMarshaler

func (j *comStack) push(m compactMarshaler) {
	(*j) = append(*j, m)
}
func (j *comStack) pop() (ret compactMarshaler) {
	end := len(*j) - 1
	ret, (*j) = (*j)[end], (*j)[:end]
	return
}
