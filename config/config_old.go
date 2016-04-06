package config

/*
import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"time"

	"github.com/andmar/fraudion/types"
	"github.com/andmar/fraudion/utils"
)

const (
	constMinimumTriggerCheckPeriod     = "1m"
	constMaximumActionChainSleepPeriod = "5m"
)

var (
	constSupportedSoftware               = []string{"*ast_elastix_2.3", "*ast_alone_1.8"}
	constSupportedCDRSources             = []string{"*db_mysql"}
	constSupportedEmailMethods           = []string{"*gmail"}
	constSupportedCallOriginationMethods = []string{"*ami"}
	constSupportedActions                = []string{"*email", "*call", "*http", "*localcommands"}
)

// CheckJSONSanityAndLoadConfigs ...
func (fraudionConfig *types.FraudionConfig) CheckJSONSanityAndLoadConfigs(ConfigsJSON *types.FraudionConfigJSON /*, validateOnly bool*/ /*) error {*/

// ** General Section
// * Monitored Software
/*if err := errorOnValueNotFound(ConfigsJSON.General.MonitoredSoftware == nil, "General", "MonitoredSoftware", "string", "\"nothing\""); err != nil {
		return utils.DebugLogAndGetError(err.Error(), true)
	}
	valueOfMonitoredSoftware, hasCorrectType := ConfigsJSON.General.MonitoredSoftware.(string)
	if err := errorOnIncorrectType(hasCorrectType, "General", "MonitoredSoftware", "string", reflect.TypeOf(ConfigsJSON.General.MonitoredSoftware)); err != nil {
		return utils.DebugLogAndGetError(err.Error(), true)
	}
	found := utils.StringInStringsSlice(valueOfMonitoredSoftware, constSupportedSoftware)
	if err := errorOnIncorrectValue(!found, "General", "MonitoredSoftware", fmt.Sprintf("One of: %s", constSupportedSoftware), valueOfMonitoredSoftware); err != nil {
		return utils.DebugLogAndGetError(err.Error(), true)
	}
	fraudionConfig.General.MonitoredSoftware = valueOfMonitoredSoftware

	// * CDRs Source
	if err := errorOnValueNotFound(ConfigsJSON.General.CDRsSource == nil, "General", "CDRsSource", "string", "\"nothing\""); err != nil {
		return utils.DebugLogAndGetError(err.Error(), true)
	}
	valueOfCDRsSource, hasCorrectType := ConfigsJSON.General.CDRsSource.(string)
	if err := errorOnIncorrectType(hasCorrectType, "General", "CDRsSource", "string", reflect.TypeOf(ConfigsJSON.General.CDRsSource)); err != nil {
		return utils.DebugLogAndGetError(err.Error(), true)
	}
	found = utils.StringInStringsSlice(valueOfCDRsSource, constSupportedCDRSources)
	if err := errorOnIncorrectValue(!found, "General", "CDRsSource", fmt.Sprintf("One of: %s", constSupportedCDRSources), valueOfCDRsSource); err != nil {
		return utils.DebugLogAndGetError(err.Error(), true)
	}
	fraudionConfig.General.CDRsSource = valueOfCDRsSource

	// * DefaultTriggerCheckPeriod
	if err := errorOnValueNotFound(ConfigsJSON.General.DefaultTriggerCheckPeriod == nil, "General", "DefaultTriggerCheckPeriod", "string", "\"nothing\""); err != nil {
		return utils.DebugLogAndGetError(err.Error(), true)
	}
	valueOfDefaultTriggerCheckPeriod, hasCorrectType := ConfigsJSON.General.DefaultTriggerCheckPeriod.(string)
	if err := errorOnIncorrectType(hasCorrectType, "General", "DefaultTriggerCheckPeriod", "string", reflect.TypeOf(ConfigsJSON.General.DefaultTriggerCheckPeriod)); err != nil {
		return utils.DebugLogAndGetError(err.Error(), true)
	}
	durationValueOfDefaultTriggerCheckPeriod, errDuration := time.ParseDuration(valueOfDefaultTriggerCheckPeriod)
	if err := errorOnIncorrectValue(errDuration != nil, "General", "DefaultTriggerCheckPeriod", "parseable Duration (e.g. 5\"s\", \"m\" or \"h\")", valueOfDefaultTriggerCheckPeriod); err != nil {
		return utils.DebugLogAndGetError(err.Error(), true)
	}
	durationValueOfConstMinimumTriggerCheckPeriod, errMinimumCheckPeriodDuration := time.ParseDuration(constMinimumTriggerCheckPeriod)
	errorMessageSuffix := fmt.Sprintf("Minimum value is \"%s\"", constMinimumTriggerCheckPeriod)
	if err := errorOnIncorrectValue(durationValueOfDefaultTriggerCheckPeriod < durationValueOfConstMinimumTriggerCheckPeriod, "General", "DefaultTriggerCheckPeriod", errorMessageSuffix, ""); err != nil {
		return utils.DebugLogAndGetError(err.Error(), true)
	}
	if errMinimumCheckPeriodDuration != nil {
		return utils.DebugLogAndGetError("Internal: There seems to be an issue with the definition of constMinimumTriggerCheckPeriod", true)
	}
	fraudionConfig.General.DefaultTriggerCheckPeriod = durationValueOfDefaultTriggerCheckPeriod

	// DefaultHitThreshold
	if err := errorOnValueNotFound(ConfigsJSON.General.DefaultHitThreshold == nil, "General", "DefaultHitThreshold", "int", "\"nothing\""); err != nil {
		return utils.DebugLogAndGetError(err.Error(), true)
	}
	valueOfDefaultHitThreshold, hasCorrectType := ConfigsJSON.General.DefaultHitThreshold.(float64)
	if err := errorOnIncorrectType(hasCorrectType, "General", "DefaultHitThreshold", "int", reflect.TypeOf(ConfigsJSON.General.DefaultHitThreshold)); err != nil {
		return utils.DebugLogAndGetError(err.Error(), true)
	}
	if err := errorOnIncorrectValue(valueOfDefaultHitThreshold < 0, "General", "DefaultHitThreshold", "Must be >= 0", strconv.FormatFloat(valueOfDefaultHitThreshold, 'f', 6, 64)); err != nil {
		return utils.DebugLogAndGetError(err.Error(), true)
	}
	fraudionConfig.General.DefaultHitThreshold = uint32(valueOfDefaultHitThreshold)

	// * DefaultMinimumDestinationNumberLength
	if err := errorOnValueNotFound(ConfigsJSON.General.DefaultMinimumDestinationNumberLength == nil, "General", "DefaultMinimumDestinationNumberLength", "string", "\"nothing\""); err != nil {
		return utils.DebugLogAndGetError(err.Error(), true)
	}
	valueOfDefaultMinimumDestinationNumberLength, hasCorrectType := ConfigsJSON.General.DefaultMinimumDestinationNumberLength.(float64)
	if err := errorOnIncorrectType(hasCorrectType, "General", "DefaultMinimumDestinationNumberLength", "int", reflect.TypeOf(ConfigsJSON.General.DefaultMinimumDestinationNumberLength)); err != nil {
		return utils.DebugLogAndGetError(err.Error(), true)
	}
	if err := errorOnIncorrectValue(valueOfDefaultMinimumDestinationNumberLength < 0, "General", "DefaultMinimumDestinationNumberLength", "Must be >= 0", strconv.FormatFloat(valueOfDefaultMinimumDestinationNumberLength, 'f', 6, 64)); err != nil {
		return utils.DebugLogAndGetError(err.Error(), true)
	}
	fraudionConfig.General.DefaultActionChainRunCount = uint32(valueOfDefaultMinimumDestinationNumberLength)

	// * DefaultActionChainSleepPeriod
	if err := errorOnValueNotFound(ConfigsJSON.General.DefaultActionChainSleepPeriod == nil, "General", "DefaultActionChainSleepPeriod", "string", "\"nothing\""); err != nil {
		return utils.DebugLogAndGetError(err.Error(), true)
	}
	valueOfDefaultActionChainSleepPeriod, hasCorrectType := ConfigsJSON.General.DefaultActionChainSleepPeriod.(string)
	if err := errorOnIncorrectType(hasCorrectType, "General", "DefaultActionChainSleePeriod", "string", reflect.TypeOf(ConfigsJSON.General.DefaultActionChainSleepPeriod)); err != nil {
		return utils.DebugLogAndGetError(err.Error(), true)
	}
	durationValueOfdefaultActionChainSleepPeriodDuration, errDuration := time.ParseDuration(valueOfDefaultActionChainSleepPeriod)
	if err := errorOnIncorrectValue(errDuration != nil, "General", "DefaultActionChainSleepPeriod", "parseable Duration (e.g. 5\"s\", \"m\" or \"h\")", valueOfDefaultActionChainSleepPeriod); err != nil {
		return utils.DebugLogAndGetError(err.Error(), true)
	}
	durationValueOfmaximumActionChainSleepPeriodDuration, errMaximumChainSleepPeriod := time.ParseDuration(constMaximumActionChainSleepPeriod)
	errorMessageSuffix = fmt.Sprintf("Maximum value is \"%s\"", constMaximumActionChainSleepPeriod)
	if err := errorOnIncorrectValue(durationValueOfdefaultActionChainSleepPeriodDuration > durationValueOfmaximumActionChainSleepPeriodDuration, "General", "DefaultActionChainSleepPeriod", errorMessageSuffix, ""); err != nil {
		return utils.DebugLogAndGetError(err.Error(), true)
	}
	if errMaximumChainSleepPeriod != nil {
		return utils.DebugLogAndGetError("Internal: There seems to be an issue with the definition of constMaximumActionChainSleepPeriod", true)
	}
	fraudionConfig.General.DefaultActionChainSleepPeriod = durationValueOfdefaultActionChainSleepPeriodDuration

	// * DefaultActionChainRunCount
	if err := errorOnValueNotFound(ConfigsJSON.General.DefaultActionChainRunCount == nil, "General", "DefaultActionChainRunCount", "string", "\"nothing\""); err != nil {
		return utils.DebugLogAndGetError(err.Error(), true)
	}
	valueOfDefaultActionChainRunCount, hasCorrectType := ConfigsJSON.General.DefaultActionChainRunCount.(float64)
	if err := errorOnIncorrectType(hasCorrectType, "General", "DefaultActionChainRunCount", "int", reflect.TypeOf(ConfigsJSON.General.DefaultActionChainRunCount)); err != nil {
		return utils.DebugLogAndGetError(err.Error(), true)
	}
	if err := errorOnIncorrectValue(valueOfDefaultActionChainRunCount <= 0, "General", "DefaultActionChainRunCount", "Must be > 0", strconv.FormatFloat(valueOfDefaultActionChainRunCount, 'f', 6, 64)); err != nil {
		return utils.DebugLogAndGetError(err.Error(), true)
	}
	fraudionConfig.General.DefaultActionChainRunCount = uint32(valueOfDefaultActionChainRunCount)

	// ** Triggers Section
	// * SimultaneousCalls
	if err := errorOnValueNotFound(ConfigsJSON.Triggers.SimultaneousCalls == nil, "Triggers", "SimultaneousCalls", "map[string]interface{}", "\"nothing\""); err != nil {

		utils.DebugLogAndGetError(err.Error(), true)
		fraudionConfig.Triggers.SimultaneousCalls.Enabled = false

	} else {

		mapOfTriggerSimultaneousCallsConfigJSON, hasCorrectType := ConfigsJSON.Triggers.SimultaneousCalls.(map[string]interface{})
		if err := errorOnIncorrectType(hasCorrectType, "Triggers", "SimultaneousCalls", "map[string]interface{}", reflect.TypeOf(ConfigsJSON.Triggers.SimultaneousCalls)); err != nil {
			return utils.DebugLogAndGetError(err.Error(), true)
		}

		valueEnabled, err := checkJSONSanityOfEnabledConfigs(mapOfTriggerSimultaneousCallsConfigJSON, "SimultaneousCalls")
		if err != nil {
			return utils.DebugLogAndGetError(err.Error(), true)
		}
		fraudionConfig.Triggers.SimultaneousCalls.Enabled = valueEnabled

		if valueEnabled {

			// TODO This is optional. So if this is not found use the default!
			durationValueOfCheckPeriod, err := checkJSONSanityOfCheckPeriodConfigs(fraudionConfig, mapOfTriggerSimultaneousCallsConfigJSON, "SimultaneousCalls", durationValueOfConstMinimumTriggerCheckPeriod)
			if err != nil {
				return utils.DebugLogAndGetError(err.Error(), true)
			}
			fraudionConfig.Triggers.SimultaneousCalls.CheckPeriod = durationValueOfCheckPeriod

			valueHitThreshold, err := checkJSONSanityOfHitThresholdConfigs(fraudionConfig, mapOfTriggerSimultaneousCallsConfigJSON, "SimultaneousCalls")
			if err != nil {
				return utils.DebugLogAndGetError(err.Error(), true)
			}
			fraudionConfig.Triggers.SimultaneousCalls.HitThreshold = valueHitThreshold

			// TODO This is optional. So if this is not found use the default!
			valueMinimumNumberLength, err := checkJSONSanityOfMinimumNumberLengthConfigs(fraudionConfig, mapOfTriggerSimultaneousCallsConfigJSON, "SimultaneousCalls")
			if err != nil {
				return utils.DebugLogAndGetError(err.Error(), true)
			}
			fraudionConfig.Triggers.SimultaneousCalls.MinimumNumberLength = valueMinimumNumberLength

		}

	}

	// * DangerousDestinations
	if err := errorOnValueNotFound(ConfigsJSON.Triggers.DangerousDestinations == nil, "Triggers", "DangerousDestinations", "map[string]interface{}", "\"nothing\""); err != nil {

		utils.DebugLogAndGetError(err.Error(), true)
		fraudionConfig.Triggers.DangerousDestinations.Enabled = false

	} else {

		mapOfTriggerDangerousDestinationsConfigJSON, hasCorrectType := ConfigsJSON.Triggers.DangerousDestinations.(map[string]interface{})
		if err := errorOnIncorrectType(hasCorrectType, "Triggers", "DangerousDestinations", "map[string]interface{}", reflect.TypeOf(ConfigsJSON.Triggers.SimultaneousCalls)); err != nil {
			return utils.DebugLogAndGetError(err.Error(), true)
		}

		valueEnabled, err := checkJSONSanityOfEnabledConfigs(mapOfTriggerDangerousDestinationsConfigJSON, "DangerousDestinations")
		if err != nil {
			return utils.DebugLogAndGetError(err.Error(), true)
		}
		fraudionConfig.Triggers.DangerousDestinations.Enabled = valueEnabled

		if valueEnabled {

			durationValueOfCheckPeriod, err := checkJSONSanityOfCheckPeriodConfigs(fraudionConfig, mapOfTriggerDangerousDestinationsConfigJSON, "DangerousDestinations", durationValueOfConstMinimumTriggerCheckPeriod)
			if err != nil {
				return utils.DebugLogAndGetError(err.Error(), true)
			}
			fraudionConfig.Triggers.DangerousDestinations.CheckPeriod = durationValueOfCheckPeriod

			valueHitThreshold, err := checkJSONSanityOfHitThresholdConfigs(fraudionConfig, mapOfTriggerDangerousDestinationsConfigJSON, "DangerousDestinations")
			if err != nil {
				return utils.DebugLogAndGetError(err.Error(), true)
			}
			fraudionConfig.Triggers.DangerousDestinations.HitThreshold = valueHitThreshold

			valueMinimumNumberLength, err := checkJSONSanityOfMinimumNumberLengthConfigs(fraudionConfig, mapOfTriggerDangerousDestinationsConfigJSON, "DangerousDestinations")
			if err != nil {
				return utils.DebugLogAndGetError(err.Error(), true)
			}
			fraudionConfig.Triggers.DangerousDestinations.MinimumNumberLength = valueMinimumNumberLength

			valuePrefixList, err := checkJSONSanityOfPrefixListConfigs(fraudionConfig, mapOfTriggerDangerousDestinationsConfigJSON, "DangerousDestinations")
			if err != nil {
				return err
			}
			fraudionConfig.Triggers.DangerousDestinations.PrefixList = valuePrefixList

		}

	}

	// * ExpectedDestinations
	if err := errorOnValueNotFound(ConfigsJSON.Triggers.ExpectedDestinations == nil, "Triggers", "ExpectedDestinations", "map[string]interface{}", "\"nothing\""); err != nil {

		utils.DebugLogAndGetError(err.Error(), true)
		fraudionConfig.Triggers.ExpectedDestinations.Enabled = false

	} else {

		mapOfTriggerExpectedDestinationsConfigJSON, hasCorrectType := ConfigsJSON.Triggers.ExpectedDestinations.(map[string]interface{})
		if err := errorOnIncorrectType(hasCorrectType, "Triggers", "ExpectedDestinations", "map[string]interface{}", reflect.TypeOf(ConfigsJSON.Triggers.ExpectedDestinations)); err != nil {
			return utils.DebugLogAndGetError(err.Error(), true)
		}

		valueEnabled, err := checkJSONSanityOfEnabledConfigs(mapOfTriggerExpectedDestinationsConfigJSON, "ExpectedDestinations")
		if err != nil {
			return utils.DebugLogAndGetError(err.Error(), true)
		}
		fraudionConfig.Triggers.ExpectedDestinations.Enabled = valueEnabled

		if valueEnabled {

			durationValueOfCheckPeriod, err := checkJSONSanityOfCheckPeriodConfigs(fraudionConfig, mapOfTriggerExpectedDestinationsConfigJSON, "ExpectedDestinations", durationValueOfConstMinimumTriggerCheckPeriod)
			if err != nil {
				return utils.DebugLogAndGetError(err.Error(), true)
			}
			fraudionConfig.Triggers.ExpectedDestinations.CheckPeriod = durationValueOfCheckPeriod

			valueHitThreshold, err := checkJSONSanityOfHitThresholdConfigs(fraudionConfig, mapOfTriggerExpectedDestinationsConfigJSON, "ExpectedDestinations")
			if err != nil {
				return utils.DebugLogAndGetError(err.Error(), true)
			}
			fraudionConfig.Triggers.ExpectedDestinations.HitThreshold = valueHitThreshold

			valueMinimumNumberLength, err := checkJSONSanityOfMinimumNumberLengthConfigs(fraudionConfig, mapOfTriggerExpectedDestinationsConfigJSON, "ExpectedDestinations")
			if err != nil {
				return utils.DebugLogAndGetError(err.Error(), true)
			}
			fraudionConfig.Triggers.ExpectedDestinations.MinimumNumberLength = valueMinimumNumberLength

			valuePrefixList, err := checkJSONSanityOfPrefixListConfigs(fraudionConfig, mapOfTriggerExpectedDestinationsConfigJSON, "ExpectedDestinations")
			if err != nil {
				return err
			}
			fraudionConfig.Triggers.ExpectedDestinations.PrefixList = valuePrefixList

		}

	}

	// * SmallDurationCalls
	if err := errorOnValueNotFound(ConfigsJSON.Triggers.SmallDurationCalls == nil, "Triggers", "SmallDurationCalls", "map[string]interface{}", "\"nothing\""); err != nil {

		utils.DebugLogAndGetError(err.Error(), true)
		fraudionConfig.Triggers.SmallDurationCalls.Enabled = false

	} else {

		mapOfTriggerSmallDurationCallsConfigJSON, hasCorrectType := ConfigsJSON.Triggers.SmallDurationCalls.(map[string]interface{})
		if err := errorOnIncorrectType(hasCorrectType, "Triggers", "SmallDurationCalls", "map[string]interface{}", reflect.TypeOf(ConfigsJSON.Triggers.SmallDurationCalls)); err != nil {
			return utils.DebugLogAndGetError(err.Error(), true)
		}

		valueEnabled, err := checkJSONSanityOfEnabledConfigs(mapOfTriggerSmallDurationCallsConfigJSON, "SmallDurationCalls")
		if err != nil {
			return utils.DebugLogAndGetError(err.Error(), true)
		}
		fraudionConfig.Triggers.SmallDurationCalls.Enabled = valueEnabled

		if valueEnabled {

			durationValueOfCheckPeriod, err := checkJSONSanityOfCheckPeriodConfigs(fraudionConfig, mapOfTriggerSmallDurationCallsConfigJSON, "SmallDurationCalls", durationValueOfConstMinimumTriggerCheckPeriod)
			if err != nil {
				return utils.DebugLogAndGetError(err.Error(), true)
			}
			fraudionConfig.Triggers.SmallDurationCalls.CheckPeriod = durationValueOfCheckPeriod

			valueHitThreshold, err := checkJSONSanityOfHitThresholdConfigs(fraudionConfig, mapOfTriggerSmallDurationCallsConfigJSON, "SmallDurationCalls")
			if err != nil {
				return utils.DebugLogAndGetError(err.Error(), true)
			}
			fraudionConfig.Triggers.SmallDurationCalls.HitThreshold = valueHitThreshold

			valueMinimumNumberLength, err := checkJSONSanityOfMinimumNumberLengthConfigs(fraudionConfig, mapOfTriggerSmallDurationCallsConfigJSON, "SmallDurationCalls")
			if err != nil {
				return utils.DebugLogAndGetError(err.Error(), true)
			}
			fraudionConfig.Triggers.SmallDurationCalls.MinimumNumberLength = valueMinimumNumberLength

			if err := errorOnValueNotFound(!utils.StringKeyInMap("DurationThreshold", mapOfTriggerSmallDurationCallsConfigJSON), "Triggers", "SmallDurationCalls/DurationThreshold", "map[string]interface{}", "\"nothing\""); err != nil {
				// TODO If this trigger is enabled, this value should be mandatory!
				utils.DebugLogAndGetError(err.Error(), true)
			}

			valueDurationString, hasCorrectType := mapOfTriggerSmallDurationCallsConfigJSON["DurationThreshold"].(string)
			if err := errorOnIncorrectType(hasCorrectType, "Triggers", "SmallDurationCalls/DurationThreshold", "time.Duration", reflect.TypeOf(mapOfTriggerSmallDurationCallsConfigJSON["DurationThreshold"])); err != nil {
				return utils.DebugLogAndGetError(err.Error(), true)
			}

			valueDuration, errDuration := time.ParseDuration(valueDurationString)
			if err := errorOnIncorrectValue(errDuration != nil, "Triggers", "SmallDurationCalls/DurationThreshold", "Must be a valid duration", ""); err != nil {
				return utils.DebugLogAndGetError(err.Error(), true)
			}

			fraudionConfig.Triggers.SmallDurationCalls.DurationThreshold = valueDuration

		}

	}

	// ** Actions Section
	// * Email
	if err := errorOnValueNotFound(ConfigsJSON.Actions.Email == nil, "Actions", "Email", "map[string]interface{}", "\"nothing\""); err != nil {

		utils.DebugLogAndGetError(err.Error(), true)
		fraudionConfig.Actions.Email.Enabled = false

	} else {

		mapOfActionEmail, hasCorrectType := ConfigsJSON.Actions.Email.(map[string]interface{})
		if err := errorOnIncorrectType(hasCorrectType, "Actions", "Email", "map[string]interface{}", reflect.TypeOf(ConfigsJSON.Actions.Email)); err != nil {
			return utils.DebugLogAndGetError(err.Error(), true)
		}

		valueEnabled, err := checkJSONSanityOfEnabledConfigs(mapOfActionEmail, "Email")
		if err != nil {
			return utils.DebugLogAndGetError(err.Error(), true)
		}
		fraudionConfig.Actions.Email.Enabled = valueEnabled

		if valueEnabled {

			// Method
			if err := errorOnValueNotFound(!utils.StringKeyInMap("Method", mapOfActionEmail), "Actions", "Email/Method", "map[string]interface{}", "\"nothing\""); err != nil {
				return utils.DebugLogAndGetError(err.Error(), true)
			}
			valueOfMethod, hasCorrectType := mapOfActionEmail["Method"].(string)
			if err := errorOnIncorrectType(hasCorrectType, "Actions", "Email/Method", "string", reflect.TypeOf(mapOfActionEmail["Method"])); err != nil {
				return utils.DebugLogAndGetError(err.Error(), true)
			}
			if err := errorOnIncorrectValue(!utils.StringInStringsSlice(valueOfMethod, constSupportedEmailMethods), "Actions", "Email/Method", fmt.Sprintf("One of: %s", constSupportedEmailMethods), valueOfMethod); err != nil {
				return utils.DebugLogAndGetError(err.Error(), true)
			}
			fraudionConfig.Actions.Email.Method = valueOfMethod

			// Username
			if err := errorOnValueNotFound(!utils.StringKeyInMap("Username", mapOfActionEmail), "Actions", "Email/Username", "map[string]interface{}", "\"nothing\""); err != nil {
				return utils.DebugLogAndGetError(err.Error(), true)
			}
			valueOfUsername, hasCorrectType := mapOfActionEmail["Username"].(string)
			if err := errorOnIncorrectType(hasCorrectType, "Actions", "Email/Username", "string", reflect.TypeOf(mapOfActionEmail["Username"])); err != nil {
				return utils.DebugLogAndGetError(err.Error(), true)
			}
			fraudionConfig.Actions.Email.Username = valueOfUsername

			// Password
			if err := errorOnValueNotFound(!utils.StringKeyInMap("Password", mapOfActionEmail), "Actions", "Email/Password", "map[string]interface{}", "\"nothing\""); err != nil {
				return utils.DebugLogAndGetError(err.Error(), true)
			}
			valueOfPassword, hasCorrectType := mapOfActionEmail["Password"].(string)
			if err := errorOnIncorrectType(hasCorrectType, "Actions", "Email/Password", "string", reflect.TypeOf(mapOfActionEmail["Password"])); err != nil {
				return utils.DebugLogAndGetError(err.Error(), true)
			}
			fraudionConfig.Actions.Email.Password = valueOfPassword

		}

	}

	// * HTTP
	if err := errorOnValueNotFound(ConfigsJSON.Actions.HTTP == nil, "Actions", "HTTP", "map[string]interface{}", "\"nothing\""); err != nil {

		utils.DebugLogAndGetError(err.Error(), true)
		fraudionConfig.Actions.HTTP.Enabled = false

	} else {

		mapOfActionHTTP, hasCorrectType := ConfigsJSON.Actions.HTTP.(map[string]interface{})
		if err := errorOnIncorrectType(hasCorrectType, "Actions", "HTTP", "map[string]interface{}", reflect.TypeOf(ConfigsJSON.Actions.HTTP)); err != nil {
			return utils.DebugLogAndGetError(err.Error(), true)
		}

		valueEnabled, err := checkJSONSanityOfEnabledConfigs(mapOfActionHTTP, "HTTP")
		if err != nil {
			return utils.DebugLogAndGetError(err.Error(), true)
		}
		fraudionConfig.Actions.HTTP.Enabled = valueEnabled

	}

	// * Call
	if err := errorOnValueNotFound(ConfigsJSON.Actions.Call == nil, "Actions", "Call", "map[string]interface{}", "\"nothing\""); err != nil {

		utils.DebugLogAndGetError(err.Error(), true)
		fraudionConfig.Actions.Call.Enabled = false

	} else {

		mapOfActionCall, hasCorrectType := ConfigsJSON.Actions.Call.(map[string]interface{})
		if err := errorOnIncorrectType(hasCorrectType, "Actions", "Email", "map[string]interface{}", reflect.TypeOf(ConfigsJSON.Actions.Call)); err != nil {
			return utils.DebugLogAndGetError(err.Error(), true)
		}

		valueEnabled, err := checkJSONSanityOfEnabledConfigs(mapOfActionCall, "Email")
		if err != nil {
			return utils.DebugLogAndGetError(err.Error(), true)
		}
		fraudionConfig.Actions.Call.Enabled = valueEnabled

		if valueEnabled {

			// OriginatedMethod [Mandatory]
			if err := errorOnValueNotFound(!utils.StringKeyInMap("OriginateMethod", mapOfActionCall), "Actions", "Call", "map[string]interface{}", "\"nothing\""); err != nil {
				return utils.DebugLogAndGetError(err.Error(), true)
			}
			valueOfOriginateMethod, hasCorrectType := mapOfActionCall["OriginateMethod"].(string)
			if err := errorOnIncorrectType(hasCorrectType, "Actions", "Call/OriginatedMethod", "string", reflect.TypeOf(mapOfActionCall["Method"])); err != nil {
				return utils.DebugLogAndGetError(err.Error(), true)
			}
			if err := errorOnIncorrectValue(!utils.StringInStringsSlice(valueOfOriginateMethod, constSupportedCallOriginationMethods), "Actions", "Call/Method", fmt.Sprintf("One of: %s", constSupportedCallOriginationMethods), valueOfOriginateMethod); err != nil {
				return utils.DebugLogAndGetError(err.Error(), true)
			}
			fraudionConfig.Actions.Call.OriginateMethod = valueOfOriginateMethod

		}

	}

	// * LocalCommands
	if err := errorOnValueNotFound(ConfigsJSON.Actions.LocalCommands == nil, "Actions", "LocalCommands", "map[string]interface{}", "\"nothing\""); err != nil {

		utils.DebugLogAndGetError(err.Error(), true)
		fraudionConfig.Actions.LocalCommands.Enabled = false

	} else {

		mapOfActionLocalCommands, hasCorrectType := ConfigsJSON.Actions.LocalCommands.(map[string]interface{})
		if err := errorOnIncorrectType(hasCorrectType, "Actions", "LocalCommands", "map[string]interface{}", reflect.TypeOf(ConfigsJSON.Actions.LocalCommands)); err != nil {
			return utils.DebugLogAndGetError(err.Error(), true)
		}

		valueEnabled, err := checkJSONSanityOfEnabledConfigs(mapOfActionLocalCommands, "LocalCommands")
		if err != nil {
			return utils.DebugLogAndGetError(err.Error(), true)
		}
		fraudionConfig.Actions.LocalCommands.Enabled = valueEnabled

	}

	// ** ActionChains
	// * List
	if err := errorOnValueNotFound(ConfigsJSON.ActionChains.List == nil, "ActionChains", "List", "interface{}", "\"nothing\""); err != nil {

		utils.DebugLogAndGetError(err.Error(), true)
		fraudionConfig.ActionChains.List = nil // Returns empty map, no ActionChains

	} else {

		mapOfActionChainsList, hasCorrectType := ConfigsJSON.ActionChains.List.(map[string]interface{})
		if err := errorOnIncorrectType(hasCorrectType, "ActionChains", "List", "map[string]interface{}", reflect.TypeOf(ConfigsJSON.ActionChains.List)); err != nil {
			return utils.DebugLogAndGetError(err.Error(), true)
		}

		list := make(map[string][]actionChainAction)

		for actionChainValueKey, actionChainValue := range mapOfActionChainsList {

			var actionsList []actionChainAction

			for _, actionValue := range actionChainValue.([]interface{}) {

				// TODO Type assertions missing here!

				action := new(actionChainAction)

				mapOfActionValue := actionValue.(map[string]interface{})

				action.ActionName = mapOfActionValue["Action"].(string)

				var contactNames []string
				for _, contactValue := range mapOfActionValue["Contacts"].([]interface{}) {
					contactNames = append(contactNames, contactValue.(string))
				}
				action.ContactNames = contactNames

				actionsList = append(actionsList, *action)

			}

			list[actionChainValueKey] = actionsList

		}

		fraudionConfig.ActionChains.List = list

	}

	// ** Contacts
	// * List
	if err := errorOnValueNotFound(ConfigsJSON.ActionChains.List == nil, "Contacts", "List", "interface{}", "\"nothing\""); err != nil {

		utils.DebugLogAndGetError(err.Error(), true)
		fraudionConfig.Contacts.List = nil // Returns empty map, no ActionChains

	} else {

		mapOfContactsList, hasCorrectType := ConfigsJSON.Contacts.List.(map[string]interface{})
		if err := errorOnIncorrectType(hasCorrectType, "Contacts", "List", "map[string]interface{}", reflect.TypeOf(ConfigsJSON.Contacts.List)); err != nil {
			return utils.DebugLogAndGetError(err.Error(), true)
		}

		list := make(map[string]contact)

		for contactName, contactData := range mapOfContactsList {

			// TODO Type assertions missing here!

			var newContact contact

			for contactDataKey, contactDataValue := range contactData.(map[string]interface{}) {

				switch contactDataKey {

				// TODO Other fields (cases for contactDataKey) missing here!

				case "forActions":
					var actionNames []string

					for _, actionName := range contactDataValue.([]interface{}) {
						actionNames = append(actionNames, actionName.(string))
					}
					newContact.ForActions = actionNames

				}

			}

			list[contactName] = newContact

		}

		fraudionConfig.Contacts.List = list

	}

	return nil
}

func errorOnIncorrectType(condition bool, sectionName string, valueNameOrPath string, typeExpected string, typeFound interface{}) error {
	if !condition {
		errorMessage := fmt.Sprintf("Value of \"%s\" in Section \"%s\" has an incorrect type (expected %s, found %s)", valueNameOrPath, sectionName, typeExpected, typeFound)
		return fmt.Errorf("%s", errorMessage)
	}
	return nil
}

func errorOnIncorrectValue(condition bool, sectionName string, valueNameOrPath string, valueExpected string, valueFound string) error {
	if condition {
		errorMessage := fmt.Sprintf("Value of \"%s\" in Section \"%s\" is incorrect (expected %s, found %s)", valueNameOrPath, sectionName, valueExpected, valueFound)
		return fmt.Errorf("%s", errorMessage)
	}
	return nil
}

func errorOnValueNotFound(condition bool, sectionName string, valueNameOrPath string, valueExpected string, valueFound string) error {
	if condition {
		errorMessage := fmt.Sprintf("Value of \"%s\" in Section \"%s\" not found", valueNameOrPath, sectionName)
		return fmt.Errorf("%s", errorMessage)
	}
	return nil
}

func checkJSONSanityOfEnabledConfigs(JSONConfig map[string]interface{}, sectionName string) (bool, error) {

	if utils.StringKeyInMap("Enabled", JSONConfig) {

		value, hasCorrectType := JSONConfig["Enabled"].(bool)
		if err := errorOnIncorrectType(hasCorrectType, sectionName, "Enabled", "bool", reflect.TypeOf(JSONConfig["enabled"])); err != nil {
			return false, err
		}

		return value, nil

	}

	// TODO:LOG Log this to Syslog
	assume := false
	fmt.Printf("WARNING: \"Enabled\" value for section \"%s\" not found, assuming %t\n", sectionName, assume)

	return assume, nil

}

func checkJSONSanityOfCheckPeriodConfigs(fraudionConfig *types.FraudionConfig, simCallsTriggerJSONConfig map[string]interface{}, triggerName string, minimumTriggerCheckPeriodDuration time.Duration) (time.Duration, error) {

	if utils.StringKeyInMap("CheckPeriod", simCallsTriggerJSONConfig) {

		valueString, hasCorrectType := simCallsTriggerJSONConfig["CheckPeriod"].(string)
		if err := errorOnIncorrectType(hasCorrectType, triggerName, "CheckPeriod", "string", reflect.TypeOf(simCallsTriggerJSONConfig["check_period"])); err != nil {
			return 0, err
		}

		valueDuration, errDuration := time.ParseDuration(valueString)
		if err := errorOnIncorrectValue(errDuration != nil, triggerName, "CheckPeriod", "Must be a valid duration", valueString); err != nil {
			return 0, err
		}

		errorMessageSuffix := fmt.Sprintf("Minimum value is \"%s\"", constMinimumTriggerCheckPeriod)
		if err := errorOnIncorrectValue(valueDuration < minimumTriggerCheckPeriodDuration, triggerName, "check_period", errorMessageSuffix, valueString); err != nil {
			return 0, err
		}

		return valueDuration, nil

	}

	// TODO:LOG Log this to Syslog
	fmt.Printf("WARNING: \"CheckPeriod\" value for Trigger \"%s\" not found, assuming %s\n", triggerName, fraudionConfig.General.DefaultTriggerCheckPeriod)
	return fraudionConfig.General.DefaultTriggerCheckPeriod, nil

}

func checkJSONSanityOfMinimumNumberLengthConfigs(fraudionConfig *types.FraudionConfig, simCallsTriggerJSONConfig map[string]interface{}, triggerName string) (uint32, error) {

	if utils.StringKeyInMap("MinimumNumberLength", simCallsTriggerJSONConfig) {

		value, hasCorrectType := simCallsTriggerJSONConfig["MinimumNumberLength"].(float64)
		if err := errorOnIncorrectType(hasCorrectType, triggerName, "MinimumNumberLength", "float64", reflect.TypeOf(simCallsTriggerJSONConfig["minimum_number_length"])); err != nil {
			return 0, err
		}

		if err := errorOnIncorrectValue(value <= 0, triggerName, "MinimumNumberLength", "Must be > 0", strconv.FormatFloat(value, 'f', 6, 64)); err != nil {
			return 0, err
		}

		return uint32(value), nil

	}

	// This value is mandatory!

	// TODO:LOG Log this to Syslog
	fmt.Printf("WARNING: \"MinimumNumberLength\" value for \"%s\" Trigger not found, assuming %d\n", triggerName, fraudionConfig.General.DefaultMinimumDestinationNumberLength)
	return fraudionConfig.General.DefaultMinimumDestinationNumberLength, nil

}

func checkJSONSanityOfHitThresholdConfigs(fraudionConfig *types.FraudionConfig, simCallsTriggerJSONConfig map[string]interface{}, triggerName string) (uint32, error) {

	if utils.StringKeyInMap("HitThreshold", simCallsTriggerJSONConfig) {

		value, hasCorrectType := simCallsTriggerJSONConfig["HitThreshold"].(float64)
		if err := errorOnIncorrectType(hasCorrectType, triggerName, "HitThreshold", "float64", reflect.TypeOf(simCallsTriggerJSONConfig["max_call_threshold"])); err != nil {
			return 0, err
		}

		if err := errorOnIncorrectValue(value <= 0, triggerName, "HitThreshold", "integer > 0", strconv.FormatFloat(value, 'f', 6, 64)); err != nil {
			return 0, err
		}

		return uint32(value), nil

	}

	return 0, fmt.Errorf("\"HitThreshold\" value for Trigger \"%s\" not found", triggerName)

}

func checkJSONSanityOfPrefixListConfigs(fraudionConfig *FraudionConfig, simCallsTriggerJSONConfig map[string]interface{}, triggerName string) ([]string, error) {

	if utils.StringKeyInMap("PrefixList", simCallsTriggerJSONConfig) {

		sliceOfInterface := simCallsTriggerJSONConfig["PrefixList"]

		// WARNING: Just checks if it's a slice basically (of anything really because interface{}) (verification of the type of the items of the slice is done below)
		valueInterfaceSlice, hasCorrectType := sliceOfInterface.([]interface{})
		if err := errorOnIncorrectType(hasCorrectType, triggerName, "PrefixList", "[]string", reflect.TypeOf(sliceOfInterface)); err != nil {
			return *new([]string), err
		}

		sliceReflectValue := reflect.ValueOf(valueInterfaceSlice)
		var finalSlice []string
		for i := 0; i < sliceReflectValue.Len(); i++ {

			// TODO This is not needed because the previous check catches this, because ints implement no interface? Don't know but it catches when the slice has ints...
			valueOfSliceItem, hasCorrectType := sliceReflectValue.Index(i).Interface().(string)
			if err := errorOnIncorrectType(hasCorrectType, triggerName, "PrefixList", "[]", reflect.TypeOf(sliceOfInterface)); err != nil {
				return *new([]string), err
			}

			isNumericString, errCheckingValue := regexp.MatchString("^[0-9]+$", valueOfSliceItem)
			if errCheckingValue != nil {
				//fmt.Printf("ERROR: Internal: There seems to be an issue with checking if the \"prefix_list\" string values are numeric\n")
				return *new([]string), fmt.Errorf("internal: there seems to be an issue with checking if the \"prefix_list\" string values are numeric")
			}
			if err := errorOnIncorrectValue(!isNumericString, triggerName, "PrefixList", "Values must be Numeric Strings", valueOfSliceItem); err != nil {
				return *new([]string), err
			}

			finalSlice = append(finalSlice, valueOfSliceItem)
		}

		return finalSlice, nil

	}

	return *new([]string), fmt.Errorf("\"PrefixList\" value for Trigger \"%s\" not found", triggerName)

}
*/
