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
	decryptEnvFlagRegexp := `^(\w*)=((?:` + perfix + `))?(.*)$`
	var reDecryptEnv = regexp.MustCompile(decryptEnvFlagRegexp)

	for _, value := range os.Environ() {
		match := reDecryptEnv.FindStringSubmatch(value)
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

func covertFlags(decrypt []Key) string {
	var result string
	for _, flag := range decrypt {
		result += fmt.Sprintf(" -%s=%s", flag.Key, flag.Value)
	}
	return strings.TrimLeft(result, " ")
}

func EnvsExporter(keys []Key) {
	for _, k := range keys {
		fmt.Printf("export %s='%s'\n", k.Key, k.Value)
	}
}

func FlagsExporter(keys []Key) {
	flags := covertFlags(keys)
	fmt.Println(strings.TrimLeft(flags, " "))
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
