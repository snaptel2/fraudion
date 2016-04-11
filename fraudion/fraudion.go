package fraudion

import (
	"time"

	"github.com/andmar/fraudion/config"
)

func init() {
	Global = new(Fraudion)
	Global.State = new(State)
}

// Global ...
var Global *Fraudion

// Fraudion The "Type"
type Fraudion struct {
	StartUpTime time.Time
	Configs     *config.FraudionConfig
	State       *State
}

// State ...
type State struct {
	Triggers StateTriggers
}

// StateTriggers ...
type StateTriggers struct {
	StateSimultaneousCalls     stateSimultaneousCalls
	StateDangerousDestinations stateDangerousDestinations
	StateExpectedDestinations  stateExpectedDestinations
	StateSmallDurationCalls    stateSmallCallDurations
}

type stateSimultaneousCalls struct {
	LastActionChainRunTime time.Time
	ActionChainRunCount    uint32
}

type stateDangerousDestinations struct {
	LastActionChainRunTime time.Time
	ActionChainRunCount    uint32
}

type stateExpectedDestinations struct {
	LastActionChainRunTime time.Time
	ActionChainRunCount    uint32
}

type stateSmallCallDurations struct {
	LastActionChainRunTime time.Time
	ActionChainRunCount    uint32
}
