package detailed

type detStack []detailedState

func (j *detStack) push(m detailedState) {
	(*j) = append(*j, m)
}

// returns the removed state
func (j *detStack) pop() (ret detailedState) {
	end := len(*j) - 1
	ret, (*j) = (*j)[end], (*j)[:end]
	return
}
