package charm

// OnExit provides a callback when a state ends
func OnExit(name string, onExit func()) State {
	return Statement("on exit", func(rune) (none State) {
		onExit()
		return
	})
}

// for the very next rune, return nil ( unhandled )
var Terminal = Statement("terminal", func(rune) (none State) {
	return
})
