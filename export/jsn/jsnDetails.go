package jsn

type jsnDetails []detailedState

func (j jsnDetails) top() detailedState {
	return j[len(j)-1]
}

func (j *jsnDetails) push(m detailedState) {
	(*j) = append(*j, m)
}
func (j *jsnDetails) pop() (ret detailedState) {
	end := len(*j) - 1
	ret, (*j) = (*j)[end], (*j)[:end]
	return
}
