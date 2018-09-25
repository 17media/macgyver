package keys

import (
	"fmt"
	"regexp"
	"strings"
)

type Type string

const (
	TypeText Type = "text"
	TypeEnv       = "env"
)

var (
	keys = map[Type]Keys{
		TypeText: &flagsKeys{},
		TypeEnv:  &envsKeys{},
	}
)

// Keys defines keys operations
type Keys interface {
	Import(input []string, prefix string) []Key
	Export(keys []Key)
}

type Key struct {
	Key         string
	Value       string
	IsEncrypted bool
}

func Get(t Type) (Keys, error) {
	k, ok := keys[t]
	if !ok {
		return nil, fmt.Errorf("Without support %+v keys.Type", t)
	}
	return k, nil
}

type envsKeys struct {
}

func (e *envsKeys) Import(input []string, prefix string) []Key {
	var k []Key
	decryptEnvFlagRegexp := `^(\w*)=((?:` + prefix + `))?(.*)$`
	var reDecryptEnv = regexp.MustCompile(decryptEnvFlagRegexp)

	for _, value := range input {
		match := reDecryptEnv.FindStringSubmatch(value)
		flag := &Key{
			Key:         match[1],
			Value:       match[3],
			IsEncrypted: match[2] == prefix,
		}
		if flag.IsEncrypted {
			k = append(k, *flag)
		}
	}
	return k
}

func (e *envsKeys) Export(keys []Key) {
	for _, k := range keys {
		fmt.Printf("export %s='%s'\n", k.Key, k.Value)
	}
}

type flagsKeys struct {
}

func (f *flagsKeys) Import(input []string, prefix string) []Key {
	var k []Key
	decryptFlagRegexp := `^\-(\w*)=((?:` + prefix + `))?(.*)$`
	var reDecryptFlag = regexp.MustCompile(decryptFlagRegexp)

	for _, value := range input {
		match := reDecryptFlag.FindStringSubmatch(value)
		flag := &Key{
			Key:         match[1],
			Value:       match[3],
			IsEncrypted: match[2] == prefix,
		}
		k = append(k, *flag)
	}
	return k
}

func (f *flagsKeys) Export(keys []Key) {
	flags := covertFlags(keys)
	fmt.Println(strings.TrimLeft(flags, " "))
}

func covertFlags(decrypt []Key) string {
	var result string
	for _, flag := range decrypt {
		result += fmt.Sprintf(" -%s=%s", flag.Key, flag.Value)
	}
	return strings.TrimLeft(result, " ")
}
