package keys

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Keys defines keys operations
type Keys interface {
	Import(string, string) ([]Key, error)
}

type Key struct {
	Key         string
	Value       string
	IsEncrypted bool
}

func EnvsImporter(perfix string) []Key {
	var k []Key
	decryptFlagRegexp := `^(\w*)=((?:` + perfix + `))?(.*)$`
	var reDecryptFlag = regexp.MustCompile(decryptFlagRegexp)

	for _, value := range os.Environ() {
		match := reDecryptFlag.FindStringSubmatch(value)
		flag := &Key{
			Key:         match[1],
			Value:       match[3],
			IsEncrypted: match[2] == perfix,
		}
		if flag.IsEncrypted {
			k = append(k, *flag)
		}
	}
	return k
}

func EnvsOutputer(keys []Key) {
	for _, k := range keys {
		fmt.Printf("export %s='%s'\n", k.Key, k.Value)
	}
}

func FlagsImporter(args, perfix string) []Key {
	var k []Key
	decryptFlagRegexp := `^\-(\w*)=((?:` + perfix + `))?(.*)$`
	var reDecryptFlag = regexp.MustCompile(decryptFlagRegexp)

	splitargs := strings.Split(args, " ")
	for _, value := range splitargs {
		match := reDecryptFlag.FindStringSubmatch(value)
		flag := &Key{
			Key:         match[1],
			Value:       match[3],
			IsEncrypted: match[2] == perfix,
		}
		k = append(k, *flag)
	}
	return k
}
