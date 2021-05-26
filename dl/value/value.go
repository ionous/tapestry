package value

func (op *Text) Value() (ret string) {
	if s := string(*op); s != "$EMPTY" {
		ret = s
	}
	return
}
