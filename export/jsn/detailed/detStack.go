package detailed

type detStack []detailedMarshaler

func (j *detStack) push(m detailedMarshaler) {
	(*j) = append(*j, m)
}

// returns the removed state
func (j *detStack) pop() (ret detailedMarshaler) {
	end := len(*j) - 1
	ret, (*j) = (*j)[end], (*j)[:end]
	return
}
