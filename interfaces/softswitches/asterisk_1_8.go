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
