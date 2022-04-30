package keys

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

var (
	keys = map[Type]Keys{
		TypeText: &flagsKeys{},
		TypeEnv:  &envsKeys{},
		TypeFile: &fileKeys{},
	}
)

type envsKeys struct {
}

func (e *envsKeys) Import(input []string, secretTag string) []Key {
	var ks []Key
	envRegexp := `^(\w+)=(.+)$`
	reEnv := regexp.MustCompile(envRegexp)
	reSecret := getSecretRegexp(secretTag)
	for _, env := range input {
		key, value, err := getKVfromInput(env, reEnv)
		if err != nil {
			log.Printf("WARN: %s\n", err)
		}
		ks = append(ks, Key{
			Key:     key,
			Value:   value,
			Secrets: reSecret.parseValueToSecrets(value),
		})
	}
	return ks
}

func (e *envsKeys) Export(keys []Key, secretTag string, writeCloser io.WriteCloser) error {
	exportStrs := ""
	reSecret := getSecretRegexp(secretTag)
	for _, k := range keys {
		newValue := reSecret.replaceSecrets(k.Value, k.Secrets)
		exportStrs += fmt.Sprintf("export %s='%s'\n", k.Key, newValue)
	}
	if _, err := writeCloser.Write([]byte(exportStrs)); err != nil {
		return err
	}
	return writeCloser.Close()
}

type flagsKeys struct {
}

func (f *flagsKeys) Import(input []string, secretTag string) []Key {
	var ks []Key
	flagRegexp := `^\-(\w+)=(.+)$`
	reFlag := regexp.MustCompile(flagRegexp)
	reSecret := getSecretRegexp(secretTag)
	for _, flag := range input {
		key, value, err := getKVfromInput(flag, reFlag)
		if err != nil {
			log.Printf("WARN: %s\n", err)
		}
		ks = append(ks, Key{
			Key:     key,
			Value:   value,
			Secrets: reSecret.parseValueToSecrets(value),
		})
	}
	return ks
}

func (f *flagsKeys) Export(keys []Key, secretTag string, writeCloser io.WriteCloser) error {
	exportFlags := ""
	reSecret := getSecretRegexp(secretTag)
	for _, k := range keys {
		newValue := reSecret.replaceSecrets(k.Value, k.Secrets)
		exportFlags += fmt.Sprintf(" -%s=%s", k.Key, newValue)
	}
	if _, err := writeCloser.Write([]byte(strings.TrimLeft(exportFlags, " ") + "\n")); err != nil {
		return err
	}
	return writeCloser.Close()
}

func parseYaml(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = parseYaml(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = parseYaml(v)
		}
	}
	return i
}

func convertYamlToMap(data []byte) map[string]interface{} {
	var body interface{}
	if err := yaml.Unmarshal(data, &body); err != nil {
		panic(err)
	}
	body = parseYaml(body)
	// Test parse to JSON
	if _, err := json.Marshal(body); err != nil {
		panic(err)
	}
	return body.(map[string]interface{})
}

type fileKeys struct {
}

func (f *fileKeys) Import(input []string, secretTag string) []Key {
	var ks []Key
	reSecret := getSecretRegexp(secretTag)
	rawData, err := ioutil.ReadFile(input[0])
	if err != nil {
		log.Panicf("ReadFile failed %s", err)
	}
	data := convertYamlToMap(rawData)
	for key, value := range data {
		switch v := value.(type) {
		case string:
			ks = append(ks, Key{
				Key:     key,
				Value:   v,
				Secrets: reSecret.parseValueToSecrets(v),
			})
		case map[interface{}]interface{}:

		case []interface{}:

		default:

		}
	}
	fmt.Println(ks)
	return ks
}

func (f *fileKeys) Export(keys []Key, secretTag string, writeCloser io.WriteCloser) error {
	exportStrs := ""
	reSecret := getSecretRegexp(secretTag)
	for _, k := range keys {
		newValue := reSecret.replaceSecrets(k.Value, k.Secrets)
		exportStrs += fmt.Sprintf("export %s='%s'\n", k.Key, newValue)
	}
	if _, err := writeCloser.Write([]byte(exportStrs)); err != nil {
		return err
	}
	return writeCloser.Close()
}

func getKVfromInput(input string, re *regexp.Regexp) (key string, value string, err error) {
	kv := re.FindStringSubmatch(input)
	if len(kv) != 3 {
		emptyRegexp := `^\-(\w+)=$`
		emptyFlag := regexp.MustCompile(emptyRegexp)
		k := emptyFlag.FindStringSubmatch(input)
		return k[1], "", fmt.Errorf(`Cannot find value for key "%s"`, k[1])
	}
	key = kv[1]
	value = kv[2]
	return key, value, nil
}

func getSecretRegexp(secretTag string) *secretRegexp {
	secretTagRegexpTemplate := `<%s>(.*?)</%s>|<%s>(.*?)</%s>`
	secretTagRegexp := fmt.Sprintf(secretTagRegexpTemplate,
		strings.ToLower(secretTag), strings.ToLower(secretTag),
		strings.ToUpper(secretTag), strings.ToUpper(secretTag),
	)
	return &secretRegexp{
		secretTag: secretTag,
		re:        regexp.MustCompile(secretTagRegexp),
	}
}

type secretRegexp struct {
	secretTag string
	re        *regexp.Regexp
}

func (s *secretRegexp) parseValueToSecrets(value string) []*Secret {
	matchedSecrets := s.re.FindAllStringSubmatch(value, -1)
	// If we don't find any tagged secrets, we regards the entire value as a plaintext secret.
	if len(matchedSecrets) == 0 {
		return []*Secret{
			{
				Text:        value,
				IsEncrypted: false,
			},
		}
	}

	// Otherwise, returns secret converted from tagged secrets.
	var secrets []*Secret
	for _, ms := range matchedSecrets {
		plaintext := ms[1]
		ciphertext := ms[2]

		text := plaintext
		isEncrypted := false
		if ciphertext != "" {
			text = ciphertext
			isEncrypted = true
		}
		secrets = append(secrets, &Secret{
			Text:        text,
			IsEncrypted: isEncrypted,
		})
	}
	return secrets
}

func (s *secretRegexp) replaceSecrets(value string, secrets []*Secret) string {
	if len(secrets) == 0 {
		return ""
	}

	if len(s.re.FindAllStringSubmatch(value, -1)) == 0 {
		return s.outputSecret(secrets[0])
	}

	i := 0
	f := func(_ string) string {
		defer func() { i++ }()
		return s.outputSecret(secrets[i])
	}
	return s.re.ReplaceAllStringFunc(value, f)
}

func (s *secretRegexp) outputSecret(secret *Secret) string {
	text := secret.Text
	if secret.IsEncrypted == false {
		return text
	}
	return fmt.Sprintf("<%s>%s</%s>", strings.ToUpper(s.secretTag), text, strings.ToUpper(s.secretTag))
}
