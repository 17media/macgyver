package keys

import (
	"fmt"
	"io"
)

// Type of the Keys implemented Keys
type Type string

const (
	TypeText Type = "text"
	TypeEnv       = "env"
	TypeFile      = "file"
)

// Get returns the Keys of a Type is exists
func Get(t Type) (Keys, error) {
	k, ok := keys[t]
	if !ok {
		return nil, fmt.Errorf("Without support %+v keys.Type", t)
	}
	return k, nil
}

// Key is one of parsed input
type Key struct {
	Key     string
	Value   string
	Secrets []*Secret
}

// Secret is one of the secrets of a key
type Secret struct {
	Text        string
	IsEncrypted bool
}

// Keys defines keys operations
type Keys interface {
	// Import parses the input into keys, the secrets of keys are parsed by secretTag.
	// The regexp pattern of secrets is `<%s>(.*?)</%s>|<%s>(.*?)</%s>`. Only one group might be captured.
	// If no secretTag is captured, the entire value of the key will be regarded as a plaintext secret. (i.e. Secret{Text:"Value of the Key", IsEncrypted: false})
	//
	// Example for secretTag:
	//   If secretTag is `kms`, the regexp pattern will be <kms>(.*?)</kms>|<KMS>(.*?)</KMS>.
	//   The first group is lower case of the secret_tag which means plaintext of the secret(i.e. all characters between `<kms>` and `</kms>`). If it is not empty string, `secret.IsEncrypted` will be set to false.
	//   The second group is upper case of the secret_tag which means ciphertext of the secret(i.e. all characters between `<KMS>` and `</KMS>`). If it is not empty string, `secret.IsEncrypted` will be set to true.
	Import(input []string, secretTag string) []Key

	// Export outputs the `keys` to a string formated by the `Keys Type` and writes the string to the `writeCloser`
	// The secrets in the `Key.Value` will be replaced by the current `Key.Secrets`.
	// If `secret.IsEncrypted` is false, the secret will be the `secret.Text` only.
	// If `secret.IsEncrypted` is true, the secret will be the <secretTag>secret.Text</secretTag>.
	Export(keys []Key, secretTag string, writeCloser io.WriteCloser) error
}
