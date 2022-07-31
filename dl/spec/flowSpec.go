package spec

func (op *FlowSpec) FriendlyLede(blockType *TypeSpec) (ret string) {
	if lede := op.Name; len(lede) > 0 {
		ret = FriendlyName(lede, false)
	} else {
		ret = FriendlyName(blockType.Name, false)
	}
	return
}
