package chart

type chartStack []State

func (j *chartStack) push(m State) {
	(*j) = append(*j, m)
}

// returns the removed state
func (j *chartStack) pop() (ret State) {
	end := len(*j) - 1
	ret, (*j) = (*j)[end], (*j)[:end]
	return
}
