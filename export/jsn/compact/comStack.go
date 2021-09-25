package compact

type comStack []compactState

func (j comStack) top() compactState {
	return j[len(j)-1]
}

func (j *comStack) push(m compactState) {
	(*j) = append(*j, m)
}
func (j *comStack) pop() (ret compactState) {
	end := len(*j) - 1
	ret, (*j) = (*j)[end], (*j)[:end]
	return
}
