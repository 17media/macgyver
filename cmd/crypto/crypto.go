package crypto

// Crypto defines crypto operations
type Crypto interface {
	Encrypt([]byte) ([]byte, error)
	Decrypt([]byte) ([]byte, error)
}

// ProviderNewFunc creates
var ProviderNewFunc = map[string]func() Crypto{}

// Providers stores Crypto implementations
var Providers = map[string]Crypto{}

// Register stores Crypto implementation's newFunc.
func Register(s string, f func() Crypto) {
	ProviderNewFunc[s] = f
}

// Init creates specified Crypto implementation and stores it in Providers
func Init(cryptoProvider string) {
	f, ok := ProviderNewFunc[cryptoProvider]
	if !ok {
		panic("Without support " + cryptoProvider + " encrypt")
	}
	Providers[cryptoProvider] = f()
}
