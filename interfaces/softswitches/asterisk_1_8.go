package softswitches

// Asterisk1_8 ...
type Asterisk1_8 struct{}

// GetSimultaneousCalls ...
func (asterisk *Asterisk1_8) GetSimultaneousCalls() (int, error) {
	return 0, nil
}

// GenerateCall ...
func (asterisk *Asterisk1_8) GenerateCall() error {
	return nil
}

// GetCDRColumnNames ...
func (asterisk *Asterisk1_8) GetCDRColumnNames() ([]string, error) {
	return *new([]string), nil
}

// NumberInDialString ...
func (asterisk *Asterisk1_8) NumberInDialString() (bool, error) {

	// Dial String Formats:
	// Asterisk ?: Dial(type1/identifier1[&type2/identifier2[&type3/identifier3... ] ], timeout, options, URL) Source: http://www.voip-info.org/wiki/view/Asterisk+cmd+Dial
	// Asterisk 8: Dial(Technology/Resource&[Technology2/Resource2[&...]],[timeout,[options,[URL]]]) Source: https://wiki.asterisk.org/wiki/display/AST/Application_Dial

	return false, nil
}
