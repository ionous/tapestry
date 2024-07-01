package frame

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/dl/game"
)

// marker interface for frame commands.
// ( assigned to those commands in the frame.idl )
type Notification interface{}

// Queues incoming events and collects text output between events.
type Collector struct {
	events []Notification
	// collects any text output since the last event
	// flushed before each event so that .events becomes a list of
	// [output], [event], [output], [event]
	buf strings.Builder
}

// returns and clears all events
func (out *Collector) GetEvents() (ret []Notification) {
	ret, out.events = out.flush(), nil
	return
}

func (out *Collector) onGameEvent(sig game.Signal) {
	out.flush()
	out.addEvent(&GameSignal{Signal: sig.String()})
}
func (out *Collector) onStartScene(domains []string) {
	out.flush()
	out.addEvent(&SceneStarted{Domains: domains})
}
func (out *Collector) onEndScene(domains []string) {
	out.flush()
	out.addEvent(&SceneEnded{Domains: domains})
}
func (out *Collector) onChangeState(noun, aspect, prev, trait string) {
	out.flush()
	out.addEvent(&StateChanged{Noun: noun, Aspect: aspect, Prev: prev, Trait: trait})
}
func (out *Collector) onChangeRel(a, b, rel string) {
	out.flush()
	out.addEvent(&PairChanged{A: a, B: b, Rel: rel})
}
func (out *Collector) addEvent(evt Notification) {
	out.events = append(out.events, evt)
}

func (out *Collector) flush() []Notification {
	if out.buf.Len() > 0 {
		str := out.buf.String()
		out.buf.Reset()
		out.addEvent(&FrameOutput{Text: str})
	}
	return out.events
}
