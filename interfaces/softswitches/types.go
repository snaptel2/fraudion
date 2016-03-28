package softswitches

// Softswitch ...
type Softswitch interface {
	GetSimultaneousCalls() (int, error)
	GenerateCall() error
}
