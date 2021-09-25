package detailed

type detStack []detailedState

func (j detStack) top() detailedState {
	return j[len(j)-1]
}

func (j *detStack) push(m detailedState) {
	(*j) = append(*j, m)
}
func (j *detStack) pop() (ret detailedState) {
	end := len(*j) - 1
	ret, (*j) = (*j)[end], (*j)[:end]
	return
}
