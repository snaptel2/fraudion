package config

import (
	"time"
)

// FraudionConfig ...
type FraudionConfig struct {
	General      General
	Triggers     Triggers
	Actions      Actions
	ActionChains ActionChains
	Contacts     Contacts
}

// General ...
type General struct {
	MonitoredSoftware                     string
	CDRsSource                            string
	DefaultTriggerCheckPeriod             time.Duration
	DefaultHitThreshold                   uint32
	DefaultMinimumDestinationNumberLength uint32
	DefaultActionChainSleepPeriod         time.Duration
	DefaultActionChainRunCount            uint32
}

// Triggers ...
type Triggers struct {
	SimultaneousCalls     triggerSimultaneousCalls
	DangerousDestinations triggerDangerousDestinations
	ExpectedDestinations  triggerExpectedDestinations
	SmallDurationCalls    triggerSmallCallDurations
}

type triggerSimultaneousCalls struct {
	Enabled             bool
	CheckPeriod         time.Duration
	HitThreshold        uint32
	MinimumNumberLength uint32
}

type triggerDangerousDestinations struct {
	Enabled             bool
	CheckPeriod         time.Duration
	HitThreshold        uint32
	MinimumNumberLength uint32
	PrefixList          []string
}

type triggerExpectedDestinations struct {
	Enabled             bool
	CheckPeriod         time.Duration
	HitThreshold        uint32
	MinimumNumberLength uint32
	PrefixList          []string
}

type triggerSmallCallDurations struct {
	Enabled             bool
	CheckPeriod         time.Duration
	HitThreshold        uint32
	MinimumNumberLength uint32
	DurationThreshold   time.Duration
}

// Actions ...
type Actions struct {
	Email         actionEmail
	Call          actionCall
	HTTP          actionHTTP
	LocalCommands map[string]string
}

type actionEmail struct {
	Enabled        bool
	DefaultMessage string
}

type actionCall struct {
	Enabled        bool
	DefaultMessage string
}

type actionHTTP struct {
	Enabled           bool
	DefaultURL        string
	DefaultMethod     string
	DefaultParameters map[string]string
}

type actionLocalCommands struct {
	Enabled bool
	List    map[string]string
}

// ActionChains ...
type ActionChains struct {
	List map[string][]actionChainAction
}

type actionChainAction struct {
	Action   string
	Contacts []string
	Command  string
}

// Contacts ...
type Contacts struct {
	List map[string]contact
}

type contact struct {
	ForActions     []string
	PhoneNumber    string
	Email          string
	Message        string
	HTTPURL        string
	HTTPMethod     string
	HTTPParameters map[string]string
}

// FraudionConfigJSON ...
type FraudionConfigJSON struct {
	General      GeneralJSON
	Triggers     TriggersJSON
	Actions      ActionsJSON
	ActionChains ActionChainsJSON
	Contacts     ContactsJSON
}

// GeneralJSON ...
type GeneralJSON struct {
	MonitoredSoftware                     interface{}
	CDRsSource                            interface{}
	DefaultTriggerCheckPeriod             interface{}
	DefaultActionChainSleepPeriod         interface{}
	DefaultActionChainRunCount            interface{}
	DefaultMinimumDestinationNumberLength interface{}
}

// TriggersJSON ...
type TriggersJSON struct {
	SimultaneousCalls     interface{}
	DangerousDestinations interface{}
	ExpectedDestinations  interface{}
	SmallDurationCalls    interface{}
}

// ActionsJSON ...
type ActionsJSON struct {
	Email         interface{}
	Call          interface{}
	HTTP          interface{}
	LocalCommands interface{}
}

// ActionChainsJSON ...
type ActionChainsJSON struct {
	List interface{}
}

// ContactsJSON ...
type ContactsJSON struct {
	List interface{}
}
