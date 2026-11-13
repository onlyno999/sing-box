package adapter

import (
	"sync"
	"sync/atomic"
)

// Dial health, NekoBox++ fork addition: for every outbound tag the
// core keeps the count of CONSECUTIVE failed dials made by real
// routed traffic. A dial that fails means a connection tried to
// go out through that outbound and its handshake never came back
// — the node behind the outbound did not answer. Any successful
// dial resets the count to zero, so isolated failures aimed at a
// dead target site (the outbound itself worked) never accumulate:
// on a busy line successes keep landing in between.
//
// The client app reads the count for the main line's tag through
// libcore instead of actively probing the node in use: the real
// traffic already knows whether the node answers.

type dialHealthState struct {
	consecutiveFailures atomic.Int64
}

var dialHealthRegistry sync.Map // string -> *dialHealthState

func dialHealthFor(tag string) *dialHealthState {
	if state, loaded := dialHealthRegistry.Load(tag); loaded {
		return state.(*dialHealthState)
	}
	state, _ := dialHealthRegistry.LoadOrStore(tag, &dialHealthState{})
	return state.(*dialHealthState)
}

// RecordDialResult notes the outcome of one real dial through the
// outbound with the given tag.
func RecordDialResult(tag string, success bool) {
	if tag == "" {
		return
	}
	state := dialHealthFor(tag)
	if success {
		state.consecutiveFailures.Store(0)
	} else {
		state.consecutiveFailures.Add(1)
	}
}

// DialConsecutiveFailures returns how many dials through the
// outbound with the given tag have failed in a row.
func DialConsecutiveFailures(tag string) int64 {
	if state, loaded := dialHealthRegistry.Load(tag); loaded {
		return state.(*dialHealthState).consecutiveFailures.Load()
	}
	return 0
}

// ResetDialHealth restarts the failure count for a tag. The app
// calls it when the member behind a line tag changes, so the new
// member is judged on its own dials only.
func ResetDialHealth(tag string) {
	dialHealthFor(tag).consecutiveFailures.Store(0)
}
