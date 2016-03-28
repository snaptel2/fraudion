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
	Enabled                bool
	CheckPeriod            time.Duration
	HitThreshold           uint32
	MinimumNumberLength    uint32
	ActionChainName        string
	LastActionChainRunTime time.Time
}

type triggerDangerousDestinations struct {
	Enabled                bool
	CheckPeriod            time.Duration
	HitThreshold           uint32
	MinimumNumberLength    uint32
	PrefixList             []string
	ActionChainName        string
	LastActionChainRunTime time.Time
}

type triggerExpectedDestinations struct {
	Enabled                bool
	CheckPeriod            time.Duration
	HitThreshold           uint32
	MinimumNumberLength    uint32
	PrefixList             []string
	ActionChainName        string
	LastActionChainRunTime time.Time
}

type triggerSmallCallDurations struct {
	Enabled                bool
	CheckPeriod            time.Duration
	HitThreshold           uint32
	MinimumNumberLength    uint32
	DurationThreshold      time.Duration
	ActionChainName        string
	LastActionChainRunTime time.Time
}

// Actions ...
type Actions struct {
	Email         actionEmail
	Call          actionCall
	HTTP          actionHTTP
	LocalCommands actionLocalCommands
}

type actionEmail struct {
	Enabled  bool
	Method   string
	Username string
	Password string
}

type actionCall struct {
	Enabled         bool
	OriginateMethod string
}

type actionHTTP struct {
	Enabled bool
}

type actionLocalCommands struct {
	Enabled bool
}

// ActionChains ...
type ActionChains struct {
	List map[string][]actionChainAction
}

type actionChainAction struct {
	ActionName   string
	ContactNames []string
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
	DefaultHitThreshold                   interface{}
	DefaultMinimumDestinationNumberLength interface{}
	DefaultActionChainSleepPeriod         interface{}
	DefaultActionChainRunCount            interface{}
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
