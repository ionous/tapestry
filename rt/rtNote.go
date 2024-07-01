package rt

// Callbacks to listen for system level changes
type Notifier struct {
	StartedScene    func(domains []string)
	EndedScene      func(domains []string)
	ChangedState    func(noun, aspect, oldState, newState string)
	ChangedRelative func(a, b, rel string)
}
