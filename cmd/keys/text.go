package keys

import (
	"regexp"
	"strings"
)

type text struct{}

func init() {
	Register("text", &text{})
}

func (im *text) Import(args, perfix string) ([]Key, error) {
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
	return k, nil
}
