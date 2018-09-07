package keys


// Keys defines keys operations
type Keys interface {
	Import(string, string) ([]Key, error)
}

type Key struct {
	Key   string
	Value string
  IsEncrypted bool
}


// Type stores Crypto implementations
var Types = map[string]Keys{}

// Register stores Keys implementation's newFunc.
func Register(s string, c Keys) {
	Types[s] = c
}
